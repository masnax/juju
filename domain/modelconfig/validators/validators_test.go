// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package validators

import (
	"github.com/juju/errors"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/testing"
)

type dummySecretBackendProviderFunc func(string) (bool, error)

type dummySpaceProviderFunc func(string) (bool, error)

type validatorsSuite struct{}

var _ = gc.Suite(&validatorsSuite{})

func (d dummySecretBackendProviderFunc) HasSecretsBackend(s string) (bool, error) {
	return d(s)
}

func (d dummySpaceProviderFunc) HasSpace(s string) (bool, error) {
	return d(s)
}

func (_ *validatorsSuite) TestCharmhubURLChange(c *gc.C) {
	oldCfg, err := config.New(config.NoDefaults, map[string]any{
		"name":         "wallyworld",
		"uuid":         testing.ModelTag.Id(),
		"type":         "sometype",
		"charmhub-url": "https://charmhub.example.com",
	})
	c.Assert(err, jc.ErrorIsNil)

	newCfg, err := config.New(config.NoDefaults, map[string]any{
		"name":         "wallyworld",
		"uuid":         testing.ModelTag.Id(),
		"type":         "sometype",
		"charmhub-url": "https://charmhub1.example.com",
	})
	c.Assert(err, jc.ErrorIsNil)

	var validationError *config.ValidationError
	_, err = CharmhubURLChange()(newCfg, oldCfg)
	c.Assert(errors.As(err, &validationError), jc.IsTrue)
	c.Assert(validationError.InvalidAttrs, gc.DeepEquals, []string{"charmhub-url"})
}

func (_ *validatorsSuite) TestCharmhubURLNoChange(c *gc.C) {
	oldCfg, err := config.New(config.NoDefaults, map[string]any{
		"name":         "wallyworld",
		"uuid":         testing.ModelTag.Id(),
		"type":         "sometype",
		"charmhub-url": "https://charmhub.example.com",
	})
	c.Assert(err, jc.ErrorIsNil)

	newCfg, err := config.New(config.NoDefaults, map[string]any{
		"name":         "wallyworld",
		"uuid":         testing.ModelTag.Id(),
		"type":         "sometype",
		"charmhub-url": "https://charmhub.example.com",
	})
	c.Assert(err, jc.ErrorIsNil)

	_, err = CharmhubURLChange()(newCfg, oldCfg)
	c.Assert(err, jc.ErrorIsNil)
}

func (_ *validatorsSuite) TestSpaceCheckerFound(c *gc.C) {
	provider := dummySpaceProviderFunc(func(s string) (bool, error) {
		c.Assert(s, gc.Equals, "foobar")
		return true, nil
	})

	oldCfg, err := config.New(config.NoDefaults, map[string]any{
		"name":          "wallyworld",
		"uuid":          testing.ModelTag.Id(),
		"type":          "sometype",
		"default-space": "foobar",
	})
	c.Assert(err, jc.ErrorIsNil)

	newCfg, err := config.New(config.NoDefaults, map[string]any{
		"name":          "wallyworld",
		"uuid":          testing.ModelTag.Id(),
		"type":          "sometype",
		"default-space": "foobar",
	})
	c.Assert(err, jc.ErrorIsNil)

	_, err = SpaceChecker(provider)(newCfg, oldCfg)
	c.Assert(err, jc.ErrorIsNil)
}

func (_ *validatorsSuite) TestSpaceCheckerNotFound(c *gc.C) {
	provider := dummySpaceProviderFunc(func(s string) (bool, error) {
		c.Assert(s, gc.Equals, "foobar")
		return false, nil
	})

	oldCfg, err := config.New(config.NoDefaults, map[string]any{
		"name":          "wallyworld",
		"uuid":          testing.ModelTag.Id(),
		"type":          "sometype",
		"default-space": "foobar",
	})
	c.Assert(err, jc.ErrorIsNil)

	newCfg, err := config.New(config.NoDefaults, map[string]any{
		"name":          "wallyworld",
		"uuid":          testing.ModelTag.Id(),
		"type":          "sometype",
		"default-space": "foobar",
	})
	c.Assert(err, jc.ErrorIsNil)

	_, err = SpaceChecker(provider)(newCfg, oldCfg)
	var validationError *config.ValidationError
	c.Assert(errors.As(err, &validationError), jc.IsTrue)
	c.Assert(validationError.InvalidAttrs, gc.DeepEquals, []string{"default-space"})
}

func (_ *validatorsSuite) TestSpaceCheckerError(c *gc.C) {
	providerErr := errors.New("some error")
	provider := dummySpaceProviderFunc(func(s string) (bool, error) {
		c.Assert(s, gc.Equals, "foobar")
		return false, providerErr
	})

	oldCfg, err := config.New(config.NoDefaults, map[string]any{
		"name":          "wallyworld",
		"uuid":          testing.ModelTag.Id(),
		"type":          "sometype",
		"default-space": "foobar",
	})
	c.Assert(err, jc.ErrorIsNil)

	newCfg, err := config.New(config.NoDefaults, map[string]any{
		"name":          "wallyworld",
		"uuid":          testing.ModelTag.Id(),
		"type":          "sometype",
		"default-space": "foobar",
	})
	c.Assert(err, jc.ErrorIsNil)

	_, err = SpaceChecker(provider)(newCfg, oldCfg)
	c.Assert(err, jc.ErrorIs, providerErr)
}

