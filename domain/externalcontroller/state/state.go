// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/canonical/sqlair"
	"github.com/juju/errors"
	"github.com/juju/names/v4"
	"github.com/juju/utils/v3"

	"github.com/juju/juju/core/crossmodel"
	"github.com/juju/juju/database"
	"github.com/juju/juju/domain"
)

type State struct {
	*domain.StateBase
}

func NewState(factory domain.TxnRunnerFactory) *State {
	return &State{
		StateBase: domain.NewStateBase(factory),
	}
}

func (st *State) Controller(
	ctx context.Context,
	controllerUUID string,
) (*crossmodel.ControllerInfo, error) {
	db, err := st.DB()
	if err != nil {
		return nil, errors.Trace(err)
	}

	q := `
SELECT (ctrl.uuid,
       alias,
       ca_cert,
       address) as &ExternalController.*,
       model.uuid as &ExternalController.model
FROM external_controller AS ctrl
LEFT JOIN external_model AS model
ON ctrl.uuid = model.controller_uuid
LEFT JOIN external_controller_address AS addrs
ON ctrl.uuid = addrs.controller_uuid
WHERE ctrl.uuid = $M.id`
	s, err := sqlair.Prepare(q, ExternalController{}, sqlair.M{})
	if err != nil {
		return nil, errors.Annotatef(err, "preparing %q", q)
	}

	var rows []ExternalController
	if err := db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		return errors.Trace(tx.Query(ctx, s, sqlair.M{"id": controllerUUID}).GetAll(&rows))
	}); err != nil {
		return nil, errors.Annotate(err, "querying external controller")
	}

	if len(rows) == 0 {
		return nil, errors.NotFoundf("external controller %q", controllerUUID)
	}

	return controllerInfoFromRows(rows), nil
}

func (st *State) ControllersForModels(ctx context.Context, modelUUIDs ...string) ([]crossmodel.ControllerInfo, error) {
	db, err := st.DB()
	if err != nil {
		return nil, errors.Trace(err)
	}

	// TODO(nvinuesa): We should use an `IN` statement query here, instead
	// of looping over the list of models and performing N queries, but
	// they are not yet supported on sqlair:
	// https://github.com/canonical/sqlair/pull/76
	q := `
SELECT (ctrl.uuid,  
       ctrl.alias,
       ctrl.ca_cert,
       addrs.address) as &ExternalController.*,
       model.uuid as &ExternalController.model
FROM external_controller AS ctrl	
JOIN external_model AS model
ON ctrl.uuid = model.controller_uuid
LEFT JOIN external_controller_address AS addrs
ON ctrl.uuid = addrs.controller_uuid
WHERE model.uuid = $M.model`

	var resultControllers []crossmodel.ControllerInfo

	uniqueControllers := make(map[string]crossmodel.ControllerInfo)
	s, err := sqlair.Prepare(q, ExternalController{}, sqlair.M{})
	if err != nil {
		return nil, errors.Annotatef(err, "preparing %q", q)
	}

	var rows []ExternalController
	if err := db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		for _, modelUUID := range modelUUIDs {
			err := tx.Query(ctx, s, sqlair.M{"model": modelUUID}).GetAll(&rows)
			if err != nil {
				return errors.Trace(err)
			}
			if len(rows) > 0 {
				uniqueControllers[rows[0].ID] = *controllerInfoFromRows(rows)
			}

			for _, controller := range uniqueControllers {
				resultControllers = append(resultControllers, controller)
			}
		}

		return nil
	}); err != nil {
		return nil, errors.Annotate(err, "querying external controller")
	}

	return resultControllers, nil
}

// controllerInfoFromRows is a convenient method for de-normalizing multiple
// rows containing the (necessary) columns from the three tables that
// represent the ControllerInfo struct (i.e. including addresses and models).
func controllerInfoFromRows(rows []ExternalController) *crossmodel.ControllerInfo {
	// We are sure that only one controller is represented by the list
	// of rows:
	controllerInfo := &crossmodel.ControllerInfo{
		ControllerTag: names.NewControllerTag(rows[0].ID),
		CACert:        rows[0].CACert,
		Alias:         rows[0].Alias.String,
	}

	// Prepare structs for unique models and addresses for each
	// controller.
	uniqueModelUUIDs := make(map[string]string)
	uniqueAddresses := make(map[string]string)
	for _, row := range rows {
		// Each row contains only one address, so it's safe
		// to access the only possible (nullable) value.
		if row.Addr.Valid {
			uniqueAddresses[row.Addr.String] = row.Addr.String
		}
		// Each row contains only one model, so it's safe
		// to access the only possible (nullable) value.
		if row.ModelUUID.Valid {
			uniqueModelUUIDs[row.ModelUUID.String] = row.ModelUUID.String
		}
	}

	// Flatten the models and addresses.
	var modelUUIDs []string
	for _, modelUUID := range uniqueModelUUIDs {
		modelUUIDs = append(modelUUIDs, modelUUID)
	}
	controllerInfo.ModelUUIDs = modelUUIDs
	var addresses []string
	for _, modelUUID := range uniqueAddresses {
		addresses = append(addresses, modelUUID)
	}
	controllerInfo.Addrs = addresses

	return controllerInfo
}

