package menu

import (
	"be-dashboard-nba/usecase/entities"
	"be-dashboard-nba/repository/menu"
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
