package usecase

import (
	"context"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/role/domain"

	"github.com/pkg/errors"
)

func (s *roleUsecase) ReadRoleUsecase(ctx context.Context, filter domain.RoleFilter) (data []domain.Role, count int, err error) {
	count, err = s.roleRepo.ReadRolesCount(ctx, filter)
	if err != nil {
		s.log.Error().Err(err).Msg("error to get role pagination count")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	data, err = s.roleRepo.ReadRolesQuery(ctx, filter)
	if err != nil {
		s.log.Error().Err(err).Msg("error to get role query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	return
}
