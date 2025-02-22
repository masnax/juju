// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// This package exists solely to avoid circular imports.

package factory

import (
	"github.com/juju/errors"

	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/internal/container"
	"github.com/juju/juju/internal/container/kvm"
	"github.com/juju/juju/internal/container/lxd"
)

// NewContainerManager creates the appropriate container.Manager for the
// specified container type.
var NewContainerManager = func(forType instance.ContainerType, conf container.ManagerConfig) (container.Manager, error) {
	switch forType {
	case instance.LXD:
		return lxd.NewContainerManager(conf, lxd.NewLocalServer)
	case instance.KVM:
		return kvm.NewContainerManager(conf)
	}
	return nil, errors.Errorf("unknown container type: %q", forType)
}
