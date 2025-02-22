// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/internal/secrets/provider (interfaces: SecretBackendProvider)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	provider "github.com/juju/juju/internal/secrets/provider"
	names "github.com/juju/names/v4"
	gomock "go.uber.org/mock/gomock"
)

// MockSecretBackendProvider is a mock of SecretBackendProvider interface.
type MockSecretBackendProvider struct {
	ctrl     *gomock.Controller
	recorder *MockSecretBackendProviderMockRecorder
}

// MockSecretBackendProviderMockRecorder is the mock recorder for MockSecretBackendProvider.
type MockSecretBackendProviderMockRecorder struct {
	mock *MockSecretBackendProvider
}

// NewMockSecretBackendProvider creates a new mock instance.
func NewMockSecretBackendProvider(ctrl *gomock.Controller) *MockSecretBackendProvider {
	mock := &MockSecretBackendProvider{ctrl: ctrl}
	mock.recorder = &MockSecretBackendProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretBackendProvider) EXPECT() *MockSecretBackendProviderMockRecorder {
	return m.recorder
}

// CleanupModel mocks base method.
func (m *MockSecretBackendProvider) CleanupModel(arg0 *provider.ModelBackendConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CleanupModel", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CleanupModel indicates an expected call of CleanupModel.
func (mr *MockSecretBackendProviderMockRecorder) CleanupModel(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CleanupModel", reflect.TypeOf((*MockSecretBackendProvider)(nil).CleanupModel), arg0)
}

// CleanupSecrets mocks base method.
func (m *MockSecretBackendProvider) CleanupSecrets(arg0 *provider.ModelBackendConfig, arg1 names.Tag, arg2 provider.SecretRevisions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CleanupSecrets", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CleanupSecrets indicates an expected call of CleanupSecrets.
func (mr *MockSecretBackendProviderMockRecorder) CleanupSecrets(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CleanupSecrets", reflect.TypeOf((*MockSecretBackendProvider)(nil).CleanupSecrets), arg0, arg1, arg2)
}

// Initialise mocks base method.
func (m *MockSecretBackendProvider) Initialise(arg0 *provider.ModelBackendConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Initialise", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Initialise indicates an expected call of Initialise.
func (mr *MockSecretBackendProviderMockRecorder) Initialise(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Initialise", reflect.TypeOf((*MockSecretBackendProvider)(nil).Initialise), arg0)
}

// NewBackend mocks base method.
func (m *MockSecretBackendProvider) NewBackend(arg0 *provider.ModelBackendConfig) (provider.SecretsBackend, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewBackend", arg0)
	ret0, _ := ret[0].(provider.SecretsBackend)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewBackend indicates an expected call of NewBackend.
func (mr *MockSecretBackendProviderMockRecorder) NewBackend(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewBackend", reflect.TypeOf((*MockSecretBackendProvider)(nil).NewBackend), arg0)
}

// RestrictedConfig mocks base method.
func (m *MockSecretBackendProvider) RestrictedConfig(arg0 *provider.ModelBackendConfig, arg1 bool, arg2 names.Tag, arg3, arg4 provider.SecretRevisions) (*provider.BackendConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RestrictedConfig", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(*provider.BackendConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RestrictedConfig indicates an expected call of RestrictedConfig.
func (mr *MockSecretBackendProviderMockRecorder) RestrictedConfig(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RestrictedConfig", reflect.TypeOf((*MockSecretBackendProvider)(nil).RestrictedConfig), arg0, arg1, arg2, arg3, arg4)
}

// Type mocks base method.
func (m *MockSecretBackendProvider) Type() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Type")
	ret0, _ := ret[0].(string)
	return ret0
}

// Type indicates an expected call of Type.
func (mr *MockSecretBackendProviderMockRecorder) Type() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Type", reflect.TypeOf((*MockSecretBackendProvider)(nil).Type))
}
