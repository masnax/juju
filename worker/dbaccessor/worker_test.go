// Copyright 2022 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package dbaccessor

import (
	"context"
	"time"

	jc "github.com/juju/testing/checkers"
	"github.com/juju/worker/v3"
	"github.com/juju/worker/v3/dependency"
	"github.com/juju/worker/v3/workertest"
	"go.uber.org/mock/gomock"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/core/database"
	"github.com/juju/juju/internal/database/dqlite"
	"github.com/juju/juju/internal/pubsub/apiserver"
	"github.com/juju/juju/testing"
)

type workerSuite struct {
	baseSuite

	trackedDB *MockTrackedDB
}

var _ = gc.Suite(&workerSuite{})

func (s *workerSuite) TestKilledGetDBErrDying(c *gc.C) {
	defer s.setupMocks(c).Finish()

	dbDone := make(chan struct{})
	s.expectClock()
	s.expectTrackedDBUpdateNodeAndKill(dbDone)

	mgrExp := s.nodeManager.EXPECT()
	mgrExp.EnsureDataDir().Return(c.MkDir(), nil)
	mgrExp.IsExistingNode().Return(true, nil).Times(1)
	mgrExp.IsLoopbackBound(gomock.Any()).Return(true, nil).Times(2)
	mgrExp.WithLogFuncOption().Return(nil)
	mgrExp.WithTracingOption().Return(nil)

	// We may or may not get this call.
	mgrExp.SetClusterToLocalNode(gomock.Any()).Return(nil).AnyTimes()

	s.expectNodeStartupAndShutdown()

	s.client.EXPECT().Cluster(gomock.Any()).Return(nil, nil)

	s.hub.EXPECT().Subscribe(apiserver.DetailsTopic, gomock.Any()).Return(func() {}, nil)

	w := s.newWorker(c)
	defer func() {
		close(dbDone)
		workertest.DirtyKill(c, w)
	}()
	dbw := w.(*dbWorker)
	ensureStartup(c, dbw)

	w.Kill()

	_, err := dbw.GetDB("anything")
	c.Assert(err, jc.ErrorIs, database.ErrDBAccessorDying)
}

func (s *workerSuite) TestStartupTimeoutSingleControllerReconfigure(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.expectClock()

	mgrExp := s.nodeManager.EXPECT()
	mgrExp.EnsureDataDir().Return(c.MkDir(), nil)
	mgrExp.IsExistingNode().Return(true, nil).Times(2)
	mgrExp.IsLoopbackBound(gomock.Any()).Return(false, nil).Times(3)
	mgrExp.WithTLSOption().Return(nil, nil)
	mgrExp.WithLogFuncOption().Return(nil)
	mgrExp.WithTracingOption().Return(nil)
	mgrExp.SetClusterToLocalNode(gomock.Any()).Return(nil)

	// App gets started, we time out waiting, then we close it.
	appExp := s.dbApp.EXPECT()
	appExp.Ready(gomock.Any()).Return(context.DeadlineExceeded)
	appExp.Close().Return(nil)

	// We expect to request API details.
	s.hub.EXPECT().Subscribe(apiserver.DetailsTopic, gomock.Any()).Return(func() {}, nil)
	s.hub.EXPECT().Publish(apiserver.DetailsRequestTopic, gomock.Any()).Return(func() {}, nil)

	w := s.newWorker(c)
	defer workertest.DirtyKill(c, w)

	// Topology is just us. We should reconfigure the node and shut down.
	select {
	case w.(*dbWorker).apiServerChanges <- apiserver.Details{
		Servers: map[string]apiserver.APIServer{"0": {ID: "0", InternalAddress: "10.6.6.6:1234"}},
	}:
	case <-time.After(testing.LongWait):
		c.Fatal("timed out waiting for cluster change to be processed")
	}

	err := workertest.CheckKilled(c, w)
	c.Assert(err, jc.ErrorIs, dependency.ErrBounce)
}

