package service

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/entities"
	"be-dashboard-nba/pkg/user/repository"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *service) ReadDetailUserService(
	ctx context.Context,
	id string,
) (data entities.User, err error) {
	r := repository.NewRepository(s.db)

	data, err = r.ReadDetailUserQuery(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warn().Str("id", id).Msg("user detail not found")
			return
		}
		s.log.Error().Err(err).Str("id", id).Msg("error reading user detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	return
}
