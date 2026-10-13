package usecase

import (
	"github.com/ABA-Developer/go-backend-template/constant"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	"context"
	"database/sql"
	"errors"
)

func (s *useCase) ReadUserProfileUseCase(ctx context.Context, userID string) (data model.User, err error) {

	r := s.newUserRepo(s.db)

	data, err = r.ReadUserProfileQry(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log(ctx).Warn().Str("id", userID).Msg("user profile not found")
			err = constant.ErrUserIdNotFound
			return
		}
		log(ctx).Error().Err(err).Str("id", userID).Msg("error reading query user pprofile")
		err = constant.ErrUnknownSource
		return
	}

	return
}
