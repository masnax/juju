// Copyright 2012-2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package openstack_test

import (
	"testing"

	gc "gopkg.in/check.v1"
)

//go:generate go run go.uber.org/mock/mockgen -package openstack -destination network_mock_test.go github.com/juju/juju/provider/openstack SSLHostnameConfig,Networking,NetworkingBase,NetworkingNeutron,NetworkingAuthenticatingClient,NetworkingNova,NetworkingEnvironConfig
//go:generate go run go.uber.org/mock/mockgen -package openstack -destination cloud_mock_test.go github.com/juju/juju/cloudconfig/cloudinit NetworkingConfig
//go:generate go run go.uber.org/mock/mockgen -package openstack -destination goose_mock_test.go github.com/go-goose/goose/v5/client AuthenticatingClient

func Test(t *testing.T) {
	gc.TestingT(t)
}
