package repository

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/pkg/entities"
	"context"
)

type Repository interface {
	CreateSessionQuery(ctx context.Context, args SessionParams) (err error)
	UpdateSessionQuery(ctx context.Context, args SessionParams) (err error)
	DeleteSessionQuery(ctx context.Context, id string) (err error)
	ReadDetailSessionQuery(ctx context.Context, id string) (data entities.Session, err error)
	ReadDetailUserByEmailQuery(ctx context.Context, email string) (data entities.User, err error)
	CreateLoginAttemp(ctx context.Context, args LoginAttempParams) (err error)
	CreateLoginRecord(ctx context.Context, args LoginRecord) (err error)
	CheckPermissionQuery(ctx context.Context, menuURL constant.MenuKey, userID string, permissionCode []string) (bool, error)
	ReadDetailUserByIdQuery(ctx context.Context, id string) (data entities.User, err error)
}

type repository struct {
	DB db.Query
}

func NewRepo(db db.Query) Repository {
	return &repository{
		DB: db,
	}
}
