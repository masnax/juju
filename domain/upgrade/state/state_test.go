// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"
	"database/sql"

	"github.com/canonical/sqlair"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/utils/v3"
	"github.com/juju/version/v2"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/core/upgrade"
	schematesting "github.com/juju/juju/domain/schema/testing"
	domainupgrade "github.com/juju/juju/domain/upgrade"
	upgradeerrors "github.com/juju/juju/domain/upgrade/errors"
	"github.com/juju/juju/internal/database"
)

type stateSuite struct {
	schematesting.ControllerSuite

	st *State

	upgradeUUID domainupgrade.UUID
}

var _ = gc.Suite(&stateSuite{})

func (s *stateSuite) SetUpTest(c *gc.C) {
	s.ControllerSuite.SetUpTest(c)
	s.st = NewState(s.TxnRunnerFactory())

	// Add a completed upgrade before tests start
	uuid, err := s.st.CreateUpgrade(context.Background(), version.MustParse("2.9.42"), version.MustParse("3.0.0"))
	c.Assert(err, jc.ErrorIsNil)
	err = s.st.SetControllerReady(context.Background(), uuid, "0")
	c.Assert(err, jc.ErrorIsNil)
	err = s.st.StartUpgrade(context.Background(), uuid)
	c.Assert(err, jc.ErrorIsNil)
	err = s.st.SetDBUpgradeCompleted(context.Background(), uuid)
	c.Assert(err, jc.ErrorIsNil)
	err = s.st.SetControllerDone(context.Background(), uuid, "0")
	c.Assert(err, jc.ErrorIsNil)

	s.upgradeUUID = uuid
}

func (s *stateSuite) TestCreateUpgrade(c *gc.C) {
	uuid, err := s.st.CreateUpgrade(context.Background(), version.MustParse("3.0.0"), version.MustParse("3.0.1"))
	c.Assert(err, jc.ErrorIsNil)

	upgradeInfo, err := s.st.UpgradeInfo(context.Background(), uuid)
	c.Assert(err, jc.ErrorIsNil)
	c.Check(upgradeInfo, gc.DeepEquals, upgrade.Info{
		UUID:            uuid.String(),
		PreviousVersion: "3.0.0",
		TargetVersion:   "3.0.1",
		State:           upgrade.Created,
	})

	nodeInfos := s.getUpgrade(c, s.st, uuid)
	c.Check(nodeInfos, gc.HasLen, 0)
}

func (s *stateSuite) TestCreateUpgradeAlreadyExists(c *gc.C) {
	_, err := s.st.CreateUpgrade(context.Background(), version.MustParse("3.0.0"), version.MustParse("3.0.1"))
	c.Assert(err, jc.ErrorIsNil)

	_, err = s.st.CreateUpgrade(context.Background(), version.MustParse("4.0.0"), version.MustParse("4.0.1"))
	c.Assert(database.IsErrConstraintUnique(err), jc.IsTrue)
}

func (s *stateSuite) TestSetControllerReady(c *gc.C) {
	uuid, err := s.st.CreateUpgrade(context.Background(), version.MustParse("3.0.0"), version.MustParse("3.0.1"))
	c.Assert(err, jc.ErrorIsNil)
	err = s.st.SetControllerReady(context.Background(), uuid, "0")
	c.Assert(err, jc.ErrorIsNil)

	nodeInfos := s.getUpgrade(c, s.st, uuid)
	c.Check(nodeInfos, gc.HasLen, 1)
	c.Check(nodeInfos[0], gc.Equals, infoControllerNode{
		ControllerNodeID: "0",
	})
}

func (s *stateSuite) TestSetControllerReadyWithoutUpgrade(c *gc.C) {
	uuid := utils.MustNewUUID().String()

	err := s.st.SetControllerReady(context.Background(), domainupgrade.UUID(uuid), "0")
	c.Assert(err, gc.NotNil)
	c.Check(database.IsErrConstraintForeignKey(err), jc.IsTrue)
}

