package domain

import (
	"be-dashboard-nba/constant"
	"context"
	"database/sql"

	"github.com/stretchr/testify/mock"
)

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) WithTx(tx *sql.Tx) AuthRepository {
	m.Called(tx)
	return m
}

func (m *MockAuthRepository) CreateSessionQuery(ctx context.Context, args Session) error {
	ret := m.Called(ctx, args)
	return ret.Error(0)
}

func (m *MockAuthRepository) UpdateSessionQuery(ctx context.Context, args Session) error {
	ret := m.Called(ctx, args)
	return ret.Error(0)
}

func (m *MockAuthRepository) DeleteSessionQuery(ctx context.Context, id string) error {
	ret := m.Called(ctx, id)
	return ret.Error(0)
}

func (m *MockAuthRepository) ReadDetailSessionQuery(ctx context.Context, id string) (Session, error) {
	ret := m.Called(ctx, id)
	return ret.Get(0).(Session), ret.Error(1)
}

func (m *MockAuthRepository) ReadDetailUserByEmailQuery(ctx context.Context, email string) (User, error) {
	ret := m.Called(ctx, email)
	return ret.Get(0).(User), ret.Error(1)
}

func (m *MockAuthRepository) CreateLoginAttemp(ctx context.Context, args LoginAttemp) error {
	ret := m.Called(ctx, args)
	return ret.Error(0)
}

func (m *MockAuthRepository) CreateLoginRecord(ctx context.Context, args LoginRecord) error {
	ret := m.Called(ctx, args)
	return ret.Error(0)
}

func (m *MockAuthRepository) CheckPermissionQuery(ctx context.Context, menuURL constant.MenuKey, userID string, permissionCode []string) (bool, error) {
	ret := m.Called(ctx, menuURL, userID, permissionCode)
	return ret.Get(0).(bool), ret.Error(1)
}

func (m *MockAuthRepository) ReadDetailUserByIdQuery(ctx context.Context, id string) (User, error) {
	ret := m.Called(ctx, id)
	return ret.Get(0).(User), ret.Error(1)
}
