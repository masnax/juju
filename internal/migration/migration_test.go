// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package migration_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/juju/charm/v11"
	"github.com/juju/description/v4"
	"github.com/juju/errors"
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/version/v2"
	"go.uber.org/mock/gomock"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/controller"
	coremigration "github.com/juju/juju/core/migration"
	"github.com/juju/juju/core/modelmigration"
	"github.com/juju/juju/core/resources"
	resourcetesting "github.com/juju/juju/core/resources/testing"
	"github.com/juju/juju/internal/migration"
	"github.com/juju/juju/internal/tools"
	"github.com/juju/juju/state"
	jujutesting "github.com/juju/juju/testing"
)

type ImportSuite struct {
	testing.IsolationSuite

	controllerConfigService *MockControllerConfigService
}

var _ = gc.Suite(&ImportSuite{})

func (s *ImportSuite) SetUpTest(c *gc.C) {
	s.IsolationSuite.SetUpTest(c)
}

func (s *ImportSuite) TestBadBytes(c *gc.C) {
	defer s.setupMocks(c).Finish()

	bytes := []byte("not a model")
	scope := modelmigration.NewScope(nil, nil)
	controller := &fakeImporter{}
	importer := migration.NewModelImporter(controller, scope, s.controllerConfigService)
	model, st, err := importer.ImportModel(context.Background(), bytes)
	c.Check(st, gc.IsNil)
	c.Check(model, gc.IsNil)
	c.Assert(err, gc.ErrorMatches, "yaml: unmarshal errors:\n.*")
}

const model = `
cloud: dev
config:
  name: foo
  type: lxd
  uuid: bd3fae18-5ea1-4bc5-8837-45400cf1f8f6
actions:
  actions: []
  version: 1
applications:
  applications: []
  version: 1
cloud-image-metadata:
  cloudimagemetadata: []
  version: 1
filesystems:
  filesystems: []
  version: 1
ip-addresses:
  ip-addresses: []
  version: 1
link-layer-devices:
  link-layer-devices: []
  version: 1
machines:
  machines: []
  version: 1
owner: admin
relations:
  relations: []
  version: 1
sequences:
  machine: 2
spaces:
  spaces: []
  version: 1
ssh-host-keys:
  ssh-host-keys: []
  version: 1
storage-pools:
  pools: []
  version: 1
storages:
  storages: []
  version: 1
subnets:
  subnets: []
  version: 1
users:
  users: []
  version: 1
volumes:
  volumes: []
  version: 1
version: 1
`

func (s *ImportSuite) exportImport(c *gc.C, leaders map[string]string) {
	bytes := []byte(model)
	st := &state.State{}
	m := &state.Model{}
	controller := &fakeImporter{st: st, m: m}
	scope := modelmigration.NewScope(nil, nil)
	importer := migration.NewModelImporter(controller, scope, s.controllerConfigService)
	gotM, gotSt, err := importer.ImportModel(context.Background(), bytes)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(controller.model.Tag().Id(), gc.Equals, "bd3fae18-5ea1-4bc5-8837-45400cf1f8f6")
	c.Assert(gotM, gc.Equals, m)
	c.Assert(gotSt, gc.Equals, st)
}

func (s *ImportSuite) TestImportModel(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.exportImport(c, map[string]string{})
}

func (s *ImportSuite) TestUploadBinariesConfigValidate(c *gc.C) {
	type T migration.UploadBinariesConfig // alias for brevity

	check := func(modify func(*T), missing string) {
		config := T{
			CharmDownloader:    struct{ migration.CharmDownloader }{},
			CharmUploader:      struct{ migration.CharmUploader }{},
			ToolsDownloader:    struct{ migration.ToolsDownloader }{},
			ToolsUploader:      struct{ migration.ToolsUploader }{},
			ResourceDownloader: struct{ migration.ResourceDownloader }{},
			ResourceUploader:   struct{ migration.ResourceUploader }{},
		}
		modify(&config)
		realConfig := migration.UploadBinariesConfig(config)
		c.Check(realConfig.Validate(), gc.ErrorMatches, fmt.Sprintf("missing %s not valid", missing))
	}

	check(func(c *T) { c.CharmDownloader = nil }, "CharmDownloader")
	check(func(c *T) { c.CharmUploader = nil }, "CharmUploader")
	check(func(c *T) { c.ToolsDownloader = nil }, "ToolsDownloader")
	check(func(c *T) { c.ToolsUploader = nil }, "ToolsUploader")
	check(func(c *T) { c.ResourceDownloader = nil }, "ResourceDownloader")
	check(func(c *T) { c.ResourceUploader = nil }, "ResourceUploader")
}

