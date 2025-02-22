// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package secrets

import (
	"context"

	"github.com/juju/errors"
	"github.com/juju/loggo"
	"github.com/juju/names/v4"

	apiservererrors "github.com/juju/juju/apiserver/errors"
	"github.com/juju/juju/apiserver/facade"
	"github.com/juju/juju/apiserver/internal"
	"github.com/juju/juju/core/leadership"
	coresecrets "github.com/juju/juju/core/secrets"
	secretsprovider "github.com/juju/juju/internal/secrets/provider"
	"github.com/juju/juju/rpc/params"
	"github.com/juju/juju/state"
)

// SecretsDrainAPI is the implementation for the SecretsDrain facade.
type SecretsDrainAPI struct {
	authTag           names.Tag
	logger            loggo.Logger
	leadershipChecker leadership.Checker
	watcherRegistry   facade.WatcherRegistry

	model           Model
	secretsState    SecretsMetaState
	secretsConsumer SecretsConsumer
}

// NewSecretsDrainAPI returns a new SecretsDrainAPI.
func NewSecretsDrainAPI(
	authTag names.Tag,
	authorizer facade.Authorizer,
	logger loggo.Logger,
	leadershipChecker leadership.Checker,
	model Model,
	secretsState SecretsMetaState,
	secretsConsumer SecretsConsumer,
	watcherRegistry facade.WatcherRegistry,
) (*SecretsDrainAPI, error) {
	if !authorizer.AuthUnitAgent() && !authorizer.AuthApplicationAgent() && !authorizer.AuthController() {
		return nil, apiservererrors.ErrPerm
	}
	return &SecretsDrainAPI{
		authTag:           authTag,
		logger:            logger,
		leadershipChecker: leadershipChecker,
		model:             model,
		secretsState:      secretsState,
		secretsConsumer:   secretsConsumer,
		watcherRegistry:   watcherRegistry,
	}, nil
}

// GetSecretsToDrain returns metadata for the secrets that need to be drained.
func (s *SecretsDrainAPI) GetSecretsToDrain(ctx context.Context) (params.ListSecretResults, error) {
	modelConfig, err := s.model.ModelConfig(ctx)
	if err != nil {
		return params.ListSecretResults{}, errors.Trace(err)
	}
	modelType := s.model.Type()
	modelUUID := s.model.UUID()
	controllerUUID := s.model.ControllerUUID()

	activeBackend := modelConfig.SecretBackend()
	if activeBackend == secretsprovider.Auto {
		activeBackend = controllerUUID
		if modelType == state.ModelTypeCAAS {
			activeBackend = modelUUID
		}
	}
	return GetSecretMetadata(
		s.authTag, s.secretsState, s.leadershipChecker,
		func(md *coresecrets.SecretMetadata, rev *coresecrets.SecretRevisionMetadata) bool {
			if rev.ValueRef == nil {
				// Only internal backend secrets have nil ValueRef.
				if activeBackend == secretsprovider.Internal {
					return false
				}
				if activeBackend == controllerUUID {
					return false
				}
				return true
			}
			return rev.ValueRef.BackendID != activeBackend
		},
	)
}

// ChangeSecretBackend updates the backend for the specified secret after migration done.
func (s *SecretsDrainAPI) ChangeSecretBackend(ctx context.Context, args params.ChangeSecretBackendArgs) (params.ErrorResults, error) {
	result := params.ErrorResults{
		Results: make([]params.ErrorResult, len(args.Args)),
	}
	for i, arg := range args.Args {
		err := s.changeSecretBackendForOne(arg)
		result.Results[i].Error = apiservererrors.ServerError(err)
	}
	return result, nil
}

func (s *SecretsDrainAPI) changeSecretBackendForOne(arg params.ChangeSecretBackendArg) (err error) {
	uri, err := coresecrets.ParseURI(arg.URI)
	if err != nil {
		return errors.Trace(err)
	}
	token, err := CanManage(s.secretsConsumer, s.leadershipChecker, s.authTag, uri)
	if err != nil {
		return errors.Trace(err)
	}
	return s.secretsState.ChangeSecretBackend(toChangeSecretBackendParams(token, uri, arg))
}

func toChangeSecretBackendParams(token leadership.Token, uri *coresecrets.URI, arg params.ChangeSecretBackendArg) state.ChangeSecretBackendParams {
	params := state.ChangeSecretBackendParams{
		Token:    token,
		URI:      uri,
		Revision: arg.Revision,
		Data:     arg.Content.Data,
	}
	if arg.Content.ValueRef != nil {
		params.ValueRef = &coresecrets.ValueRef{
			BackendID:  arg.Content.ValueRef.BackendID,
			RevisionID: arg.Content.ValueRef.RevisionID,
		}
	}
	return params
}

// WatchSecretBackendChanged sets up a watcher to notify of changes to the secret backend.
func (s *SecretsDrainAPI) WatchSecretBackendChanged(ctx context.Context) (params.NotifyWatchResult, error) {
	stateWatcher := s.model.WatchForModelConfigChanges()
	w, err := newSecretBackendModelConfigWatcher(ctx, s.model, stateWatcher, s.logger)
	if err != nil {
		return params.NotifyWatchResult{
			Error: apiservererrors.ServerError(err),
		}, nil
	}
	id, _, err := internal.EnsureRegisterWatcher[struct{}](ctx, s.watcherRegistry, w)
	if err != nil {
		return params.NotifyWatchResult{
			Error: apiservererrors.ServerError(err),
		}, nil
	}
	return params.NotifyWatchResult{
		NotifyWatcherId: id,
	}, nil
}
