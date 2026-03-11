package usecase

import (
	"context"

	"be-dashboard-nba/constant"

	"github.com/pkg/errors"
)

func (s *roleUsecase) DeleteRoleUsecase(ctx context.Context, roleID int) (err error) {
	_, err = s.roleRepo.ReadRoleByIDQuery(ctx, roleID)
	if err != nil {
		if errors.Is(err, constant.ErrDataNotFound) {
			s.log.Warn().Int("id", roleID).Msg("role detail not found")
			err = constant.ErrRoleIdNotFound
			return
		}
		s.log.Error().Err(err).Int("id", roleID).Msg("error reading role detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	err = s.roleRepo.DeleteRoleQuery(ctx, roleID)
	if err != nil {
		s.log.Error().Err(err).Int("id", roleID).Msg("error to delete role")
		err = errors.WithStack(constant.ErrUnknownSource)
	}
	return
}
