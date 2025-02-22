// Copyright 2018 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package lease

import (
	"github.com/juju/errors"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/worker/v3"
	gc "gopkg.in/check.v1"

	coredatabase "github.com/juju/juju/core/database"
	"github.com/juju/juju/core/lease"
)

type manifoldSuite struct {
	baseSuite
}

var _ = gc.Suite(&manifoldSuite{})

func (s *manifoldSuite) TestValidateConfig(c *gc.C) {
	defer s.setupMocks(c).Finish()

	cfg := s.getConfig()
	c.Check(cfg.Validate(), jc.ErrorIsNil)

	cfg.AgentName = ""
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)

	cfg.ClockName = ""
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)

	cfg.DBAccessorName = ""
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)

	cfg.TraceName = ""
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)

	cfg.Logger = nil
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)

	cfg.PrometheusRegisterer = nil
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)

	cfg.NewWorker = nil
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)

	cfg.NewStore = nil
	c.Check(cfg.Validate(), jc.ErrorIs, errors.NotValid)
}

func (s *manifoldSuite) getConfig() ManifoldConfig {
	return ManifoldConfig{
		AgentName:      "agent",
		ClockName:      "clock",
		DBAccessorName: "dbaccessor",
		TraceName:      "trace",

		Logger:               s.logger,
		PrometheusRegisterer: s.prometheusRegisterer,
		NewWorker: func(mc ManagerConfig) (worker.Worker, error) {
			return nil, nil
		},
		NewStore: func(coredatabase.DBGetter, Logger) lease.Store {
			return nil
		},
	}
}
