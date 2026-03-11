package usecase

import (
	"be-dashboard-nba/internal/modules/menu/domain"

	"github.com/rs/zerolog"
)

type menuUsecase struct {
	menuRepo domain.MenuRepository
	log      *zerolog.Logger
}

func NewMenuUsecase(menuRepo domain.MenuRepository, log *zerolog.Logger) domain.MenuUsecase {
	return &menuUsecase{
		menuRepo: menuRepo,
		log:      log,
	}
}
