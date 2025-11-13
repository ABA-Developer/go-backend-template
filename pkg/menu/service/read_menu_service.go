package service

import (
	menuPresenter "be-dashboard-nba/api/presenter/menu"
	"be-dashboard-nba/pkg/entities"
	"be-dashboard-nba/pkg/menu/repository"
	"context"
)

func (s *service) ReadListMenuService(
	ctx context.Context,
	request menuPresenter.ReadMenuListRequest,
) (data []entities.Menu, err error) {
	r := repository.NewRepository(s.db)
	params := request.ToParams()
	data, err = r.ReadListMenuQuery(ctx, params)
	if err != nil {
		s.log.Error().Err(err).Msg("error query read list menu")
		return
	}
	return
}
