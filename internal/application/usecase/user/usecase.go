package user

import (
	contractRepo "be-dashboard-nba/internal/domain/contract/repository"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	db "be-dashboard-nba/internal/infrastructure/database"
	roleRepo "be-dashboard-nba/internal/infrastructure/repository/role"
	userRepo "be-dashboard-nba/internal/infrastructure/repository/user"

	"github.com/rs/zerolog"
)

type useCase struct {
	db    db.DB
	log   *zerolog.Logger

	newUserRepo func(q db.Query) contractRepo.UserRepository
	newRoleRepo func(q db.Query) contractRepo.RoleRepository
}

func NewUseCase(db db.DB, log *zerolog.Logger) contract.UserUseCase {
	return &useCase{
		db:  db,
		log: log,

		newUserRepo: userRepo.NewRepository,
		newRoleRepo: roleRepo.NewRepository,
	}
}
