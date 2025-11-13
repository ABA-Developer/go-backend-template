package presenter

import "be-dashboard-nba/pkg/menu_permission/repository"

type CreateMenuPermissionRequest struct {
	Code       string `json:"code" validate:"required,min=1,max=50"`
	ActionName string `json:"action_name" validate:"required,min=1,max=50"`
}

func (req *CreateMenuPermissionRequest) ToParams(userID string, menuID int) (payload repository.CreateMenuPermissionPayload) {
	payload = repository.CreateMenuPermissionPayload{
		Code:       req.Code,
		ActionName: req.ActionName,
		MenuID:     menuID,
		CreatedBy:  userID,
	}

	return
}
