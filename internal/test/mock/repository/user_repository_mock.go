package mock

import (
	"context"

	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func NewUserRepository() *UserRepositoryMock { return &UserRepositoryMock{} }

func (m *UserRepositoryMock) CreateUserQry(ctx context.Context, args repoParams.CreateUserParams) (string, error) {
	call := m.Called(ctx, args)

	var r0 string
	raw0 := call.Get(0)
	if casted, ok := raw0.(string); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *UserRepositoryMock) CreateUserRoleQry(ctx context.Context, roleID int, userID string) error {
	call := m.Called(ctx, roleID, userID)
	return call.Error(0)
}

func (m *UserRepositoryMock) UpdateUserQry(ctx context.Context, args repoParams.UpdateUserParams) error {
	call := m.Called(ctx, args)
	return call.Error(0)
}

func (m *UserRepositoryMock) UpdateUserRoleQry(ctx context.Context, roleID int, userID string) error {
	call := m.Called(ctx, roleID, userID)
	return call.Error(0)
}

func (m *UserRepositoryMock) DeleteUserQry(ctx context.Context, id string) error {
	call := m.Called(ctx, id)
	return call.Error(0)
}

func (m *UserRepositoryMock) ReadUsersQry(ctx context.Context, args repoParams.ReadListUserParams) ([]model.User, error) {
	call := m.Called(ctx, args)

	var r0 []model.User
	raw0 := call.Get(0)
	if casted, ok := raw0.([]model.User); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *UserRepositoryMock) ReadCountUserQry(ctx context.Context, args repoParams.ReadListUserParams) (int, error) {
	call := m.Called(ctx, args)

	var r0 int
	raw0 := call.Get(0)
	if casted, ok := raw0.(int); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *UserRepositoryMock) ReadUserByIDQry(ctx context.Context, id string) (model.User, error) {
	call := m.Called(ctx, id)

	var r0 model.User
	raw0 := call.Get(0)
	if casted, ok := raw0.(model.User); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *UserRepositoryMock) ReadUserProfileQry(ctx context.Context, id string) (model.User, error) {
	call := m.Called(ctx, id)

	var r0 model.User
	raw0 := call.Get(0)
	if casted, ok := raw0.(model.User); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *UserRepositoryMock) IsUserEmailExistsQry(ctx context.Context, email string) (bool, error) {
	call := m.Called(ctx, email)
	return call.Bool(0), call.Error(1)
}

func (m *UserRepositoryMock) IsUpdateUserEmailExistsQry(ctx context.Context, email, id string) (bool, error) {
	call := m.Called(ctx, email, id)
	return call.Bool(0), call.Error(1)
}
