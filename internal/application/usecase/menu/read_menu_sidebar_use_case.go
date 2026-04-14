package menu

import (
	"be-dashboard-nba/internal/domain/model"
	"context"
)

func (s *useCase) ReadSidebarMenuUseCase(
	ctx context.Context,
	userID string,
) (data []model.Menu, err error) {
	r := s.newMenuRepo(s.db)
	data, err = r.ReadSidebarMenuQuery(ctx, userID)
	if err != nil {
		s.log.Error().Err(err).Msg("error query read list menu")
		return
	}
	return
}
