// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package modelworkermanager

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/juju/clock"
	"github.com/juju/errors"
	"github.com/juju/loggo"
	"github.com/juju/names/v4"
	"github.com/juju/worker/v3"
	"github.com/juju/worker/v3/catacomb"

	agentengine "github.com/juju/juju/agent/engine"
	"github.com/juju/juju/apiserver/apiserverhttp"
	"github.com/juju/juju/controller"
	corelogger "github.com/juju/juju/core/logger"
	"github.com/juju/juju/internal/pki"
	"github.com/juju/juju/state"
)

// ModelWatcher provides an interface for watching the additiona and
// removal of models.
type ModelWatcher interface {
	WatchModels() state.StringsWatcher
}

// ControllerConfigGetter is an interface that returns the controller config.
type ControllerConfigGetter interface {
	ControllerConfig(context.Context) (controller.Config, error)
}

// Controller provides an interface for getting models by UUID,
// and other details needed to pass into the function to start workers for a model.
// Once a model is no longer required, the returned function must
// be called to dispose of the model.
type Controller interface {
	Model(modelUUID string) (Model, func(), error)
	RecordLogger(modelUUID string) (RecordLogger, error)
}

// Model represents a model.
type Model interface {
	MigrationMode() state.MigrationMode
	Type() state.ModelType
	Name() string
	Owner() names.UserTag
}

// RecordLogger writes logs to backing store.
type RecordLogger interface {
	io.Closer
	// Log writes the given log records to the logger's storage.
	Log([]corelogger.LogRecord) error
}

// ModelLogger is a database backed loggo Writer.
type ModelLogger interface {
	loggo.Writer
	Close() error
}

// MetricSink describes a way to unregister a model metrics collector. This
// ensures that we correctly tidy up after the removal of a model.
type MetricSink = agentengine.MetricSink

// ModelMetrics defines a way to create metrics for a model.
type ModelMetrics interface {
	ForModel(names.ModelTag) MetricSink
}

// NewModelConfig holds the information required by the NewModelWorkerFunc
// to start the workers for the specified model
type NewModelConfig struct {
	Authority        pki.Authority
	ModelName        string // Use a fully qualified name "<namespace>-<name>"
	ModelUUID        string
	ModelType        state.ModelType
	ModelLogger      ModelLogger
	ModelMetrics     MetricSink
	Mux              *apiserverhttp.Mux
	ControllerConfig controller.Config
}

// NewModelWorkerFunc should return a worker responsible for running
// all a model's required workers; and for returning nil when there's
// no more model to manage.
type NewModelWorkerFunc func(config NewModelConfig) (worker.Worker, error)

// Config holds the dependencies and configuration necessary to run
// a model worker manager.
type Config struct {
	Authority              pki.Authority
	Clock                  clock.Clock
	Logger                 Logger
	MachineID              string
	ModelWatcher           ModelWatcher
	ModelMetrics           ModelMetrics
	Mux                    *apiserverhttp.Mux
	Controller             Controller
	ControllerConfigGetter ControllerConfigGetter
	NewModelWorker         NewModelWorkerFunc
	ErrorDelay             time.Duration
}

// Validate returns an error if config cannot be expected to drive
// a functional model worker manager.
func (config Config) Validate() error {
	if config.Authority == nil {
		return errors.NotValidf("nil authority")
	}
	if config.Clock == nil {
		return errors.NotValidf("nil Clock")
	}
	if config.Logger == nil {
		return errors.NotValidf("nil Logger")
	}
	if config.MachineID == "" {
		return errors.NotValidf("empty MachineID")
	}
	if config.ModelWatcher == nil {
		return errors.NotValidf("nil ModelWatcher")
	}
	if config.ModelMetrics == nil {
		return errors.NotValidf("nil ModelMetrics")
	}
	if config.Controller == nil {
		return errors.NotValidf("nil Controller")
	}
	if config.ControllerConfigGetter == nil {
		return errors.NotValidf("nil ControllerConfigGetter")
	}
	if config.NewModelWorker == nil {
		return errors.NotValidf("nil NewModelWorker")
	}
	if config.ErrorDelay <= 0 {
		return errors.NotValidf("non-positive ErrorDelay")
	}
	return nil
}

// New starts a new model worker manager.
func New(config Config) (worker.Worker, error) {
	if err := config.Validate(); err != nil {
		return nil, errors.Trace(err)
	}
	m := &modelWorkerManager{
		config: config,
	}

	err := catacomb.Invoke(catacomb.Plan{
		Site: &m.catacomb,
		Work: m.loop,
	})
	if err != nil {
		return nil, errors.Trace(err)
	}
	return m, nil
}