func (s *workerSuite) TestStartupTimeoutMultipleControllerRetry(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.expectClock()

	mgrExp := s.nodeManager.EXPECT()
	mgrExp.EnsureDataDir().Return(c.MkDir(), nil).Times(2)
	mgrExp.IsExistingNode().Return(true, nil).Times(2)
	mgrExp.IsLoopbackBound(gomock.Any()).Return(false, nil).Times(4)

	// We expect 2 attempts to start.
	mgrExp.WithTLSOption().Return(nil, nil).Times(2)
	mgrExp.WithLogFuncOption().Return(nil).Times(2)
	mgrExp.WithTracingOption().Return(nil).Times(2)

	// App gets started, we time out waiting, then we close it both times.
	appExp := s.dbApp.EXPECT()
	appExp.Ready(gomock.Any()).Return(context.DeadlineExceeded).Times(2)
	appExp.Close().Return(nil).Times(2)

	// We expect to request API details.
	s.hub.EXPECT().Subscribe(apiserver.DetailsTopic, gomock.Any()).Return(func() {}, nil)
	s.hub.EXPECT().Publish(apiserver.DetailsRequestTopic, gomock.Any()).Return(func() {}, nil).Times(2)

	w := s.newWorker(c)
	defer workertest.CleanKill(c, w)
	dbw := w.(*dbWorker)

	// If there are multiple servers reported, we can't reason about our
	// current state in a discrete fashion. The worker throws an error.
	select {
	case dbw.apiServerChanges <- apiserver.Details{
		Servers: map[string]apiserver.APIServer{
			"0": {ID: "0", InternalAddress: "10.6.6.6:1234"},
			"1": {ID: "1", InternalAddress: "10.6.6.7:1234"},
		},
	}:
	case <-time.After(testing.LongWait):
		c.Fatal("timed out waiting for cluster change to be processed")
	}

	// At this point, the Dqlite node is not started.
	// The worker is waiting for legitimate server detail messages.
	select {
	case <-dbw.dbReady:
		c.Fatal("Dqlite node should not be started yet.")
	case <-time.After(testing.ShortWait):
	}
}

func (s *workerSuite) TestStartupNotExistingNodeThenCluster(c *gc.C) {
	defer s.setupMocks(c).Finish()

	dbDone := make(chan struct{})
	s.expectClock()
	s.expectTrackedDBUpdateNodeAndKill(dbDone)

	mgrExp := s.nodeManager.EXPECT()
	mgrExp.EnsureDataDir().Return(c.MkDir(), nil)
	mgrExp.IsExistingNode().Return(false, nil).Times(4)
	mgrExp.WithAddressOption("10.6.6.6").Return(nil)
	mgrExp.WithClusterOption([]string{"10.6.6.7"}).Return(nil)
	mgrExp.WithLogFuncOption().Return(nil)
	mgrExp.WithTLSOption().Return(nil, nil)
	mgrExp.WithTracingOption().Return(nil)
	mgrExp.IsLoopbackBound(gomock.Any()).Return(false, nil)

	s.client.EXPECT().Cluster(gomock.Any()).Return(nil, nil)

	s.expectNodeStartupAndShutdown()
	s.dbApp.EXPECT().Handover(gomock.Any()).Return(nil)

	// When we are starting up as a new node,
	// we request details immediately.
	s.hub.EXPECT().Subscribe(apiserver.DetailsTopic, gomock.Any()).Return(func() {}, nil)
	s.hub.EXPECT().Publish(apiserver.DetailsRequestTopic, gomock.Any()).Return(func() {}, nil)

	w := s.newWorker(c)
	defer func() {
		close(dbDone)
		workertest.CleanKill(c, w)
	}()
	dbw := w.(*dbWorker)

	// Without a bind address for ourselves we keep waiting.
	select {
	case dbw.apiServerChanges <- apiserver.Details{
		Servers: map[string]apiserver.APIServer{
			"0": {ID: "0"},
			"1": {ID: "1", InternalAddress: "10.6.6.7:1234"},
		},
	}:
	case <-time.After(testing.LongWait):
		c.Fatal("timed out waiting for cluster change to be processed")
	}

	// Without other cluster members we keep waiting.
	select {
	case w.(*dbWorker).apiServerChanges <- apiserver.Details{
		Servers: map[string]apiserver.APIServer{
			"0": {ID: "0", InternalAddress: "10.6.6.6:1234"},
		},
	}:
	case <-time.After(testing.LongWait):
		c.Fatal("timed out waiting for cluster change to be processed")
	}

	// At this point, the Dqlite node is not started.
	// The worker is waiting for legitimate server detail messages.
	select {
	case <-dbw.dbReady:
		c.Fatal("Dqlite node should not be started yet.")
	case <-time.After(testing.ShortWait):
	}

	// Push a message onto the API details channel,
	// enabling node startup as a cluster member.
	select {
	case w.(*dbWorker).apiServerChanges <- apiserver.Details{
		Servers: map[string]apiserver.APIServer{
			"0": {ID: "0", InternalAddress: "10.6.6.6:1234"},
			"1": {ID: "1", InternalAddress: "10.6.6.7:1234"},
		},
	}:
	case <-time.After(testing.LongWait):
		c.Fatal("timed out waiting for cluster change to be processed")
	}

	ensureStartup(c, dbw)

	s.client.EXPECT().Leader(gomock.Any()).Return(&dqlite.NodeInfo{
		ID:      1,
		Address: "10.10.1.1",
	}, nil)
	report := w.(interface{ Report() map[string]any }).Report()
	c.Assert(report, MapHasKeys, []string{
		"leader",
		"leader-id",
		"leader-role",
	})
}

