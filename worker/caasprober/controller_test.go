// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package caasprober_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	k8sconstants "github.com/juju/juju/caas/kubernetes/provider/constants"
	"github.com/juju/juju/internal/observability/probe"
	"github.com/juju/juju/worker/caasprober"
)

type ControllerSuite struct {
}

type dummyMux struct {
	AddHandlerFunc    func(string, string, http.Handler) error
	RemoveHandlerFunc func(string, string)
}

var _ = gc.Suite(&ControllerSuite{})

func (d *dummyMux) AddHandler(i, j string, h http.Handler) error {
	if d.AddHandlerFunc == nil {
		return nil
	}
	return d.AddHandlerFunc(i, j, h)
}

func (d *dummyMux) RemoveHandler(i, j string) {
	if d.RemoveHandlerFunc != nil {
		d.RemoveHandlerFunc(i, j)
	}
}

func (s *ControllerSuite) TestControllerMuxRegistration(c *gc.C) {
	var (
		livenessRegistered    = false
		livenessDeRegistered  = false
		readinessRegistered   = false
		readinessDeRegistered = false
		startupRegistered     = false
		startupDeRegistered   = false
		waitGroup             = sync.WaitGroup{}
	)

	waitGroup.Add(3)
	mux := dummyMux{
		AddHandlerFunc: func(m, p string, _ http.Handler) error {
			c.Check(m, gc.Equals, http.MethodGet)
			switch p {
			case k8sconstants.AgentHTTPPathLiveness:
				c.Check(livenessRegistered, jc.IsFalse)
				livenessRegistered = true
				waitGroup.Done()
			case k8sconstants.AgentHTTPPathReadiness:
				c.Check(readinessRegistered, jc.IsFalse)
				readinessRegistered = true
				waitGroup.Done()
			case k8sconstants.AgentHTTPPathStartup:
				c.Check(startupRegistered, jc.IsFalse)
				startupRegistered = true
				waitGroup.Done()
			default:
				c.Errorf("unknown path registered in controller: %s", p)
			}
			return nil
		},
		RemoveHandlerFunc: func(m, p string) {
			c.Check(m, gc.Equals, http.MethodGet)
			switch p {
			case k8sconstants.AgentHTTPPathLiveness:
				c.Check(livenessDeRegistered, jc.IsFalse)
				livenessDeRegistered = true
				waitGroup.Done()
			case k8sconstants.AgentHTTPPathReadiness:
				c.Check(readinessDeRegistered, jc.IsFalse)
				readinessDeRegistered = true
				waitGroup.Done()
			case k8sconstants.AgentHTTPPathStartup:
				c.Check(startupDeRegistered, jc.IsFalse)
				startupDeRegistered = true
				waitGroup.Done()
			default:
				c.Errorf("unknown path registered in controller: %s", p)
			}
		},
	}

	probes := caasprober.NewCAASProbes()
	probes.Liveness.Probes["test"] = probe.NotImplemented
	probes.Readiness.Probes["test"] = probe.NotImplemented
	probes.Startup.Probes["test"] = probe.NotImplemented

	controller, err := caasprober.NewController(probes, &mux)
	c.Assert(err, jc.ErrorIsNil)

	waitGroup.Wait()
	waitGroup.Add(3)
	controller.Kill()

	waitGroup.Wait()
	err = controller.Wait()
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(livenessRegistered, jc.IsTrue)
	c.Assert(livenessDeRegistered, jc.IsTrue)
	c.Assert(readinessRegistered, jc.IsTrue)
	c.Assert(readinessDeRegistered, jc.IsTrue)
	c.Assert(startupRegistered, jc.IsTrue)
	c.Assert(startupDeRegistered, jc.IsTrue)
}

func (s *ControllerSuite) TestControllerNotImplemented(c *gc.C) {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(3)

	mux := dummyMux{
		AddHandlerFunc: func(m, p string, h http.Handler) error {
			req := httptest.NewRequest(m, p, nil)
			recorder := httptest.NewRecorder()
			h.ServeHTTP(recorder, req)
			c.Check(recorder.Result().StatusCode, gc.Equals, http.StatusNotImplemented)
			waitGroup.Done()
			return nil
		},
		RemoveHandlerFunc: func(m, p string) {},
	}

	probes := caasprober.NewCAASProbes()
	probes.Liveness.Probes["test"] = probe.NotImplemented
	probes.Readiness.Probes["test"] = probe.NotImplemented
	probes.Startup.Probes["test"] = probe.NotImplemented

	controller, err := caasprober.NewController(probes, &mux)
	c.Assert(err, jc.ErrorIsNil)

	waitGroup.Wait()
	controller.Kill()
	err = controller.Wait()
	c.Assert(err, jc.ErrorIsNil)
}