func (s *stateSuite) TestSetControllerReadyMultipleTimes(c *gc.C) {
	err := s.st.SetControllerReady(context.Background(), s.upgradeUUID, "0")
	c.Assert(err, gc.NotNil)
	c.Check(database.IsErrConstraintUnique(err), jc.IsTrue)
}

func (s *stateSuite) TestAllProvisionedControllersReadyTrue(c *gc.C) {
	db := s.DB()

	_, err := db.Exec("INSERT INTO controller_node (controller_id, dqlite_node_id) VALUES ('1', 1)")
	c.Assert(err, jc.ErrorIsNil)

	uuid, err := s.st.CreateUpgrade(context.Background(), version.MustParse("3.0.0"), version.MustParse("3.0.1"))
	c.Assert(err, jc.ErrorIsNil)
	err = s.st.SetControllerReady(context.Background(), uuid, "0")
	c.Assert(err, jc.ErrorIsNil)
	err = s.st.SetControllerReady(context.Background(), uuid, "1")
	c.Assert(err, jc.ErrorIsNil)

	allProvisioned, err := s.st.AllProvisionedControllersReady(context.Background(), uuid)
	c.Assert(err, jc.ErrorIsNil)
	c.Check(allProvisioned, jc.IsTrue)
}

func (s *stateSuite) TestAllProvisionedControllersReadyFalse(c *gc.C) {
	db := s.DB()

	_, err := db.Exec("INSERT INTO controller_node (controller_id, dqlite_node_id) VALUES ('1', 1), ('2', 2)")
	c.Assert(err, jc.ErrorIsNil)

	uuid, err := s.st.CreateUpgrade(context.Background(), version.MustParse("3.0.0"), version.MustParse("3.0.1"))
	c.Assert(err, jc.ErrorIsNil)
	err = s.st.SetControllerReady(context.Background(), uuid, "0")
	c.Assert(err, jc.ErrorIsNil)
	err = s.st.SetControllerReady(context.Background(), uuid, "1")
	c.Assert(err, jc.ErrorIsNil)

	allProvisioned, err := s.st.AllProvisionedControllersReady(context.Background(), uuid)
	c.Assert(err, jc.ErrorIsNil)
	c.Check(allProvisioned, jc.IsFalse)
}

func (s *stateSuite) TestAllProvisionedControllersReadyMultipleControllers(c *gc.C) {
	db := s.DB()

	_, err := db.Exec("INSERT INTO controller_node (controller_id, dqlite_node_id) VALUES ('1', 1)")
	c.Assert(err, jc.ErrorIsNil)
	_, err = db.Exec("INSERT INTO controller_node (controller_id, dqlite_node_id) VALUES ('2', 2)")
	c.Assert(err, jc.ErrorIsNil)

	uuid, err := s.st.CreateUpgrade(context.Background(), version.MustParse("3.0.0"), version.MustParse("3.0.1"))
	c.Assert(err, jc.ErrorIsNil)
	err = s.st.SetControllerReady(context.Background(), uuid, "0")
	c.Assert(err, jc.ErrorIsNil)
	err = s.st.SetControllerReady(context.Background(), uuid, "1")
	c.Assert(err, jc.ErrorIsNil)
	err = s.st.SetControllerReady(context.Background(), uuid, "2")
	c.Assert(err, jc.ErrorIsNil)

	allProvisioned, err := s.st.AllProvisionedControllersReady(context.Background(), uuid)
	c.Assert(err, jc.ErrorIsNil)
	c.Check(allProvisioned, jc.IsTrue)
}

func (s *stateSuite) TestAllProvisionedControllersReadyMultipleControllersWithoutAllBeingReady(c *gc.C) {
	db := s.DB()

	_, err := db.Exec("INSERT INTO controller_node (controller_id, dqlite_node_id) VALUES ('1', 1)")
	c.Assert(err, jc.ErrorIsNil)
	_, err = db.Exec("INSERT INTO controller_node (controller_id, dqlite_node_id) VALUES ('2', 2)")
	c.Assert(err, jc.ErrorIsNil)

	uuid, err := s.st.CreateUpgrade(context.Background(), version.MustParse("3.0.0"), version.MustParse("3.0.1"))
	c.Assert(err, jc.ErrorIsNil)
	err = s.st.SetControllerReady(context.Background(), uuid, "0")
	c.Assert(err, jc.ErrorIsNil)
	err = s.st.SetControllerReady(context.Background(), uuid, "1")
	c.Assert(err, jc.ErrorIsNil)

	allProvisioned, err := s.st.AllProvisionedControllersReady(context.Background(), uuid)
	c.Assert(err, jc.ErrorIsNil)
	c.Check(allProvisioned, jc.IsFalse)
}