func (s *workerSuite) TestWorkerStartupExistingNode(c *gc.C) {
	defer s.setupMocks(c).Finish()

	dbDone := make(chan struct{})
	s.expectClock()
	s.expectTrackedDBUpdateNodeAndKill(dbDone)

	mgrExp := s.nodeManager.EXPECT()
	mgrExp.EnsureDataDir().Return(c.MkDir(), nil)

	// If this is an existing node, we do not invoke the address or cluster
	// options, but if the node is not as bootstrapped, we do assume it is
	// part of a cluster, and uses the TLS option.
	// IsBootstrapped node is called twice - once to check the startup
	// conditions and then again upon worker shutdown.
	mgrExp.IsExistingNode().Return(true, nil)
	mgrExp.IsLoopbackBound(gomock.Any()).Return(false, nil).Times(2)
	mgrExp.WithLogFuncOption().Return(nil)
	mgrExp.WithTLSOption().Return(nil, nil)
	mgrExp.WithTracingOption().Return(nil)

	s.client.EXPECT().Cluster(gomock.Any()).Return(nil, nil)

	s.expectNodeStartupAndShutdown()
	s.dbApp.EXPECT().Handover(gomock.Any()).Return(nil)

	s.hub.EXPECT().Subscribe(apiserver.DetailsTopic, gomock.Any()).Return(func() {}, nil)

	w := s.newWorker(c)
	defer func() {
		close(dbDone)
		workertest.CleanKill(c, w)
	}()

	ensureStartup(c, w.(*dbWorker))
}

func (s *workerSuite) TestWorkerStartupAsBootstrapNodeSingleServerNoRebind(c *gc.C) {
	defer s.setupMocks(c).Finish()

	dbDone := make(chan struct{})
	s.expectClock()
	s.expectTrackedDBUpdateNodeAndKill(dbDone)

	dataDir := c.MkDir()
	mgrExp := s.nodeManager.EXPECT()
	mgrExp.EnsureDataDir().Return(dataDir, nil).MinTimes(1)

	// If this is an existing node, we do not
	// invoke the address or cluster options.
	mgrExp.IsExistingNode().Return(true, nil).Times(3)
	mgrExp.IsLoopbackBound(gomock.Any()).Return(true, nil).Times(4)
	mgrExp.WithLogFuncOption().Return(nil)
	mgrExp.WithTracingOption().Return(nil)

	s.client.EXPECT().Cluster(gomock.Any()).Return(nil, nil)

	s.expectNodeStartupAndShutdown()

	s.hub.EXPECT().Subscribe(apiserver.DetailsTopic, gomock.Any()).Return(func() {}, nil)

	w := s.newWorker(c)
	defer func() {
		close(dbDone)
		workertest.CleanKill(c, w)
	}()
	dbw := w.(*dbWorker)

	ensureStartup(c, dbw)

	// At this point we have started successfully.
	// Push a message onto the API details channel.
	// A single server does not cause a binding change.
	select {
	case dbw.apiServerChanges <- apiserver.Details{
		Servers: map[string]apiserver.APIServer{
			"0": {ID: "0", InternalAddress: "10.6.6.6:1234"},
		},
	}:
	case <-time.After(testing.LongWait):
		c.Fatal("timed out waiting for cluster change to be processed")
	}

	// Multiple servers still do not cause a binding change
	// if there is no internal address to bind to.
	select {
	case dbw.apiServerChanges <- apiserver.Details{
		Servers: map[string]apiserver.APIServer{
			"0": {ID: "0"},
			"1": {ID: "1", InternalAddress: "10.6.6.7:1234"},
			"2": {ID: "2", InternalAddress: "10.6.6.8:1234"},
		},
	}:
	case <-time.After(testing.LongWait):
		c.Fatal("timed out waiting for cluster change to be processed")
	}
}

