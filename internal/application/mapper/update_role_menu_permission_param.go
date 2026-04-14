package mapper

import (
	"be-dashboard-nba/internal/application/dto"
	rolePresenter "be-dashboard-nba/internal/presentation/presenter/role"
)

func ToUpdateRoleMenuPermissionParam(item *rolePresenter.UpdateRoleAccessItem, roleID int) dto.UpdateRoleMenuPermission {
	return dto.UpdateRoleMenuPermission{
		MenuPermissionID: item.AccessID,
		RoleID:           roleID,
	}
}
