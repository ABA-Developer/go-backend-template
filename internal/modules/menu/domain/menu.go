package domain

import (
	"context"
	"time"
)

// Menu represents the core entity for application menus without SQL driver leakages.
type Menu struct {
	ID          int        `json:"id"`
	ParentID    *int32     `json:"parent_id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	URL         *string    `json:"url"`
	Sort        int        `json:"sort"`
	Group       string     `json:"group"`
	Icon        *string    `json:"icon"`
	Active      bool       `json:"active"`
	Display     bool       `json:"display"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedBy   *string    `json:"updated_by"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

// =======================
// Payload / Filter Structs
// =======================

type MenuFilter struct {
	SetSearch bool
	Search    string
}

type MenuCreatePayload struct {
	ParentID    *int32
	Name        string
	Description *string
	URL         *string
	Sort        int
	Group       string
	Icon        *string
	Active      bool
	Display     bool
	CreatedBy   string
}

type MenuUpdatePayload struct {
	ID          int
	ParentID    *int32
	Name        string
	Description *string
	URL         *string
	Sort        int
	Group       string
	Icon        *string
	Active      bool
	Display     bool
	UpdatedBy   *string
}

type MenuUpdateSortItemPayload struct {
	ID        int
	Sort      int
	UpdatedBy string
	ParentID  *int32
	Group     string
}

type MenuUpdateSortPayload struct {
	Group     string
	ParentID  *int
	SortedIDs []int
	UpdatedBy string
}

// =======================
// Read Models (Views)
// =======================

type MenuListItem struct {
	ID       int            `json:"id"`
	Name     string         `json:"name"`
	Icon     *string        `json:"icon"`
	Url      *string        `json:"url"`
	Children []MenuListItem `json:"children"`
	Sort     int            `json:"sort"`
}

type MenuDetail struct {
	ID          int     `json:"id"`
	ParentID    *int32  `json:"parent_id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	URL         *string `json:"url"`
	Sort        int     `json:"sort"`
	Group       string  `json:"group"`
	Icon        *string `json:"icon"`
	Active      bool    `json:"active"`
	Display     bool    `json:"display"`
}

type MenuParent struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Group string `json:"group"`
}

type MenuRepository interface {
	// =======================
	// Queries (Read)
	// =======================
	ReadSidebarMenuQuery(ctx context.Context, userID string) (data []Menu, err error)
	ReadListMenuQuery(ctx context.Context, params MenuFilter) (data []Menu, err error)
	ReadCountMenuQuery(ctx context.Context, params MenuFilter) (count int64, err error)
	ReadParentMenuQuery(ctx context.Context) (data []MenuParent, err error)
	ReadMenuByIDQuery(ctx context.Context, menuID int) (data Menu, err error)
	CountMenuChildren(ctx context.Context, menuID int) (int, error)

	// =======================
	// Commands (Write / Tx)
	// =======================
	DeleteMenuQuery(ctx context.Context, menuID int) (err error)

	CreateMenuQuery(ctx context.Context, params MenuCreatePayload, userID string) (err error)
	UpdateMenuQuery(ctx context.Context, params MenuUpdatePayload, updateChildrenGroup bool) (err error)
	UpdateMenuOrderQuery(ctx context.Context, payloads []MenuUpdateSortItemPayload) error
}

type MenuUsecase interface {
	CreateMenuUsecase(ctx context.Context, payload MenuCreatePayload, userID string) (err error)
	UpdateMenuUsecase(ctx context.Context, payload MenuUpdatePayload) (err error)
	ReadListMenuUsecase(ctx context.Context, filter MenuFilter) (data []Menu, err error)
	ReadSidebarMenuUsecase(ctx context.Context, userID string) (data []Menu, err error)
	DeleteMenuUsecase(ctx context.Context, menuID int) (err error)
	UpdateMenuOrderUsecase(ctx context.Context, payload MenuUpdateSortPayload) (err error)
	ReadMenuDetailUsecase(ctx context.Context, menuID int) (data MenuDetail, err error)
	ReadMenuParentUsecase(ctx context.Context) (data []MenuParent, err error)
}
