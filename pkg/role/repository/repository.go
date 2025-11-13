package repository

import (
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/pkg/entities"
	"context"
)

type Repository interface {
	ReadRolesQuery(ctx context.Context, args ReadRolesParams) (data []entities.Role, err error)
	ReadRolesCount(ctx context.Context, args ReadRolesParams) (count int, err error)
	ReadRoleByIDQuery(ctx context.Context, roleID int) (data entities.Role, err error)
	CreateRoleQuery(ctx context.Context, payload CreateRolePayload) (err error)
	UpdateRoleQuery(ctx context.Context, payload UpdateRolePayload) (err error)
	DeleteRoleQuery(ctx context.Context, roleID int) (err error)
	ReadRoleAccess(ctx context.Context, roleID int) (data []entities.RoleAccessResponse, err error)
	DeleteRoleAccess(ctx context.Context, payload UpdateRoleMenuPermission) (err error)
	CreateRoleAccess(ctx context.Context, payload UpdateRoleMenuPermission) (err error)
}

type repository struct {
	db db.Query
}

func NewRepository(db db.Query) Repository {
	return &repository{
		db: db,
	}
}
