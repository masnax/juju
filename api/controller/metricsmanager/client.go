// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package metricsmanager

import (
	"context"

	"github.com/juju/errors"
	"github.com/juju/names/v4"

	"github.com/juju/juju/api/base"
	"github.com/juju/juju/rpc/params"
)

// Option is a function that can be used to configure a Client.
type Option = base.Option

// WithTracer returns an Option that configures the Client to use the
// supplied tracer.
var WithTracer = base.WithTracer

// Client provides access to the metrics manager api
type Client struct {
	modelTag names.ModelTag
	facade   base.FacadeCaller
}

// MetricsManagerClient defines the methods on the metricsmanager API end point.
type MetricsManagerClient interface {
	CleanupOldMetrics() error
	SendMetrics() error
}

var _ MetricsManagerClient = (*Client)(nil)

// NewClient creates a new client for accessing the metricsmanager api
func NewClient(apiCaller base.APICaller, options ...Option) (*Client, error) {
	modelTag, ok := apiCaller.ModelTag()
	if !ok {
		return nil, errors.New("metricsmanager client is not appropriate for controller-only API")

	}
	facade := base.NewFacadeCaller(apiCaller, "MetricsManager", options...)
	return &Client{
		modelTag: modelTag,
		facade:   facade,
	}, nil
}

// CleanupOldMetrics looks for metrics that are 24 hours old (or older)
// and have been sent. Any metrics it finds are deleted.
func (c *Client) CleanupOldMetrics() error {
	p := params.Entities{Entities: []params.Entity{
		{c.modelTag.String()},
	}}
	var results params.ErrorResults
	err := c.facade.FacadeCall(context.TODO(), "CleanupOldMetrics", p, &results)
	if err != nil {
		return errors.Trace(err)
	}
	return results.OneError()
}

// SendMetrics will send any unsent metrics to the collection service.
func (c *Client) SendMetrics() error {
	p := params.Entities{Entities: []params.Entity{
		{c.modelTag.String()},
	}}
	var results params.ErrorResults
	err := c.facade.FacadeCall(context.TODO(), "SendMetrics", p, &results)
	if err != nil {
		return errors.Trace(err)
	}
	return results.OneError()
}
