package repository

import (
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/pkg/entities"
	"context"
)

type Repository interface {
	CreateUserQuery(ctx context.Context, args CreateUserParams) (err error)
	UpdateUserQuery(ctx context.Context, args UpdateUserParams) (err error)
	DeleteUserQuery(ctx context.Context, id string) (err error)
	ReadListUserQuery(ctx context.Context, args ReadListUserParams) (data []entities.User, err error)
	GetCountUserQuery(ctx context.Context, args ReadListUserParams) (count int64, err error)
	ReadDetailUserQuery(ctx context.Context, id string) (data entities.User, err error)
	IsUserEmailExistsQuery(ctx context.Context, email string) (exists bool, err error)
	IsUpdateUserEmailExistsQuery(ctx context.Context, email, id string) (exists bool, err error)
}

type repository struct {
	db db.Query
}

func NewRepository(db db.Query) Repository {
	return &repository{
		db: db,
	}
}
