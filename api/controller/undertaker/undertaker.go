// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package undertaker

import (
	"context"

	"github.com/juju/errors"
	"github.com/juju/names/v4"

	"github.com/juju/juju/api/base"
	"github.com/juju/juju/api/common"
	"github.com/juju/juju/api/common/cloudspec"
	"github.com/juju/juju/core/status"
	"github.com/juju/juju/core/watcher"
	"github.com/juju/juju/rpc/params"
)

// Option is a function that can be used to configure a Client.
type Option = base.Option

// WithTracer returns an Option that configures the Client to use the
// supplied tracer.
var WithTracer = base.WithTracer

// NewWatcherFunc exists to let us test Watch properly.
type NewWatcherFunc func(base.APICaller, params.NotifyWatchResult) watcher.NotifyWatcher

// Client provides access to the undertaker API
type Client struct {
	*cloudspec.CloudSpecAPI
	*common.ModelWatcher
	modelTag   names.ModelTag
	caller     base.FacadeCaller
	newWatcher NewWatcherFunc
}

// NewClient creates a new client for accessing the undertaker API.
func NewClient(caller base.APICaller, newWatcher NewWatcherFunc, options ...Option) (*Client, error) {
	modelTag, ok := caller.ModelTag()
	if !ok {
		return nil, errors.New("undertaker client is not appropriate for controller-only API")
	}
	facadeCaller := base.NewFacadeCaller(caller, "Undertaker", options...)
	return &Client{
		modelTag:     modelTag,
		caller:       facadeCaller,
		newWatcher:   newWatcher,
		CloudSpecAPI: cloudspec.NewCloudSpecAPI(facadeCaller, modelTag),
		ModelWatcher: common.NewModelWatcher(facadeCaller),
	}, nil
}

// ModelInfo returns information on the model needed by the undertaker worker.
func (c *Client) ModelInfo() (params.UndertakerModelInfoResult, error) {
	result := params.UndertakerModelInfoResult{}
	err := c.entityFacadeCall("ModelInfo", &result)
	return result, errors.Trace(err)
}

// ProcessDyingModel checks if a dying model has any machines or applications.
// If there are none, the model's life is changed from dying to dead.
func (c *Client) ProcessDyingModel() error {
	return c.entityFacadeCall("ProcessDyingModel", nil)
}

// RemoveModel removes any records of this model from Juju.
func (c *Client) RemoveModel() error {
	return c.entityFacadeCall("RemoveModel", nil)
}

// SetStatus sets the status of the model.
func (c *Client) SetStatus(status status.Status, message string, data map[string]interface{}) error {
	args := params.SetStatus{
		Entities: []params.EntityStatusArgs{
			{c.modelTag.String(), status.String(), message, data},
		},
	}
	var results params.ErrorResults
	if err := c.caller.FacadeCall(context.TODO(), "SetStatus", args, &results); err != nil {
		return errors.Trace(err)
	}
	if len(results.Results) != 1 {
		return errors.Errorf("expected 1 result, got %d", len(results.Results))
	}
	if results.Results[0].Error != nil {
		return errors.Trace(results.Results[0].Error)
	}
	return nil
}

func (c *Client) entityFacadeCall(name string, results interface{}) error {
	args := params.Entities{
		Entities: []params.Entity{{c.modelTag.String()}},
	}
	return c.caller.FacadeCall(context.TODO(), name, args, results)
}

// WatchModelResources starts a watcher for changes to the model's
// machines and applications.
func (c *Client) WatchModelResources() (watcher.NotifyWatcher, error) {
	var results params.NotifyWatchResults
	err := c.entityFacadeCall("WatchModelResources", &results)
	if err != nil {
		return nil, err
	}
	if len(results.Results) != 1 {
		return nil, errors.Errorf("expected 1 result, got %d", len(results.Results))
	}
	result := results.Results[0]
	if result.Error != nil {
		return nil, result.Error
	}
	w := c.newWatcher(c.caller.RawAPICaller(), result)
	return w, nil
}

// WatchModel starts a watcher for changes to the model.
func (c *Client) WatchModel() (watcher.NotifyWatcher, error) {
	var results params.NotifyWatchResults
	err := c.entityFacadeCall("WatchModel", &results)
	if err != nil {
		return nil, err
	}
	if len(results.Results) != 1 {
		return nil, errors.Errorf("expected 1 result, got %d", len(results.Results))
	}
	result := results.Results[0]
	if result.Error != nil {
		return nil, result.Error
	}
	w := c.newWatcher(c.caller.RawAPICaller(), result)
	return w, nil
}
