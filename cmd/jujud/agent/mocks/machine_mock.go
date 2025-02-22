// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/cmd/jujud/agent (interfaces: CommandRunner)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	exec "github.com/juju/utils/v3/exec"
	gomock "go.uber.org/mock/gomock"
)

// MockCommandRunner is a mock of CommandRunner interface.
type MockCommandRunner struct {
	ctrl     *gomock.Controller
	recorder *MockCommandRunnerMockRecorder
}

// MockCommandRunnerMockRecorder is the mock recorder for MockCommandRunner.
type MockCommandRunnerMockRecorder struct {
	mock *MockCommandRunner
}

// NewMockCommandRunner creates a new mock instance.
func NewMockCommandRunner(ctrl *gomock.Controller) *MockCommandRunner {
	mock := &MockCommandRunner{ctrl: ctrl}
	mock.recorder = &MockCommandRunnerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommandRunner) EXPECT() *MockCommandRunnerMockRecorder {
	return m.recorder
}

// RunCommands mocks base method.
func (m *MockCommandRunner) RunCommands(arg0 exec.RunParams) (*exec.ExecResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunCommands", arg0)
	ret0, _ := ret[0].(*exec.ExecResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunCommands indicates an expected call of RunCommands.
func (mr *MockCommandRunnerMockRecorder) RunCommands(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunCommands", reflect.TypeOf((*MockCommandRunner)(nil).RunCommands), arg0)
}
