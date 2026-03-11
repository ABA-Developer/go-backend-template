package usecase

import (
	"context"

	"be-dashboard-nba/constant"

	"github.com/pkg/errors"
)

func (s *menuPermissionUsecase) DeleteMenuPermissionUsecase(ctx context.Context, id int) (err error) {
	_, err = s.menuPermissionRepo.ReadMenuPermissionByIDQuery(ctx, id)
	if err != nil {
		if errors.Is(err, constant.ErrDataNotFound) {
			s.log.Warn().Int("id", id).Msg("menu permission not found for delete")
			return constant.ErrMenuPermissionIdNotFound
		}

		s.log.Error().Err(err).Int("id", id).Msg("error reading menu permission detail query")
		return errors.WithStack(constant.ErrUnknownSource)
	}

	err = s.menuPermissionRepo.DeleteMenuPermissionQuery(ctx, id)
	if err != nil {
		s.log.Error().Err(err).Int("id", id).Msg("error to delete menu permission")
		return errors.WithStack(constant.ErrUnknownSource)
	}

	return nil
}
