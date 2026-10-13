package mock

import (
	"context"

	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"

	"github.com/stretchr/testify/mock"
)

type MenuPermissionRepositoryMock struct {
	mock.Mock
}

func NewMenuPermissionRepository() *MenuPermissionRepositoryMock {
	return &MenuPermissionRepositoryMock{}
}

func (m *MenuPermissionRepositoryMock) ReadMenuPermissionListQry(ctx context.Context, args repoParams.ReadMenuPermissionParams) ([]model.MenuPermission, error) {
	call := m.Called(ctx, args)

	var r0 []model.MenuPermission
	raw0 := call.Get(0)
	if casted, ok := raw0.([]model.MenuPermission); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *MenuPermissionRepositoryMock) ReadMenuPermissionByIdQry(ctx context.Context, id int) (model.MenuPermission, error) {
	call := m.Called(ctx, id)

	var r0 model.MenuPermission
	raw0 := call.Get(0)
	if casted, ok := raw0.(model.MenuPermission); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *MenuPermissionRepositoryMock) ReadMenuPermissionCountQry(ctx context.Context, args repoParams.ReadMenuPermissionParams) (int, error) {
	call := m.Called(ctx, args)

	var r0 int
	raw0 := call.Get(0)
	if casted, ok := raw0.(int); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *MenuPermissionRepositoryMock) UpdateMenuPermissionQry(ctx context.Context, params repoParams.UpdateMenuPermissionParams) error {
	call := m.Called(ctx, params)
	return call.Error(0)
}

func (m *MenuPermissionRepositoryMock) CreateMenuPermissionQry(ctx context.Context, params repoParams.CreateMenuPermissionParams) error {
	call := m.Called(ctx, params)
	return call.Error(0)
}

func (m *MenuPermissionRepositoryMock) DeleteMenuPermissionQry(ctx context.Context, id int) error {
	call := m.Called(ctx, id)
	return call.Error(0)
}
