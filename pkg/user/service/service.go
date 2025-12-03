package service

import (
	presenter "be-dashboard-nba/api/presenter/user"
	userPresenter "be-dashboard-nba/api/presenter/user"
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/internal/utils"
	"be-dashboard-nba/pkg/entities"
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
)

type Service interface {
	ReadDetailUserService(ctx context.Context, id string) (data entities.User, err error)
	UpdateUserService(ctx context.Context, request presenter.UpdateUserRequest, updatedBy string, userID string) (err error)
	ReadUsersService(ctx context.Context, req userPresenter.ReadUserRequest) (data entities.UserPaginationResponse, err error)
	CreateUserService(ctx context.Context, req userPresenter.CreateUserRequest, userID string) (err error)
	DeleteUserService(ctx context.Context, userID string, deletedBy string) (err error)
	ReadUserProfile(ctx context.Context, userID string) (data entities.User, err error)
}

type service struct {
	db  db.DB
	log *zerolog.Logger
}

func NewService(db db.DB, log *zerolog.Logger) Service {
	return &service{
		db:  db,
		log: log,
	}
}

func newTestService(t *testing.T) (*service, sqlmock.Sqlmock, func()) {
	mockDB, mock, cleanup := utils.NewMockDB(t)
	svc := NewService(mockDB, &zerolog.Logger{}).(*service)
	return svc, mock, cleanup
}
