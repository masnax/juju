// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/core/lease (interfaces: ExpiryStore)

// Package leaseexpiry_test is a generated GoMock package.
package leaseexpiry_test

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockExpiryStore is a mock of ExpiryStore interface.
type MockExpiryStore struct {
	ctrl     *gomock.Controller
	recorder *MockExpiryStoreMockRecorder
}

// MockExpiryStoreMockRecorder is the mock recorder for MockExpiryStore.
type MockExpiryStoreMockRecorder struct {
	mock *MockExpiryStore
}

// NewMockExpiryStore creates a new mock instance.
func NewMockExpiryStore(ctrl *gomock.Controller) *MockExpiryStore {
	mock := &MockExpiryStore{ctrl: ctrl}
	mock.recorder = &MockExpiryStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExpiryStore) EXPECT() *MockExpiryStoreMockRecorder {
	return m.recorder
}

// ExpireLeases mocks base method.
func (m *MockExpiryStore) ExpireLeases(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExpireLeases", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExpireLeases indicates an expected call of ExpireLeases.
func (mr *MockExpiryStoreMockRecorder) ExpireLeases(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExpireLeases", reflect.TypeOf((*MockExpiryStore)(nil).ExpireLeases), arg0)
}
