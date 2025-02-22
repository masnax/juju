// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package changestreampruner

import (
	"github.com/juju/clock"
	"github.com/juju/errors"
	"github.com/juju/worker/v3"
	"github.com/juju/worker/v3/dependency"
)

// Logger represents the logging methods called.
type Logger interface {
	Errorf(message string, args ...interface{})
	Warningf(message string, args ...interface{})
	Infof(message string, args ...interface{})
	Debugf(message string, args ...interface{})
	Tracef(message string, args ...interface{})
	IsTraceEnabled() bool
}

// NewWorkerFn is an alias function that allows the creation of
// EventQueueWorker.
type NewWorkerFn = func(WorkerConfig) (worker.Worker, error)

// ManifoldConfig defines the names of the manifolds on which a Manifold will
// depend.
type ManifoldConfig struct {
	DBAccessor string

	Clock     clock.Clock
	Logger    Logger
	NewWorker NewWorkerFn
}

func (cfg ManifoldConfig) Validate() error {
	if cfg.DBAccessor == "" {
		return errors.NotValidf("empty DBAccessorName")
	}
	if cfg.Clock == nil {
		return errors.NotValidf("nil Clock")
	}
	if cfg.Logger == nil {
		return errors.NotValidf("nil Logger")
	}
	if cfg.NewWorker == nil {
		return errors.NotValidf("nil NewWorker")
	}
	return nil
}

// Manifold returns a dependency manifold that runs the changestream
// worker, using the resource names defined in the supplied config.
func Manifold(config ManifoldConfig) dependency.Manifold {
	return dependency.Manifold{
		Inputs: []string{
			config.DBAccessor,
		},
		Start: func(context dependency.Context) (worker.Worker, error) {
			if err := config.Validate(); err != nil {
				return nil, errors.Trace(err)
			}

			var dbGetter DBGetter
			if err := context.Get(config.DBAccessor, &dbGetter); err != nil {
				return nil, errors.Trace(err)
			}

			cfg := WorkerConfig{
				DBGetter: dbGetter,
				Clock:    config.Clock,
				Logger:   config.Logger,
			}

			w, err := config.NewWorker(cfg)
			if err != nil {
				return nil, errors.Trace(err)
			}
			return w, nil
		},
	}
}

func NewWorker(cfg WorkerConfig) (worker.Worker, error) {
	return newWorker(cfg)
}
