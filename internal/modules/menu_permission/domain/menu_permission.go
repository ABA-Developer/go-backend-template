package domain

import (
	"be-dashboard-nba/api/presenter"
	"context"
	"time"
)

type MenuPermission struct {
	ID         int        `json:"id"`
	MenuID     int        `json:"menu_id"`
	Code       string     `json:"code"`
	ActionName string     `json:"action_name"`
	CreatedBy  string     `json:"created_by"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedBy  *string    `json:"updated_by"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type MenuPermissionFilter struct {
	SetSearch bool
	Search    string
	Order     string
	Limit     int
	Offset    int
	Page      int
	MenuID    int
}

type MenuPermissionCreatePayload struct {
	MenuID     int
	Code       string
	ActionName string
	CreatedBy  string
}

type MenuPermissionUpdatePayload struct {
	ID         int
	Code       string
	ActionName string
	UpdatedBy  string
}

type MenuPermissionDetail struct {
	ID         int    `json:"id"`
	MenuID     int    `json:"menu_id"`
	Code       string `json:"code"`
	ActionName string `json:"action_name"`
}

type MenuPermissionPaginationResponse struct {
	Data       []MenuPermissionDetail
	Pagination presenter.Pagination
}

// =======================
// Interfaces
// =======================

type MenuPermissionRepository interface {
	// Queries (Read)
	ReadMenuPermissionListQuery(ctx context.Context, params MenuPermissionFilter) (data []MenuPermissionDetail, err error)
	ReadMenuPermissionCountQuery(ctx context.Context, params MenuPermissionFilter) (count int, err error)
	ReadMenuPermissionByIDQuery(ctx context.Context, id int) (data MenuPermissionDetail, err error)

	// Commands (Write)
	CreateMenuPermissionQuery(ctx context.Context, payload MenuPermissionCreatePayload) (err error)
	UpdateMenuPermissionQuery(ctx context.Context, payload MenuPermissionUpdatePayload) (err error)
	DeleteMenuPermissionQuery(ctx context.Context, id int) (err error)
}

type MenuPermissionUsecase interface {
	CreateMenuPermissionUsecase(ctx context.Context, payload MenuPermissionCreatePayload) (err error)
	UpdateMenuPermissionUsecase(ctx context.Context, payload MenuPermissionUpdatePayload) (err error)
	DeleteMenuPermissionUsecase(ctx context.Context, id int) (err error)

	ReadListMenuPermissionUsecase(ctx context.Context, filter MenuPermissionFilter) (data MenuPermissionPaginationResponse, err error)
	ReadMenuPermissionDetailUsecase(ctx context.Context, id int) (data MenuPermissionDetail, err error)
}
