package repository

import (
	"context"

	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
)

type RoleRepository interface {
	ReadRolesQry(ctx context.Context, args repoParams.ReadRolesParams) (data []model.Role, err error)
	ReadRolesCountQry(ctx context.Context, args repoParams.ReadRolesParams) (count int, err error)
	ReadRoleByIDQry(ctx context.Context, roleID int) (data model.Role, err error)
	CreateRoleQry(ctx context.Context, params repoParams.CreateRoleParams) (err error)
	UpdateRoleQry(ctx context.Context, params repoParams.UpdateRoleParams) (err error)
	DeleteRoleQry(ctx context.Context, roleID int) (err error)
	ReadRoleAccessQry(ctx context.Context, args repoParams.ReadRoleAccessParams) (data []model.RoleAccessResponse, err error)
	ReadRoleAccessCountQry(ctx context.Context, args repoParams.ReadRoleAccessParams) (count int, err error)
	DeleteRoleAccessQry(ctx context.Context, payload repoParams.UpdateRoleMenuPermission) (err error)
	CreateRoleAccessQry(ctx context.Context, payload repoParams.UpdateRoleMenuPermission) (err error)
}
