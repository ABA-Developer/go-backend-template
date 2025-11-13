package service

import (
	rolePresenter "be-dashboard-nba/api/presenter/role"
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/pkg/entities"
	"context"

	"github.com/rs/zerolog"
)

type Service interface {
	ReadRolesService(ctx context.Context, args rolePresenter.ReadRolesRequest) (data entities.RolePaginationResponse, err error)
	ReadRoleDetail(ctx context.Context, roleID int) (data entities.Role, err error)
	CreateRoleService(ctx context.Context, paylaod rolePresenter.CreateRoleRequest, userID string) (err error)
	UpdateRoleService(ctx context.Context, payload rolePresenter.UpdateRoleRequest, userID string, roleID int) (err error)
	DeleteRoleService(ctx context.Context, roleID int) (err error)
	ReadRoleAccessService(ctx context.Context, roleID int) (data []entities.RoleAccessResponse, err error)
	UpdateRoleAccessService(ctx context.Context, roleID int, request rolePresenter.UpdateRoleAccessRequest) (err error)
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
