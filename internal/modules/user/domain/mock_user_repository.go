package domain

import (
	"be-dashboard-nba/internal/core/db"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) WithTx(tx db.DB) UserRepository {
	args := m.Called(tx)
	return args.Get(0).(UserRepository)
}

func (m *MockUserRepository) CreateUserWithRoleTx(ctx context.Context, user User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUserWithRoleTx(ctx context.Context, user User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUserQuery(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) ReadUsersQuery(ctx context.Context, filter UserFilter) ([]UserListRow, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]UserListRow), args.Error(1)
}

func (m *MockUserRepository) ReadCountUserQuery(ctx context.Context, filter UserFilter) (int, error) {
	args := m.Called(ctx, filter)
	return args.Int(0), args.Error(1)
}

func (m *MockUserRepository) ReadUserByIDQuery(ctx context.Context, id string) (UserDetailRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(UserDetailRow), args.Error(1)
}

func (m *MockUserRepository) ReadUserProfileQuery(ctx context.Context, id string) (UserProfileRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(UserProfileRow), args.Error(1)
}

func (m *MockUserRepository) IsUserEmailExistsQuery(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) IsUpdateUserEmailExistsQuery(ctx context.Context, email, id string) (bool, error) {
	args := m.Called(ctx, email, id)
	return args.Bool(0), args.Error(1)
}
