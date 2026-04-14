package mapper

import (
	"database/sql"

	"be-dashboard-nba/internal/application/dto"
	rolePresenter "be-dashboard-nba/internal/presentation/presenter/role"
)

func ToCreateRoleParams(req *rolePresenter.CreateRoleRequest, createdBy string) dto.CreateRoleParams {
	params := dto.CreateRoleParams{
		Name:      req.Name,
		Code:      req.Code,
		CreatedBy: createdBy,
	}

	if req.Description != nil {
		params.Description = sql.NullString{String: *req.Description, Valid: true}
	}

	return params
}
