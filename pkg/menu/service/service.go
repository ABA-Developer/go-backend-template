package service

import (
	menuPresenter "be-dashboard-nba/api/presenter/menu"
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/pkg/entities"
	"context"

	"github.com/rs/zerolog"
)

type Service interface {
	CreateMenuService(ctx context.Context, request menuPresenter.CreateMenuRequest, userID string) (err error)
	UpdateMenuService(ctx context.Context, request menuPresenter.UpdateMenuRequest, userID string, MenuID int) (err error)
	ReadListMenuService(ctx context.Context, request menuPresenter.ReadMenuListRequest) (data []entities.Menu, err error)
	ReadSidebarMenuService(ctx context.Context, userID string) (data []entities.Menu, err error)
	DeleteMenuService(ctx context.Context, menuID int) (err error)
	UpdateMenuOrderService(ctx context.Context, request menuPresenter.UpdateMenuOrderRequest, userID string) (err error)
	ReadMenuDetailService(ctx context.Context, menuID int) (data entities.Menu, err error)
	ReadMenuParentService(ctx context.Context) (data []entities.Menu, err error)
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
