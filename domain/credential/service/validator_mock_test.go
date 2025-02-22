// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/domain/credential/service (interfaces: CredentialValidator)

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	cloud "github.com/juju/juju/cloud"
	credential "github.com/juju/juju/domain/credential"
	gomock "go.uber.org/mock/gomock"
)

// MockCredentialValidator is a mock of CredentialValidator interface.
type MockCredentialValidator struct {
	ctrl     *gomock.Controller
	recorder *MockCredentialValidatorMockRecorder
}

// MockCredentialValidatorMockRecorder is the mock recorder for MockCredentialValidator.
type MockCredentialValidatorMockRecorder struct {
	mock *MockCredentialValidator
}

// NewMockCredentialValidator creates a new mock instance.
func NewMockCredentialValidator(ctrl *gomock.Controller) *MockCredentialValidator {
	mock := &MockCredentialValidator{ctrl: ctrl}
	mock.recorder = &MockCredentialValidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCredentialValidator) EXPECT() *MockCredentialValidatorMockRecorder {
	return m.recorder
}

// Validate mocks base method.
func (m *MockCredentialValidator) Validate(arg0 context.Context, arg1 CredentialValidationContext, arg2 credential.ID, arg3 *cloud.Credential, arg4 bool) ([]error, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]error)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Validate indicates an expected call of Validate.
func (mr *MockCredentialValidatorMockRecorder) Validate(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockCredentialValidator)(nil).Validate), arg0, arg1, arg2, arg3, arg4)
}
