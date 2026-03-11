package usecase

import (
	"context"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/user/domain"
	"be-dashboard-nba/internal/core/utils"

	"github.com/pkg/errors"
)

func (s *userUsecase) ReadUsersUsecase(ctx context.Context, filter domain.UserFilter) (data domain.UserPaginationResponse, err error) {
	totalItems, err := s.userRepo.ReadCountUserQuery(ctx, filter)
	if err != nil {
		s.log.Error().Err(err).Msg("error query get user count")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	pagination := utils.GeneratePagination(totalItems, filter.Page, filter.Limit)

	userData, err := s.userRepo.ReadUsersQuery(ctx, filter)
	if err != nil {
		s.log.Error().Err(err).Msg("error query read roles")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	data = domain.UserPaginationResponse{
		Data:       userData,
		Pagination: pagination,
	}

	return

}
