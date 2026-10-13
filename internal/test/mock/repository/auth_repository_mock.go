package mock

import (
	"context"

	"github.com/ABA-Developer/go-backend-template/constant"
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"

	"github.com/stretchr/testify/mock"
)

type AuthRepositoryMock struct {
	mock.Mock
}

func NewAuthRepository() *AuthRepositoryMock { return &AuthRepositoryMock{} }

func (m *AuthRepositoryMock) CreateSessionQry(ctx context.Context, args repoParams.SessionParams) error {
	call := m.Called(ctx, args)
	return call.Error(0)
}

func (m *AuthRepositoryMock) UpdateSessionQry(ctx context.Context, args repoParams.SessionParams) error {
	call := m.Called(ctx, args)
	return call.Error(0)
}

func (m *AuthRepositoryMock) DeleteSessionQry(ctx context.Context, id string) error {
	call := m.Called(ctx, id)
	return call.Error(0)
}

func (m *AuthRepositoryMock) ReadDetailSessionQry(ctx context.Context, id string) (model.Session, error) {
	call := m.Called(ctx, id)

	var r0 model.Session
	raw0 := call.Get(0)
	if casted, ok := raw0.(model.Session); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *AuthRepositoryMock) ReadDetailUserByEmailQry(ctx context.Context, email string) (model.User, error) {
	call := m.Called(ctx, email)

	var r0 model.User
	raw0 := call.Get(0)
	if casted, ok := raw0.(model.User); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *AuthRepositoryMock) CreateLoginAttempQry(ctx context.Context, args repoParams.LoginAttempParams) error {
	call := m.Called(ctx, args)
	return call.Error(0)
}

func (m *AuthRepositoryMock) CreateLoginRecordQry(ctx context.Context, args repoParams.LoginRecord) error {
	call := m.Called(ctx, args)
	return call.Error(0)
}

func (m *AuthRepositoryMock) CheckPermissionQry(ctx context.Context, menuURL constant.MenuKey, userID string, permissionCode []string) (bool, error) {
	call := m.Called(ctx, menuURL, userID, permissionCode)
	return call.Bool(0), call.Error(1)
}

func (m *AuthRepositoryMock) ReadDetailUserByIdQry(ctx context.Context, id string) (model.User, error) {
	call := m.Called(ctx, id)

	var r0 model.User
	raw0 := call.Get(0)
	if casted, ok := raw0.(model.User); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}