func (st *State) UpdateExternalController(
	ctx context.Context,
	ci crossmodel.ControllerInfo,
) error {
	db, err := st.DB()
	if err != nil {
		return errors.Trace(err)
	}

	err = db.StdTxn(ctx, func(ctx context.Context, tx *sql.Tx) error {
		return st.updateExternalControllerTx(ctx, tx, ci)
	})

	return errors.Trace(err)
}

// ImportExternalControllers imports the list of ControllerInfo
// external controllers on one single transaction.
func (st *State) ImportExternalControllers(ctx context.Context, infos []crossmodel.ControllerInfo) error {
	db, err := st.DB()
	if err != nil {
		return errors.Trace(err)
	}

	err = db.StdTxn(ctx, func(ctx context.Context, tx *sql.Tx) error {
		for _, ci := range infos {
			err := st.updateExternalControllerTx(
				ctx,
				tx,
				ci,
			)
			if err != nil {
				return errors.Trace(err)
			}
		}
		return nil
	})

	return errors.Trace(err)
}

func (st *State) updateExternalControllerTx(
	ctx context.Context,
	tx *sql.Tx,
	ci crossmodel.ControllerInfo,
) error {
	cID := ci.ControllerTag.Id()
	q := `
INSERT INTO external_controller (uuid, alias, ca_cert)
VALUES (?, ?, ?)
  ON CONFLICT(uuid) DO UPDATE SET alias=excluded.alias, ca_cert=excluded.ca_cert`[1:]

	if _, err := tx.ExecContext(ctx, q, cID, ci.Alias, ci.CACert); err != nil {
		return errors.Trace(err)
	}

	addrsBinds, addrsAnyVals := database.SliceToPlaceholder(ci.Addrs)
	q = fmt.Sprintf(`
DELETE FROM external_controller_address
WHERE  controller_uuid = ?
AND    address NOT IN (%s)`[1:], addrsBinds)

	args := append([]any{cID}, addrsAnyVals...)
	if _, err := tx.ExecContext(ctx, q, args...); err != nil {
		return errors.Trace(err)
	}

	for _, addr := range ci.Addrs {
		q := `
INSERT INTO external_controller_address (uuid, controller_uuid, address)
VALUES (?, ?, ?)
  ON CONFLICT(controller_uuid, address) DO NOTHING`[1:]

		if _, err := tx.ExecContext(ctx, q, utils.MustNewUUID().String(), cID, addr); err != nil {
			return errors.Trace(err)
		}
	}

	// TODO (manadart 2023-05-13): Check current implementation and see if
	// we need to delete models as we do for addresses, or whether this
	// (additive) approach is what we have now.
	for _, modelUUID := range ci.ModelUUIDs {
		q := `
INSERT INTO external_model (uuid, controller_uuid)
VALUES (?, ?)
  ON CONFLICT(uuid) DO UPDATE SET controller_uuid=excluded.controller_uuid`[1:]

		if _, err := tx.ExecContext(ctx, q, modelUUID, cID); err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}

func (st *State) ModelsForController(
	ctx context.Context,
	controllerUUID string,
) ([]string, error) {
	db, err := st.DB()
	if err != nil {
		return nil, errors.Trace(err)
	}

	q := `
SELECT uuid 
FROM external_model 
WHERE controller_uuid = ?`

	var modelUUIDs []string
	err = db.StdTxn(ctx, func(ctx context.Context, tx *sql.Tx) error {
		rows, err := tx.QueryContext(ctx, q, controllerUUID)
		if err != nil {
			return errors.Trace(err)
		}

		for rows.Next() {
			var modelUUID string
			if err := rows.Scan(&modelUUID); err != nil {
				_ = rows.Close()
				return errors.Trace(err)
			}
			modelUUIDs = append(modelUUIDs, modelUUID)
		}

		return nil
	})

	return modelUUIDs, err
}
