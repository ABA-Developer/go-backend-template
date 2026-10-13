package mock

import (
	"context"

	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"

	"github.com/stretchr/testify/mock"
)

type RoleRepositoryMock struct {
	mock.Mock
}

func NewRoleRepository() *RoleRepositoryMock { return &RoleRepositoryMock{} }

func (m *RoleRepositoryMock) ReadRolesQry(ctx context.Context, args repoParams.ReadRolesParams) ([]model.Role, error) {
	call := m.Called(ctx, args)

	var r0 []model.Role
	raw0 := call.Get(0)
	if casted, ok := raw0.([]model.Role); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *RoleRepositoryMock) ReadRolesCountQry(ctx context.Context, args repoParams.ReadRolesParams) (int, error) {
	call := m.Called(ctx, args)

	var r0 int
	raw0 := call.Get(0)
	if casted, ok := raw0.(int); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *RoleRepositoryMock) ReadRoleByIDQry(ctx context.Context, roleID int) (model.Role, error) {
	call := m.Called(ctx, roleID)

	var r0 model.Role
	raw0 := call.Get(0)
	if casted, ok := raw0.(model.Role); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *RoleRepositoryMock) CreateRoleQry(ctx context.Context, params repoParams.CreateRoleParams) error {
	call := m.Called(ctx, params)
	return call.Error(0)
}

func (m *RoleRepositoryMock) UpdateRoleQry(ctx context.Context, params repoParams.UpdateRoleParams) error {
	call := m.Called(ctx, params)
	return call.Error(0)
}

func (m *RoleRepositoryMock) DeleteRoleQry(ctx context.Context, roleID int) error {
	call := m.Called(ctx, roleID)
	return call.Error(0)
}

func (m *RoleRepositoryMock) ReadRoleAccessQry(ctx context.Context, args repoParams.ReadRoleAccessParams) ([]model.RoleAccessResponse, error) {
	call := m.Called(ctx, args)

	var r0 []model.RoleAccessResponse
	raw0 := call.Get(0)
	if casted, ok := raw0.([]model.RoleAccessResponse); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *RoleRepositoryMock) ReadRoleAccessCountQry(ctx context.Context, args repoParams.ReadRoleAccessParams) (int, error) {
	call := m.Called(ctx, args)

	var r0 int
	raw0 := call.Get(0)
	if casted, ok := raw0.(int); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *RoleRepositoryMock) DeleteRoleAccessQry(ctx context.Context, payload repoParams.UpdateRoleMenuPermission) error {
	call := m.Called(ctx, payload)
	return call.Error(0)
}

func (m *RoleRepositoryMock) CreateRoleAccessQry(ctx context.Context, payload repoParams.UpdateRoleMenuPermission) error {
	call := m.Called(ctx, payload)
	return call.Error(0)
}
