package presenter

import (
	"be-dashboard-nba/pkg/role/repository"
	"database/sql"
)

type CreateRoleRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Code        string  `json:"code" validate:"required,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty"`
}

func (req *CreateRoleRequest) ToParams(userID string) (params repository.CreateRolePayload) {
	params.Name = req.Name
	params.Code = req.Code
	params.CreatedBy = userID

	if req.Description != nil {
		params.Description = sql.NullString{
			String: *req.Description,
			Valid:  true,
		}
	}

	return
}
