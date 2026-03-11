package usecase

import (
	"context"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/role/domain"

	"github.com/pkg/errors"
)

func (s *roleUsecase) CreateRoleUsecase(ctx context.Context, payload domain.CreateRolePayload) (err error) {
	err = s.roleRepo.CreateRoleQuery(ctx, payload)
	if err != nil {
		s.log.Error().Err(err).Interface("request_payload", payload).Msg("error to create role")
		err = errors.WithStack(constant.ErrUnknownSource)
	}
	return
}
