// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package cloudspec_test

import (
	"context"
	"errors"

	"github.com/juju/names/v4"
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/utils/v3"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/apiserver/common"
	"github.com/juju/juju/apiserver/common/cloudspec"
	apiservertesting "github.com/juju/juju/apiserver/testing"
	"github.com/juju/juju/cloud"
	"github.com/juju/juju/core/watcher"
	environscloudspec "github.com/juju/juju/environs/cloudspec"
	"github.com/juju/juju/rpc/params"
	"github.com/juju/juju/state"
	coretesting "github.com/juju/juju/testing"
)

type CloudSpecSuite struct {
	testing.IsolationSuite
	testing.Stub
	result   environscloudspec.CloudSpec
	authFunc common.AuthFunc
	api      cloudspec.CloudSpecAPI
}

var _ = gc.Suite(&CloudSpecSuite{})

func (s *CloudSpecSuite) SetUpTest(c *gc.C) {
	s.IsolationSuite.SetUpTest(c)
	s.Stub.ResetCalls()

	s.authFunc = func(tag names.Tag) bool {
		s.AddCall("Auth", tag)
		return tag == coretesting.ModelTag
	}
	s.api = s.getTestCloudSpec(
		apiservertesting.NewFakeNotifyWatcher(),
		apiservertesting.NewFakeNotifyWatcher())
	credential := cloud.NewCredential(
		"auth-type",
		map[string]string{"k": "v"},
	)
	s.result = environscloudspec.CloudSpec{
		Type:             "type",
		Name:             "name",
		Region:           "region",
		Endpoint:         "endpoint",
		IdentityEndpoint: "identity-endpoint",
		StorageEndpoint:  "storage-endpoint",
		Credential:       &credential,
		CACertificates:   []string{coretesting.CACert},
		SkipTLSVerify:    true,
	}
}

func (s *CloudSpecSuite) getTestCloudSpec(cloudWatcher, credentialContentWatcher state.NotifyWatcher) cloudspec.CloudSpecAPI {
	return cloudspec.NewCloudSpec(
		common.NewResources(),
		func(tag names.ModelTag) (environscloudspec.CloudSpec, error) {
			s.AddCall("CloudSpec", tag)
			return s.result, s.NextErr()
		},
		func(_ context.Context, tag names.ModelTag) (watcher.NotifyWatcher, error) {
			s.AddCall("WatchCloudSpec", tag)
			return cloudWatcher, s.NextErr()
		},
		func(tag names.ModelTag) (state.NotifyWatcher, error) {
			s.AddCall("WatchCredentialReference", tag)
			return apiservertesting.NewFakeNotifyWatcher(), s.NextErr()
		},
		func(_ context.Context, tag names.ModelTag) (watcher.NotifyWatcher, error) {
			s.AddCall("WatchCredentialContent", tag)
			return credentialContentWatcher, s.NextErr()
		},
		func() (common.AuthFunc, error) {
			s.AddCall("GetAuthFunc")
			return s.authFunc, s.NextErr()
		})
}

func (s *CloudSpecSuite) TestCloudSpec(c *gc.C) {
	otherModelTag := names.NewModelTag(utils.MustNewUUID().String())
	machineTag := names.NewMachineTag("42")
	result, err := s.api.CloudSpec(
		context.Background(),
		params.Entities{Entities: []params.Entity{
			{coretesting.ModelTag.String()},
			{otherModelTag.String()},
			{machineTag.String()},
		}},
	)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Results, jc.DeepEquals, []params.CloudSpecResult{{
		Result: &params.CloudSpec{
			Type:             "type",
			Name:             "name",
			Region:           "region",
			Endpoint:         "endpoint",
			IdentityEndpoint: "identity-endpoint",
			StorageEndpoint:  "storage-endpoint",
			Credential: &params.CloudCredential{
				AuthType:   "auth-type",
				Attributes: map[string]string{"k": "v"},
			},
			CACertificates: []string{coretesting.CACert},
			SkipTLSVerify:  true,
		},
	}, {
		Error: &params.Error{
			Code:    params.CodeUnauthorized,
			Message: "permission denied",
		},
	}, {
		Error: &params.Error{
			Message: `"machine-42" is not a valid model tag`,
		},
	}})
	s.CheckCalls(c, []testing.StubCall{
		{"GetAuthFunc", nil},
		{"Auth", []interface{}{coretesting.ModelTag}},
		{"CloudSpec", []interface{}{coretesting.ModelTag}},
		{"Auth", []interface{}{otherModelTag}},
	})
}

