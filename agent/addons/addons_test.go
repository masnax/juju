// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package addons_test

import (
	"runtime"
	"time"

	"github.com/juju/clock"
	"github.com/juju/errors"
	"github.com/juju/loggo"
	"github.com/juju/names/v4"
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/worker/v3"
	"github.com/juju/worker/v3/dependency"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/mock/gomock"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/agent/addons"
	coretesting "github.com/juju/juju/testing"
	"github.com/juju/juju/worker/introspection"
)

type introspectionSuite struct {
	testing.IsolationSuite
}

var _ = gc.Suite(&introspectionSuite{})

func (s *introspectionSuite) TestStartNonLinux(c *gc.C) {
	if runtime.GOOS == "linux" {
		c.Skip("testing for non-linux")
	}
	var started bool

	cfg := addons.IntrospectionConfig{
		WorkerFunc: func(_ introspection.Config) (worker.Worker, error) {
			started = true
			return nil, errors.New("shouldn't call start")
		},
		Logger: coretesting.CheckLogger{Log: c},
	}

	err := addons.StartIntrospection(cfg)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(started, jc.IsFalse)
}

func (s *introspectionSuite) TestStartError(c *gc.C) {
	if runtime.GOOS != "linux" {
		c.Skip("introspection worker not supported on non-linux")
	}

	cfg := addons.IntrospectionConfig{
		AgentTag:      names.NewMachineTag("42"),
		NewSocketName: addons.DefaultIntrospectionSocketName,
		WorkerFunc: func(_ introspection.Config) (worker.Worker, error) {
			return nil, errors.New("boom")
		},
		Logger: coretesting.CheckLogger{Log: c},
	}

	err := addons.StartIntrospection(cfg)
	c.Check(err, gc.ErrorMatches, "boom")
}

func (s *introspectionSuite) TestStartSuccess(c *gc.C) {
	if runtime.GOOS != "linux" {
		c.Skip("introspection worker not supported on non-linux")
	}
	fake := &dummyWorker{
		done: make(chan struct{}),
	}

	config := dependency.EngineConfig{
		IsFatal:    func(err error) bool { return false },
		WorstError: func(err1, err2 error) error { return err1 },
		Clock:      clock.WallClock,
		Metrics:    dependency.DefaultMetrics(),
		Logger:     loggo.GetLogger("juju.worker.dependency"),
	}
	engine, err := dependency.NewEngine(config)
	c.Assert(err, jc.ErrorIsNil)

	cfg := addons.IntrospectionConfig{
		AgentTag:      names.NewMachineTag("42"),
		Engine:        engine,
		NewSocketName: func(tag names.Tag) string { return "bananas" },
		WorkerFunc: func(cfg introspection.Config) (worker.Worker, error) {
			fake.config = cfg
			return fake, nil
		},
		Logger: coretesting.CheckLogger{Log: c},
	}

	err = addons.StartIntrospection(cfg)
	c.Assert(err, jc.ErrorIsNil)

	c.Check(fake.config.DepEngine, gc.Equals, engine)
	c.Check(fake.config.SocketName, gc.Equals, "bananas")

	// Stopping the engine causes the introspection worker to stop.
	engine.Kill()

	select {
	case <-fake.done:
	case <-time.After(coretesting.LongWait):
		c.Fatalf("worker did not get stopped")
	}
}

func (s *introspectionSuite) TestDefaultIntrospectionSocketName(c *gc.C) {
	name := addons.DefaultIntrospectionSocketName(names.NewMachineTag("42"))
	c.Assert(name, gc.Equals, "jujud-machine-42")
}

type dummyWorker struct {
	config introspection.Config
	done   chan struct{}
}

func (d *dummyWorker) Kill() {
	close(d.done)
}

func (d *dummyWorker) Wait() error {
	<-d.done
	return nil
}

type registerSuite struct {
	testing.IsolationSuite
}

var _ = gc.Suite(&registerSuite{})

func (s *registerSuite) TestRegisterEngineMetrics(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	done := make(chan struct{}, 1)

	collector := dummyCollector{}

	registry := NewMockRegisterer(ctrl)
	registry.EXPECT().Register(collector)
	registry.EXPECT().Unregister(collector).Do(func(_ prometheus.Collector) {
		close(done)
	})
	sink := NewMockMetricSink(ctrl)
	sink.EXPECT().Unregister()

	worker := &dummyWorker{
		done: make(chan struct{}, 1),
	}

	err := addons.RegisterEngineMetrics(registry, collector, worker, sink)
	c.Assert(err, jc.ErrorIsNil)

	worker.Kill()

	select {
	case <-done:
	case <-time.After(testing.ShortWait):
	}
}

type dummyCollector struct{}

// Describe is part of the prometheus.Collector interface.
func (dummyCollector) Describe(ch chan<- *prometheus.Desc) {
}

// Collect is part of the prometheus.Collector interface.
func (dummyCollector) Collect(ch chan<- prometheus.Metric) {
}
