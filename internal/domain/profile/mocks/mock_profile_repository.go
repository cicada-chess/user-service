// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/domain/profile/interfaces/repository.go

// Package mock_interfaces is a generated GoMock package.
package mock_interfaces

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/profile/entity"
)

// MockProfileRepository is a mock of ProfileRepository interface.
type MockProfileRepository struct {
	ctrl     *gomock.Controller
	recorder *MockProfileRepositoryMockRecorder
}

// MockProfileRepositoryMockRecorder is the mock recorder for MockProfileRepository.
type MockProfileRepositoryMockRecorder struct {
	mock *MockProfileRepository
}

// NewMockProfileRepository creates a new mock instance.
func NewMockProfileRepository(ctrl *gomock.Controller) *MockProfileRepository {
	mock := &MockProfileRepository{ctrl: ctrl}
	mock.recorder = &MockProfileRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProfileRepository) EXPECT() *MockProfileRepositoryMockRecorder {
	return m.recorder
}

// CheckProfileExists mocks base method.
func (m *MockProfileRepository) CheckProfileExists(ctx context.Context, userID string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckProfileExists", ctx, userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckProfileExists indicates an expected call of CheckProfileExists.
func (mr *MockProfileRepositoryMockRecorder) CheckProfileExists(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckProfileExists", reflect.TypeOf((*MockProfileRepository)(nil).CheckProfileExists), ctx, userID)
}

// CreateProfile mocks base method.
func (m *MockProfileRepository) CreateProfile(ctx context.Context, profile *entity.Profile) (*entity.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProfile", ctx, profile)
	ret0, _ := ret[0].(*entity.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProfile indicates an expected call of CreateProfile.
func (mr *MockProfileRepositoryMockRecorder) CreateProfile(ctx, profile interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProfile", reflect.TypeOf((*MockProfileRepository)(nil).CreateProfile), ctx, profile)
}

// GetByUserID mocks base method.
func (m *MockProfileRepository) GetByUserID(ctx context.Context, userID string) (*entity.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUserID", ctx, userID)
	ret0, _ := ret[0].(*entity.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUserID indicates an expected call of GetByUserID.
func (mr *MockProfileRepositoryMockRecorder) GetByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUserID", reflect.TypeOf((*MockProfileRepository)(nil).GetByUserID), ctx, userID)
}

// UpdateProfile mocks base method.
func (m *MockProfileRepository) UpdateProfile(ctx context.Context, profile *entity.Profile) (*entity.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProfile", ctx, profile)
	ret0, _ := ret[0].(*entity.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateProfile indicates an expected call of UpdateProfile.
func (mr *MockProfileRepositoryMockRecorder) UpdateProfile(ctx, profile interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProfile", reflect.TypeOf((*MockProfileRepository)(nil).UpdateProfile), ctx, profile)
}