func (s *stateSuite) TestAllProvisionedControllersReadyUnprovisionedController(c *gc.C) {
	db := s.DB()

	_, err := db.Exec("INSERT INTO controller_node (controller_id, dqlite_node_id) VALUES ('1', 1)")
	c.Assert(err, jc.ErrorIsNil)
	_, err = db.Exec("INSERT INTO controller_node (controller_id) VALUES ('2')")
	c.Assert(err, jc.ErrorIsNil)

	uuid, err := s.st.CreateUpgrade(context.Background(), version.MustParse("3.0.0"), version.MustParse("3.0.1"))
	c.Assert(err, jc.ErrorIsNil)
	err = s.st.SetControllerReady(context.Background(), uuid, "0")
	c.Assert(err, jc.ErrorIsNil)
	err = s.st.SetControllerReady(context.Background(), uuid, "1")
	c.Assert(err, jc.ErrorIsNil)

	allProvisioned, err := s.st.AllProvisionedControllersReady(context.Background(), uuid)
	c.Assert(err, jc.ErrorIsNil)
	c.Check(allProvisioned, jc.IsTrue)
}

func (s *stateSuite) TestStartUpgrade(c *gc.C) {
	uuid, err := s.st.CreateUpgrade(context.Background(), version.MustParse("3.0.0"), version.MustParse("3.0.1"))
	c.Assert(err, jc.ErrorIsNil)
	err = s.st.SetControllerReady(context.Background(), uuid, "0")
	c.Assert(err, jc.ErrorIsNil)

	err = s.st.StartUpgrade(context.Background(), uuid)
	c.Assert(err, jc.ErrorIsNil)

	s.ensureUpgradeInfoState(c, uuid, upgrade.Started)
}

func (s *stateSuite) TestStartUpgradeCalledMultipleTimes(c *gc.C) {
	uuid, err := s.st.CreateUpgrade(context.Background(), version.MustParse("3.0.0"), version.MustParse("3.0.1"))
	c.Assert(err, jc.ErrorIsNil)
	err = s.st.SetControllerReady(context.Background(), uuid, "0")
	c.Assert(err, jc.ErrorIsNil)

	err = s.st.StartUpgrade(context.Background(), uuid)
	c.Assert(err, jc.ErrorIsNil)

	s.ensureUpgradeInfoState(c, uuid, upgrade.Started)

	err = s.st.StartUpgrade(context.Background(), uuid)
	c.Assert(err, jc.ErrorIs, upgradeerrors.ErrUpgradeAlreadyStarted)

	s.ensureUpgradeInfoState(c, uuid, upgrade.Started)
}

func (s *stateSuite) TestStartUpgradeBeforeCreated(c *gc.C) {
	uuid := utils.MustNewUUID().String()
	err := s.st.StartUpgrade(context.Background(), domainupgrade.UUID(uuid))
	c.Assert(err, jc.ErrorIs, sql.ErrNoRows)
}

