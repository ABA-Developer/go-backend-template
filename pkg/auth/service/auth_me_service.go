package service

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/auth/repository"
	"be-dashboard-nba/pkg/entities"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *service) AuthMeService(ctx context.Context, id string) (data entities.User, err error) {
	r := repository.NewRepo(s.db)
	data, err = r.ReadDetailUserByIdQuery(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warn().Str("id", id).Msg("user detail not found")
			err = constant.ErrUserIdNotFound
			return
		}
		s.log.Error().Err(err).Str("id", id).Msg("error reading user detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	return
}
