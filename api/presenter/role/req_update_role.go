package presenter

import (
	"be-dashboard-nba/pkg/role/repository"
	"database/sql"
)

type UpdateRoleRequest struct {
	Name        string  `json:"name" validate:"required"`
	Code        string  `json:"code" validate:"required"`
	Description *string `json:"description" validate:"omitempty"`
}

func (req *UpdateRoleRequest) ToParams(userID string, roleID int) (params repository.UpdateRolePayload) {
	params.Name = req.Name
	params.Code = req.Code
	params.UpdatedBy = userID
	params.RoleID = roleID

	if req.Description != nil {
		params.Description = sql.NullString{
			String: *req.Description,
			Valid:  true,
		}
	}

	return
}
