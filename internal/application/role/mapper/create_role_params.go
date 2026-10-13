package mapper

import (
	"database/sql"

	"github.com/ABA-Developer/go-backend-template/internal/application/role/dto"
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
)

func ToCreateRoleParams(req *dto.CreateRolePayload, createdBy string) repoParams.CreateRoleParams {
	params := repoParams.CreateRoleParams{
		Name:      req.Name,
		Code:      req.Code,
		CreatedBy: createdBy,
	}

	if req.Description != nil {
		params.Description = sql.NullString{String: *req.Description, Valid: true}
	}

	return params
}
