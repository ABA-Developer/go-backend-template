package usecase

import (
	"context"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/core/utils"
	"be-dashboard-nba/internal/modules/user/domain"

	"github.com/pkg/errors"
)

func (s *userUsecase) CreateUserUsecase(ctx context.Context, userDomain domain.User) (err error) {
	hashedPassword, err := utils.GenerateHashPassword(userDomain.Password)
	if err != nil {
		s.log.Error().Err(err).Msg("error hash password")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}
	userDomain.Password = hashedPassword

	err = s.userRepo.CreateUserWithRoleTx(ctx, userDomain)
	if err != nil {
		if errors.Is(err, constant.ErrEmailAlreadyExists) || errors.Is(err, constant.ErrRoleIdNotFound) {
			s.log.Warn().Err(err).Msg("validation error on create user")
			return err
		}

		s.log.Error().Err(err).Interface("request_payload", userDomain).Msg("error to create user via transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	return nil
}
