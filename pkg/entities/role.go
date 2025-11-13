package entities

import (
	"be-dashboard-nba/api/presenter"
	"database/sql"
	"time"
)

type Role struct {
	ID          int
	Code        string
	Name        string
	Description sql.NullString
	CreatedBy   sql.NullTime
	CreatedAt   time.Time
	UpdatedBy   sql.NullString
	UpdatedAt   sql.NullTime
}

type RoleAccessResponse struct {
	RoleID         int
	RoleName       string
	MenuID         int
	MenuName       string
	PermissionID   int
	PermissionName string
	PermissionCode string
	HasAccess      bool
}

type RolePaginationResponse struct {
	Data       []Role
	Pagination presenter.Pagination
}
