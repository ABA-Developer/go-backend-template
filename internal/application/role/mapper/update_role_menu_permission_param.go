package mapper

import (
	"github.com/ABA-Developer/go-backend-template/internal/application/role/dto"
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
)

func ToUpdateRoleMenuPermissionParam(item *dto.UpdateRoleAccessPayloadItem, roleID int) repoParams.UpdateRoleMenuPermission {
	return repoParams.UpdateRoleMenuPermission{
		MenuPermissionID: item.AccessID,
		RoleID:           roleID,
	}
}
