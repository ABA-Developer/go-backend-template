package service

import (
	"context"
	"database/sql"

	"github.com/gofiber/fiber/v2/log"
	"github.com/pkg/errors"

	"be-dashboard-nba/api/handlers/user/payload"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/utils"
	"be-dashboard-nba/pkg/user/repository"
)

func (s *Service) UpdateUserService(
	ctx context.Context,
	request payload.UpdateUserRequest,
	userID int64,
) (err error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		log.WithContext(ctx).Error(err, "error to begin transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				log.WithContext(ctx).Error(errRollback, "error to rollback transaction", "original_error", err)
				err = errors.WithStack(constant.ErrUnknownSource)
				return
			}
		}
	}()

	q := repository.NewQuery(tx)

	var password string

	if request.Password != "" {
		password, err = utils.GenerateHashPassword(request.Password)
		if err != nil {
			log.WithContext(ctx).Error(err, "error generate hash password : "+request.Password)
			return
		}
	}

	err = q.UpdateUserQuery(ctx, request.ToParams(userID, password))
	if err != nil {
		log.WithContext(ctx).Error(err, "error to update user", "request", request)
		return
	}

	if err = tx.Commit(); err != nil {
		log.WithContext(ctx).Error(err, "error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}

	return
}
