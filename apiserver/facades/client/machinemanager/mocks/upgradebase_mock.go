// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/client/machinemanager (interfaces: Authorizer,UpgradeSeries,UpgradeSeriesState,UpgradeBaseValidator)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	machinemanager "github.com/juju/juju/apiserver/facades/client/machinemanager"
	base "github.com/juju/juju/core/base"
	gomock "go.uber.org/mock/gomock"
)

// MockAuthorizer is a mock of Authorizer interface.
type MockAuthorizer struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizerMockRecorder
}

// MockAuthorizerMockRecorder is the mock recorder for MockAuthorizer.
type MockAuthorizerMockRecorder struct {
	mock *MockAuthorizer
}

// NewMockAuthorizer creates a new mock instance.
func NewMockAuthorizer(ctrl *gomock.Controller) *MockAuthorizer {
	mock := &MockAuthorizer{ctrl: ctrl}
	mock.recorder = &MockAuthorizerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthorizer) EXPECT() *MockAuthorizerMockRecorder {
	return m.recorder
}

// AuthClient mocks base method.
func (m *MockAuthorizer) AuthClient() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthClient")
	ret0, _ := ret[0].(bool)
	return ret0
}

// AuthClient indicates an expected call of AuthClient.
func (mr *MockAuthorizerMockRecorder) AuthClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthClient", reflect.TypeOf((*MockAuthorizer)(nil).AuthClient))
}

// CanRead mocks base method.
func (m *MockAuthorizer) CanRead() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CanRead")
	ret0, _ := ret[0].(error)
	return ret0
}

// CanRead indicates an expected call of CanRead.
func (mr *MockAuthorizerMockRecorder) CanRead() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CanRead", reflect.TypeOf((*MockAuthorizer)(nil).CanRead))
}

// CanWrite mocks base method.
func (m *MockAuthorizer) CanWrite() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CanWrite")
	ret0, _ := ret[0].(error)
	return ret0
}

// CanWrite indicates an expected call of CanWrite.
func (mr *MockAuthorizerMockRecorder) CanWrite() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CanWrite", reflect.TypeOf((*MockAuthorizer)(nil).CanWrite))
}

// MockUpgradeSeries is a mock of UpgradeSeries interface.
type MockUpgradeSeries struct {
	ctrl     *gomock.Controller
	recorder *MockUpgradeSeriesMockRecorder
}

// MockUpgradeSeriesMockRecorder is the mock recorder for MockUpgradeSeries.
type MockUpgradeSeriesMockRecorder struct {
	mock *MockUpgradeSeries
}

// NewMockUpgradeSeries creates a new mock instance.
func NewMockUpgradeSeries(ctrl *gomock.Controller) *MockUpgradeSeries {
	mock := &MockUpgradeSeries{ctrl: ctrl}
	mock.recorder = &MockUpgradeSeriesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUpgradeSeries) EXPECT() *MockUpgradeSeriesMockRecorder {
	return m.recorder
}

// Complete mocks base method.
func (m *MockUpgradeSeries) Complete(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Complete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Complete indicates an expected call of Complete.
func (mr *MockUpgradeSeriesMockRecorder) Complete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Complete", reflect.TypeOf((*MockUpgradeSeries)(nil).Complete), arg0)
}

// Prepare mocks base method.
func (m *MockUpgradeSeries) Prepare(arg0 context.Context, arg1, arg2 string, arg3 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Prepare", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// Prepare indicates an expected call of Prepare.
func (mr *MockUpgradeSeriesMockRecorder) Prepare(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Prepare", reflect.TypeOf((*MockUpgradeSeries)(nil).Prepare), arg0, arg1, arg2, arg3)
}

// Validate mocks base method.
func (m *MockUpgradeSeries) Validate(arg0 context.Context, arg1 []machinemanager.ValidationEntity) ([]machinemanager.ValidationResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", arg0, arg1)
	ret0, _ := ret[0].([]machinemanager.ValidationResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Validate indicates an expected call of Validate.
func (mr *MockUpgradeSeriesMockRecorder) Validate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockUpgradeSeries)(nil).Validate), arg0, arg1)
}

// MockUpgradeSeriesState is a mock of UpgradeSeriesState interface.
type MockUpgradeSeriesState struct {
	ctrl     *gomock.Controller
	recorder *MockUpgradeSeriesStateMockRecorder
}