func (_ *validatorsSuite) TestLoggincTracePermissionNoTrace(c *gc.C) {
	oldCfg, err := config.New(config.NoDefaults, map[string]any{
		"name": "wallyworld",
		"uuid": testing.ModelTag.Id(),
		"type": "sometype",
	})
	c.Assert(err, jc.ErrorIsNil)

	newCfg, err := config.New(config.NoDefaults, map[string]any{
		"name":           "wallyworld",
		"uuid":           testing.ModelTag.Id(),
		"type":           "sometype",
		"logging-config": "root=DEBUG",
	})
	c.Assert(err, jc.ErrorIsNil)

	_, err = LoggingTracePermissionChecker(false)(newCfg, oldCfg)
	c.Assert(err, jc.ErrorIsNil)
}

func (_ *validatorsSuite) TestLoggincTracePermissionTrace(c *gc.C) {
	oldCfg, err := config.New(config.NoDefaults, map[string]any{
		"name": "wallyworld",
		"uuid": testing.ModelTag.Id(),
		"type": "sometype",
	})
	c.Assert(err, jc.ErrorIsNil)

	newCfg, err := config.New(config.NoDefaults, map[string]any{
		"name":           "wallyworld",
		"uuid":           testing.ModelTag.Id(),
		"type":           "sometype",
		"logging-config": "root=TRACE",
	})
	c.Assert(err, jc.ErrorIsNil)

	_, err = LoggingTracePermissionChecker(false)(newCfg, oldCfg)
	c.Assert(err, jc.ErrorIs, ErrorLogTracingPermission)

	var validationError *config.ValidationError
	c.Assert(errors.As(err, &validationError), jc.IsTrue)
	c.Assert(validationError.InvalidAttrs, gc.DeepEquals, []string{"logging-config"})
}

func (_ *validatorsSuite) TestLoggincTracePermissionTraceAllow(c *gc.C) {
	oldCfg, err := config.New(config.NoDefaults, map[string]any{
		"name": "wallyworld",
		"uuid": testing.ModelTag.Id(),
		"type": "sometype",
	})
	c.Assert(err, jc.ErrorIsNil)

	newCfg, err := config.New(config.NoDefaults, map[string]any{
		"name":           "wallyworld",
		"uuid":           testing.ModelTag.Id(),
		"type":           "sometype",
		"logging-config": "root=TRACE",
	})
	c.Assert(err, jc.ErrorIsNil)

	_, err = LoggingTracePermissionChecker(true)(newCfg, oldCfg)
	c.Assert(err, jc.ErrorIsNil)
}

func (_ *validatorsSuite) TestSecretsBackendChecker(c *gc.C) {
	provider := dummySecretBackendProviderFunc(func(s string) (bool, error) {
		c.Assert(s, gc.Equals, "vault")
		return true, nil
	})

	oldCfg, err := config.New(config.NoDefaults, map[string]any{
		"name":           "wallyworld",
		"uuid":           testing.ModelTag.Id(),
		"type":           "sometype",
		"secret-backend": "default",
	})
	c.Assert(err, jc.ErrorIsNil)

	newCfg, err := config.New(config.NoDefaults, map[string]any{
		"name":           "wallyworld",
		"uuid":           testing.ModelTag.Id(),
		"type":           "sometype",
		"secret-backend": "vault",
	})
	c.Assert(err, jc.ErrorIsNil)

	_, err = SecretBackendChecker(provider)(newCfg, oldCfg)
	c.Assert(err, jc.ErrorIsNil)
}

func (_ *validatorsSuite) TestSecretsBackendCheckerNoExist(c *gc.C) {
	provider := dummySecretBackendProviderFunc(func(s string) (bool, error) {
		c.Assert(s, gc.Equals, "vault")
		return false, nil
	})

	oldCfg, err := config.New(config.NoDefaults, map[string]any{
		"name":           "wallyworld",
		"uuid":           testing.ModelTag.Id(),
		"type":           "sometype",
		"secret-backend": "default",
	})
	c.Assert(err, jc.ErrorIsNil)

	newCfg, err := config.New(config.NoDefaults, map[string]any{
		"name":           "wallyworld",
		"uuid":           testing.ModelTag.Id(),
		"type":           "sometype",
		"secret-backend": "vault",
	})
	c.Assert(err, jc.ErrorIsNil)

	_, err = SecretBackendChecker(provider)(newCfg, oldCfg)
	var validationError *config.ValidationError
	c.Assert(errors.As(err, &validationError), jc.IsTrue)
	c.Assert(validationError.InvalidAttrs, gc.DeepEquals, []string{"secret-backend"})
}

func (_ *validatorsSuite) TestSecretsBackendCheckerProviderError(c *gc.C) {
	providerErr := errors.New("some error")
	provider := dummySecretBackendProviderFunc(func(s string) (bool, error) {
		c.Assert(s, gc.Equals, "vault")
		return false, providerErr
	})

	oldCfg, err := config.New(config.NoDefaults, map[string]any{
		"name":           "wallyworld",
		"uuid":           testing.ModelTag.Id(),
		"type":           "sometype",
		"secret-backend": "default",
	})
	c.Assert(err, jc.ErrorIsNil)

	newCfg, err := config.New(config.NoDefaults, map[string]any{
		"name":           "wallyworld",
		"uuid":           testing.ModelTag.Id(),
		"type":           "sometype",
		"secret-backend": "vault",
	})
	c.Assert(err, jc.ErrorIsNil)

	_, err = SecretBackendChecker(provider)(newCfg, oldCfg)
	c.Assert(err, jc.ErrorIs, providerErr)
}
