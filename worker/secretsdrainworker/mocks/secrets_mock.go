// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/internal/secrets (interfaces: BackendsClient)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	secrets "github.com/juju/juju/core/secrets"
	provider "github.com/juju/juju/internal/secrets/provider"
	gomock "go.uber.org/mock/gomock"
)

// MockBackendsClient is a mock of BackendsClient interface.
type MockBackendsClient struct {
	ctrl     *gomock.Controller
	recorder *MockBackendsClientMockRecorder
}

// MockBackendsClientMockRecorder is the mock recorder for MockBackendsClient.
type MockBackendsClientMockRecorder struct {
	mock *MockBackendsClient
}

// NewMockBackendsClient creates a new mock instance.
func NewMockBackendsClient(ctrl *gomock.Controller) *MockBackendsClient {
	mock := &MockBackendsClient{ctrl: ctrl}
	mock.recorder = &MockBackendsClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBackendsClient) EXPECT() *MockBackendsClientMockRecorder {
	return m.recorder
}

// DeleteContent mocks base method.
func (m *MockBackendsClient) DeleteContent(arg0 *secrets.URI, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteContent", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteContent indicates an expected call of DeleteContent.
func (mr *MockBackendsClientMockRecorder) DeleteContent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteContent", reflect.TypeOf((*MockBackendsClient)(nil).DeleteContent), arg0, arg1)
}

// DeleteExternalContent mocks base method.
func (m *MockBackendsClient) DeleteExternalContent(arg0 secrets.ValueRef) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteExternalContent", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteExternalContent indicates an expected call of DeleteExternalContent.
func (mr *MockBackendsClientMockRecorder) DeleteExternalContent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteExternalContent", reflect.TypeOf((*MockBackendsClient)(nil).DeleteExternalContent), arg0)
}

// GetBackend mocks base method.
func (m *MockBackendsClient) GetBackend(arg0 *string, arg1 bool) (provider.SecretsBackend, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBackend", arg0, arg1)
	ret0, _ := ret[0].(provider.SecretsBackend)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetBackend indicates an expected call of GetBackend.
func (mr *MockBackendsClientMockRecorder) GetBackend(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBackend", reflect.TypeOf((*MockBackendsClient)(nil).GetBackend), arg0, arg1)
}

// GetContent mocks base method.
func (m *MockBackendsClient) GetContent(arg0 *secrets.URI, arg1 string, arg2, arg3 bool) (secrets.SecretValue, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContent", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(secrets.SecretValue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContent indicates an expected call of GetContent.
func (mr *MockBackendsClientMockRecorder) GetContent(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContent", reflect.TypeOf((*MockBackendsClient)(nil).GetContent), arg0, arg1, arg2, arg3)
}

// GetRevisionContent mocks base method.
func (m *MockBackendsClient) GetRevisionContent(arg0 *secrets.URI, arg1 int) (secrets.SecretValue, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRevisionContent", arg0, arg1)
	ret0, _ := ret[0].(secrets.SecretValue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRevisionContent indicates an expected call of GetRevisionContent.
func (mr *MockBackendsClientMockRecorder) GetRevisionContent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRevisionContent", reflect.TypeOf((*MockBackendsClient)(nil).GetRevisionContent), arg0, arg1)
}

// SaveContent mocks base method.
func (m *MockBackendsClient) SaveContent(arg0 *secrets.URI, arg1 int, arg2 secrets.SecretValue) (secrets.ValueRef, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveContent", arg0, arg1, arg2)
	ret0, _ := ret[0].(secrets.ValueRef)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveContent indicates an expected call of SaveContent.
func (mr *MockBackendsClientMockRecorder) SaveContent(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveContent", reflect.TypeOf((*MockBackendsClient)(nil).SaveContent), arg0, arg1, arg2)
}
