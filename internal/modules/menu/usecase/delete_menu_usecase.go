package usecase

import (
	"context"

	"be-dashboard-nba/constant"

	"github.com/pkg/errors"
)

func (s *menuUsecase) DeleteMenuUsecase(ctx context.Context, menuID int) error {

	err := s.menuRepo.DeleteMenuQuery(ctx, menuID)

	if err != nil {
		if errors.Is(err, constant.ErrDataNotFound) {
			s.log.Warn().Int("id", menuID).Msg("menu detail not found for delete")
			return constant.ErrMenuIdNotFound
		}

		s.log.Error().Err(err).Int("id", menuID).Msg("error delete menu")
		return errors.WithStack(constant.ErrUnknownSource)
	}

	return nil
}
