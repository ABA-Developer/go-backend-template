package repository

import (
	"context"

	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
)

type MenuPermissionRepository interface {
	ReadMenuPermissionListQry(ctx context.Context, args repoParams.ReadMenuPermissionParams) (data []model.MenuPermission, err error)
	ReadMenuPermissionByIdQry(ctx context.Context, MenuPermissionID int) (data model.MenuPermission, err error)
	ReadMenuPermissionCountQry(ctx context.Context, args repoParams.ReadMenuPermissionParams) (count int, err error)
	UpdateMenuPermissionQry(ctx context.Context, params repoParams.UpdateMenuPermissionParams) (err error)
	CreateMenuPermissionQry(ctx context.Context, params repoParams.CreateMenuPermissionParams) (err error)
	DeleteMenuPermissionQry(ctx context.Context, MenuPermissionID int) (err error)
}
