package menu

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/domain/model"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *useCase) ReadMenuDetailUseCase(ctx context.Context, menuID int) (data model.Menu, err error) {

	r := s.newMenuRepo(s.db)
	data, err = r.ReadMenuByIDQuery(ctx, menuID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warn().Int("id", menuID).Msg("menu detail not found")
			err = constant.ErrMenuIdNotFound
			return
		}
		s.log.Error().Err(err).Int("id", menuID).Msg("error reading menu detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}
	return
}
