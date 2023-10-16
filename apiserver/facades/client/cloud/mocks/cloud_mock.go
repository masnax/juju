// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/client/cloud (interfaces: UserService,User,ModelCredentialService,CredentialService,CloudService,CloudPermissionService)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	cloud "github.com/juju/juju/apiserver/facades/client/cloud"
	cloud0 "github.com/juju/juju/cloud"
	permission "github.com/juju/juju/core/permission"
	watcher "github.com/juju/juju/core/watcher"
	credential "github.com/juju/juju/domain/credential"
	names "github.com/juju/names/v4"
	gomock "go.uber.org/mock/gomock"
)

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// User mocks base method.
func (m *MockUserService) User(arg0 names.UserTag) (cloud.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "User", arg0)
	ret0, _ := ret[0].(cloud.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// User indicates an expected call of User.
func (mr *MockUserServiceMockRecorder) User(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "User", reflect.TypeOf((*MockUserService)(nil).User), arg0)
}

// MockUser is a mock of User interface.
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
}

// MockUserMockRecorder is the mock recorder for MockUser.
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance.
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUser) EXPECT() *MockUserMockRecorder {
	return m.recorder
}

// DisplayName mocks base method.
func (m *MockUser) DisplayName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DisplayName")
	ret0, _ := ret[0].(string)
	return ret0
}

// DisplayName indicates an expected call of DisplayName.
func (mr *MockUserMockRecorder) DisplayName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisplayName", reflect.TypeOf((*MockUser)(nil).DisplayName))
}

// MockModelCredentialService is a mock of ModelCredentialService interface.
type MockModelCredentialService struct {
	ctrl     *gomock.Controller
	recorder *MockModelCredentialServiceMockRecorder
}

// MockModelCredentialServiceMockRecorder is the mock recorder for MockModelCredentialService.
type MockModelCredentialServiceMockRecorder struct {
	mock *MockModelCredentialService
}

// NewMockModelCredentialService creates a new mock instance.
func NewMockModelCredentialService(ctrl *gomock.Controller) *MockModelCredentialService {
	mock := &MockModelCredentialService{ctrl: ctrl}
	mock.recorder = &MockModelCredentialServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockModelCredentialService) EXPECT() *MockModelCredentialServiceMockRecorder {
	return m.recorder
}

