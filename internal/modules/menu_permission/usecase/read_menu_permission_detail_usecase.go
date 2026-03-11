package usecase

import (
	"context"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/menu_permission/domain"

	"github.com/pkg/errors"
)

func (s *menuPermissionUsecase) ReadMenuPermissionDetailUsecase(
	ctx context.Context,
	id int,
) (data domain.MenuPermissionDetail, err error) {

	data, err = s.menuPermissionRepo.ReadMenuPermissionByIDQuery(ctx, id)

	if err != nil {
		if errors.Is(err, constant.ErrDataNotFound) {
			s.log.Warn().Int("id", id).Msg("menu permission detail not found")
			err = constant.ErrMenuPermissionIdNotFound
			return
		}

		s.log.Error().Err(err).Int("id", id).Msg("error reading menu permission detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	return data, nil
}
