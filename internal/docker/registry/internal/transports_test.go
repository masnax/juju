// Copyright 2021 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package internal_test

import (
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/juju/errors"
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	"go.uber.org/mock/gomock"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/internal/docker/registry/internal"
	"github.com/juju/juju/internal/docker/registry/mocks"
)

type transportSuite struct {
	testing.IsolationSuite
}

var _ = gc.Suite(&transportSuite{})

func (s *transportSuite) TestErrorTransport(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()
	mockRoundTripper := mocks.NewMockRoundTripper(ctrl)

	url, err := url.Parse(`https://example.com`)
	c.Assert(err, jc.ErrorIsNil)

	mockRoundTripper.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
		resps := &http.Response{
			Request:    req,
			StatusCode: http.StatusForbidden,
			Body:       io.NopCloser(strings.NewReader(`invalid input`)),
		}
		return resps, nil
	})
	t := internal.NewErrorTransport(mockRoundTripper)
	_, err = t.RoundTrip(&http.Request{URL: url})
	c.Assert(err, gc.ErrorMatches, `non-successful response status=403`)
}

func (s *transportSuite) TestBasicTransport(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()
	mockRoundTripper := mocks.NewMockRoundTripper(ctrl)

	url, err := url.Parse(`https://example.com`)
	c.Assert(err, jc.ErrorIsNil)

	// username + password.
	mockRoundTripper.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(
		func(req *http.Request) (*http.Response, error) {
			c.Assert(req.Header, jc.DeepEquals, http.Header{"Authorization": []string{"Basic " + base64.StdEncoding.EncodeToString([]byte("username:pwd"))}})
			return &http.Response{
				Request:    req,
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(``)),
			}, nil
		},
	)
	t := internal.NewBasicTransport(mockRoundTripper, "username", "pwd", "")
	_, err = t.RoundTrip(&http.Request{
		Header: http.Header{},
		URL:    url,
	})
	c.Assert(err, jc.ErrorIsNil)

	// auth token.
	mockRoundTripper.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(
		func(req *http.Request) (*http.Response, error) {
			c.Assert(req.Header, jc.DeepEquals, http.Header{"Authorization": []string{"Basic " + `dXNlcm5hbWU6cHdkMQ==`}})
			return &http.Response{
				Request:    req,
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(``)),
			}, nil
		},
	)
	t = internal.NewBasicTransport(mockRoundTripper, "", "", "dXNlcm5hbWU6cHdkMQ==")
	_, err = t.RoundTrip(&http.Request{
		Header: http.Header{},
		URL:    url,
	})
	c.Assert(err, jc.ErrorIsNil)

	// no credentials.
	mockRoundTripper.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(
		func(req *http.Request) (*http.Response, error) {
			c.Assert(req.Header, jc.DeepEquals, http.Header{})
			return &http.Response{
				Request:    req,
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(``)),
			}, nil
		},
	)
	t = internal.NewBasicTransport(mockRoundTripper, "", "", "")
	_, err = t.RoundTrip(&http.Request{
		Header: http.Header{},
		URL:    url,
	})
	c.Assert(err, jc.ErrorIsNil)
}

func (s *transportSuite) TestTokenTransportOAuthTokenProvided(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()
	mockRoundTripper := mocks.NewMockRoundTripper(ctrl)

	url, err := url.Parse(`https://example.com`)
	c.Assert(err, jc.ErrorIsNil)

	gomock.InOrder(
		mockRoundTripper.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(
			func(req *http.Request) (*http.Response, error) {
				c.Assert(req.Header, jc.DeepEquals, http.Header{"Authorization": []string{"Bearer " + `OAuth-jwt-token`}})
				c.Assert(req.URL.String(), gc.Equals, `https://example.com`)
				return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(nil)}, nil
			},
		),
	)
	t := internal.NewTokenTransport(mockRoundTripper, "", "", "", "OAuth-jwt-token", false)
	_, err = t.RoundTrip(&http.Request{
		Header: http.Header{},
		URL:    url,
	})
	c.Assert(err, jc.ErrorIsNil)
}

