package mapper

import (
	"database/sql"

	"github.com/ABA-Developer/go-backend-template/internal/application/role/dto"
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
)

func ToUpdateRoleParams(req *dto.UpdateRolePayload, updatedBy string, roleID int) repoParams.UpdateRoleParams {
	params := repoParams.UpdateRoleParams{
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
