// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/common (interfaces: CredentialService,CloudService)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	cloud "github.com/juju/juju/cloud"
	watcher "github.com/juju/juju/core/watcher"
	credential "github.com/juju/juju/domain/credential"
	gomock "go.uber.org/mock/gomock"
)

// MockCredentialService is a mock of CredentialService interface.
type MockCredentialService struct {
	ctrl     *gomock.Controller
	recorder *MockCredentialServiceMockRecorder
}

// MockCredentialServiceMockRecorder is the mock recorder for MockCredentialService.
type MockCredentialServiceMockRecorder struct {
	mock *MockCredentialService
}

// NewMockCredentialService creates a new mock instance.
func NewMockCredentialService(ctrl *gomock.Controller) *MockCredentialService {
	mock := &MockCredentialService{ctrl: ctrl}
	mock.recorder = &MockCredentialServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCredentialService) EXPECT() *MockCredentialServiceMockRecorder {
	return m.recorder
}

// CloudCredential mocks base method.
func (m *MockCredentialService) CloudCredential(arg0 context.Context, arg1 credential.ID) (cloud.Credential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudCredential", arg0, arg1)
	ret0, _ := ret[0].(cloud.Credential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloudCredential indicates an expected call of CloudCredential.
func (mr *MockCredentialServiceMockRecorder) CloudCredential(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudCredential", reflect.TypeOf((*MockCredentialService)(nil).CloudCredential), arg0, arg1)
}

// WatchCredential mocks base method.
func (m *MockCredentialService) WatchCredential(arg0 context.Context, arg1 credential.ID) (watcher.Watcher[struct{}], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchCredential", arg0, arg1)
	ret0, _ := ret[0].(watcher.Watcher[struct{}])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchCredential indicates an expected call of WatchCredential.
func (mr *MockCredentialServiceMockRecorder) WatchCredential(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchCredential", reflect.TypeOf((*MockCredentialService)(nil).WatchCredential), arg0, arg1)
}

// MockCloudService is a mock of CloudService interface.
type MockCloudService struct {
	ctrl     *gomock.Controller
	recorder *MockCloudServiceMockRecorder
}

// MockCloudServiceMockRecorder is the mock recorder for MockCloudService.
type MockCloudServiceMockRecorder struct {
	mock *MockCloudService
}

// NewMockCloudService creates a new mock instance.
func NewMockCloudService(ctrl *gomock.Controller) *MockCloudService {
	mock := &MockCloudService{ctrl: ctrl}
	mock.recorder = &MockCloudServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCloudService) EXPECT() *MockCloudServiceMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockCloudService) Get(arg0 context.Context, arg1 string) (*cloud.Cloud, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*cloud.Cloud)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCloudServiceMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCloudService)(nil).Get), arg0, arg1)
}

// WatchCloud mocks base method.
func (m *MockCloudService) WatchCloud(arg0 context.Context, arg1 string) (watcher.Watcher[struct{}], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchCloud", arg0, arg1)
	ret0, _ := ret[0].(watcher.Watcher[struct{}])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchCloud indicates an expected call of WatchCloud.
func (mr *MockCloudServiceMockRecorder) WatchCloud(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchCloud", reflect.TypeOf((*MockCloudService)(nil).WatchCloud), arg0, arg1)
}
