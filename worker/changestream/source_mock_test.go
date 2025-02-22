// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/core/changestream (interfaces: EventSource)

// Package changestream is a generated GoMock package.
package changestream

import (
	reflect "reflect"

	changestream "github.com/juju/juju/core/changestream"
	gomock "go.uber.org/mock/gomock"
)

// MockEventSource is a mock of EventSource interface.
type MockEventSource struct {
	ctrl     *gomock.Controller
	recorder *MockEventSourceMockRecorder
}

// MockEventSourceMockRecorder is the mock recorder for MockEventSource.
type MockEventSourceMockRecorder struct {
	mock *MockEventSource
}

// NewMockEventSource creates a new mock instance.
func NewMockEventSource(ctrl *gomock.Controller) *MockEventSource {
	mock := &MockEventSource{ctrl: ctrl}
	mock.recorder = &MockEventSourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventSource) EXPECT() *MockEventSourceMockRecorder {
	return m.recorder
}

// Subscribe mocks base method.
func (m *MockEventSource) Subscribe(arg0 ...changestream.SubscriptionOption) (changestream.Subscription, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Subscribe", varargs...)
	ret0, _ := ret[0].(changestream.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockEventSourceMockRecorder) Subscribe(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockEventSource)(nil).Subscribe), arg0...)
}
