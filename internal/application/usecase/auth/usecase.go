package auth

import (
	contractRepo "be-dashboard-nba/internal/domain/contract/repository"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	db "be-dashboard-nba/internal/infrastructure/database"
	authRepo "be-dashboard-nba/internal/infrastructure/repository/auth"

	"github.com/rs/zerolog"
)

type useCase struct {
	db    db.DB
	log   *zerolog.Logger
	newAuthRepo func(q db.Query) contractRepo.AuthRepository
}

func NewUseCase(db db.DB, log *zerolog.Logger) contract.AuthUseCase {
	return &useCase{
		db:  db,
		log: log,

		newAuthRepo: authRepo.NewRepo,
	}
}
