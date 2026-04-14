package menu

import (
	"be-dashboard-nba/internal/application/mapper"
	"be-dashboard-nba/internal/domain/model"
	menuPresenter "be-dashboard-nba/internal/presentation/presenter/menu"
	"context"
)

func (s *useCase) ReadListMenuUseCase(
	ctx context.Context,
	request menuPresenter.ReadMenuListRequest,
) (data []model.Menu, err error) {
	r := s.newMenuRepo(s.db)
	data, err = r.ReadListMenuQuery(ctx, mapper.ToReadListMenuParams(&request))
	if err != nil {
		s.log.Error().Err(err).Msg("error query read list menu")
		return
	}
	return
}
