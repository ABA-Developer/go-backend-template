package presenter

import "be-dashboard-nba/pkg/role/repository"

type UpdateRoleAccessItem struct {
	AccessID  int   `json:"access_id" validate:"required"`
	HasAccess *bool `json:"has_access" validate:"required"`
}

type UpdateRoleAccessRequest struct {
	AccessItem []UpdateRoleAccessItem `json:"access_item" validate:"required,min=1,dive"`
}

func (req *UpdateRoleAccessItem) ToParams(roleID int) (params repository.UpdateRoleMenuPermission) {
	params.MenuPermissionID = req.AccessID
	params.RoleID = roleID

	return
}
