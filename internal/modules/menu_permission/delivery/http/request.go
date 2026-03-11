package http

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/core/utils"
	"be-dashboard-nba/internal/modules/menu_permission/domain"
)

type CreateMenuPermissionRequest struct {
	Code       string `json:"code" validate:"required,min=1,max=50"`
	ActionName string `json:"action_name" validate:"required,min=1,max=50"`
}

type UpdateMenuPermissionRequest struct {
	Code       string `json:"code" validate:"min=1,max=50"`
	ActionName string `json:"action_name" validate:"min=1,max=50"`
}

func (req *CreateMenuPermissionRequest) ToDomainCreate(userID string, menuID int) (payload domain.MenuPermissionCreatePayload) {
	payload = domain.MenuPermissionCreatePayload{
		Code:       req.Code,
		ActionName: req.ActionName,
		MenuID:     menuID,
		CreatedBy:  userID,
	}

	return
}

func (req *UpdateMenuPermissionRequest) ToDomainUpdate(userID string, menuPermissionID int) (payload domain.MenuPermissionUpdatePayload) {
	payload = domain.MenuPermissionUpdatePayload{
		ID:         menuPermissionID,
		Code:       req.Code,
		ActionName: req.ActionName,
		UpdatedBy:  userID,
	}

	return
}

type ReadMenuPermissionRequest struct {
	utils.PaginationPayload
}

func (req *ReadMenuPermissionRequest) ToDomainFilter(menuID int) (filter domain.MenuPermissionFilter) {
	req.Init()
	if req.Limit <= 0 {
		req.Limit = constant.DefaultLimit
	}

	if req.Page <= 0 {
		req.Page = constant.DefaultPage
	}
	filter = domain.MenuPermissionFilter{
		SetSearch: req.SetSearch,
		Search:    req.Search,
		Order:     req.Order,
		Limit:     req.Limit,
		Offset:    req.Offset,
		Page:      req.Page,
		MenuID:    menuID,
	}

	return
}
