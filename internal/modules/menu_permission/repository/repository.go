package repository

import (
	"be-dashboard-nba/internal/core/db"
	"be-dashboard-nba/internal/modules/menu_permission/domain"
)

type menuPermissionRepository struct {
	db db.DB
}

func NewMenuRepository(db db.DB) domain.MenuPermissionRepository {
	return &menuPermissionRepository{
		db: db,
	}
}
