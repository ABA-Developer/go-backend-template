package service

import (
	"be-dashboard-nba/pkg/entities"
	"be-dashboard-nba/pkg/menu/repository"
	"context"
)

func (s *service) ReadMenuParentService(ctx context.Context) (data []entities.Menu, err error) {
	r := repository.NewRepository(s.db)
	data, err = r.ReadParentMenuQuery(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("error query read  menu parent")
		return
	}
	return
}
