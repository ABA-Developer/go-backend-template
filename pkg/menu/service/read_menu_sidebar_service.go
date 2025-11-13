package service

import (
	"be-dashboard-nba/pkg/entities"
	"be-dashboard-nba/pkg/menu/repository"
	"context"
)

func (s *service) ReadSidebarMenuService(
	ctx context.Context,
	userID string,
) (data []entities.Menu, err error) {
	r := repository.NewRepository(s.db)
	data, err = r.ReadSidebarMenuQuery(ctx, userID)
	if err != nil {
		s.log.Error().Err(err).Msg("error query read list menu")
		return
	}
	return
}
