package service

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/entities"
	"be-dashboard-nba/pkg/menu_permission/repository"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *service) ReadMenuPermissionDetail(ctx context.Context, MenuPermissionID int) (data entities.MenuPermission, err error) {
	r := repository.NewRepository(s.db)

	data, err = r.ReadMenuPermissionByIdQuery(ctx, MenuPermissionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warn().Int("id", MenuPermissionID).Msg("menu permission detail not found")
			err = constant.ErrMenuPermissionIdNotFound
			return
		}
		s.log.Error().Err(err).Int("id", MenuPermissionID).Msg("error reading menu permission detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	return
}
