package usecase

import (
	"context"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/menu_permission/domain"

	"github.com/pkg/errors"
)

func (s *menuPermissionUsecase) UpdateMenuPermissionUsecase(
	ctx context.Context,
	payload domain.MenuPermissionUpdatePayload,
) (err error) {

	_, err = s.menuPermissionRepo.ReadMenuPermissionByIDQuery(ctx, payload.ID)
	if err != nil {
		if errors.Is(err, constant.ErrDataNotFound) {
			s.log.Warn().Int("id", payload.ID).Msg("menu permission detail not found for update")
			return constant.ErrMenuPermissionIdNotFound
		}

		s.log.Error().Err(err).Int("id", payload.ID).Msg("error reading menu permission detail query")
		return errors.WithStack(constant.ErrUnknownSource)
	}

	err = s.menuPermissionRepo.UpdateMenuPermissionQuery(ctx, payload)
	if err != nil {
		s.log.Error().Err(err).Interface("payload", payload).Msg("error to update menu permission")
		return errors.WithStack(constant.ErrUnknownSource)
	}

	return nil
}
