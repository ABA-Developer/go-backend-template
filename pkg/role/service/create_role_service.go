package service

import (
	rolePresenter "be-dashboard-nba/api/presenter/role"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/role/repository"
	"context"

	"github.com/pkg/errors"
)

func (s *service) CreateRoleService(ctx context.Context, paylaod rolePresenter.CreateRoleRequest, userID string) (err error) {
	r := repository.NewRepository(s.db)

	err = r.CreateRoleQuery(ctx, paylaod.ToParams(userID))
	if err != nil {
		s.log.Error().Err(err).Interface("request_payload", paylaod).Msg("error to create role")
		err = errors.WithStack(constant.ErrUnknownSource)
	}
	return err
}
