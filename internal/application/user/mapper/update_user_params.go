package mapper

import (
	"database/sql"

	"github.com/ABA-Developer/go-backend-template/internal/application/user/dto"
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
)

func ToUpdateUserParams(req *dto.UpdateUserPayload, userID, updatedBy string, existingUser model.User) repoParams.UpdateUserParams {
	params := repoParams.UpdateUserParams{
		ID:        userID,
		Name:      req.Name,
		FullName:  req.FullName,
		Email:     req.Email,
		UpdatedBy: updatedBy,
		RoleID:    req.RoleID,
		Active:    existingUser.Active,
		Phone:     existingUser.Phone,
	}

	if req.Active != nil {
		params.Active = *req.Active
	}

	if req.Phone != nil {
		// Preserve prior behavior: empty string means NULL.
		if *req.Phone == "" {
			params.Phone = sql.NullString{Valid: false}
		} else {
			params.Phone = sql.NullString{String: *req.Phone, Valid: true}
		}
	}

	return params
}