// MockUpgradeSeriesStateMockRecorder is the mock recorder for MockUpgradeSeriesState.
type MockUpgradeSeriesStateMockRecorder struct {
	mock *MockUpgradeSeriesState
}

// NewMockUpgradeSeriesState creates a new mock instance.
func NewMockUpgradeSeriesState(ctrl *gomock.Controller) *MockUpgradeSeriesState {
	mock := &MockUpgradeSeriesState{ctrl: ctrl}
	mock.recorder = &MockUpgradeSeriesStateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUpgradeSeriesState) EXPECT() *MockUpgradeSeriesStateMockRecorder {
	return m.recorder
}

// ApplicationsFromMachine mocks base method.
func (m *MockUpgradeSeriesState) ApplicationsFromMachine(arg0 machinemanager.Machine) ([]machinemanager.Application, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ApplicationsFromMachine", arg0)
	ret0, _ := ret[0].([]machinemanager.Application)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ApplicationsFromMachine indicates an expected call of ApplicationsFromMachine.
func (mr *MockUpgradeSeriesStateMockRecorder) ApplicationsFromMachine(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplicationsFromMachine", reflect.TypeOf((*MockUpgradeSeriesState)(nil).ApplicationsFromMachine), arg0)
}

// MachineFromTag mocks base method.
func (m *MockUpgradeSeriesState) MachineFromTag(arg0 string) (machinemanager.Machine, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MachineFromTag", arg0)
	ret0, _ := ret[0].(machinemanager.Machine)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MachineFromTag indicates an expected call of MachineFromTag.
func (mr *MockUpgradeSeriesStateMockRecorder) MachineFromTag(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MachineFromTag", reflect.TypeOf((*MockUpgradeSeriesState)(nil).MachineFromTag), arg0)
}

// MockUpgradeBaseValidator is a mock of UpgradeBaseValidator interface.
type MockUpgradeBaseValidator struct {
	ctrl     *gomock.Controller
	recorder *MockUpgradeBaseValidatorMockRecorder
}

// MockUpgradeBaseValidatorMockRecorder is the mock recorder for MockUpgradeBaseValidator.
type MockUpgradeBaseValidatorMockRecorder struct {
	mock *MockUpgradeBaseValidator
}

// NewMockUpgradeBaseValidator creates a new mock instance.
func NewMockUpgradeBaseValidator(ctrl *gomock.Controller) *MockUpgradeBaseValidator {
	mock := &MockUpgradeBaseValidator{ctrl: ctrl}
	mock.recorder = &MockUpgradeBaseValidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUpgradeBaseValidator) EXPECT() *MockUpgradeBaseValidatorMockRecorder {
	return m.recorder
}

// ValidateApplications mocks base method.
func (m *MockUpgradeBaseValidator) ValidateApplications(arg0 context.Context, arg1 []machinemanager.Application, arg2 base.Base, arg3 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateApplications", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateApplications indicates an expected call of ValidateApplications.
func (mr *MockUpgradeBaseValidatorMockRecorder) ValidateApplications(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateApplications", reflect.TypeOf((*MockUpgradeBaseValidator)(nil).ValidateApplications), arg0, arg1, arg2, arg3)
}

// ValidateBase mocks base method.
func (m *MockUpgradeBaseValidator) ValidateBase(arg0, arg1 base.Base, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateBase", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateBase indicates an expected call of ValidateBase.
func (mr *MockUpgradeBaseValidatorMockRecorder) ValidateBase(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateBase", reflect.TypeOf((*MockUpgradeBaseValidator)(nil).ValidateBase), arg0, arg1, arg2)
}

// ValidateMachine mocks base method.
func (m *MockUpgradeBaseValidator) ValidateMachine(arg0 machinemanager.Machine) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateMachine", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateMachine indicates an expected call of ValidateMachine.
func (mr *MockUpgradeBaseValidatorMockRecorder) ValidateMachine(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateMachine", reflect.TypeOf((*MockUpgradeBaseValidator)(nil).ValidateMachine), arg0)
}

// ValidateUnits mocks base method.
func (m *MockUpgradeBaseValidator) ValidateUnits(arg0 []machinemanager.Unit) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateUnits", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateUnits indicates an expected call of ValidateUnits.
func (mr *MockUpgradeBaseValidatorMockRecorder) ValidateUnits(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateUnits", reflect.TypeOf((*MockUpgradeBaseValidator)(nil).ValidateUnits), arg0)
}
