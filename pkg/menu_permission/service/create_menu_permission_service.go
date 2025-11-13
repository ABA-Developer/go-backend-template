package service

import (
	menuPermissionPresenter "be-dashboard-nba/api/presenter/menu_permission"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/menu_permission/repository"
	"context"

	"github.com/pkg/errors"
)

func (s *service) CreateMenuPermissionService(ctx context.Context, payload menuPermissionPresenter.CreateMenuPermissionRequest, userID string, menuID int) (err error) {
	r := repository.NewRepository(s.db)

	err = r.CreateMenuPermissionQuery(ctx, payload.ToParams(userID, menuID))
	if err != nil {
		s.log.Error().Err(err).Interface("request_payload", payload).Msg("error to create menu permission")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}
	return
}
