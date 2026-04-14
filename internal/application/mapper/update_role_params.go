package mapper

import (
	"database/sql"

	"be-dashboard-nba/internal/application/dto"
	rolePresenter "be-dashboard-nba/internal/presentation/presenter/role"
)

func ToUpdateRoleParams(req *rolePresenter.UpdateRoleRequest, updatedBy string, roleID int) dto.UpdateRoleParams {
	params := dto.UpdateRoleParams{
		RoleID:    roleID,
		Name:      req.Name,
		Code:      req.Code,
		UpdatedBy: updatedBy,
	}

	if req.Description != nil {
		params.Description = sql.NullString{String: *req.Description, Valid: true}
	}

	return params
}
