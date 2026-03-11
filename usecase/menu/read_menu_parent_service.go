package menu

import (
	"be-dashboard-nba/usecase/entities"
	"be-dashboard-nba/repository/menu"
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
