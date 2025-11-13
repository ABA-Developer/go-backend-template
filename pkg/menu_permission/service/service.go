package service

import (
	menuPermissionPresenter "be-dashboard-nba/api/presenter/menu_permission"
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/pkg/entities"
	"context"

	"github.com/rs/zerolog"
)

type Service interface {
	ReadMenuPermissionListParams(ctx context.Context, args menuPermissionPresenter.ReadMenuPermissionListRequest, MenuID int) (data entities.MenuPermissionPaginationResponse, err error)
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
