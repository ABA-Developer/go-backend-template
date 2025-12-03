package service

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/entities"
	"be-dashboard-nba/pkg/user/repository"
	"context"
	"database/sql"
	"errors"
)

func (s *service) ReadUserProfile(ctx context.Context, userID string) (data entities.User, err error) {

	r := repository.NewRepository(s.db)

	data, err = r.ReadUserProfileQuery(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warn().Str("id", userID).Msg("user profile not found")
			err = constant.ErrUserIdNotFound
			return
		}
		s.log.Err(err).Err(err).Str("id", userID).Msg("error reading query user pprofile")
		err = constant.ErrUnknownSource
		return
	}

	return
}
