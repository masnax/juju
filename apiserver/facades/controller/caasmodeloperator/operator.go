// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package caasmodeloperator

import (
	"context"
	"fmt"

	"github.com/juju/errors"
	"github.com/juju/loggo"
	"github.com/juju/names/v4"

	"github.com/juju/juju/apiserver/common"
	apiservererrors "github.com/juju/juju/apiserver/errors"
	"github.com/juju/juju/apiserver/facade"
	"github.com/juju/juju/cloudconfig/podcfg"
	"github.com/juju/juju/controller"
	"github.com/juju/juju/internal/docker"
	"github.com/juju/juju/rpc/params"
	"github.com/juju/juju/state/watcher"
)

// TODO (manadart 2020-10-21): Remove the ModelUUID method
// from the next version of this facade.

// API represents the controller model operator facade.
type API struct {
	*common.APIAddresser
	*common.PasswordChanger

	auth      facade.Authorizer
	ctrlState CAASControllerState
	state     CAASModelOperatorState
	logger    loggo.Logger

	resources facade.Resources
}

// NewAPI is alternative means of constructing a controller model facade.
func NewAPI(
	authorizer facade.Authorizer,
	resources facade.Resources,
	ctrlSt CAASControllerState,
	st CAASModelOperatorState,
	logger loggo.Logger,
) (*API, error) {

	if !authorizer.AuthController() {
		return nil, apiservererrors.ErrPerm
	}

	return &API{
		auth:            authorizer,
		APIAddresser:    common.NewAPIAddresser(ctrlSt, resources),
		PasswordChanger: common.NewPasswordChanger(st, common.AuthFuncForTagKind(names.ModelTagKind)),
		ctrlState:       ctrlSt,
		state:           st,
		logger:          logger,
		resources:       resources,
	}, nil
}

// WatchModelOperatorProvisioningInfo provides a watcher for changes that affect the
// information returned by ModelOperatorProvisioningInfo.
func (a *API) WatchModelOperatorProvisioningInfo() (params.NotifyWatchResult, error) {
	result := params.NotifyWatchResult{}

	model, err := a.state.Model()
	if err != nil {
		return result, errors.Trace(err)
	}

	controllerConfigWatcher := a.ctrlState.WatchControllerConfig()
	controllerAPIHostPortsWatcher := a.ctrlState.WatchAPIHostPortsForAgents()
	modelConfigWatcher := model.WatchForModelConfigChanges()

	multiWatcher := common.NewMultiNotifyWatcher(controllerConfigWatcher, controllerAPIHostPortsWatcher, modelConfigWatcher)

	if _, ok := <-multiWatcher.Changes(); ok {
		result.NotifyWatcherId = a.resources.Register(multiWatcher)
	} else {
		return result, watcher.EnsureErr(multiWatcher)
	}

	return result, nil
}

// ModelOperatorProvisioningInfo returns the information needed for provisioning
// a new model operator into a caas cluster.
func (a *API) ModelOperatorProvisioningInfo(ctx context.Context) (params.ModelOperatorInfo, error) {
	var result params.ModelOperatorInfo
	controllerConf, err := a.ctrlState.ControllerConfig()
	if err != nil {
		return result, err
	}

	model, err := a.state.Model()
	if err != nil {
		return result, errors.Trace(err)
	}
	modelConfig, err := model.ModelConfig(ctx)
	if err != nil {
		return result, errors.Trace(err)
	}

	vers, ok := modelConfig.AgentVersion()
	if !ok {
		return result, errors.NewNotValid(nil,
			fmt.Sprintf("agent version is missing in the model config %q",
				modelConfig.Name()))
	}

	apiAddresses, err := a.APIAddresses(context.Background())
	if err != nil && apiAddresses.Error != nil {
		err = apiAddresses.Error
	}
	if err != nil {
		return result, errors.Annotate(err, "getting api addresses")
	}

	registryPath, err := podcfg.GetJujuOCIImagePath(controllerConf, vers)
	if err != nil {
		return result, errors.Trace(err)
	}

	imageRepoDetails, err := docker.NewImageRepoDetails(controllerConf.CAASImageRepo())
	if err != nil {
		return result, errors.Annotatef(err, "parsing %s", controller.CAASImageRepo)
	}
	imageInfo := params.NewDockerImageInfo(docker.ConvertToResourceImageDetails(imageRepoDetails), registryPath)
	a.logger.Tracef("image info %v", imageInfo)

	result = params.ModelOperatorInfo{
		APIAddresses: apiAddresses.Result,
		ImageDetails: imageInfo,
		Version:      vers,
	}
	return result, nil
}

// ModelUUID returns the model UUID that this facade is used to operate.
// It is implemented here directly as a result of removing it from
// embedded APIAddresser *without* bumping the facade version.
// It should be blanked when this facade version is next incremented.
func (a *API) ModelUUID(ctx context.Context) params.StringResult {
	return params.StringResult{Result: a.state.ModelUUID()}
}

// APIHostPorts returns the API server addresses.
func (u *API) APIHostPorts(ctx context.Context) (result params.APIHostPortsResult, err error) {
	controllerConfig, err := u.ctrlState.ControllerConfig()
	if err != nil {
		return result, errors.Trace(err)
	}

	return u.APIAddresser.APIHostPorts(ctx, controllerConfig)
}

// APIAddresses returns the list of addresses used to connect to the API.
func (u *API) APIAddresses(ctx context.Context) (result params.StringsResult, err error) {
	controllerConfig, err := u.ctrlState.ControllerConfig()
	if err != nil {
		return result, errors.Trace(err)
	}

	return u.APIAddresser.APIAddresses(ctx, controllerConfig)
}
