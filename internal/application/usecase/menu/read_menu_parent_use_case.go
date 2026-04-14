package menu

import (
	"be-dashboard-nba/internal/domain/model"
	"context"
)

func (s *useCase) ReadMenuParentUseCase(ctx context.Context) (data []model.Menu, err error) {
	r := s.newMenuRepo(s.db)
	data, err = r.ReadParentMenuQuery(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("error query read  menu parent")
		return
	}
	return
}