func (s *CloudSpecSuite) TestWatchCloudSpecsChanges(c *gc.C) {
	otherModelTag := names.NewModelTag(utils.MustNewUUID().String())
	machineTag := names.NewMachineTag("42")
	result, err := s.api.WatchCloudSpecsChanges(
		context.Background(),
		params.Entities{Entities: []params.Entity{
			{coretesting.ModelTag.String()},
			{otherModelTag.String()},
			{machineTag.String()},
		}},
	)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Results, jc.DeepEquals, []params.NotifyWatchResult{{
		NotifyWatcherId: "1",
	}, {
		Error: &params.Error{
			Code:    params.CodeUnauthorized,
			Message: "permission denied",
		},
	}, {
		Error: &params.Error{
			Message: `"machine-42" is not a valid model tag`,
		},
	}})
	s.CheckCalls(c, []testing.StubCall{
		{"GetAuthFunc", nil},
		{"Auth", []interface{}{coretesting.ModelTag}},
		{"WatchCloudSpec", []interface{}{coretesting.ModelTag}},
		{"WatchCredentialReference", []interface{}{coretesting.ModelTag}},
		{"WatchCredentialContent", []interface{}{coretesting.ModelTag}},
		{"Auth", []interface{}{otherModelTag}},
	})
}

func (s *CloudSpecSuite) TestWatchCloudSpecsNoCredentialContentToWatch(c *gc.C) {
	s.api = s.getTestCloudSpec(apiservertesting.NewFakeNotifyWatcher(), nil)
	result, err := s.api.WatchCloudSpecsChanges(
		context.Background(),
		params.Entities{Entities: []params.Entity{
			{coretesting.ModelTag.String()},
		}},
	)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Results, jc.DeepEquals, []params.NotifyWatchResult{{
		NotifyWatcherId: "1",
	}})
	s.CheckCalls(c, []testing.StubCall{
		{"GetAuthFunc", nil},
		{"Auth", []interface{}{coretesting.ModelTag}},
		{"WatchCloudSpec", []interface{}{coretesting.ModelTag}},
		{"WatchCredentialReference", []interface{}{coretesting.ModelTag}},
		{"WatchCredentialContent", []interface{}{coretesting.ModelTag}},
	})
}

func (s *CloudSpecSuite) TestCloudSpecNilCredential(c *gc.C) {
	s.result.Credential = nil
	result, err := s.api.CloudSpec(
		context.Background(),
		params.Entities{
			Entities: []params.Entity{{coretesting.ModelTag.String()}},
		},
	)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Results, jc.DeepEquals, []params.CloudSpecResult{{
		Result: &params.CloudSpec{
			Type:             "type",
			Name:             "name",
			Region:           "region",
			Endpoint:         "endpoint",
			IdentityEndpoint: "identity-endpoint",
			StorageEndpoint:  "storage-endpoint",
			Credential:       nil,
			CACertificates:   []string{coretesting.CACert},
			SkipTLSVerify:    true,
		},
	}})
}

func (s *CloudSpecSuite) TestCloudSpecGetAuthFuncError(c *gc.C) {
	expect := errors.New("bewm")
	s.SetErrors(expect)
	result, err := s.api.CloudSpec(
		context.Background(),
		params.Entities{
			Entities: []params.Entity{{coretesting.ModelTag.String()}},
		},
	)
	c.Assert(err, gc.Equals, expect)
	c.Assert(result, jc.DeepEquals, params.CloudSpecResults{})
}

func (s *CloudSpecSuite) TestCloudSpecCloudSpecError(c *gc.C) {
	s.SetErrors(nil, errors.New("bewm"))
	result, err := s.api.CloudSpec(
		context.Background(),
		params.Entities{
			Entities: []params.Entity{{coretesting.ModelTag.String()}},
		},
	)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result, jc.DeepEquals, params.CloudSpecResults{Results: []params.CloudSpecResult{{
		Error: &params.Error{Message: "bewm"},
	}}})
}
