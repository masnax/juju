// Copyright 2022 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package backups

import (
	"reflect"

	"github.com/juju/juju/apiserver/facade"
)

// Register is called to expose a package of facades onto a given registry.
func Register(registry facade.FacadeRegistry) {
	registry.MustRegister("Backups", 3, func(ctx facade.Context) (facade.Facade, error) {
		return newFacade(ctx)
	}, reflect.TypeOf((*API)(nil)))
}

// newFacade provides the required signature for facade registration.
func newFacade(ctx facade.Context) (*API, error) {
	return NewAPI(
		ctx.ServiceFactory().ControllerConfig(),
		ctx.Auth(),
		ctx.MachineTag(),
		ctx.DataDir(),
		ctx.LogDir(),
	)
}
