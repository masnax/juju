// Copyright 2018 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	stdcontext "context"
	"fmt"

	"github.com/juju/collections/set"
	"github.com/juju/errors"

	"github.com/juju/juju/caas"
	"github.com/juju/juju/cloud"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/domain/credential"
	"github.com/juju/juju/domain/model"
	"github.com/juju/juju/environs"
	environscloudspec "github.com/juju/juju/environs/cloudspec"
	"github.com/juju/juju/environs/config"
	envcontext "github.com/juju/juju/environs/envcontext"
	"github.com/juju/juju/environs/instances"
)

// MachineService provides access to all machines.
type MachineService interface {
	// AllMachines returns all machines in the model.
	AllMachines() ([]Machine, error)
}

// Machine defines machine methods needed for the check.
type Machine interface {
	// IsManual returns true if the machine was manually provisioned.
	IsManual() (bool, error)

	// IsContainer returns true if the machine is a container.
	IsContainer() bool

	// InstanceId returns the provider specific instance id for this
	// machine, or a NotProvisionedError, if not set.
	InstanceId() (instance.Id, error)

	// Id returns the machine id.
	Id() string
}

// CloudProvider defines methods needed from the cloud provider to perform the check.
type CloudProvider interface {
	// AllInstances returns all instances currently known to the cloud provider.
	AllInstances(ctx envcontext.ProviderCallContext) ([]instances.Instance, error)
}

// CredentialValidationContext provides access to artefacts needed to
// validate a credential for a given model.
type CredentialValidationContext struct {
	ControllerUUID string

	Config         *config.Config
	MachineService MachineService

	ModelType model.Type
	Cloud     cloud.Cloud
	Region    string
}

// CredentialValidator instances check that a given credential is
// valid for any models which want to use it.
type CredentialValidator interface {
	Validate(
		ctx stdcontext.Context,
		validationContext CredentialValidationContext,
		credentialID credential.ID,
		credential *cloud.Credential,
		checkCloudInstances bool,
	) ([]error, error)
}

type defaultCredentialValidator struct{}

// NewCredentialValidator returns the credential validator used in production.
func NewCredentialValidator() CredentialValidator {
	return defaultCredentialValidator{}
}

// Validate checks if a new cloud credential could be valid for a model whose
// details are defined in the context.
func (v defaultCredentialValidator) Validate(
	ctx stdcontext.Context,
	validationContext CredentialValidationContext,
	id credential.ID,
	cred *cloud.Credential,
	checkCloudInstances bool,
) (machineErrors []error, err error) {
	if err := id.Validate(); err != nil {
		return nil, fmt.Errorf("credential %w", err)
	}

	openParams, err := v.buildOpenParams(validationContext, id, cred)
	if err != nil {
		return nil, errors.Trace(err)
	}
	switch validationContext.ModelType {
	case model.TypeCAAS:
		return checkCAASModelCredential(ctx, openParams)
	case model.TypeIAAS:
		return checkIAASModelCredential(ctx, validationContext.MachineService, openParams, checkCloudInstances)
	default:
		return nil, errors.NotSupportedf("model type %q", validationContext.ModelType)
	}
}

func checkCAASModelCredential(ctx stdcontext.Context, brokerParams environs.OpenParams) ([]error, error) {
	broker, err := newCAASBroker(ctx, brokerParams)
	if err != nil {
		return nil, errors.Trace(err)
	}

	if err = broker.CheckCloudCredentials(); err != nil {
		return nil, errors.Trace(err)
	}
	return nil, nil
}

func checkIAASModelCredential(ctx stdcontext.Context, machineService MachineService, openParams environs.OpenParams, checkCloudInstances bool) ([]error, error) {
	env, err := newEnv(ctx, openParams)
	if err != nil {
		return nil, errors.Trace(err)
	}
	// We only check persisted machines vs known cloud instances.
	// In the future, this check may be extended to other cloud resources,
	// entities and operation-level authorisations such as interfaces,
	// ability to CRUD storage, etc.
	return checkMachineInstances(ctx, machineService, env, checkCloudInstances)
}

