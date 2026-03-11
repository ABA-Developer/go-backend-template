package menu_permission

import (
	menuPermissionPresenter "be-dashboard-nba/api/presenter/menu_permission"
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/internal/utils"
	"be-dashboard-nba/usecase/entities"
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
)

type Service interface {
	ReadMenuPermissionService(ctx context.Context, args menuPermissionPresenter.ReadMenuPermissionListRequest, MenuID int) (data entities.MenuPermissionPaginationResponse, err error)
	CreateMenuPermissionService(ctx context.Context, payload menuPermissionPresenter.CreateMenuPermissionRequest, userID string, menuID int) (err error)
	UpdateMenuPermissionService(ctx context.Context, payload menuPermissionPresenter.UpdateMenuPermissionRequest, userID string, menuPermissionID int) (err error)
	ReadMenuPermissionDetail(ctx context.Context, MenuPermissionID int) (data entities.MenuPermission, err error)
	DeleteMenuPermissionService(ctx context.Context, menuPermissionID int) (err error)
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
