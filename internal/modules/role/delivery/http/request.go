package http

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/core/utils"
	"be-dashboard-nba/internal/modules/role/domain"
)

type CreateRoleRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Code        string  `json:"code" validate:"required,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty"`
}

func (req *CreateRoleRequest) ToDomain(userID string) (payload domain.CreateRolePayload) {
	payload.Name = req.Name
	payload.Code = req.Code
	payload.CreatedBy = userID
	payload.Description = req.Description
	return
}

type UpdateRoleRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Code        string  `json:"code" validate:"required,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty"`
}

func (req *UpdateRoleRequest) ToDomain(userID string, roleID int) (payload domain.UpdateRolePayload) {
	payload.RoleID = roleID
	payload.Name = req.Name
	payload.Code = req.Code
	payload.UpdatedBy = &userID
	payload.Description = req.Description
	return
}

type ReadRolesRequest struct {
	utils.PaginationPayload
}

func (req *ReadRolesRequest) ToDomainFilter() (filter domain.RoleFilter) {
	req.Init()

	if req.Limit <= 0 {
		req.Limit = constant.DefaultLimit
	}
	if req.Page <= 0 {
		req.Page = constant.DefaultPage
	}
	filter = domain.RoleFilter{
		SetSearch: req.SetSearch,
		Search:    req.Search,
		Order:     req.Order,
		Offset:    req.Offset,
		Limit:     req.Limit,
		Page:      req.Page,
	}
	return
}

type ReadRoleAccessesRequest struct {
	utils.PaginationPayload
}

func (req *ReadRoleAccessesRequest) ToDomainFilter(roleID int) (filter domain.RoleAccessFilter) {
	req.Init()

	if req.Limit <= 0 {
		req.Limit = constant.DefaultLimit
	}
	if req.Page <= 0 {
		req.Page = constant.DefaultPage
	}
	filter = domain.RoleAccessFilter{
		RoleID:    roleID,
		SetSearch: req.SetSearch,
		Search:    req.Search,
		Order:     req.Order,
		Offset:    req.Offset,
		Limit:     req.Limit,
		Page:      req.Page,
	}
	return
}

type UpdateRoleAccessItem struct {
	AccessID  int   `json:"access_id" validate:"required"`
	HasAccess *bool `json:"has_access" validate:"required"`
}

type UpdateRoleAccessRequest struct {
	AccessItem []UpdateRoleAccessItem `json:"access_item" validate:"required,min=1,dive"`
}

func (req *UpdateRoleAccessRequest) ToDomainPayloads(roleID int) (payloads []domain.UpdateRoleMenuPermission) {
	for _, item := range req.AccessItem {
		payloads = append(payloads, domain.UpdateRoleMenuPermission{
			MenuPermissionID: item.AccessID,
			RoleID:           roleID,
			HasAccess:        item.HasAccess,
		})
	}
	return
}
