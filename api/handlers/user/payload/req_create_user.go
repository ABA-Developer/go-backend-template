package payload

import "be-dashboard-nba/pkg/user/repository"

type CreateUserRequest struct {
	FirstName  string `json:"first_name" validate:"required"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8"`
	Role       string `json:"role" validate:"required"`
	IsActive   *bool  `json:"is_active"`
}

func (req *CreateUserRequest) ToParams(userID int64, password string) (params repository.CreateUserParams) {
	params = repository.CreateUserParams{
		FirstName:  req.FirstName,
		MiddleName: req.MiddleName,
		LastName:   req.LastName,
		Email:      req.Email,
		Password:   password,
		Role:       req.Role,
		IsActive:   *req.IsActive,
		CreatedBy:  userID,
	}

	return
}
