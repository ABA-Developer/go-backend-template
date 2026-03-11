package usecase

import (
	"context"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/user/domain"

	"github.com/pkg/errors"
)

func (s *userUsecase) ReadDetailUserUsecase(
	ctx context.Context,
	id string,
) (data domain.UserDetailRow, err error) {
	data, err = s.userRepo.ReadUserByIDQuery(ctx, id)
	if err != nil {
		if errors.Is(err, constant.ErrDataNotFound) {
			s.log.Warn().Str("id", id).Msg("user detail not found")
			err = constant.ErrUserIdNotFound
			return
		}
		s.log.Error().Err(err).Str("id", id).Msg("error reading user detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	return
}
