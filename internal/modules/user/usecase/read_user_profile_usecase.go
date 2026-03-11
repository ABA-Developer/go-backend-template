package usecase

import (
	"context"
	"errors"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/user/domain"
)

func (s *userUsecase) ReadUserProfileUsecase(ctx context.Context, userID string) (data domain.UserProfileRow, err error) {
	data, err = s.userRepo.ReadUserProfileQuery(ctx, userID)
	if err != nil {
		if errors.Is(err, constant.ErrDataNotFound) {
			s.log.Warn().Str("id", userID).Msg("user profile not found")
			err = constant.ErrUserIdNotFound
			return
		}
		s.log.Err(err).Err(err).Str("id", userID).Msg("error reading query user pprofile")
		err = constant.ErrUnknownSource
		return
	}

	return
}
