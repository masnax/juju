// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package usersecretsdrain

import (
	"github.com/juju/errors"
	"github.com/juju/loggo"

	commonsecrets "github.com/juju/juju/apiserver/common/secrets"
	apiservererrors "github.com/juju/juju/apiserver/errors"
	coresecrets "github.com/juju/juju/core/secrets"
	"github.com/juju/juju/internal/secrets"
	"github.com/juju/juju/internal/secrets/provider"
	"github.com/juju/juju/rpc/params"
)

var logger = loggo.GetLogger("juju.apiserver.usersecretsdrain")

// SecretsDrainAPI is the implementation for the SecretsDrain facade.
type SecretsDrainAPI struct {
	*commonsecrets.SecretsDrainAPI

	secretsState SecretsState

	drainConfigGetter   commonsecrets.BackendDrainConfigGetter
	backendConfigGetter commonsecrets.BackendConfigGetter
}

// GetSecretBackendConfigs gets the config needed to create a client to secret backends for the drain worker.
func (s *SecretsDrainAPI) GetSecretBackendConfigs(arg params.SecretBackendArgs) (params.SecretBackendConfigResults, error) {
	if len(arg.BackendIDs) > 1 {
		return params.SecretBackendConfigResults{}, errors.Errorf("Maximumly only one backend ID can be specified for drain")
	}
	var backendID string
	if len(arg.BackendIDs) == 1 {
		backendID = arg.BackendIDs[0]
	}
	results := params.SecretBackendConfigResults{
		Results: make(map[string]params.SecretBackendConfigResult, 1),
	}
	cfgInfo, err := s.drainConfigGetter(backendID)
	if err != nil {
		return results, errors.Trace(err)
	}
	if len(cfgInfo.Configs) == 0 {
		return results, errors.NotFoundf("no secret backends available")
	}
	results.ActiveID = cfgInfo.ActiveID
	for id, cfg := range cfgInfo.Configs {
		results.Results[id] = params.SecretBackendConfigResult{
			ControllerUUID: cfg.ControllerUUID,
			ModelUUID:      cfg.ModelUUID,
			ModelName:      cfg.ModelName,
			Draining:       true,
			Config: params.SecretBackendConfig{
				BackendType: cfg.BackendType,
				Params:      cfg.Config,
			},
		}
	}
	return results, nil
}

// GetSecretContentInfo returns the secret values for the specified secrets.
func (s *SecretsDrainAPI) GetSecretContentInfo(args params.GetSecretContentArgs) (params.SecretContentResults, error) {
	result := params.SecretContentResults{
		Results: make([]params.SecretContentResult, len(args.Args)),
	}
	for i, arg := range args.Args {
		content, backend, draining, err := s.getSecretContent(arg)
		if err != nil {
			result.Results[i].Error = apiservererrors.ServerError(err)
			continue
		}
		contentParams := params.SecretContentParams{}
		if content.ValueRef != nil {
			contentParams.ValueRef = &params.SecretValueRef{
				BackendID:  content.ValueRef.BackendID,
				RevisionID: content.ValueRef.RevisionID,
			}
		}
		if content.SecretValue != nil {
			contentParams.Data = content.SecretValue.EncodedValues()
		}
		result.Results[i].Content = contentParams
		if backend != nil {
			result.Results[i].BackendConfig = &params.SecretBackendConfigResult{
				ControllerUUID: backend.ControllerUUID,
				ModelUUID:      backend.ModelUUID,
				ModelName:      backend.ModelName,
				Draining:       draining,
				Config: params.SecretBackendConfig{
					BackendType: backend.BackendType,
					Params:      backend.Config,
				},
			}
		}
	}
	return result, nil
}

func (s *SecretsDrainAPI) getSecretContent(arg params.GetSecretContentArg) (
	*secrets.ContentParams, *provider.ModelBackendConfig, bool, error,
) {
	if arg.URI == "" {
		return nil, nil, false, errors.NewNotValid(nil, "empty URI")
	}

	uri, err := coresecrets.ParseURI(arg.URI)
	if err != nil {
		return nil, nil, false, errors.Trace(err)
	}
	logger.Debugf("getting secret content for: %s", uri)

	md, err := s.secretsState.GetSecret(uri)
	if err != nil {
		return nil, nil, false, errors.Trace(err)
	}

	val, valueRef, err := s.secretsState.GetSecretValue(md.URI, md.LatestRevision)
	if err != nil {
		return nil, nil, false, errors.Trace(err)
	}
	content := &secrets.ContentParams{SecretValue: val, ValueRef: valueRef}
	if content.ValueRef == nil {
		// Internal secret.
		return content, nil, false, errors.Trace(err)
	}
	// Get backend config for external secret.
	backend, draining, err := s.getBackend(content.ValueRef.BackendID)
	return content, backend, draining, errors.Trace(err)
}

func (s *SecretsDrainAPI) getBackend(backendID string) (*provider.ModelBackendConfig, bool, error) {
	cfgInfo, err := s.backendConfigGetter([]string{backendID}, false)
	if err != nil {
		return nil, false, errors.Trace(err)
	}
	cfg, ok := cfgInfo.Configs[backendID]
	if ok {
		return &provider.ModelBackendConfig{
			ControllerUUID: cfg.ControllerUUID,
			ModelUUID:      cfg.ModelUUID,
			ModelName:      cfg.ModelName,
			BackendConfig: provider.BackendConfig{
				BackendType: cfg.BackendType,
				Config:      cfg.Config,
			},
		}, backendID != cfgInfo.ActiveID, nil
	}
	return nil, false, errors.NotFoundf("secret backend %q", backendID)
}

// GetSecretRevisionContentInfo returns the secret values for the specified secret revisions.
func (s *SecretsDrainAPI) GetSecretRevisionContentInfo(arg params.SecretRevisionArg) (params.SecretContentResults, error) {
	result := params.SecretContentResults{
		Results: make([]params.SecretContentResult, len(arg.Revisions)),
	}
	uri, err := coresecrets.ParseURI(arg.URI)
	if err != nil {
		return params.SecretContentResults{}, errors.Trace(err)
	}

	for i, rev := range arg.Revisions {
		val, valueRef, err := s.secretsState.GetSecretValue(uri, rev)
		if err != nil {
			result.Results[i].Error = apiservererrors.ServerError(err)
			continue
		}
		contentParams := params.SecretContentParams{}
		if valueRef != nil {
			contentParams.ValueRef = &params.SecretValueRef{
				BackendID:  valueRef.BackendID,
				RevisionID: valueRef.RevisionID,
			}
			backend, draining, err := s.getBackend(valueRef.BackendID)
			if err != nil {
				result.Results[i].Error = apiservererrors.ServerError(err)
				continue
			}
			result.Results[i].BackendConfig = &params.SecretBackendConfigResult{
				ControllerUUID: backend.ControllerUUID,
				ModelUUID:      backend.ModelUUID,
				ModelName:      backend.ModelName,
				Draining:       draining,
				Config: params.SecretBackendConfig{
					BackendType: backend.BackendType,
					Params:      backend.Config,
				},
			}
		}
		if val != nil {
			contentParams.Data = val.EncodedValues()
		}
		result.Results[i].Content = contentParams
	}
	return result, nil
}
