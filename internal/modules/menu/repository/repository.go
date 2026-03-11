package repository

import (
	"be-dashboard-nba/internal/core/db"
	"be-dashboard-nba/internal/modules/menu/domain"
)

type menuRepository struct {
	db db.DB
}

func NewMenuRepository(db db.DB) domain.MenuRepository {
	return &menuRepository{
		db: db,
	}
}