func (s *stateSuite) TestSetControllerDone(c *gc.C) {
	db := s.DB()

	_, err := db.Exec("INSERT INTO controller_node (controller_id, dqlite_node_id) VALUES ('1', 1), ('2', 2)")
	c.Assert(err, jc.ErrorIsNil)
	uuid, err := s.st.CreateUpgrade(context.Background(), version.MustParse("3.0.0"), version.MustParse("3.0.1"))
	c.Assert(err, jc.ErrorIsNil)
	err = s.st.SetControllerReady(context.Background(), uuid, "0")
	c.Assert(err, jc.ErrorIsNil)
	err = s.st.SetControllerReady(context.Background(), uuid, "1")
	c.Assert(err, jc.ErrorIsNil)

	err = s.st.SetControllerDone(context.Background(), uuid, "1")
	c.Assert(err, jc.ErrorIsNil)

	nodeInfos := s.getUpgrade(c, s.st, uuid)
	c.Assert(err, jc.ErrorIsNil)
	c.Check(nodeInfos, gc.HasLen, 2)
	c.Check(nodeInfos[0], gc.Equals, infoControllerNode{ControllerNodeID: "0"})
	c.Check(nodeInfos[1], gc.Equals, infoControllerNode{ControllerNodeID: "1", NodeUpgradeCompletedAt: nodeInfos[1].NodeUpgradeCompletedAt})
	c.Check(nodeInfos[1].NodeUpgradeCompletedAt.Valid, jc.IsTrue)
}

func (s *stateSuite) TestSetControllerDoneNotExists(c *gc.C) {
	uuid, err := s.st.CreateUpgrade(context.Background(), version.MustParse("3.0.0"), version.MustParse("3.0.1"))
	c.Assert(err, jc.ErrorIsNil)

	err = s.st.SetControllerDone(context.Background(), uuid, "0")
	c.Assert(err, gc.ErrorMatches, `controller node "0" not ready`)
}

func (s *stateSuite) TestSetControllerDoneCompleteUpgrade(c *gc.C) {
	db := s.DB()

	_, err := db.Exec("INSERT INTO controller_node (controller_id, dqlite_node_id) VALUES ('1', 1)")
	c.Assert(err, jc.ErrorIsNil)

	// Start the upgrade
	uuid := s.startUpgrade(c)

	// Set the nodes to ready.
	err = s.st.SetControllerReady(context.Background(), uuid, "0")
	c.Assert(err, jc.ErrorIsNil)
	err = s.st.SetControllerReady(context.Background(), uuid, "1")
	c.Assert(err, jc.ErrorIsNil)

	// Set the db upgrade completed.
	err = s.st.SetDBUpgradeCompleted(context.Background(), uuid)
	c.Assert(err, jc.ErrorIsNil)

	// Ensure that all the steps have been completed.
	err = s.st.SetControllerDone(context.Background(), uuid, "0")
	c.Assert(err, jc.ErrorIsNil)

	// Check that the upgrade hasn't been completed for just one node.
	activeUpgrade, err := s.st.ActiveUpgrade(context.Background())
	c.Assert(err, jc.ErrorIsNil)
	c.Check(activeUpgrade, gc.Equals, uuid)

	// Set the last node.
	err = s.st.SetControllerDone(context.Background(), uuid, "1")
	c.Assert(err, jc.ErrorIsNil)

	// The active upgrade should be done.
	_, err = s.st.ActiveUpgrade(context.Background())
	c.Assert(err, jc.ErrorIs, sql.ErrNoRows)
}

func (s *stateSuite) TestSetControllerDoneCompleteUpgradeEmptyCompletedAt(c *gc.C) {
	db := s.DB()

	_, err := db.Exec("INSERT INTO controller_node (controller_id, dqlite_node_id) VALUES ('1', 1)")
	c.Assert(err, jc.ErrorIsNil)

	uuid := s.startUpgrade(c)

	err = s.st.SetControllerReady(context.Background(), uuid, "0")
	c.Assert(err, jc.ErrorIsNil)
	err = s.st.SetControllerReady(context.Background(), uuid, "1")
	c.Assert(err, jc.ErrorIsNil)

	_, err = db.Exec(`
UPDATE upgrade_info_controller_node 
SET    node_upgrade_completed_at = ''
WHERE  upgrade_info_uuid = ?
       AND controller_node_id = 0`, uuid)
	c.Assert(err, jc.ErrorIsNil)

	activeUpgrade, err := s.st.ActiveUpgrade(context.Background())
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(activeUpgrade, gc.Equals, uuid)

	// Set the db upgrade completed.
	err = s.st.SetDBUpgradeCompleted(context.Background(), uuid)
	c.Assert(err, jc.ErrorIsNil)

	err = s.st.SetControllerDone(context.Background(), uuid, "1")
	c.Assert(err, jc.ErrorIsNil)

	_, err = s.st.ActiveUpgrade(context.Background())
	c.Assert(err, jc.ErrorIs, sql.ErrNoRows)
}

