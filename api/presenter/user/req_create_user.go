package presenter

import (
	"be-dashboard-nba/pkg/user/repository"
	"database/sql"
)

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Phone    string `json:"phone"`
	Active   *bool  `json:"active"`
	ImgPath  string `json:"img_path"`
	ImgName  string `json:"img_name"`
}

func (req *CreateUserRequest) ToParams(userID string, password string) (params repository.CreateUserParams) {

	params = repository.CreateUserParams{
		Name:      req.Name,
		FullName:  req.FullName,
		Email:     req.Email,
		Password:  password,
		CreatedBy: userID,
	}

	if req.Active != nil {
		params.Active = *req.Active
	} else {
		params.Active = true
	}

	if req.Phone != "" {
		params.Phone = sql.NullString{String: req.Phone, Valid: true}
	}
	if req.ImgPath != "" {
		params.ImgPath = sql.NullString{String: req.ImgPath, Valid: true}
	}
	if req.ImgName != "" {
		params.ImgName = sql.NullString{String: req.ImgName, Valid: true}
	}

	return
}
