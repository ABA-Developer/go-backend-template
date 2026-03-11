package usecase

import (
	"context"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/role/domain"

	"github.com/pkg/errors"
)

func (s *roleUsecase) UpdateRoleUsecase(ctx context.Context, payload domain.UpdateRolePayload) (err error) {
	_, err = s.roleRepo.ReadRoleByIDQuery(ctx, payload.RoleID)
	if err != nil {
		if errors.Is(err, constant.ErrDataNotFound) {
			s.log.Warn().Int("id", payload.RoleID).Msg("role detail not found")
			err = constant.ErrRoleIdNotFound
			return
		}
		s.log.Error().Err(err).Int("id", payload.RoleID).Msg("error reading role detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	err = s.roleRepo.UpdateRoleQuery(ctx, payload)
	if err != nil {
		s.log.Error().Err(err).Interface("request_payload", payload).Msg("error to update role")
		err = errors.WithStack(constant.ErrUnknownSource)
	}
	return
}
