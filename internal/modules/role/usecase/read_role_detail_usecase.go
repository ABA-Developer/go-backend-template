package usecase

import (
	"context"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/role/domain"

	"github.com/pkg/errors"
)

func (s *roleUsecase) ReadDetailRoleUsecase(ctx context.Context, roleID int) (data domain.Role, err error) {
	data, err = s.roleRepo.ReadRoleByIDQuery(ctx, roleID)
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

	return
}
