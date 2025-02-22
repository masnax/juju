// Copyright 2022 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package remoterelations

import (
	"reflect"

	"github.com/juju/errors"

	"github.com/juju/juju/apiserver/common"
	commoncrossmodel "github.com/juju/juju/apiserver/common/crossmodel"
	"github.com/juju/juju/apiserver/facade"
	corelogger "github.com/juju/juju/core/logger"
)

// Register is called to expose a package of facades onto a given registry.
func Register(registry facade.FacadeRegistry) {
	registry.MustRegister("RemoteRelations", 2, func(ctx facade.Context) (facade.Facade, error) {
		return newAPI(ctx) // Adds UpdateControllersForModels and WatchLocalRelationChanges.
	}, reflect.TypeOf((*API)(nil)))
}

// newAPI creates a new server-side API facade backed by global state.
func newAPI(ctx facade.Context) (*API, error) {
	systemState, err := ctx.StatePool().SystemState()
	if err != nil {
		return nil, errors.Trace(err)
	}
	serviceFactory := ctx.ServiceFactory()
	controllerConfigService := serviceFactory.ControllerConfig()
	externalControllerService := serviceFactory.ExternalController()
	return NewRemoteRelationsAPI(
		stateShim{st: ctx.State(), Backend: commoncrossmodel.GetBackend(ctx.State())},
		externalControllerService,
		common.NewControllerConfigAPI(systemState, controllerConfigService, externalControllerService),
		ctx.Resources(), ctx.Auth(),
		ctx.Logger().ChildWithLabels("remoterelations", corelogger.CMR),
	)
}
