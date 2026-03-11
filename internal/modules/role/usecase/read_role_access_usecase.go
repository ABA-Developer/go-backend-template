package usecase

import (
	"context"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/role/domain"

	"github.com/pkg/errors"
)

func (s *roleUsecase) ReadRoleAccessUsecase(ctx context.Context, filter domain.RoleAccessFilter) (data []domain.RoleAccessResponse, count int, err error) {
	_, err = s.roleRepo.ReadRoleByIDQuery(ctx, filter.RoleID)
	if err != nil {
		if errors.Is(err, constant.ErrDataNotFound) {
			s.log.Warn().Int("id", filter.RoleID).Msg("role detail not found")
			err = constant.ErrRoleIdNotFound
			return
		}
		s.log.Error().Err(err).Int("id", filter.RoleID).Msg("error reading role detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	count, err = s.roleRepo.ReadRoleAccessCount(ctx, filter)
	if err != nil {
		s.log.Error().Err(err).Msg("error to get role access count")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	data, err = s.roleRepo.ReadRoleAccessQuery(ctx, filter)
	if err != nil {
		s.log.Error().Err(err).Msg("error to get role access query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	return
}