// CloudCredentialUpdated mocks base method.
func (m *MockModelCredentialService) CloudCredentialUpdated(arg0 names.CloudCredentialTag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudCredentialUpdated", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CloudCredentialUpdated indicates an expected call of CloudCredentialUpdated.
func (mr *MockModelCredentialServiceMockRecorder) CloudCredentialUpdated(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudCredentialUpdated", reflect.TypeOf((*MockModelCredentialService)(nil).CloudCredentialUpdated), arg0)
}

// CredentialModels mocks base method.
func (m *MockModelCredentialService) CredentialModels(arg0 names.CloudCredentialTag) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CredentialModels", arg0)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CredentialModels indicates an expected call of CredentialModels.
func (mr *MockModelCredentialServiceMockRecorder) CredentialModels(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CredentialModels", reflect.TypeOf((*MockModelCredentialService)(nil).CredentialModels), arg0)
}

// CredentialModelsAndOwnerAccess mocks base method.
func (m *MockModelCredentialService) CredentialModelsAndOwnerAccess(arg0 names.CloudCredentialTag) ([]cloud0.CredentialOwnerModelAccess, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CredentialModelsAndOwnerAccess", arg0)
	ret0, _ := ret[0].([]cloud0.CredentialOwnerModelAccess)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CredentialModelsAndOwnerAccess indicates an expected call of CredentialModelsAndOwnerAccess.
func (mr *MockModelCredentialServiceMockRecorder) CredentialModelsAndOwnerAccess(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CredentialModelsAndOwnerAccess", reflect.TypeOf((*MockModelCredentialService)(nil).CredentialModelsAndOwnerAccess), arg0)
}

// RemoveModelsCredential mocks base method.
func (m *MockModelCredentialService) RemoveModelsCredential(arg0 names.CloudCredentialTag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveModelsCredential", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveModelsCredential indicates an expected call of RemoveModelsCredential.
func (mr *MockModelCredentialServiceMockRecorder) RemoveModelsCredential(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveModelsCredential", reflect.TypeOf((*MockModelCredentialService)(nil).RemoveModelsCredential), arg0)
}

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

// AllCloudCredentials mocks base method.
func (m *MockCredentialService) AllCloudCredentials(arg0 context.Context, arg1 string) ([]credential.CloudCredential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllCloudCredentials", arg0, arg1)
	ret0, _ := ret[0].([]credential.CloudCredential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllCloudCredentials indicates an expected call of AllCloudCredentials.
func (mr *MockCredentialServiceMockRecorder) AllCloudCredentials(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllCloudCredentials", reflect.TypeOf((*MockCredentialService)(nil).AllCloudCredentials), arg0, arg1)
}

// CloudCredential mocks base method.
func (m *MockCredentialService) CloudCredential(arg0 context.Context, arg1 names.CloudCredentialTag) (cloud0.Credential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudCredential", arg0, arg1)
	ret0, _ := ret[0].(cloud0.Credential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloudCredential indicates an expected call of CloudCredential.
func (mr *MockCredentialServiceMockRecorder) CloudCredential(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudCredential", reflect.TypeOf((*MockCredentialService)(nil).CloudCredential), arg0, arg1)
}

// CloudCredentials mocks base method.
func (m *MockCredentialService) CloudCredentials(arg0 context.Context, arg1, arg2 string) (map[string]cloud0.Credential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudCredentials", arg0, arg1, arg2)
	ret0, _ := ret[0].(map[string]cloud0.Credential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloudCredentials indicates an expected call of CloudCredentials.
func (mr *MockCredentialServiceMockRecorder) CloudCredentials(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudCredentials", reflect.TypeOf((*MockCredentialService)(nil).CloudCredentials), arg0, arg1, arg2)
}

// RemoveCloudCredential mocks base method.
func (m *MockCredentialService) RemoveCloudCredential(arg0 context.Context, arg1 names.CloudCredentialTag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveCloudCredential", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveCloudCredential indicates an expected call of RemoveCloudCredential.
func (mr *MockCredentialServiceMockRecorder) RemoveCloudCredential(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveCloudCredential", reflect.TypeOf((*MockCredentialService)(nil).RemoveCloudCredential), arg0, arg1)
}

// UpdateCloudCredential mocks base method.
func (m *MockCredentialService) UpdateCloudCredential(arg0 context.Context, arg1 names.CloudCredentialTag, arg2 cloud0.Credential) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCloudCredential", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCloudCredential indicates an expected call of UpdateCloudCredential.
func (mr *MockCredentialServiceMockRecorder) UpdateCloudCredential(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCloudCredential", reflect.TypeOf((*MockCredentialService)(nil).UpdateCloudCredential), arg0, arg1, arg2)
}

// WatchCredential mocks base method.
func (m *MockCredentialService) WatchCredential(arg0 context.Context, arg1 names.CloudCredentialTag) (watcher.Watcher[struct{}], error) {
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

// Delete mocks base method.
func (m *MockCloudService) Delete(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCloudServiceMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCloudService)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *MockCloudService) Get(arg0 context.Context, arg1 string) (*cloud0.Cloud, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*cloud0.Cloud)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCloudServiceMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCloudService)(nil).Get), arg0, arg1)
}

// ListAll mocks base method.
func (m *MockCloudService) ListAll(arg0 context.Context) ([]cloud0.Cloud, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAll", arg0)
	ret0, _ := ret[0].([]cloud0.Cloud)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAll indicates an expected call of ListAll.
func (mr *MockCloudServiceMockRecorder) ListAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAll", reflect.TypeOf((*MockCloudService)(nil).ListAll), arg0)
}

// Save mocks base method.
func (m *MockCloudService) Save(arg0 context.Context, arg1 cloud0.Cloud) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockCloudServiceMockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockCloudService)(nil).Save), arg0, arg1)
}

// MockCloudPermissionService is a mock of CloudPermissionService interface.
type MockCloudPermissionService struct {
	ctrl     *gomock.Controller
	recorder *MockCloudPermissionServiceMockRecorder
}

// MockCloudPermissionServiceMockRecorder is the mock recorder for MockCloudPermissionService.
type MockCloudPermissionServiceMockRecorder struct {
	mock *MockCloudPermissionService
}

// NewMockCloudPermissionService creates a new mock instance.
func NewMockCloudPermissionService(ctrl *gomock.Controller) *MockCloudPermissionService {
	mock := &MockCloudPermissionService{ctrl: ctrl}
	mock.recorder = &MockCloudPermissionServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCloudPermissionService) EXPECT() *MockCloudPermissionServiceMockRecorder {
	return m.recorder
}

// CloudsForUser mocks base method.
func (m *MockCloudPermissionService) CloudsForUser(arg0 names.UserTag) ([]cloud0.CloudAccess, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudsForUser", arg0)
	ret0, _ := ret[0].([]cloud0.CloudAccess)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloudsForUser indicates an expected call of CloudsForUser.
func (mr *MockCloudPermissionServiceMockRecorder) CloudsForUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudsForUser", reflect.TypeOf((*MockCloudPermissionService)(nil).CloudsForUser), arg0)
}

// CreateCloudAccess mocks base method.
func (m *MockCloudPermissionService) CreateCloudAccess(arg0 string, arg1 names.UserTag, arg2 permission.Access) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCloudAccess", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateCloudAccess indicates an expected call of CreateCloudAccess.
func (mr *MockCloudPermissionServiceMockRecorder) CreateCloudAccess(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCloudAccess", reflect.TypeOf((*MockCloudPermissionService)(nil).CreateCloudAccess), arg0, arg1, arg2)
}

// GetCloudAccess mocks base method.
func (m *MockCloudPermissionService) GetCloudAccess(arg0 string, arg1 names.UserTag) (permission.Access, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCloudAccess", arg0, arg1)
	ret0, _ := ret[0].(permission.Access)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCloudAccess indicates an expected call of GetCloudAccess.
func (mr *MockCloudPermissionServiceMockRecorder) GetCloudAccess(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCloudAccess", reflect.TypeOf((*MockCloudPermissionService)(nil).GetCloudAccess), arg0, arg1)
}

// GetCloudUsers mocks base method.
func (m *MockCloudPermissionService) GetCloudUsers(arg0 string) (map[string]permission.Access, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCloudUsers", arg0)
	ret0, _ := ret[0].(map[string]permission.Access)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCloudUsers indicates an expected call of GetCloudUsers.
func (mr *MockCloudPermissionServiceMockRecorder) GetCloudUsers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCloudUsers", reflect.TypeOf((*MockCloudPermissionService)(nil).GetCloudUsers), arg0)
}

// RemoveCloudAccess mocks base method.
func (m *MockCloudPermissionService) RemoveCloudAccess(arg0 string, arg1 names.UserTag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveCloudAccess", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveCloudAccess indicates an expected call of RemoveCloudAccess.
func (mr *MockCloudPermissionServiceMockRecorder) RemoveCloudAccess(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveCloudAccess", reflect.TypeOf((*MockCloudPermissionService)(nil).RemoveCloudAccess), arg0, arg1)
}

// UpdateCloudAccess mocks base method.
func (m *MockCloudPermissionService) UpdateCloudAccess(arg0 string, arg1 names.UserTag, arg2 permission.Access) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCloudAccess", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCloudAccess indicates an expected call of UpdateCloudAccess.
func (mr *MockCloudPermissionServiceMockRecorder) UpdateCloudAccess(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCloudAccess", reflect.TypeOf((*MockCloudPermissionService)(nil).UpdateCloudAccess), arg0, arg1, arg2)
}
