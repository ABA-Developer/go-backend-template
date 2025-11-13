package repository

import (
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/pkg/entities"
	"context"
)

type Repository interface {
	ReadSidebarMenuQuery(ctx context.Context, userID string) (data []entities.Menu, err error)
	ReadListMenuQuery(ctx context.Context, params ReadListMenuParams) (data []entities.Menu, err error)
	ReadCountMenuQuery(ctx context.Context, params ReadListMenuParams) (count int64, err error)
	ReadParentMenuQuery(ctx context.Context) (data []entities.Menu, err error)
	CreateMenuQuery(ctx context.Context, params CreateMenuParams) (err error)
	UpdateMenuQuery(ctx context.Context, params UpdateMenuParams) (err error)
	DeleteMenuQuery(ctx context.Context, menuID int) (err error)
	ReadMenuByIDQuery(ctx context.Context, menuID int) (data entities.Menu, err error)
	ReadSortForGroup(ctx context.Context, group string) (int, error)
	ReadNextSortForParent(ctx context.Context, parentID int32) (int, error)
	UpdateMenuSortQuery(ctx context.Context, params UpdateMenuSortParams,
	) (err error)
}

type repository struct {
	db db.Query
}

func NewRepository(db db.Query) Repository {
	return &repository{
		db: db,
	}
}
