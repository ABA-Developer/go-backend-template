package presenter

import "be-dashboard-nba/pkg/entities"

type ReadMenuPermissionResponse struct {
	ID         int    `json:"id"`
	MenuID     int    `json:"menu_id"`
	Code       string `json:"code"`
	ActionName string `json:"action_name"`
}

func ToReadMenuPermissionListResponse(entity entities.MenuPermission) (response ReadMenuPermissionResponse) {
	response.ID = entity.ID
	response.MenuID = entity.MenuID
	response.Code = entity.Code
	response.ActionName = entity.ActionName

	return
}

func ToReadMenuPermissionListResponses(entities []entities.MenuPermission) (responses []ReadMenuPermissionResponse) {
	responses = make([]ReadMenuPermissionResponse, len(entities))

	for i := range entities {
		responses[i] = ToReadMenuPermissionListResponse(entities[i])
	}

	return
}
