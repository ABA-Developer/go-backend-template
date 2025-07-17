package payload

import "be-dashboard-nba/pkg/user/repository"

type UpdateUserRequest struct {
	ID         int64  `json:"id"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	IsActive   *bool  `json:"is_active"`
}

func (req *UpdateUserRequest) ToParams(userID int64, password string) (params repository.UpdateUserParams) {
	params = repository.UpdateUserParams{
		ID:         req.ID,
		FirstName:  req.FirstName,
		MiddleName: req.MiddleName,
		LastName:   req.LastName,
		Email:      req.Email,
		Password:   password,
		Role:       req.Role,
		IsActive:   *req.IsActive,
		UpdatedBy:  userID,
	}

	return
}
