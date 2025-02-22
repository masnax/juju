// Code generated by MockGen. DO NOT EDIT.
// Source: ./consume.go

// Package eventsource is a generated GoMock package.
package eventsource

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockWatcher is a mock of Watcher interface.
type MockWatcher[T any] struct {
	ctrl     *gomock.Controller
	recorder *MockWatcherMockRecorder[T]
}

// MockWatcherMockRecorder is the mock recorder for MockWatcher.
type MockWatcherMockRecorder[T any] struct {
	mock *MockWatcher[T]
}

// NewMockWatcher creates a new mock instance.
func NewMockWatcher[T any](ctrl *gomock.Controller) *MockWatcher[T] {
	mock := &MockWatcher[T]{ctrl: ctrl}
	mock.recorder = &MockWatcherMockRecorder[T]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWatcher[T]) EXPECT() *MockWatcherMockRecorder[T] {
	return m.recorder
}

// Changes mocks base method.
func (m *MockWatcher[T]) Changes() <-chan T {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Changes")
	ret0, _ := ret[0].(<-chan T)
	return ret0
}

// Changes indicates an expected call of Changes.
func (mr *MockWatcherMockRecorder[T]) Changes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Changes", reflect.TypeOf((*MockWatcher[T])(nil).Changes))
}

// Kill mocks base method.
func (m *MockWatcher[T]) Kill() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Kill")
}

// Kill indicates an expected call of Kill.
func (mr *MockWatcherMockRecorder[T]) Kill() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Kill", reflect.TypeOf((*MockWatcher[T])(nil).Kill))
}

// Wait mocks base method.
func (m *MockWatcher[T]) Wait() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Wait")
	ret0, _ := ret[0].(error)
	return ret0
}

// Wait indicates an expected call of Wait.
func (mr *MockWatcherMockRecorder[T]) Wait() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Wait", reflect.TypeOf((*MockWatcher[T])(nil).Wait))
}
