package repository

import (
	"context"

	"github.com/ABA-Developer/go-backend-template/constant"
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
)

type AuthRepository interface {
	CreateSessionQry(ctx context.Context, args repoParams.SessionParams) (err error)
	UpdateSessionQry(ctx context.Context, args repoParams.SessionParams) (err error)
	DeleteSessionQry(ctx context.Context, id string) (err error)
	ReadDetailSessionQry(ctx context.Context, id string) (data model.Session, err error)
	ReadDetailUserByEmailQry(ctx context.Context, email string) (data model.User, err error)
	CreateLoginAttempQry(ctx context.Context, args repoParams.LoginAttempParams) (err error)
	CreateLoginRecordQry(ctx context.Context, args repoParams.LoginRecord) (err error)
	CheckPermissionQry(ctx context.Context, menuURL constant.MenuKey, userID string, permissionCode []string) (bool, error)
	ReadDetailUserByIdQry(ctx context.Context, id string) (data model.User, err error)
}