func (s *stateSuite) TestActiveUpgradesNoUpgrades(c *gc.C) {
	_, err := s.st.ActiveUpgrade(context.Background())
	c.Assert(err, jc.ErrorIs, sql.ErrNoRows)
}

func (s *stateSuite) TestActiveUpgradesSingular(c *gc.C) {
	uuid, err := s.st.CreateUpgrade(context.Background(), version.MustParse("3.0.0"), version.MustParse("3.0.1"))
	c.Assert(err, jc.ErrorIsNil)

	activeUpgrade, err := s.st.ActiveUpgrade(context.Background())
	c.Assert(err, jc.ErrorIsNil)
	c.Check(activeUpgrade, gc.Equals, uuid)
}

func (s *stateSuite) TestSetDBUpgradeCompleted(c *gc.C) {
	uuid := s.startUpgrade(c)

	err := s.st.SetDBUpgradeCompleted(context.Background(), uuid)
	c.Assert(err, jc.ErrorIsNil)
	err = s.st.SetDBUpgradeCompleted(context.Background(), uuid)
	c.Assert(err, gc.ErrorMatches, `expected to set db upgrade completed.*`)

	s.ensureUpgradeInfoState(c, uuid, upgrade.DBCompleted)
}

func (s *stateSuite) TestUpgradeInfo(c *gc.C) {
	uuid, err := s.st.CreateUpgrade(context.Background(), version.MustParse("3.0.0"), version.MustParse("3.0.1"))
	c.Assert(err, jc.ErrorIsNil)

	upgradeInfo, err := s.st.UpgradeInfo(context.Background(), uuid)
	c.Assert(err, jc.ErrorIsNil)
	c.Check(upgradeInfo, gc.Equals, upgrade.Info{
		UUID:            uuid.String(),
		PreviousVersion: "3.0.0",
		TargetVersion:   "3.0.1",
		State:           upgrade.Created,
	})
}

func (s *stateSuite) ensureUpgradeInfoState(c *gc.C, uuid domainupgrade.UUID, state upgrade.State) {
	upgradeInfo, err := s.st.UpgradeInfo(context.Background(), uuid)
	c.Assert(err, jc.ErrorIsNil)
	c.Check(upgradeInfo.State, gc.Equals, state)
}

func (s *stateSuite) startUpgrade(c *gc.C) domainupgrade.UUID {
	_, err := s.st.ActiveUpgrade(context.Background())
	c.Assert(err, jc.ErrorIs, sql.ErrNoRows)

	uuid, err := s.st.CreateUpgrade(context.Background(), version.MustParse("3.0.0"), version.MustParse("3.0.1"))
	c.Assert(err, jc.ErrorIsNil)

	err = s.st.StartUpgrade(context.Background(), uuid)
	c.Assert(err, jc.ErrorIsNil)

	return uuid
}

func (s *stateSuite) getUpgrade(c *gc.C, st *State, upgradeUUID domainupgrade.UUID) []infoControllerNode {
	db, err := s.st.DB()
	c.Assert(err, jc.ErrorIsNil)

	nodeInfosQ := `
SELECT * AS &infoControllerNode.* FROM upgrade_info_controller_node
WHERE upgrade_info_uuid = $M.info_uuid`
	nodeInfosS, err := sqlair.Prepare(nodeInfosQ, infoControllerNode{}, sqlair.M{})
	c.Assert(err, jc.ErrorIsNil)

	var (
		nodeInfos []infoControllerNode
	)
	err = db.Txn(context.Background(), func(ctx context.Context, tx *sqlair.TX) error {
		err = tx.Query(ctx, nodeInfosS, sqlair.M{"info_uuid": upgradeUUID}).GetAll(&nodeInfos)
		c.Assert(err, jc.ErrorIsNil)
		return nil
	})
	c.Assert(err, jc.ErrorIsNil)
	return nodeInfos
}
