// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/controller/externalcontrollerupdater (interfaces: ECService)

// Package externalcontrollerupdater_test is a generated GoMock package.
package externalcontrollerupdater_test

import (
	context "context"
	reflect "reflect"

	crossmodel "github.com/juju/juju/core/crossmodel"
	watcher "github.com/juju/juju/core/watcher"
	gomock "go.uber.org/mock/gomock"
)

// MockECService is a mock of ECService interface.
type MockECService struct {
	ctrl     *gomock.Controller
	recorder *MockECServiceMockRecorder
}

// MockECServiceMockRecorder is the mock recorder for MockECService.
type MockECServiceMockRecorder struct {
	mock *MockECService
}

// NewMockECService creates a new mock instance.
func NewMockECService(ctrl *gomock.Controller) *MockECService {
	mock := &MockECService{ctrl: ctrl}
	mock.recorder = &MockECServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockECService) EXPECT() *MockECServiceMockRecorder {
	return m.recorder
}

// Controller mocks base method.
func (m *MockECService) Controller(arg0 context.Context, arg1 string) (*crossmodel.ControllerInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Controller", arg0, arg1)
	ret0, _ := ret[0].(*crossmodel.ControllerInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Controller indicates an expected call of Controller.
func (mr *MockECServiceMockRecorder) Controller(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Controller", reflect.TypeOf((*MockECService)(nil).Controller), arg0, arg1)
}

// UpdateExternalController mocks base method.
func (m *MockECService) UpdateExternalController(arg0 context.Context, arg1 crossmodel.ControllerInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateExternalController", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateExternalController indicates an expected call of UpdateExternalController.
func (mr *MockECServiceMockRecorder) UpdateExternalController(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateExternalController", reflect.TypeOf((*MockECService)(nil).UpdateExternalController), arg0, arg1)
}

// Watch mocks base method.
func (m *MockECService) Watch() (watcher.Watcher[[]string], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Watch")
	ret0, _ := ret[0].(watcher.Watcher[[]string])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Watch indicates an expected call of Watch.
func (mr *MockECServiceMockRecorder) Watch() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockECService)(nil).Watch))
}
