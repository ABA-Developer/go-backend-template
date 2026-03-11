package usecase

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/menu_permission/domain"
	"context"

	"github.com/pkg/errors"
)

func (s *menuPermissionUsecase) CreateMenuPermissionUsecase(
	ctx context.Context,
	payload domain.MenuPermissionCreatePayload,
) (err error) {

	_, err = s.menuRepo.ReadMenuByIDQuery(ctx, payload.MenuID)
	if err != nil {
		if errors.Is(err, constant.ErrDataNotFound) {
			s.log.Warn().Int("menu_id", payload.MenuID).Msg("menu detail not found for permission creation")
			return constant.ErrMenuIdNotFound
		}

		s.log.Error().Err(err).Int("menu_id", payload.MenuID).Msg("error reading menu detail query")
		return errors.WithStack(constant.ErrUnknownSource)
	}

	err = s.menuPermissionRepo.CreateMenuPermissionQuery(ctx, payload)
	if err != nil {
		s.log.Error().Err(err).Interface("payload", payload).Msg("error to create menu permission")
		return errors.WithStack(constant.ErrUnknownSource)
	}

	return nil
}
