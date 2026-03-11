package usecase

import (
	"be-dashboard-nba/internal/modules/user/domain"

	"github.com/rs/zerolog"
)

type userUsecase struct {
	userRepo domain.UserRepository
	log      *zerolog.Logger
}

func NewUserUsecase(userRepo domain.UserRepository, log *zerolog.Logger) domain.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		log:      log,
	}
}
