package usecase

import (
	"context"
	"fmt"

	"be-dashboard-nba/constant"

	"github.com/pkg/errors"
)

func (s *userUsecase) DeleteUserUsecase(ctx context.Context, userID string, deletedBy string) (err error) {

	existingUser, err := s.userRepo.ReadUserByIDQuery(ctx, userID)
	if err != nil {
		if errors.Is(err, constant.ErrDataNotFound) {
			s.log.Warn().Str("id", userID).Msg("user detail not found")
			err = constant.ErrUserIdNotFound
			return
		}
		s.log.Error().Err(err).Str("id", userID).Msg("error reading user detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}
	fmt.Printf("user id: %s, deletedBy: %s", existingUser.ID, deletedBy)
	if existingUser.ID == deletedBy {
		err = constant.ErrForbiddenSelfDelete
		return
	}

	err = s.userRepo.DeleteUserQuery(ctx, userID)
	if err != nil {
		s.log.Error().Err(err).Str("id", userID).Msg("error to delete user")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	return

}
