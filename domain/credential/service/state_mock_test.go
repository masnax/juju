// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/domain/credential/service (interfaces: State,WatcherFactory)

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	cloud "github.com/juju/juju/cloud"
	changestream "github.com/juju/juju/core/changestream"
	watcher "github.com/juju/juju/core/watcher"
	credential "github.com/juju/juju/domain/credential"
	model "github.com/juju/juju/domain/model"
	gomock "go.uber.org/mock/gomock"
)

// MockState is a mock of State interface.
type MockState struct {
	ctrl     *gomock.Controller
	recorder *MockStateMockRecorder
}

// MockStateMockRecorder is the mock recorder for MockState.
type MockStateMockRecorder struct {
	mock *MockState
}

// NewMockState creates a new mock instance.
func NewMockState(ctrl *gomock.Controller) *MockState {
	mock := &MockState{ctrl: ctrl}
	mock.recorder = &MockStateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockState) EXPECT() *MockStateMockRecorder {
	return m.recorder
}

// AllCloudCredentials mocks base method.
func (m *MockState) AllCloudCredentials(arg0 context.Context, arg1 string) ([]credential.CloudCredential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllCloudCredentials", arg0, arg1)
	ret0, _ := ret[0].([]credential.CloudCredential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllCloudCredentials indicates an expected call of AllCloudCredentials.
func (mr *MockStateMockRecorder) AllCloudCredentials(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllCloudCredentials", reflect.TypeOf((*MockState)(nil).AllCloudCredentials), arg0, arg1)
}

// CloudCredential mocks base method.
func (m *MockState) CloudCredential(arg0 context.Context, arg1, arg2, arg3 string) (cloud.Credential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudCredential", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(cloud.Credential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloudCredential indicates an expected call of CloudCredential.
func (mr *MockStateMockRecorder) CloudCredential(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudCredential", reflect.TypeOf((*MockState)(nil).CloudCredential), arg0, arg1, arg2, arg3)
}

// CloudCredentials mocks base method.
func (m *MockState) CloudCredentials(arg0 context.Context, arg1, arg2 string) (map[string]cloud.Credential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudCredentials", arg0, arg1, arg2)
	ret0, _ := ret[0].(map[string]cloud.Credential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloudCredentials indicates an expected call of CloudCredentials.
func (mr *MockStateMockRecorder) CloudCredentials(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudCredentials", reflect.TypeOf((*MockState)(nil).CloudCredentials), arg0, arg1, arg2)
}

// GetCloud mocks base method.
func (m *MockState) GetCloud(arg0 context.Context, arg1 string) (cloud.Cloud, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCloud", arg0, arg1)
	ret0, _ := ret[0].(cloud.Cloud)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCloud indicates an expected call of GetCloud.
func (mr *MockStateMockRecorder) GetCloud(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCloud", reflect.TypeOf((*MockState)(nil).GetCloud), arg0, arg1)
}

// InvalidateCloudCredential mocks base method.
func (m *MockState) InvalidateCloudCredential(arg0 context.Context, arg1, arg2, arg3, arg4 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InvalidateCloudCredential", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// InvalidateCloudCredential indicates an expected call of InvalidateCloudCredential.
func (mr *MockStateMockRecorder) InvalidateCloudCredential(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InvalidateCloudCredential", reflect.TypeOf((*MockState)(nil).InvalidateCloudCredential), arg0, arg1, arg2, arg3, arg4)
}

// ModelsUsingCloudCredential mocks base method.
func (m *MockState) ModelsUsingCloudCredential(arg0 context.Context, arg1 credential.ID) (map[model.UUID]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelsUsingCloudCredential", arg0, arg1)
	ret0, _ := ret[0].(map[model.UUID]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ModelsUsingCloudCredential indicates an expected call of ModelsUsingCloudCredential.
func (mr *MockStateMockRecorder) ModelsUsingCloudCredential(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelsUsingCloudCredential", reflect.TypeOf((*MockState)(nil).ModelsUsingCloudCredential), arg0, arg1)
}

// RemoveCloudCredential mocks base method.
func (m *MockState) RemoveCloudCredential(arg0 context.Context, arg1, arg2, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveCloudCredential", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveCloudCredential indicates an expected call of RemoveCloudCredential.
func (mr *MockStateMockRecorder) RemoveCloudCredential(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveCloudCredential", reflect.TypeOf((*MockState)(nil).RemoveCloudCredential), arg0, arg1, arg2, arg3)
}

// UpsertCloudCredential mocks base method.
func (m *MockState) UpsertCloudCredential(arg0 context.Context, arg1, arg2, arg3 string, arg4 cloud.Credential) (*bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertCloudCredential", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(*bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpsertCloudCredential indicates an expected call of UpsertCloudCredential.
func (mr *MockStateMockRecorder) UpsertCloudCredential(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertCloudCredential", reflect.TypeOf((*MockState)(nil).UpsertCloudCredential), arg0, arg1, arg2, arg3, arg4)
}

// WatchCredential mocks base method.
func (m *MockState) WatchCredential(arg0 context.Context, arg1 func(string, string, changestream.ChangeType) (watcher.Watcher[struct{}], error), arg2, arg3, arg4 string) (watcher.Watcher[struct{}], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchCredential", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(watcher.Watcher[struct{}])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchCredential indicates an expected call of WatchCredential.
func (mr *MockStateMockRecorder) WatchCredential(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchCredential", reflect.TypeOf((*MockState)(nil).WatchCredential), arg0, arg1, arg2, arg3, arg4)
}

// MockWatcherFactory is a mock of WatcherFactory interface.
type MockWatcherFactory struct {
	ctrl     *gomock.Controller
	recorder *MockWatcherFactoryMockRecorder
}

// MockWatcherFactoryMockRecorder is the mock recorder for MockWatcherFactory.
type MockWatcherFactoryMockRecorder struct {
	mock *MockWatcherFactory
}

// NewMockWatcherFactory creates a new mock instance.
func NewMockWatcherFactory(ctrl *gomock.Controller) *MockWatcherFactory {
	mock := &MockWatcherFactory{ctrl: ctrl}
	mock.recorder = &MockWatcherFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWatcherFactory) EXPECT() *MockWatcherFactoryMockRecorder {
	return m.recorder
}

// NewValueWatcher mocks base method.
func (m *MockWatcherFactory) NewValueWatcher(arg0, arg1 string, arg2 changestream.ChangeType) (watcher.Watcher[struct{}], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewValueWatcher", arg0, arg1, arg2)
	ret0, _ := ret[0].(watcher.Watcher[struct{}])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewValueWatcher indicates an expected call of NewValueWatcher.
func (mr *MockWatcherFactoryMockRecorder) NewValueWatcher(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewValueWatcher", reflect.TypeOf((*MockWatcherFactory)(nil).NewValueWatcher), arg0, arg1, arg2)
}
