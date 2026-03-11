package usecase

import (
	"context"

	"be-dashboard-nba/internal/modules/menu/domain"
)

func (s *menuUsecase) ReadListMenuUsecase(
	ctx context.Context,
	params domain.MenuFilter,
) (data []domain.Menu, err error) {
	data, err = s.menuRepo.ReadListMenuQuery(ctx, params)
	if err != nil {
		s.log.Error().Err(err).Msg("error query read list menu")
		return
	}
	return
}
