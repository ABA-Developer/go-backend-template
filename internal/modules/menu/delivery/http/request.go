package http

import (
	"be-dashboard-nba/internal/core/utils"
	"be-dashboard-nba/internal/modules/menu/domain"
	"strings"
)

type CreateMenuRequest struct {
	ParentID    *int    `json:"parent_id" validate:"omitempty,min=1"`
	Name        string  `json:"name" validate:"required,min=1,max=50"`
	Description *string `json:"description" validate:"omitempty,max=100"`
	URL         *string `json:"url" validate:"required,max=100"`
	Group       string  `json:"group" validate:"required,min=1,max=50"`
	Icon        *string `json:"icon" validate:"omitempty,max=50"`
	Active      *bool   `json:"active" validate:"required,boolean"`
	Display     *bool   `json:"display" validate:"required,boolean"`
}

type UpdateMenuRequest struct {
	ParentID    *int    `json:"parent_id" validate:"omitempty,min=1"`
	Name        string  `json:"name" validate:"required,min=1,max=50"`
	Description *string `json:"description" validate:"omitempty,max=100"`
	URL         *string `json:"url" validate:"omitempty,max=100,uri"`
	Group       string  `json:"group" validate:"required,min=1,max=50"`
	Icon        *string `json:"icon" validate:"omitempty,max=50"`
	Active      *bool   `json:"active" validate:"required,boolean"`
	Display     *bool   `json:"display" validate:"required,boolean"`
}

func (req *CreateMenuRequest) ToDomainCreateMenu(userID string) domain.MenuCreatePayload {
	params := domain.MenuCreatePayload{
		Name:      req.Name,
		CreatedBy: userID,
		Group:     req.Group,
	}

	if req.Active != nil {
		params.Active = *req.Active
	}
	if req.Display != nil {
		params.Display = *req.Display
	}

	if req.ParentID != nil {
		val := int32(*req.ParentID)
		params.ParentID = &val
	}

	params.Description = req.Description
	params.Icon = req.Icon

	if req.URL != nil {
		urlToStore := *req.URL
		if !strings.HasPrefix(urlToStore, "/") {
			urlToStore = "/" + urlToStore
		}
		params.URL = &urlToStore
	}

	return params
}

type ReadMenuRequest struct {
	utils.PaginationPayload
}

func (req *ReadMenuRequest) ToDomainFilter() (filter domain.MenuFilter) {
	req.Init()

	filter = domain.MenuFilter{
		SetSearch: req.SetSearch,
		Search:    req.Search,
	}

	return
}

func (req *UpdateMenuRequest) ToDomainUpdatePayload(userID string, menuID int) (payload domain.MenuUpdatePayload) {
	params := domain.MenuUpdatePayload{
		Name:      req.Name,
		Group:     req.Group,
		UpdatedBy: &userID,
		ID:        menuID,
	}

	if req.Active != nil {
		params.Active = *req.Active
	}
	if req.Display != nil {
		params.Display = *req.Display
	}

	if req.ParentID != nil {
		val := int32(*req.ParentID)
		params.ParentID = &val
	}

	params.Description = req.Description
	params.Icon = req.Icon

	if req.URL != nil {
		urlToStore := *req.URL
		if !strings.HasPrefix(urlToStore, "/") {
			urlToStore = "/" + urlToStore
		}
		params.URL = &urlToStore
	}

	return params

}

type UpdateMenuOrderRequest struct {
	Group     string `json:"group" validate:"required"`
	ParentID  *int   `json:"parent_id" validate:"omitempty,min=1"`
	SortedIDs []int  `json:"sorted_ids" validate:"required,min=1,dive,min=1"`
}

func (req *UpdateMenuOrderRequest) ToDomainUpdateOrder(userID string) domain.MenuUpdateSortPayload {
	return domain.MenuUpdateSortPayload{
		Group:     req.Group,
		ParentID:  req.ParentID,
		SortedIDs: req.SortedIDs,
		UpdatedBy: userID,
	}
}
