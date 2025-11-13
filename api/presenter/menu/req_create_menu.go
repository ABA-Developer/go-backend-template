package presenter

import (
	"be-dashboard-nba/pkg/menu/repository"
	"database/sql"
	"strings"
)

type CreateMenuRequest struct {
	ParentID    *int    `json:"parent_id" validate:"omitempty,min=1"`
	Name        string  `json:"name" validate:"required,min=1,max=50"`
	Description *string `json:"description" validate:"omitempty,max=100"`
	URL         *string `json:"url" validate:"required,max=100,uri"`
	Group       string  `json:"group" validate:"required,min=1,max=50"`
	Icon        *string `json:"icon" validate:"omitempty,max=50"`
	Active      *bool   `json:"active" validate:"required,boolean"`
	Display     *bool   `json:"display" validate:"required,boolean"`
}

func (req *CreateMenuRequest) ToParams(userID string) (params repository.CreateMenuParams) {

	params = repository.CreateMenuParams{
		Name:      req.Name,
		CreatedBy: userID,
		Group:     req.Group,
		Active:    *req.Active,
		Display:   *req.Display,
	}

	if req.ParentID != nil {
		params.ParentID = sql.NullInt32{Int32: int32(*req.ParentID), Valid: true}
	}

	if req.Description != nil {
		params.Description = sql.NullString{String: *req.Description, Valid: true}
	}

	if req.URL != nil {
		urlToStore := *req.URL
		if !strings.HasPrefix(urlToStore, "/") {
			urlToStore = "/" + urlToStore
		}
		params.URL = sql.NullString{String: urlToStore, Valid: true}
	}

	if req.Icon != nil {
		params.Icon = sql.NullString{String: *req.Icon, Valid: true}
	}

	return
}
