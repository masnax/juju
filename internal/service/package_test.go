// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service_test

import (
	"testing"

	gc "gopkg.in/check.v1"
)

//go:generate go run go.uber.org/mock/mockgen -package mocks -destination mocks/service.go github.com/juju/juju/internal/service Service

func TestPackage(t *testing.T) {
	gc.TestingT(t)
}
