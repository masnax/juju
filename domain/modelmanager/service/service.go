// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"context"

	"github.com/juju/errors"

	"github.com/juju/juju/core/database"
	"github.com/juju/juju/domain"
	"github.com/juju/juju/domain/model"
)

// State defines a interface for interacting with the underlying state.
type State interface {
	Create(context.Context, model.UUID) error
	List(context.Context) ([]model.UUID, error)
	Delete(context.Context, model.UUID) error
}

// Service defines a service for interacting with the underlying state.
type Service struct {
	st        State
	dbDeleter database.DBDeleter
}

// NewService returns a new Service for interacting with the underlying state.
func NewService(st State, dbDeleter database.DBDeleter) *Service {
	return &Service{
		st:        st,
		dbDeleter: dbDeleter,
	}
}

// Create takes a model UUID and creates a new model.
func (s *Service) Create(ctx context.Context, uuid model.UUID) error {
	if err := uuid.Validate(); err != nil {
		return errors.Annotatef(err, "validating model uuid %q", uuid)
	}

	err := s.st.Create(ctx, uuid)
	return errors.Annotatef(domain.CoerceError(err), "creating model %q", uuid)
}

// ModelList returns a list of all model UUIDs.
// This only includes active models. Either a model is within the model manager
// list or it's not.
// Note: This shouldn't be used as a proxy for alive models. This hasn't got
// the same guarantees. Instead this should only be used for managing models
// from a dqlite perspective.
func (s *Service) ModelList(ctx context.Context) ([]model.UUID, error) {
	uuids, err := s.st.List(ctx)
	if err != nil {
		return nil, errors.Annotatef(err, "retrieving model list")
	}
	return uuids, nil
}

// Delete takes a model UUID and deletes the model if it exists.
func (s *Service) Delete(ctx context.Context, uuid model.UUID) error {
	if err := uuid.Validate(); err != nil {
		return errors.Annotatef(err, "validating model uuid %q", uuid)
	}

	// Deletion of the model in state should prevent any future requests to
	// acquire the tracked db from the db manager.
	if err := s.st.Delete(ctx, uuid); err != nil {
		if errors.Is(err, domain.ErrNoRecord) {
			return nil
		}
		return errors.Annotatef(domain.CoerceError(err), "deleting model %q", uuid)
	}

	err := s.dbDeleter.DeleteDB(uuid.String())
	return errors.Annotatef(err, "stopping model %q", uuid)
}
