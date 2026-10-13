package mapper

import (
	"database/sql"

	"github.com/ABA-Developer/go-backend-template/internal/application/user/dto"
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
)

func ToCreateUserParams(req *dto.CreateUserPayload, createdBy, hashedPassword string) repoParams.CreateUserParams {
	params := repoParams.CreateUserParams{
		Name:      req.Name,
		FullName:  req.FullName,
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedBy: createdBy,
		RoleID:    req.RoleID,
	}

	if req.Active != nil {
		params.Active = *req.Active
	} else {
		params.Active = true
	}

	if req.Phone != nil {
		params.Phone = sql.NullString{String: *req.Phone, Valid: true}
	}

	return params
}