// checkMachineInstances compares model machines from state with
// the ones reported by the provider using supplied credential.
// This only makes sense for non-k8s providers.
func checkMachineInstances(ctx stdcontext.Context, machineService MachineService, provider CloudProvider, checkCloudInstances bool) ([]error, error) {
	// Get machines from state
	machines, err := machineService.AllMachines()
	if err != nil {
		return nil, errors.Trace(err)
	}

	var results []error

	machinesByInstance := make(map[string]string)
	for _, machine := range machines {
		if machine.IsContainer() {
			// Containers don't correspond to instances at the
			// provider level.
			continue
		}
		if manual, err := machine.IsManual(); err != nil {
			return nil, errors.Trace(err)
		} else if manual {
			continue
		}
		instanceId, err := machine.InstanceId()
		if errors.Is(err, errors.NotProvisioned) {
			// Skip over this machine; we wouldn't expect the cloud
			// to know about it.
			continue
		} else if err != nil {
			results = append(results, errors.Annotatef(err, "getting instance id for machine %s", machine.Id()))
			continue
		}
		machinesByInstance[string(instanceId)] = machine.Id()
	}

	// Check that we can see all machines' instances regardless of their state as perceived by the cloud, i.e.
	// this call will return all non-terminated instances.
	callCtx := envcontext.WithoutCredentialInvalidator(ctx)
	instances, err := provider.AllInstances(callCtx)
	if err != nil {
		return nil, errors.Trace(err)
	}

	// From here, there 2 ways of checking whether the credential is valid:
	// 1. Can we reach all cloud instances that machines know about?
	// 2. Can we cross examine all machines we know about with all the instances we can reach
	// and ensure that they correspond 1:1.
	// Second check (2) is more useful for model migration, for example, since we want to know if
	// we have moved the known universe correctly. However, it is a but redundant if we just care about
	// credential validity since the first check (1) addresses all our concerns.

	instanceIds := set.NewStrings()
	for _, instance := range instances {
		id := string(instance.Id())
		instanceIds.Add(id)
		if checkCloudInstances {
			if _, found := machinesByInstance[id]; !found {
				results = append(results, errors.Errorf("no machine with instance %q", id))
			}
		}
	}

	for instanceId, name := range machinesByInstance {
		if !instanceIds.Contains(instanceId) {
			results = append(results, errors.Errorf("couldn't find instance %q for machine %s", instanceId, name))
		}
	}

	return results, nil
}

var (
	newEnv        = environs.New
	newCAASBroker = caas.New
)

func (v defaultCredentialValidator) buildOpenParams(
	ctx CredentialValidationContext, credentialID credential.ID, credential *cloud.Credential,
) (environs.OpenParams, error) {
	fail := func(original error) (environs.OpenParams, error) {
		return environs.OpenParams{}, original
	}

	err := v.validateCloudCredential(ctx.Cloud, credentialID)
	if err != nil {
		return fail(errors.Trace(err))
	}

	tempCloudSpec, err := environscloudspec.MakeCloudSpec(ctx.Cloud, ctx.Region, credential)
	if err != nil {
		return fail(errors.Trace(err))
	}

	return environs.OpenParams{
		ControllerUUID: ctx.ControllerUUID,
		Cloud:          tempCloudSpec,
		Config:         ctx.Config,
	}, nil
}

// validateCloudCredential validates the given cloud credential
// name against the provided cloud definition and credentials.
func (v defaultCredentialValidator) validateCloudCredential(
	cld cloud.Cloud,
	credentialID credential.ID,
) error {
	if !credentialID.IsZero() {
		if credentialID.Cloud != cld.Name {
			return errors.NotValidf("credential %q", credentialID)
		}
		return nil
	}
	var hasEmptyAuth bool
	for _, authType := range cld.AuthTypes {
		if authType != cloud.EmptyAuthType {
			continue
		}
		hasEmptyAuth = true
		break
	}
	if !hasEmptyAuth {
		return errors.NotValidf("missing CloudCredential")
	}
	return nil
}
