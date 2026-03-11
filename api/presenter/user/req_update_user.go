package presenter

import (
	"be-dashboard-nba/usecase/entities"
	"be-dashboard-nba/repository/user"
	"database/sql"
)

type UpdateUserRequest struct {
	Name     string  `json:"name" validate:"required"`
	FullName string  `json:"full_name"  validate:"required"`
	Email    string  `json:"email"  validate:"required,email"`
	RoleID   int     `json:"role_id" validate:"required"`
	Phone    *string `json:"phone"`
	Active   *bool   `json:"active"  validate:"required"`
}

func (req *UpdateUserRequest) ToParams(userID string, updatedBy string, existingUser entities.User) (params repository.UpdateUserParams) {
	params = repository.UpdateUserParams{
		ID:        userID,
		Name:      req.Name,
		FullName:  req.FullName,
		Email:     req.Email,
		UpdatedBy: updatedBy,
		RoleID:    req.RoleID,
	}

	if req.Active != nil {
		params.Active = *req.Active
	} else {

		params.Active = existingUser.Active
	}

	if req.Phone != nil {
		if *req.Phone == "" {
			params.Phone = sql.NullString{Valid: false}
		} else {
			params.Phone = sql.NullString{String: *req.Phone, Valid: true}
		}
	} else {
		params.Phone = existingUser.Phone
	}

	return
}
