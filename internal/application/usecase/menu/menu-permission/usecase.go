package menu_permission

import (
	contractRepo "be-dashboard-nba/internal/domain/contract/repository"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	db "be-dashboard-nba/internal/infrastructure/database"
	menuRepo "be-dashboard-nba/internal/infrastructure/repository/menu"
	menuPermissionRepo "be-dashboard-nba/internal/infrastructure/repository/menu/menu-permission"

	"github.com/rs/zerolog"
)

type useCase struct {
	db                    db.DB
	log                   *zerolog.Logger
	newMenuRepo           func(q db.Query) contractRepo.MenuRepository
	newMenuPermissionRepo func(q db.Query) contractRepo.MenuPermissionRepository
}

func NewUseCase(db db.DB, log *zerolog.Logger) contract.MenuPermissionUseCase {
	return &useCase{
		db:  db,
		log: log,

		newMenuRepo:           menuRepo.NewRepository,
		newMenuPermissionRepo: menuPermissionRepo.NewRepository,
	}
}