func (s *ImportSuite) TestBinariesMigration(c *gc.C) {
	downloader := &fakeDownloader{}
	uploader := &fakeUploader{
		tools:     make(map[version.Binary]string),
		resources: make(map[string]string),
	}

	toolsMap := map[version.Binary]string{
		version.MustParseBinary("2.1.0-ubuntu-amd64"): "/tools/0",
		version.MustParseBinary("2.0.0-ubuntu-amd64"): "/tools/1",
	}

	app0Res := resourcetesting.NewResource(c, nil, "blob0", "app0", "blob0").Resource
	app1Res := resourcetesting.NewResource(c, nil, "blob1", "app1", "blob1").Resource
	app1UnitRes := app1Res
	app1UnitRes.Revision = 1
	app2Res := resourcetesting.NewPlaceholderResource(c, "blob2", "app2")
	resources := []coremigration.SerializedModelResource{
		{ApplicationRevision: app0Res},
		{
			ApplicationRevision: app1Res,
			UnitRevisions:       map[string]resources.Resource{"app1/99": app1UnitRes},
		},
		{ApplicationRevision: app2Res},
	}

	config := migration.UploadBinariesConfig{
		Charms: []string{
			// These 2 are out of order. Rev 2 must be uploaded first.
			"local:trusty/magic-10",
			"local:trusty/magic-2",
			"ch:trusty/postgresql-42",
		},
		CharmDownloader:    downloader,
		CharmUploader:      uploader,
		Tools:              toolsMap,
		ToolsDownloader:    downloader,
		ToolsUploader:      uploader,
		Resources:          resources,
		ResourceDownloader: downloader,
		ResourceUploader:   uploader,
	}
	err := migration.UploadBinaries(config)
	c.Assert(err, jc.ErrorIsNil)

	expectedCharms := []string{
		// Note ordering.
		"ch:trusty/postgresql-42",
		"local:trusty/magic-2",
		"local:trusty/magic-10",
	}
	c.Assert(downloader.charms, jc.DeepEquals, expectedCharms)
	c.Assert(uploader.charms, jc.DeepEquals, expectedCharms)

	c.Assert(downloader.uris, jc.SameContents, []string{
		"/tools/0",
		"/tools/1",
	})
	c.Assert(uploader.tools, jc.DeepEquals, toolsMap)

	c.Assert(downloader.resources, jc.SameContents, []string{
		"app0/blob0",
		"app1/blob1",
	})
	c.Assert(uploader.resources, jc.DeepEquals, map[string]string{
		"app0/blob0": "blob0",
		"app1/blob1": "blob1",
	})
	c.Assert(uploader.unitResources, jc.SameContents, []string{"app1/99-blob1"})
}

func (s *ImportSuite) TestWrongCharmURLAssigned(c *gc.C) {
	downloader := &fakeDownloader{}
	uploader := &fakeUploader{
		reassignCharmURL: true,
	}

	config := migration.UploadBinariesConfig{
		Charms:             []string{"local:foo/bar-2"},
		CharmDownloader:    downloader,
		CharmUploader:      uploader,
		ToolsDownloader:    downloader,
		ToolsUploader:      uploader,
		ResourceDownloader: downloader,
		ResourceUploader:   uploader,
	}
	err := migration.UploadBinaries(config)
	c.Assert(err, gc.ErrorMatches,
		"cannot upload charms: charm local:foo/bar-2 unexpectedly assigned local:foo/bar-1")
}

func (s *ImportSuite) setupMocks(c *gc.C) *gomock.Controller {
	ctrl := gomock.NewController(c)

	s.controllerConfigService = NewMockControllerConfigService(ctrl)
	s.controllerConfigService.EXPECT().ControllerConfig(gomock.Any()).Return(jujutesting.FakeControllerConfig(), nil).AnyTimes()

	return ctrl
}

type fakeImporter struct {
	model            description.Model
	st               *state.State
	m                *state.Model
	controllerConfig controller.Config
}

func (i *fakeImporter) Import(model description.Model, controllerConfig controller.Config) (*state.Model, *state.State, error) {
	i.model = model
	i.controllerConfig = controllerConfig
	return i.m, i.st, nil
}

type fakeDownloader struct {
	charms    []string
	uris      []string
	resources []string
}

func (d *fakeDownloader) OpenCharm(curl string) (io.ReadCloser, error) {
	d.charms = append(d.charms, curl)
	// Return the charm URL string as the fake charm content
	return io.NopCloser(bytes.NewReader([]byte(curl + " content"))), nil
}

func (d *fakeDownloader) OpenURI(uri string, query url.Values) (io.ReadCloser, error) {
	if query != nil {
		panic("query should be empty")
	}
	d.uris = append(d.uris, uri)
	// Return the URI string as fake content
	return io.NopCloser(bytes.NewReader([]byte(uri))), nil
}

func (d *fakeDownloader) OpenResource(app, name string) (io.ReadCloser, error) {
	d.resources = append(d.resources, app+"/"+name)
	// Use the resource name as the content.
	return io.NopCloser(bytes.NewReader([]byte(name))), nil
}

type fakeUploader struct {
	tools            map[version.Binary]string
	charms           []string
	resources        map[string]string
	unitResources    []string
	reassignCharmURL bool
}

func (f *fakeUploader) UploadTools(r io.ReadSeeker, v version.Binary) (tools.List, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, errors.Trace(err)
	}
	f.tools[v] = string(data)
	return tools.List{&tools.Tools{Version: v}}, nil
}

func (f *fakeUploader) UploadCharm(u *charm.URL, r io.ReadSeeker) (*charm.URL, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if string(data) != u.String()+" content" {
		panic(fmt.Sprintf("unexpected charm body for %s: %s", u.String(), data))
	}
	f.charms = append(f.charms, u.String())

	outU := *u
	if f.reassignCharmURL {
		outU.Revision--
	}
	return &outU, nil
}

func (f *fakeUploader) UploadResource(res resources.Resource, r io.ReadSeeker) error {
	body, err := io.ReadAll(r)
	if err != nil {
		return errors.Trace(err)
	}
	f.resources[res.ApplicationID+"/"+res.Name] = string(body)
	return nil
}

func (f *fakeUploader) SetPlaceholderResource(res resources.Resource) error {
	f.resources[res.ApplicationID+"/"+res.Name] = "<placeholder>"
	return nil
}

func (f *fakeUploader) SetUnitResource(unit string, res resources.Resource) error {
	f.unitResources = append(f.unitResources, unit+"-"+res.Name)
	return nil
}
