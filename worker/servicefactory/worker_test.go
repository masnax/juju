// Copyright 2022 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package servicefactory

import (
	"github.com/juju/errors"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/worker/v3"
	"github.com/juju/worker/v3/workertest"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/core/changestream"
	coredatabase "github.com/juju/juju/core/database"
	"github.com/juju/juju/domain/model"
	"github.com/juju/juju/internal/servicefactory"
)

type workerSuite struct {
	baseSuite
}

var _ = gc.Suite(&workerSuite{})

func (s *workerSuite) TestValidateConfig(c *gc.C) {
	defer s.setupMocks(c).Finish()

	cfg := s.getConfig()
	c.Check(cfg.Validate(), jc.ErrorIsNil)

	cfg = s.getConfig()
	cfg.Logger = nil
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)

	cfg = s.getConfig()
	cfg.DBDeleter = nil
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)

	cfg = s.getConfig()
	cfg.DBGetter = nil
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

func (s *workerSuite) getConfig() Config {
	return Config{
		DBGetter:  s.dbGetter,
		DBDeleter: s.dbDeleter,
		Logger:    s.logger,
		NewServiceFactoryGetter: func(servicefactory.ControllerServiceFactory, changestream.WatchableDBGetter, Logger, ModelServiceFactoryFn) servicefactory.ServiceFactoryGetter {
			return s.serviceFactoryGetter
		},
		NewControllerServiceFactory: func(changestream.WatchableDBGetter, coredatabase.DBDeleter, Logger) servicefactory.ControllerServiceFactory {
			return s.controllerServiceFactory
		},
		NewModelServiceFactory: func(model.UUID, changestream.WatchableDBGetter, Logger) servicefactory.ModelServiceFactory {
			return s.modelServiceFactory
		},
	}
}

func (s *workerSuite) TestWorkerControllerFactory(c *gc.C) {
	defer s.setupMocks(c).Finish()

	w := s.newWorker(c)
	defer workertest.CleanKill(c, w)

	srvFact, ok := w.(*serviceFactoryWorker)
	c.Assert(ok, jc.IsTrue, gc.Commentf("worker does not implement serviceFactoryWorker"))

	factory := srvFact.ControllerFactory()
	c.Assert(factory, gc.NotNil)

	workertest.CleanKill(c, w)
}

func (s *workerSuite) TestWorkerFactoryGetter(c *gc.C) {
	defer s.setupMocks(c).Finish()

	w := s.newWorker(c)
	defer workertest.CleanKill(c, w)

	srvFact, ok := w.(*serviceFactoryWorker)
	c.Assert(ok, jc.IsTrue, gc.Commentf("worker does not implement serviceFactoryWorker"))

	factory := srvFact.FactoryGetter()
	c.Assert(factory, gc.NotNil)

	workertest.CleanKill(c, w)
}

func (s *workerSuite) newWorker(c *gc.C) worker.Worker {
	w, err := NewWorker(s.getConfig())
	c.Assert(err, jc.ErrorIsNil)
	return w
}
