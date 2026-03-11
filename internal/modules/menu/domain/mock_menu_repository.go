package domain

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockMenuRepository struct {
	mock.Mock
}

// =======================
// Queries (Read)
// =======================

func (m *MockMenuRepository) ReadSidebarMenuQuery(ctx context.Context, userID string) ([]Menu, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]Menu), args.Error(1)
}

func (m *MockMenuRepository) ReadListMenuQuery(ctx context.Context, params MenuFilter) ([]Menu, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]Menu), args.Error(1)
}

func (m *MockMenuRepository) ReadCountMenuQuery(ctx context.Context, params MenuFilter) (int64, error) {
	args := m.Called(ctx, params)
	// Gunakan args.Get(0).(int64) karena testify tidak memiliki arg.Int64() bawaan
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockMenuRepository) ReadParentMenuQuery(ctx context.Context) ([]MenuParent, error) {
	args := m.Called(ctx)
	return args.Get(0).([]MenuParent), args.Error(1)
}

func (m *MockMenuRepository) ReadMenuByIDQuery(ctx context.Context, menuID int) (Menu, error) {
	args := m.Called(ctx, menuID)
	return args.Get(0).(Menu), args.Error(1)
}

func (m *MockMenuRepository) CountMenuChildren(ctx context.Context, menuID int) (int, error) {
	args := m.Called(ctx, menuID)
	return args.Int(0), args.Error(1)
}

// =======================
// Commands (Write / Tx)
// =======================

func (m *MockMenuRepository) DeleteMenuQuery(ctx context.Context, menuID int) error {
	args := m.Called(ctx, menuID)
	return args.Error(0)
}

func (m *MockMenuRepository) CreateMenuQuery(ctx context.Context, params MenuCreatePayload, userID string) error {
	args := m.Called(ctx, params, userID)
	return args.Error(0)
}

func (m *MockMenuRepository) UpdateMenuQuery(ctx context.Context, params MenuUpdatePayload, updateChildrenGroup bool) error {
	args := m.Called(ctx, params, updateChildrenGroup)
	return args.Error(0)
}

func (m *MockMenuRepository) UpdateMenuOrderQuery(ctx context.Context, payloads []MenuUpdateSortItemPayload) error {
	args := m.Called(ctx, payloads)
	return args.Error(0)
}
