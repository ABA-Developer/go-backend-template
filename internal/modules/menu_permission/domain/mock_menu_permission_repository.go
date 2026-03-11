package domain

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockMenuPermissionRepository adalah mock implementation dari MenuPermissionRepository
type MockMenuPermissionRepository struct {
	mock.Mock
}

// =======================
// Queries (Read)
// =======================

func (m *MockMenuPermissionRepository) ReadMenuPermissionListQuery(ctx context.Context, params MenuPermissionFilter) ([]MenuPermissionDetail, error) {
	args := m.Called(ctx, params)

	var data []MenuPermissionDetail
	if args.Get(0) != nil {
		data = args.Get(0).([]MenuPermissionDetail)
	}

	return data, args.Error(1)
}

func (m *MockMenuPermissionRepository) ReadMenuPermissionCountQuery(ctx context.Context, params MenuPermissionFilter) (int, error) {
	args := m.Called(ctx, params)
	return args.Int(0), args.Error(1)
}

func (m *MockMenuPermissionRepository) ReadMenuPermissionByIDQuery(ctx context.Context, id int) (MenuPermissionDetail, error) {
	args := m.Called(ctx, id)

	var data MenuPermissionDetail
	if args.Get(0) != nil {
		data = args.Get(0).(MenuPermissionDetail)
	}

	return data, args.Error(1)
}

// =======================
// Commands (Write)
// =======================

func (m *MockMenuPermissionRepository) CreateMenuPermissionQuery(ctx context.Context, payload MenuPermissionCreatePayload) error {
	args := m.Called(ctx, payload)
	return args.Error(0)
}

func (m *MockMenuPermissionRepository) UpdateMenuPermissionQuery(ctx context.Context, payload MenuPermissionUpdatePayload) error {
	args := m.Called(ctx, payload)
	return args.Error(0)
}

func (m *MockMenuPermissionRepository) DeleteMenuPermissionQuery(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
