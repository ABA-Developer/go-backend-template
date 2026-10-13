package repository

import (
	"context"

	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
)

type UserRepository interface {
	CreateUserQry(ctx context.Context, args repoParams.CreateUserParams) (userID string, err error)
	CreateUserRoleQry(ctx context.Context, roleID int, userID string) (err error)
	UpdateUserQry(ctx context.Context, args repoParams.UpdateUserParams) (err error)
	UpdateUserRoleQry(ctx context.Context, roleID int, userID string) (err error)
	DeleteUserQry(ctx context.Context, id string) (err error)
	ReadUsersQry(ctx context.Context, args repoParams.ReadListUserParams) (data []model.User, err error)
	ReadCountUserQry(ctx context.Context, args repoParams.ReadListUserParams) (count int, err error)
	ReadUserByIDQry(ctx context.Context, id string) (data model.User, err error)
	ReadUserProfileQry(ctx context.Context, id string) (data model.User, err error)
	IsUserEmailExistsQry(ctx context.Context, email string) (exists bool, err error)
	IsUpdateUserEmailExistsQry(ctx context.Context, email, id string) (exists bool, err error)
}
