package presenter

import "be-dashboard-nba/pkg/menu_permission/repository"

type UpdateMenuPermissionRequest struct {
	Code       string `json:"code" validate:"required,min=1,max=50"`
	ActionName string `json:"action_name" validate:"required,min=1,max=50"`
}

func (req *UpdateMenuPermissionRequest) ToParams(userID string, menuPermissionID int) (payload repository.UpdateMenuPermissionPayload) {
	payload = repository.UpdateMenuPermissionPayload{
		Code:             req.Code,
		ActionName:       req.ActionName,
		UpdatedBy:        userID,
		MenuPermissionID: menuPermissionID,
	}

	return
}
