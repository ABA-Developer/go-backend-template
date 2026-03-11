package usecase

import (
	"context"

	"be-dashboard-nba/internal/modules/menu/domain"
)

func (s *menuUsecase) ReadMenuParentUsecase(ctx context.Context) (data []domain.MenuParent, err error) {
	data, err = s.menuRepo.ReadParentMenuQuery(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("error query read  menu parent")
		return
	}
	return
}
