package role

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/mapper"
	rolePresenter "be-dashboard-nba/internal/presentation/presenter/role"
	"context"

	"github.com/pkg/errors"
)

func (s *useCase) CreateRoleUseCase(ctx context.Context, paylaod rolePresenter.CreateRoleRequest, userID string) (err error) {
	r := s.newRoleRepo(s.db)

	err = r.CreateRoleQuery(ctx, mapper.ToCreateRoleParams(&paylaod, userID))
	if err != nil {
		s.log.Error().Err(err).Interface("request_payload", paylaod).Msg("error to create role")
		err = errors.WithStack(constant.ErrUnknownSource)
	}
	return err
}
