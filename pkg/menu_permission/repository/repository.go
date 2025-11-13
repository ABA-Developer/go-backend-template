package repository

import (
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/pkg/entities"
	"context"
)

type Repository interface {
	ReadMenuPermissionListQuery(ctx context.Context, args ReadMenuPermissionListParams) (data []entities.MenuPermission, err error)
	ReadMenuPermissionByIdQuery(ctx context.Context, MenuPermissionID int) (data entities.MenuPermission, err error)
	ReadMenuPermissionCount(ctx context.Context, args ReadMenuPermissionListParams) (count int, err error)
	UpdateMenuPermissionQuery(ctx context.Context, payload UpdateMenuPermissionPayload) (err error)
	CreateMenuPermissionQuery(ctx context.Context, payload CreateMenuPermissionPayload) (err error)
	DeleteMenuPermissionQuery(ctx context.Context, MenuPermissionID int) (err error)
}

type repository struct {
	db db.Query
}

func NewRepository(db db.Query) Repository {
	return &repository{
		db: db,
	}
}
