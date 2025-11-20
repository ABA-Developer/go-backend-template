package repository

import (
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/pkg/entities"
	"context"
)

type Repository interface {
	CreateUserQuery(ctx context.Context, args CreateUserParams) (userID string, err error)
	CreateUserRoleQuery(ctx context.Context, roleID int, userID string) (err error)
	UpdateUserQuery(ctx context.Context, args UpdateUserParams) (err error)
	UpdateUserRoleQuery(ctx context.Context, roleID int, userID string) (err error)
	DeleteUserQuery(ctx context.Context, id string) (err error)
	ReadUsersQuery(ctx context.Context, args ReadListUserParams) (data []entities.User, err error)
	ReadCountUserQuery(ctx context.Context, args ReadListUserParams) (count int, err error)
	ReadUserByIDQuery(ctx context.Context, id string) (data entities.User, err error)
	ReadUserProfileQuery(ctx context.Context, id string) (data entities.User, err error)
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
