package domain

import (
	"context"
	"time"
)

// Role defines the domain entity for application roles.
// DB Nullable fields are represented as pointers.
type Role struct {
	ID          int
	Code        string
	Name        string
	Description *string
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedBy   *string
	UpdatedAt   *time.Time
}

// RoleAccessResponse is a read model for fetching role access.
type RoleAccessResponse struct {
	RoleID         int    `json:"role_id"`
	RoleName       string `json:"role_name"`
	MenuID         int    `json:"menu_id"`
	MenuName       string `json:"menu_name"`
	PermissionID   int    `json:"permission_id"`
	PermissionName string `json:"permission_name"`
	PermissionCode string `json:"permission_code"`
	HasAccess      bool   `json:"has_access"`
}

// RoleAccessPaginationResponse is for paginated listing of Role Accesses.
type RoleAccessPaginationResponse struct {
	Data       []RoleAccessResponse
	Pagination interface{} // Placeholder for presenter.Pagination, will fix mapping later
}

// Payloads for Create / Update
type CreateRolePayload struct {
	Code        string
	Name        string
	Description *string
	CreatedBy   string
}

type UpdateRolePayload struct {
	RoleID      int
	Code        string
	Name        string
	Description *string
	UpdatedBy   *string
}

type UpdateRoleMenuPermission struct {
	MenuPermissionID int
	RoleID           int
	HasAccess        *bool
}

// Filters for Listing
type RoleFilter struct {
	SetSearch bool
	Search    string
	Order     string
	Limit     int
	Offset    int
	Page      int
}

type RoleAccessFilter struct {
	RoleID    int
	SetSearch bool
	Search    string
	Order     string
	Limit     int
	Offset    int
	Page      int
}

// RoleRepository defines database operations for the Role entity.
type RoleRepository interface {
	ReadRolesQuery(ctx context.Context, filter RoleFilter) (data []Role, err error)
	ReadRolesCount(ctx context.Context, filter RoleFilter) (count int, err error)
	ReadRoleByIDQuery(ctx context.Context, roleID int) (data Role, err error)
	CreateRoleQuery(ctx context.Context, payload CreateRolePayload) (err error)
	UpdateRoleQuery(ctx context.Context, payload UpdateRolePayload) (err error)
	DeleteRoleQuery(ctx context.Context, roleID int) (err error)
	ReadRoleAccessQuery(ctx context.Context, filter RoleAccessFilter) (data []RoleAccessResponse, err error)
	ReadRoleAccessCount(ctx context.Context, filter RoleAccessFilter) (count int, err error)
	UpdateRoleAccessTx(ctx context.Context, roleID int, payloads []UpdateRoleMenuPermission) (err error)
}

// RoleUsecase defines the business logic operations for the Role entity.
type RoleUsecase interface {
	CreateRoleUsecase(ctx context.Context, payload CreateRolePayload) (err error)
	UpdateRoleUsecase(ctx context.Context, payload UpdateRolePayload) (err error)
	DeleteRoleUsecase(ctx context.Context, roleID int) (err error)
	ReadRoleUsecase(ctx context.Context, filter RoleFilter) (data []Role, count int, err error)
	ReadDetailRoleUsecase(ctx context.Context, roleID int) (data Role, err error)
	ReadRoleAccessUsecase(ctx context.Context, filter RoleAccessFilter) (data []RoleAccessResponse, count int, err error)
	UpdateRoleAccessUsecase(ctx context.Context, roleID int, requestPayloads []UpdateRoleMenuPermission) (err error)
}
