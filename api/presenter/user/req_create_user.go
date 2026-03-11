package presenter

import (
	"be-dashboard-nba/repository/user"
	"database/sql"
)

type CreateUserRequest struct {
	Name     string  `json:"name" validate:"required"`
	FullName string  `json:"full_name" validate:"required"`
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=8"`
	RoleID   int     `json:"role_id" validate:"required"`
	Phone    *string `json:"phone"`
	Active   *bool   `json:"active"  validate:"required"`
}

func (req *CreateUserRequest) ToParams(userID string, password string) (params repository.CreateUserParams) {

	params = repository.CreateUserParams{
		Name:      req.Name,
		FullName:  req.FullName,
		Email:     req.Email,
		Password:  password,
		CreatedBy: userID,
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

	return
}