func (s *workerSuite) TestWorkerStartupAsBootstrapNodeThenReconfigure(c *gc.C) {
	defer s.setupMocks(c).Finish()

	dbDone := make(chan struct{})
	s.expectClock()
	s.expectTrackedDBUpdateNodeAndKill(dbDone)

	dataDir := c.MkDir()
	mgrExp := s.nodeManager.EXPECT()
	mgrExp.EnsureDataDir().Return(dataDir, nil).MinTimes(1)

	// If this is an existing node, we do not
	// invoke the address or cluster options.
	mgrExp.IsExistingNode().Return(true, nil).Times(2)
	gomock.InOrder(
		mgrExp.IsLoopbackBound(gomock.Any()).Return(true, nil).Times(2),
		// This is the check at shutdown.
		mgrExp.IsLoopbackBound(gomock.Any()).Return(false, nil))
	mgrExp.WithLogFuncOption().Return(nil)
	mgrExp.WithTracingOption().Return(nil)

	// These are the expectations around reconfiguring
	// the cluster and local node.
	mgrExp.ClusterServers(gomock.Any()).Return([]dqlite.NodeInfo{
		{
			ID:      3297041220608546238,
			Address: "127.0.0.1:17666",
			Role:    0,
		},
	}, nil)
	mgrExp.SetClusterServers(gomock.Any(), []dqlite.NodeInfo{
		{
			ID:      3297041220608546238,
			Address: "10.6.6.6:17666",
			Role:    0,
		},
	}).Return(nil)
	mgrExp.SetNodeInfo(dqlite.NodeInfo{
		ID:      3297041220608546238,
		Address: "10.6.6.6:17666",
		Role:    0,
	}).Return(nil)

	s.client.EXPECT().Cluster(gomock.Any()).Return(nil, nil)

	// Although the shut-down check for IsLoopbackBound returns false,
	// this call to shut-down is actually run before reconfiguring the node.
	// When the loop exits, the node is already set to nil.
	s.expectNodeStartupAndShutdown()

	s.hub.EXPECT().Subscribe(apiserver.DetailsTopic, gomock.Any()).Return(func() {}, nil)

	w := s.newWorker(c)
	defer func() {
		close(dbDone)
		err := workertest.CheckKilled(c, w)
		c.Assert(err, jc.ErrorIs, dependency.ErrBounce)
	}()
	dbw := w.(*dbWorker)

	ensureStartup(c, dbw)

	// At this point we have started successfully.
	// Push a message onto the API details channel to simulate a move into HA.
	select {
	case dbw.apiServerChanges <- apiserver.Details{
		Servers: map[string]apiserver.APIServer{
			"0": {ID: "0", InternalAddress: "10.6.6.6:1234"},
			"1": {ID: "1", InternalAddress: "10.6.6.7:1234"},
			"2": {ID: "2", InternalAddress: "10.6.6.8:1234"},
		},
	}:
	case <-time.After(testing.LongWait):
		c.Fatal("timed out waiting for cluster change to be processed")
	}
}

func (s *workerSuite) newWorker(c *gc.C) worker.Worker {
	return s.newWorkerWithDB(c, s.trackedDB)
}

func (s *workerSuite) setupMocks(c *gc.C) *gomock.Controller {
	ctrl := s.baseSuite.setupMocks(c)

	s.trackedDB = NewMockTrackedDB(ctrl)

	return ctrl
}

// expectTrackedDBKillUpdateNode encompasses:
// - Use by the controller node service to update the node info.
// - Kill and wait upon termination of the worker.
// The input channel is used to ensure that the runner call to Wait does not
// return until we are ready.
// The Kill expectation is soft, because this can be done via parent catacomb,
// rather than a direct call.
func (s *workerSuite) expectTrackedDBUpdateNodeAndKill(done chan struct{}) {
	s.trackedDB.EXPECT().StdTxn(gomock.Any(), gomock.Any())
	s.trackedDB.EXPECT().Kill().AnyTimes()
	s.trackedDB.EXPECT().Wait().DoAndReturn(func() error {
		<-done
		return nil
	})
}
