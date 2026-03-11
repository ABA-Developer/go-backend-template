package usecase

import (
	menuDomain "be-dashboard-nba/internal/modules/menu/domain"
	"be-dashboard-nba/internal/modules/menu_permission/domain"

	"github.com/rs/zerolog"
)

type menuPermissionUsecase struct {
	menuPermissionRepo domain.MenuPermissionRepository
	menuRepo           menuDomain.MenuRepository
	log                *zerolog.Logger
}

func NewMenuPermissionUsecase(
	mpr domain.MenuPermissionRepository,
	mr menuDomain.MenuRepository,
	log *zerolog.Logger,
) domain.MenuPermissionUsecase {
	return &menuPermissionUsecase{
		menuPermissionRepo: mpr,
		menuRepo:           mr,
		log:                log,
	}
}
