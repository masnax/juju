// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package modelmigration

import (
	"context"
	"regexp"

	"github.com/juju/description/v4"
	"github.com/juju/errors"
	"github.com/juju/names/v4"
	jc "github.com/juju/testing/checkers"
	"go.uber.org/mock/gomock"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/cloud"
	"github.com/juju/juju/domain/credential"
)

type importSuite struct {
	coordinator *MockCoordinator
	service     *MockImportService
}

var _ = gc.Suite(&importSuite{})

func (s *importSuite) setupMocks(c *gc.C) *gomock.Controller {
	ctrl := gomock.NewController(c)

	s.coordinator = NewMockCoordinator(ctrl)
	s.service = NewMockImportService(ctrl)

	return ctrl
}

func (s *importSuite) newImportOperation() *importOperation {
	return &importOperation{
		service: s.service,
	}
}

func (s *importSuite) TestRegisterImport(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.coordinator.EXPECT().Add(gomock.Any())

	RegisterImport(s.coordinator)
}

func (s *importSuite) TestEmptyCredential(c *gc.C) {
	defer s.setupMocks(c).Finish()

	// Empty model.
	model := description.NewModel(description.ModelArgs{})

	op := s.newImportOperation()
	err := op.Execute(context.Background(), model)
	c.Assert(err, jc.ErrorIsNil)
	// No import executed.
	s.service.EXPECT().UpdateCloudCredential(gomock.All(), gomock.Any(), gomock.Any()).Times(0)
}

func (s *importSuite) TestImport(c *gc.C) {
	defer s.setupMocks(c).Finish()

	// Model with 2 external controllers.
	model := description.NewModel(description.ModelArgs{})
	model.SetCloudCredential(
		description.CloudCredentialArgs{
			Owner:      names.NewUserTag("fred"),
			Cloud:      names.NewCloudTag("cirrus"),
			Name:       "foo",
			AuthType:   string(cloud.UserPassAuthType),
			Attributes: map[string]string{"hello": "world"},
		},
	)
	cred := cloud.NewCredential(cloud.UserPassAuthType, map[string]string{"hello": "world"})
	id := credential.ID{Cloud: "cirrus", Owner: "fred", Name: "foo"}
	s.service.EXPECT().CloudCredential(gomock.All(), id).Times(1).Return(cloud.Credential{}, errors.NotFound)
	s.service.EXPECT().UpdateCloudCredential(gomock.Any(), id, cred).Times(1)

	op := s.newImportOperation()
	err := op.Execute(context.Background(), model)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *importSuite) TestImportExistingMatches(c *gc.C) {
	defer s.setupMocks(c).Finish()

	// Model with 2 external controllers.
	model := description.NewModel(description.ModelArgs{})
	model.SetCloudCredential(
		description.CloudCredentialArgs{
			Owner:      names.NewUserTag("fred"),
			Cloud:      names.NewCloudTag("cirrus"),
			Name:       "foo",
			AuthType:   string(cloud.UserPassAuthType),
			Attributes: map[string]string{"hello": "world"},
		},
	)
	cred := cloud.NewCredential(cloud.UserPassAuthType, map[string]string{"hello": "world"})
	id := credential.ID{Cloud: "cirrus", Owner: "fred", Name: "foo"}
	s.service.EXPECT().CloudCredential(gomock.All(), id).Times(1).Return(cred, nil)

	op := s.newImportOperation()
	err := op.Execute(context.Background(), model)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *importSuite) TestImportExistingAuthTypeMisMatch(c *gc.C) {
	defer s.setupMocks(c).Finish()

	// Model with 2 external controllers.
	model := description.NewModel(description.ModelArgs{})
	model.SetCloudCredential(
		description.CloudCredentialArgs{
			Owner:      names.NewUserTag("fred"),
			Cloud:      names.NewCloudTag("cirrus"),
			Name:       "foo",
			AuthType:   string(cloud.UserPassAuthType),
			Attributes: map[string]string{"hello": "world"},
		},
	)
	cred := cloud.NewCredential(cloud.AccessKeyAuthType, map[string]string{"hello": "world"})
	id := credential.ID{Cloud: "cirrus", Owner: "fred", Name: "foo"}
	s.service.EXPECT().CloudCredential(gomock.All(), id).Times(1).Return(cred, nil)

	op := s.newImportOperation()
	err := op.Execute(context.Background(), model)
	c.Assert(err, gc.ErrorMatches, `credential auth type mismatch: "access-key" != "userpass"`)
}

func (s *importSuite) TestImportExistingAttributesMisMatch(c *gc.C) {
	defer s.setupMocks(c).Finish()

	// Model with 2 external controllers.
	model := description.NewModel(description.ModelArgs{})
	model.SetCloudCredential(
		description.CloudCredentialArgs{
			Owner:      names.NewUserTag("fred"),
			Cloud:      names.NewCloudTag("cirrus"),
			Name:       "foo",
			AuthType:   string(cloud.UserPassAuthType),
			Attributes: map[string]string{"hello": "world"},
		},
	)
	cred := cloud.NewCredential(cloud.UserPassAuthType, map[string]string{"goodbye": "world"})
	id := credential.ID{Cloud: "cirrus", Owner: "fred", Name: "foo"}
	s.service.EXPECT().CloudCredential(gomock.All(), id).Times(1).Return(cred, nil)

	op := s.newImportOperation()
	err := op.Execute(context.Background(), model)
	c.Assert(err, gc.ErrorMatches, regexp.QuoteMeta(`credential attribute mismatch: map[goodbye:world] != map[hello:world]`))
}
