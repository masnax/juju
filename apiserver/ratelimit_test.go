// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package apiserver_test

import (
	"time"

	"github.com/juju/clock/testclock"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/api"
	corecontroller "github.com/juju/juju/controller"
	jujutesting "github.com/juju/juju/juju/testing"
	"github.com/juju/juju/testing/factory"
)

type rateLimitSuite struct {
	jujutesting.ApiServerSuite
}

var _ = gc.Suite(&rateLimitSuite{})

func (s *rateLimitSuite) SetUpTest(c *gc.C) {
	s.Clock = testclock.NewDilatedWallClock(time.Second)
	s.ControllerConfigAttrs = map[string]interface{}{
		corecontroller.AgentRateLimitMax:  1,
		corecontroller.AgentRateLimitRate: (60 * time.Second).String(),
	}
	s.ApiServerSuite.SetUpTest(c)
}

func (s *rateLimitSuite) infoForNewMachine(c *gc.C, info *api.Info) *api.Info {
	// Make a copy
	newInfo := *info

	f, release := s.NewFactory(c, info.ModelTag.Id())
	defer release()
	machine, password := f.MakeMachineReturningPassword(
		c, &factory.MachineParams{Nonce: "fake_nonce"})

	newInfo.Tag = machine.Tag()
	newInfo.Password = password
	newInfo.Nonce = "fake_nonce"
	return &newInfo
}

func (s *rateLimitSuite) infoForNewUser(c *gc.C, info *api.Info) *api.Info {
	// Make a copy
	newInfo := *info

	f, release := s.NewFactory(c, info.ModelTag.Id())
	defer release()
	password := "shhh..."
	user := f.MakeUser(c, &factory.UserParams{
		Password: password,
	})

	newInfo.Tag = user.Tag()
	newInfo.Password = password
	return &newInfo
}

func (s *rateLimitSuite) TestRateLimitAgents(c *gc.C) {
	c.Assert(s.Server.Report(), jc.DeepEquals, map[string]interface{}{
		"agent-ratelimit-max":  1,
		"agent-ratelimit-rate": 60 * time.Second,
	})

	info := s.ControllerModelApiInfo()
	// First agent connection is fine.
	machine1 := s.infoForNewMachine(c, info)
	conn1, err := api.Open(machine1, fastDialOpts)
	c.Assert(err, jc.ErrorIsNil)
	defer conn1.Close()

	// Second machine in the same minute gets told to go away and try again.
	machine2 := s.infoForNewMachine(c, info)
	_, err = api.Open(machine2, fastDialOpts)
	c.Assert(err, gc.ErrorMatches, `try again \(try again\)`)

	// If we wait a minute and try again, it is fine.
	s.Clock.Advance(time.Minute)
	conn2, err := api.Open(machine2, fastDialOpts)
	c.Assert(err, jc.ErrorIsNil)
	defer conn2.Close()

	// And the next one is limited.
	machine3 := s.infoForNewMachine(c, info)
	_, err = api.Open(machine3, fastDialOpts)
	c.Assert(err, gc.ErrorMatches, `try again \(try again\)`)
}

func (s *rateLimitSuite) TestRateLimitNotApplicableToUsers(c *gc.C) {
	info := s.ControllerModelApiInfo()

	// First agent connection is fine.
	machine1 := s.infoForNewMachine(c, info)
	conn1, err := api.Open(machine1, fastDialOpts)
	c.Assert(err, jc.ErrorIsNil)
	defer conn1.Close()

	// User connections are fine.
	user := s.infoForNewUser(c, info)
	conn2, err := api.Open(user, fastDialOpts)
	c.Assert(err, jc.ErrorIsNil)
	defer conn2.Close()

	user2 := s.infoForNewUser(c, info)
	conn3, err := api.Open(user2, fastDialOpts)
	c.Assert(err, jc.ErrorIsNil)
	defer conn3.Close()
}
