// Copyright 2022 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package servicefactory

import (
	"github.com/juju/errors"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/worker/v3"
	dt "github.com/juju/worker/v3/dependency/testing"
	"github.com/juju/worker/v3/workertest"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/core/changestream"
	coredatabase "github.com/juju/juju/core/database"
	"github.com/juju/juju/domain/model"
	"github.com/juju/juju/internal/servicefactory"
)

type manifoldSuite struct {
	baseSuite
}

var _ = gc.Suite(&manifoldSuite{})

func (s *manifoldSuite) TestValidateConfig(c *gc.C) {
	defer s.setupMocks(c).Finish()

	cfg := s.getConfig()
	c.Check(cfg.Validate(), jc.ErrorIsNil)

	cfg = s.getConfig()
	cfg.Logger = nil
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)

	cfg = s.getConfig()
	cfg.DBAccessorName = ""
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)

	cfg = s.getConfig()
	cfg.ChangeStreamName = ""
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)

	cfg = s.getConfig()
	cfg.NewWorker = nil
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)

	cfg = s.getConfig()
	cfg.NewServiceFactoryGetter = nil
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)

	cfg = s.getConfig()
	cfg.NewControllerServiceFactory = nil
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)

	cfg = s.getConfig()
	cfg.NewModelServiceFactory = nil
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)
}

func (s *manifoldSuite) TestStart(c *gc.C) {
	context := map[string]any{
		"dbaccessor":   s.dbDeleter,
		"changestream": s.dbGetter,
	}

	manifold := Manifold(ManifoldConfig{
		DBAccessorName:              "dbaccessor",
		ChangeStreamName:            "changestream",
		Logger:                      s.logger,
		NewWorker:                   NewWorker,
		NewServiceFactoryGetter:     NewServiceFactoryGetter,
		NewControllerServiceFactory: NewControllerServiceFactory,
		NewModelServiceFactory:      NewModelServiceFactory,
	})
	w, err := manifold.Start(dt.StubContext(nil, context))
	c.Assert(err, jc.ErrorIsNil)
	defer workertest.DirtyKill(c, w)

	workertest.CheckAlive(c, w)
}

func (s *manifoldSuite) TestOutputControllerServiceFactory(c *gc.C) {
	w, err := NewWorker(Config{
		DBDeleter:                   s.dbDeleter,
		DBGetter:                    s.dbGetter,
		Logger:                      s.logger,
		NewServiceFactoryGetter:     NewServiceFactoryGetter,
		NewControllerServiceFactory: NewControllerServiceFactory,
		NewModelServiceFactory:      NewModelServiceFactory,
	})
	c.Assert(err, jc.ErrorIsNil)
	defer workertest.DirtyKill(c, w)

	manifold := ManifoldConfig{}

	var factory servicefactory.ControllerServiceFactory
	err = manifold.output(w, &factory)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *manifoldSuite) TestOutputServiceFactoryGetter(c *gc.C) {
	w, err := NewWorker(Config{
		DBDeleter:                   s.dbDeleter,
		DBGetter:                    s.dbGetter,
		Logger:                      s.logger,
		NewServiceFactoryGetter:     NewServiceFactoryGetter,
		NewControllerServiceFactory: NewControllerServiceFactory,
		NewModelServiceFactory:      NewModelServiceFactory,
	})
	c.Assert(err, jc.ErrorIsNil)
	defer workertest.DirtyKill(c, w)

	manifold := ManifoldConfig{}

	var factory servicefactory.ServiceFactoryGetter
	err = manifold.output(w, &factory)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *manifoldSuite) TestOutputInvalid(c *gc.C) {
	w, err := NewWorker(Config{
		DBDeleter:                   s.dbDeleter,
		DBGetter:                    s.dbGetter,
		Logger:                      s.logger,
		NewServiceFactoryGetter:     NewServiceFactoryGetter,
		NewControllerServiceFactory: NewControllerServiceFactory,
		NewModelServiceFactory:      NewModelServiceFactory,
	})
	c.Assert(err, jc.ErrorIsNil)
	defer workertest.DirtyKill(c, w)

	manifold := ManifoldConfig{}

	var factory struct{}
	err = manifold.output(w, &factory)
	c.Assert(err, gc.ErrorMatches, `unsupported output type .*`)
}

func (s *manifoldSuite) TestNewControllerServiceFactory(c *gc.C) {
	factory := NewControllerServiceFactory(s.dbGetter, s.dbDeleter, s.logger)
	c.Assert(factory, gc.NotNil)
}

func (s *manifoldSuite) TestNewModelServiceFactory(c *gc.C) {
	factory := NewModelServiceFactory(
		"model",
		s.dbGetter,
		s.logger,
	)
	c.Assert(factory, gc.NotNil)
}

func (s *manifoldSuite) TestNewServiceFactoryGetter(c *gc.C) {
	ctrlFactory := NewControllerServiceFactory(s.dbGetter, s.dbDeleter, s.logger)
	factory := NewServiceFactoryGetter(ctrlFactory, s.dbGetter, s.logger, NewModelServiceFactory)
	c.Assert(factory, gc.NotNil)

	modelFactory := factory.FactoryForModel("model")
	c.Assert(modelFactory, gc.NotNil)
}

func (s *manifoldSuite) getConfig() ManifoldConfig {
	return ManifoldConfig{
		DBAccessorName:   "dbaccessor",
		ChangeStreamName: "changestream",
		Logger:           s.logger,
		NewWorker: func(Config) (worker.Worker, error) {
			return nil, nil
		},
		NewServiceFactoryGetter: func(servicefactory.ControllerServiceFactory, changestream.WatchableDBGetter, Logger, ModelServiceFactoryFn) servicefactory.ServiceFactoryGetter {
			return nil
		},
		NewControllerServiceFactory: func(changestream.WatchableDBGetter, coredatabase.DBDeleter, Logger) servicefactory.ControllerServiceFactory {
			return nil
		},
		NewModelServiceFactory: func(model.UUID, changestream.WatchableDBGetter, Logger) servicefactory.ModelServiceFactory {
			return nil
		},
	}
}
