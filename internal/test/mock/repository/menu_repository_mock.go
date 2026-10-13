package mock

import (
	"context"

	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"

	"github.com/stretchr/testify/mock"
)

type MenuRepositoryMock struct {
	mock.Mock
}

func NewMenuRepository() *MenuRepositoryMock { return &MenuRepositoryMock{} }

func (m *MenuRepositoryMock) ReadSidebarMenuQry(ctx context.Context, userID string) ([]model.Menu, error) {
	call := m.Called(ctx, userID)

	var r0 []model.Menu
	raw0 := call.Get(0)
	if casted, ok := raw0.([]model.Menu); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *MenuRepositoryMock) ReadListMenuQry(ctx context.Context, params repoParams.ReadListMenuParams) ([]model.Menu, error) {
	call := m.Called(ctx, params)

	var r0 []model.Menu
	raw0 := call.Get(0)
	if casted, ok := raw0.([]model.Menu); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *MenuRepositoryMock) ReadCountMenuQry(ctx context.Context, params repoParams.ReadListMenuParams) (int64, error) {
	call := m.Called(ctx, params)

	var r0 int64
	raw0 := call.Get(0)
	if casted, ok := raw0.(int64); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *MenuRepositoryMock) ReadParentMenuQry(ctx context.Context) ([]model.Menu, error) {
	call := m.Called(ctx)

	var r0 []model.Menu
	raw0 := call.Get(0)
	if casted, ok := raw0.([]model.Menu); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *MenuRepositoryMock) CreateMenuQry(ctx context.Context, params repoParams.CreateMenuParams) error {
	call := m.Called(ctx, params)
	return call.Error(0)
}

func (m *MenuRepositoryMock) UpdateMenuQry(ctx context.Context, params repoParams.UpdateMenuParams) error {
	call := m.Called(ctx, params)
	return call.Error(0)
}

func (m *MenuRepositoryMock) DeleteMenuQry(ctx context.Context, menuID int) error {
	call := m.Called(ctx, menuID)
	return call.Error(0)
}

func (m *MenuRepositoryMock) ReadMenuByIDQry(ctx context.Context, menuID int) (model.Menu, error) {
	call := m.Called(ctx, menuID)

	var r0 model.Menu
	raw0 := call.Get(0)
	if casted, ok := raw0.(model.Menu); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *MenuRepositoryMock) ReadSortForGroupQry(ctx context.Context, group string) (int, error) {
	call := m.Called(ctx, group)

	var r0 int
	raw0 := call.Get(0)
	if casted, ok := raw0.(int); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *MenuRepositoryMock) ReadNextSortForParentQry(ctx context.Context, parentID int32) (int, error) {
	call := m.Called(ctx, parentID)

	var r0 int
	raw0 := call.Get(0)
	if casted, ok := raw0.(int); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *MenuRepositoryMock) UpdateMenuSortQry(ctx context.Context, params repoParams.UpdateMenuSortParams) error {
	call := m.Called(ctx, params)
	return call.Error(0)
}

func (m *MenuRepositoryMock) ReadNextSortForParentAndGroupQry(ctx context.Context, parentID int32, group string) (int, error) {
	call := m.Called(ctx, parentID, group)

	var r0 int
	raw0 := call.Get(0)
	if casted, ok := raw0.(int); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *MenuRepositoryMock) CountMenuChildrenQry(ctx context.Context, menuID int) (int, error) {
	call := m.Called(ctx, menuID)

	var r0 int
	raw0 := call.Get(0)
	if casted, ok := raw0.(int); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *MenuRepositoryMock) UpdateChildrenGroupQry(ctx context.Context, parentID int, newGroup string) error {
	call := m.Called(ctx, parentID, newGroup)
	return call.Error(0)
}
