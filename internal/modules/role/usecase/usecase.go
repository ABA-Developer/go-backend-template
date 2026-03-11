package usecase

import (
	"be-dashboard-nba/internal/modules/role/domain"

	"github.com/rs/zerolog"
)

type roleUsecase struct {
	roleRepo domain.RoleRepository
	log      *zerolog.Logger
}

func NewRoleUsecase(roleRepo domain.RoleRepository, log *zerolog.Logger) domain.RoleUsecase {
	return &roleUsecase{
		roleRepo: roleRepo,
		log:      log,
	}
}