func (s *transportSuite) TestTokenTransportTokenRefresh(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()
	mockRoundTripper := mocks.NewMockRoundTripper(ctrl)

	url, err := url.Parse(`https://example.com`)
	c.Assert(err, jc.ErrorIsNil)

	gomock.InOrder(
		// 1st try failed - bearer token was missing.
		mockRoundTripper.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(
			func(req *http.Request) (*http.Response, error) {
				c.Assert(req.Header, jc.DeepEquals, http.Header{})
				c.Assert(req.URL.String(), gc.Equals, `https://example.com`)
				return &http.Response{
					Request:    req,
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(nil),
					Header: http.Header{
						http.CanonicalHeaderKey("WWW-Authenticate"): []string{
							`Bearer realm="https://auth.example.com/token",service="registry.example.com",scope="repository:jujuqa/jujud-operator:pull"`,
						},
					},
				}, nil
			},
		),
		// Refresh OAuth Token.
		mockRoundTripper.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(
			func(req *http.Request) (*http.Response, error) {
				c.Assert(req.Header, jc.DeepEquals, http.Header{"Authorization": []string{"Basic " + `dXNlcm5hbWU6cHdkMQ==`}})
				c.Assert(req.URL.String(), gc.Equals, `https://auth.example.com/token?scope=repository%3Ajujuqa%2Fjujud-operator%3Apull&service=registry.example.com`)
				return &http.Response{
					Request:    req,
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"token": "OAuth-jwt-token", "access_token": "OAuth-jwt-token","expires_in": 300}`)),
				}, nil
			},
		),
		// retry.
		mockRoundTripper.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(
			func(req *http.Request) (*http.Response, error) {
				c.Assert(req.Header, jc.DeepEquals, http.Header{"Authorization": []string{"Bearer " + `OAuth-jwt-token`}})
				c.Assert(req.URL.String(), gc.Equals, `https://example.com`)
				return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(nil)}, nil
			},
		),
	)
	t := internal.NewTokenTransport(mockRoundTripper, "", "", "dXNlcm5hbWU6cHdkMQ==", "", false)
	_, err = t.RoundTrip(&http.Request{
		Header: http.Header{},
		URL:    url,
	})
	c.Assert(err, jc.ErrorIsNil)
}

func (s *transportSuite) TestTokenTransportTokenRefreshFailedRealmMissing(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()
	mockRoundTripper := mocks.NewMockRoundTripper(ctrl)

	url, err := url.Parse(`https://example.com`)
	c.Assert(err, jc.ErrorIsNil)

	gomock.InOrder(
		mockRoundTripper.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(
			func(req *http.Request) (*http.Response, error) {
				c.Assert(req.Header, jc.DeepEquals, http.Header{})
				c.Assert(req.URL.String(), gc.Equals, `https://example.com`)
				return &http.Response{
					Request:    req,
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(nil),
					Header: http.Header{
						http.CanonicalHeaderKey("WWW-Authenticate"): []string{
							`Bearer service="registry.example.com",scope="repository:jujuqa/jujud-operator:pull"`,
						},
					},
				}, nil
			},
		),
	)
	t := internal.NewTokenTransport(mockRoundTripper, "", "", "dXNlcm5hbWU6cHdkMQ==", "", false)
	_, err = t.RoundTrip(&http.Request{
		Header: http.Header{},
		URL:    url,
	})
	c.Assert(err, gc.ErrorMatches, `refreshing OAuth token: no realm specified for token auth challenge`)
}

func (s *transportSuite) TestTokenTransportTokenRefreshFailedServiceMissing(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()
	mockRoundTripper := mocks.NewMockRoundTripper(ctrl)

	url, err := url.Parse(`https://example.com`)
	c.Assert(err, jc.ErrorIsNil)

	gomock.InOrder(
		mockRoundTripper.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(
			func(req *http.Request) (*http.Response, error) {
				c.Assert(req.Header, jc.DeepEquals, http.Header{})
				c.Assert(req.URL.String(), gc.Equals, `https://example.com`)
				return &http.Response{
					Request:    req,
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(nil),
					Header: http.Header{
						http.CanonicalHeaderKey("WWW-Authenticate"): []string{
							`Bearer realm="https://auth.example.com/token",scope="repository:jujuqa/jujud-operator:pull"`,
						},
					},
				}, nil
			},
		),
	)
	t := internal.NewTokenTransport(mockRoundTripper, "", "", "dXNlcm5hbWU6cHdkMQ==", "", false)
	_, err = t.RoundTrip(&http.Request{
		Header: http.Header{},
		URL:    url,
	})
	c.Assert(err, gc.ErrorMatches, `refreshing OAuth token: no service specified for token auth challenge`)
}

func (s *transportSuite) TestUnwrapNetError(c *gc.C) {
	originalErr := errors.NotFoundf("jujud-operator:2.6.6")
	c.Assert(originalErr, jc.ErrorIs, errors.NotFound)
	var urlErr error = &url.Error{
		Op:  "Get",
		URL: "https://example.com",
		Err: originalErr,
	}
	unwrapedErr := internal.UnwrapNetError(urlErr)
	c.Assert(unwrapedErr, gc.NotNil)
	c.Assert(unwrapedErr, jc.ErrorIs, errors.NotFound)
	c.Assert(unwrapedErr, gc.ErrorMatches, `Get "https://example.com": jujud-operator:2.6.6 not found`)
}
