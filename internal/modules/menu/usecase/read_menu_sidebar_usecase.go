package usecase

import (
	"context"

	"be-dashboard-nba/internal/modules/menu/domain"
)

func (s *menuUsecase) ReadSidebarMenuUsecase(
	ctx context.Context,
	userID string,
) (data []domain.Menu, err error) {
	data, err = s.menuRepo.ReadSidebarMenuQuery(ctx, userID)
	if err != nil {
		s.log.Error().Err(err).Msg("error query read list menu")
		return
	}
	return
}
