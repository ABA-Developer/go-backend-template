package domain

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) ReadRolesQuery(ctx context.Context, filter RoleFilter) ([]Role, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]Role), args.Error(1)
}

func (m *MockRoleRepository) ReadRolesCount(ctx context.Context, filter RoleFilter) (int, error) {
	args := m.Called(ctx, filter)
	return args.Int(0), args.Error(1)
}

func (m *MockRoleRepository) ReadRoleByIDQuery(ctx context.Context, roleID int) (Role, error) {
	args := m.Called(ctx, roleID)
	return args.Get(0).(Role), args.Error(1)
}

func (m *MockRoleRepository) CreateRoleQuery(ctx context.Context, payload CreateRolePayload) error {
	args := m.Called(ctx, payload)
	return args.Error(0)
}

func (m *MockRoleRepository) UpdateRoleQuery(ctx context.Context, payload UpdateRolePayload) error {
	args := m.Called(ctx, payload)
	return args.Error(0)
}

func (m *MockRoleRepository) DeleteRoleQuery(ctx context.Context, roleID int) error {
	args := m.Called(ctx, roleID)
	return args.Error(0)
}

func (m *MockRoleRepository) ReadRoleAccessQuery(ctx context.Context, filter RoleAccessFilter) ([]RoleAccessResponse, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]RoleAccessResponse), args.Error(1)
}

func (m *MockRoleRepository) ReadRoleAccessCount(ctx context.Context, filter RoleAccessFilter) (int, error) {
	args := m.Called(ctx, filter)
	return args.Int(0), args.Error(1)
}

func (m *MockRoleRepository) UpdateRoleAccessTx(ctx context.Context, roleID int, payloads []UpdateRoleMenuPermission) error {
	args := m.Called(ctx, roleID, payloads)
	return args.Error(0)
}
