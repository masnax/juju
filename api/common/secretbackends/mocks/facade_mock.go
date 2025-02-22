// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/api/base (interfaces: FacadeCaller)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	base "github.com/juju/juju/api/base"
	gomock "go.uber.org/mock/gomock"
)

// MockFacadeCaller is a mock of FacadeCaller interface.
type MockFacadeCaller struct {
	ctrl     *gomock.Controller
	recorder *MockFacadeCallerMockRecorder
}

// MockFacadeCallerMockRecorder is the mock recorder for MockFacadeCaller.
type MockFacadeCallerMockRecorder struct {
	mock *MockFacadeCaller
}

// NewMockFacadeCaller creates a new mock instance.
func NewMockFacadeCaller(ctrl *gomock.Controller) *MockFacadeCaller {
	mock := &MockFacadeCaller{ctrl: ctrl}
	mock.recorder = &MockFacadeCallerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFacadeCaller) EXPECT() *MockFacadeCallerMockRecorder {
	return m.recorder
}

// BestAPIVersion mocks base method.
func (m *MockFacadeCaller) BestAPIVersion() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BestAPIVersion")
	ret0, _ := ret[0].(int)
	return ret0
}

// BestAPIVersion indicates an expected call of BestAPIVersion.
func (mr *MockFacadeCallerMockRecorder) BestAPIVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BestAPIVersion", reflect.TypeOf((*MockFacadeCaller)(nil).BestAPIVersion))
}

// FacadeCall mocks base method.
func (m *MockFacadeCaller) FacadeCall(arg0 context.Context, arg1 string, arg2, arg3 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FacadeCall", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// FacadeCall indicates an expected call of FacadeCall.
func (mr *MockFacadeCallerMockRecorder) FacadeCall(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FacadeCall", reflect.TypeOf((*MockFacadeCaller)(nil).FacadeCall), arg0, arg1, arg2, arg3)
}

// Name mocks base method.
func (m *MockFacadeCaller) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockFacadeCallerMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockFacadeCaller)(nil).Name))
}

// RawAPICaller mocks base method.
func (m *MockFacadeCaller) RawAPICaller() base.APICaller {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RawAPICaller")
	ret0, _ := ret[0].(base.APICaller)
	return ret0
}

// RawAPICaller indicates an expected call of RawAPICaller.
func (mr *MockFacadeCallerMockRecorder) RawAPICaller() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RawAPICaller", reflect.TypeOf((*MockFacadeCaller)(nil).RawAPICaller))
}
