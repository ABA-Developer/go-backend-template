package http

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/core/utils"
	"be-dashboard-nba/internal/modules/user/domain"
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

func (req *CreateUserRequest) ToDomain(userID string) (user domain.User) {
	user = domain.User{
		Name:      req.Name,
		FullName:  req.FullName,
		Email:     req.Email,
		Password:  req.Password,
		CreatedBy: userID,
		RoleID:    req.RoleID,
		Active:    true, // default
	}

	if req.Active != nil {
		user.Active = *req.Active
	}

	if req.Phone != nil {
		user.Phone = req.Phone
	}

	return
}

type ReadUserRequest struct {
	utils.PaginationPayload
}

func (req *ReadUserRequest) ToDomainFilter() (filter domain.UserFilter) {
	req.Init()

	if req.Limit <= 0 {
		req.Limit = constant.DefaultLimit
	}

	if req.Page <= 0 {
		req.Page = constant.DefaultPage
	}
	filter = domain.UserFilter{
		SetSearch: req.SetSearch,
		Search:    req.Search,
		Order:     req.Order,
		Offset:    req.Offset,
		Limit:     req.Limit,
		Page:      req.Page,
	}

	return
}

type UpdateUserRequest struct {
	Name     string  `json:"name" validate:"required"`
	FullName string  `json:"full_name"  validate:"required"`
	Email    string  `json:"email"  validate:"required,email"`
	RoleID   int     `json:"role_id" validate:"required"`
	Phone    *string `json:"phone"`
	Active   *bool   `json:"active"  validate:"required"`
}

func (req *UpdateUserRequest) ToDomainPayload(userID string, updatedBy string) (payload domain.UpdateUserPayload) {
	payload = domain.UpdateUserPayload{
		ID:        userID,
		Name:      req.Name,
		FullName:  req.FullName,
		Email:     req.Email,
		UpdatedBy: &updatedBy,
		RoleID:    req.RoleID,
		Active:    req.Active,
		Phone:     req.Phone,
	}

	return
}
