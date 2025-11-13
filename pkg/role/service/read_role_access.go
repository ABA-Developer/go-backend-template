package service

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/entities"
	"be-dashboard-nba/pkg/role/repository"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *service) ReadRoleAccessService(ctx context.Context, roleID int) (data []entities.RoleAccessResponse, err error) {
	r := repository.NewRepository(s.db)

	_, err = r.ReadRoleByIDQuery(ctx, roleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warn().Int("Id", roleID).Msg("role detail not found")
			err = constant.ErrRoleIdNotFound
			return
		}
		s.log.Error().Err(err).Int("id", roleID).Msg("error reading role detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	data, err = r.ReadRoleAccess(ctx, roleID)
	if err != nil {
		s.log.Error().Err(err).Int("id", roleID).Msg("error reading role detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	return
}
