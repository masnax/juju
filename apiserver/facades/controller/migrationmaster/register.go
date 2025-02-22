// Copyright 2022 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package migrationmaster

import (
	"reflect"

	"github.com/juju/errors"

	"github.com/juju/juju/apiserver/common/cloudspec"
	"github.com/juju/juju/apiserver/facade"
	"github.com/juju/juju/core/changestream"
	"github.com/juju/juju/core/modelmigration"
	"github.com/juju/juju/internal/migration"
)

// Register is called to expose a package of facades onto a given registry.
func Register(registry facade.FacadeRegistry) {
	registry.MustRegister("MigrationMaster", 3, func(ctx facade.Context) (facade.Facade, error) {
		return newMigrationMasterFacade(ctx) // Adds MinionReportTimeout.
	}, reflect.TypeOf((*API)(nil)))
}

// newMigrationMasterFacade exists to provide the required signature for API
// registration, converting st to backend.
func newMigrationMasterFacade(ctx facade.Context) (*API, error) {
	pool := ctx.StatePool()
	modelState := ctx.State()

	controllerState, err := pool.SystemState()
	if err != nil {
		return nil, errors.Trace(err)
	}

	preCheckBackend, err := migration.PrecheckShim(modelState, controllerState)
	if err != nil {
		return nil, errors.Annotate(err, "creating precheck backend")
	}

	leadership, err := ctx.LeadershipReader(modelState.ModelUUID())
	if err != nil {
		return nil, errors.Trace(err)
	}

	backend := newBacked(modelState)

	serviceFactory := ctx.ServiceFactory()
	credentialService := serviceFactory.Credential()
	return NewAPI(
		controllerState,
		backend,
		migration.NewModelExporter(
			backend,
			modelmigration.NewScope(changestream.NewTxnRunnerFactory(ctx.ControllerDB), nil),
		),
		ctx.ObjectStore(),
		preCheckBackend,
		migration.PoolShim(pool),
		ctx.Resources(),
		ctx.Auth(),
		ctx.Presence(),
		cloudspec.MakeCloudSpecGetter(pool, serviceFactory.Cloud(), credentialService),
		leadership,
		credentialService,
		serviceFactory.Upgrade(),
	)
}