type modelWorkerManager struct {
	catacomb catacomb.Catacomb
	config   Config
	runner   *worker.Runner
}

// Kill satisfies the Worker interface.
func (m *modelWorkerManager) Kill() {
	m.catacomb.Kill(nil)
}

// Wait satisfies the Worker interface.
func (m *modelWorkerManager) Wait() error {
	return m.catacomb.Wait()
}

func (m *modelWorkerManager) loop() error {
	m.runner = worker.NewRunner(worker.RunnerParams{
		IsFatal:       neverFatal,
		MoreImportant: neverImportant,
		RestartDelay:  m.config.ErrorDelay,
		Logger:        m.config.Logger,
	})
	if err := m.catacomb.Add(m.runner); err != nil {
		return errors.Trace(err)
	}
	watcher := m.config.ModelWatcher.WatchModels()
	if err := m.catacomb.Add(watcher); err != nil {
		return errors.Trace(err)
	}

	modelChanged := func(modelUUID string) error {
		model, release, err := m.config.Controller.Model(modelUUID)
		if errors.Is(err, errors.NotFound) {
			// Model was removed, ignore it.
			// The reason we ignore it here is that one of the embedded
			// workers is also responding to the model life changes and
			// when it returns a NotFound error, which is determined as a
			// fatal error for the model worker engine. This causes it to be
			// removed from the runner above. However since the runner itself
			// has neverFatal as an error handler, the runner itself doesn't
			// propagate the error.
			return nil
		} else if err != nil {
			return errors.Trace(err)
		}
		defer release()

		if !isModelActive(model) {
			// Ignore this model until it's activated - we
			// never want to run workers for an importing
			// model.
			// https://bugs.launchpad.net/juju/+bug/1646310
			return nil
		}

		cfg := NewModelConfig{
			Authority:    m.config.Authority,
			ModelName:    fmt.Sprintf("%s-%s", model.Owner().Id(), model.Name()),
			ModelUUID:    modelUUID,
			ModelType:    model.Type(),
			ModelMetrics: m.config.ModelMetrics.ForModel(names.NewModelTag(modelUUID)),
			Mux:          m.config.Mux,
		}
		return errors.Trace(m.ensure(cfg))
	}

	for {
		select {
		case <-m.catacomb.Dying():
			return m.catacomb.ErrDying()
		case uuids, ok := <-watcher.Changes():
			if !ok {
				return errors.New("changes stopped")
			}
			for _, modelUUID := range uuids {
				if err := modelChanged(modelUUID); err != nil {
					return errors.Trace(err)
				}
			}
		}
	}
}

func (m *modelWorkerManager) ensure(cfg NewModelConfig) error {
	starter := m.starter(cfg)
	if err := m.runner.StartWorker(cfg.ModelUUID, starter); !errors.Is(err, errors.AlreadyExists) {
		return errors.Trace(err)
	}
	return nil
}

func (m *modelWorkerManager) starter(cfg NewModelConfig) func() (worker.Worker, error) {
	return func() (worker.Worker, error) {
		modelUUID := cfg.ModelUUID
		modelName := fmt.Sprintf("%q (%s)", cfg.ModelName, cfg.ModelUUID)
		m.config.Logger.Debugf("starting workers for model %s", modelName)

		// Get the controller config for the model worker so that we correctly
		// handle the case where the controller config changes between model
		// worker restarts.
		ctx := m.catacomb.Context(context.Background())
		controllerConfig, err := m.config.ControllerConfigGetter.ControllerConfig(ctx)
		if err != nil {
			return nil, errors.Annotate(err, "unable to get controller config")
		}
		cfg.ControllerConfig = controllerConfig

		recordLogger, err := m.config.Controller.RecordLogger(modelUUID)
		if err != nil {
			return nil, errors.Annotatef(err, "unable to create db logger for %s", modelName)
		}

		cfg.ModelLogger = newModelLogger(
			"controller-"+m.config.MachineID,
			modelUUID,
			recordLogger,
			m.config.Clock,
			m.config.Logger,
		)
		worker, err := m.config.NewModelWorker(cfg)
		if err != nil {
			cfg.ModelLogger.Close()
			return nil, errors.Annotatef(err, "cannot manage model %s", modelName)
		}
		return worker, nil
	}
}

func neverFatal(error) bool {
	return false
}

func neverImportant(error, error) bool {
	return false
}

func isModelActive(m Model) bool {
	return m.MigrationMode() != state.MigrationModeImporting
}

// Report shows up in the dependency engine report.
func (m *modelWorkerManager) Report() map[string]any {
	if m.runner == nil {
		return nil
	}
	return m.runner.Report()
}