func (s *ControllerSuite) TestControllerProbeError(c *gc.C) {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(3)

	mux := dummyMux{
		AddHandlerFunc: func(m, p string, h http.Handler) error {
			req := httptest.NewRequest(m, p, nil)
			recorder := httptest.NewRecorder()
			h.ServeHTTP(recorder, req)
			c.Check(recorder.Result().StatusCode, gc.Equals, http.StatusInternalServerError)
			waitGroup.Done()
			return nil
		},
		RemoveHandlerFunc: func(m, p string) {},
	}

	probeErr := probe.ProberFn(func() (bool, error) {
		return false, errors.New("test error")
	})

	probes := caasprober.NewCAASProbes()
	probes.Liveness.Probes["test"] = probeErr
	probes.Readiness.Probes["test"] = probeErr
	probes.Startup.Probes["test"] = probeErr

	controller, err := caasprober.NewController(probes, &mux)
	c.Assert(err, jc.ErrorIsNil)

	waitGroup.Wait()
	controller.Kill()
	err = controller.Wait()
	c.Assert(err, jc.ErrorIsNil)
}

func (s *ControllerSuite) TestControllerProbeFail(c *gc.C) {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(3)

	mux := dummyMux{
		AddHandlerFunc: func(m, p string, h http.Handler) error {
			req := httptest.NewRequest(m, p, nil)
			recorder := httptest.NewRecorder()
			h.ServeHTTP(recorder, req)
			c.Check(recorder.Result().StatusCode, gc.Equals, http.StatusTeapot)
			waitGroup.Done()
			return nil
		},
		RemoveHandlerFunc: func(m, p string) {},
	}

	probeFail := probe.ProberFn(func() (bool, error) {
		return false, nil
	})

	probes := caasprober.NewCAASProbes()
	probes.Liveness.Probes["test"] = probeFail
	probes.Readiness.Probes["test"] = probeFail
	probes.Startup.Probes["test"] = probeFail

	controller, err := caasprober.NewController(probes, &mux)
	c.Assert(err, jc.ErrorIsNil)

	waitGroup.Wait()
	controller.Kill()
	err = controller.Wait()
	c.Assert(err, jc.ErrorIsNil)
}

func (s *ControllerSuite) TestControllerProbePass(c *gc.C) {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(3)

	mux := dummyMux{
		AddHandlerFunc: func(m, p string, h http.Handler) error {
			req := httptest.NewRequest(m, p, nil)
			recorder := httptest.NewRecorder()
			h.ServeHTTP(recorder, req)
			c.Check(recorder.Result().StatusCode, gc.Equals, http.StatusOK)
			waitGroup.Done()
			return nil
		},
		RemoveHandlerFunc: func(m, p string) {},
	}

	probes := caasprober.NewCAASProbes()
	probes.Liveness.Probes["test"] = probe.Success
	probes.Readiness.Probes["test"] = probe.Success
	probes.Startup.Probes["test"] = probe.Success

	controller, err := caasprober.NewController(probes, &mux)
	c.Assert(err, jc.ErrorIsNil)

	waitGroup.Wait()
	controller.Kill()
	err = controller.Wait()
	c.Assert(err, jc.ErrorIsNil)
}

func (s *ControllerSuite) TestControllerProbePassDetailed(c *gc.C) {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(3)

	mux := dummyMux{
		AddHandlerFunc: func(m, p string, h http.Handler) error {
			req := httptest.NewRequest(m, p+"?detailed=true", nil)
			recorder := httptest.NewRecorder()
			h.ServeHTTP(recorder, req)
			c.Check(recorder.Result().StatusCode, gc.Equals, http.StatusOK)
			c.Check(recorder.Body.String(), jc.HasPrefix, `+ test
OK: probe`)
			waitGroup.Done()
			return nil
		},
		RemoveHandlerFunc: func(m, p string) {},
	}

	probes := caasprober.NewCAASProbes()
	probes.Liveness.Probes["test"] = probe.Success
	probes.Readiness.Probes["test"] = probe.Success
	probes.Startup.Probes["test"] = probe.Success

	controller, err := caasprober.NewController(probes, &mux)
	c.Assert(err, jc.ErrorIsNil)

	waitGroup.Wait()
	controller.Kill()
	err = controller.Wait()
	c.Assert(err, jc.ErrorIsNil)
}

func (s *ControllerSuite) TestControllerProbeFailDetailed(c *gc.C) {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(3)

	mux := dummyMux{
		AddHandlerFunc: func(m, p string, h http.Handler) error {
			req := httptest.NewRequest(m, p+"?detailed=true", nil)
			recorder := httptest.NewRecorder()
			h.ServeHTTP(recorder, req)
			c.Check(recorder.Result().StatusCode, gc.Equals, http.StatusOK)
			c.Check(recorder.Body.String(), jc.HasPrefix, `- test: test error
Internal Server Error: probe`)
			waitGroup.Done()
			return nil
		},
		RemoveHandlerFunc: func(m, p string) {},
	}

	probeFail := probe.ProberFn(func() (bool, error) {
		return false, errors.New("test error")
	})

	probes := caasprober.NewCAASProbes()
	probes.Liveness.Probes["test"] = probeFail
	probes.Readiness.Probes["test"] = probeFail
	probes.Startup.Probes["test"] = probeFail

	controller, err := caasprober.NewController(probes, &mux)
	c.Assert(err, jc.ErrorIsNil)

	waitGroup.Wait()
	controller.Kill()
	err = controller.Wait()
	c.Assert(err, jc.ErrorIsNil)
}
