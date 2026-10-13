package usecase

import (
	"github.com/ABA-Developer/go-backend-template/constant"
	appdto "github.com/ABA-Developer/go-backend-template/internal/application/menu/dto"
	"github.com/ABA-Developer/go-backend-template/internal/application/menu/mapper"
	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/logger"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *useCase) CreateMenuUseCase(
	ctx context.Context,
	request appdto.CreateMenuPayload,
	userID string,
) (err error) {

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		logger.WithContext(ctx).Error().Err(err).Msg("error to begin transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				logger.WithContext(ctx).Error().Err(errRollback).AnErr("original_error", err).Msg("error to rollback transaction")
				err = errors.WithStack(constant.ErrUnknownSource)
				return
			}
		}
	}()

	q := s.newMenuRepo(tx)

	params := mapper.ToCreateMenuParams(&request, userID)

	var nextSort int

	if params.ParentID.Valid {
		nextSort, err = q.ReadNextSortForParentQry(ctx, params.ParentID.Int32)
	} else {
		nextSort, err = q.ReadSortForGroupQry(ctx, params.Group)
	}

	if err != nil {
		logger.WithContext(ctx).Error().Err(err).Msg("error getting next sort value for menu")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	params.Sort = nextSort

	err = q.CreateMenuQry(ctx, params)
	if err != nil {
		logger.WithContext(ctx).Error().Err(err).Interface("request_payload", params).Msg("error to create menu")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	if err = tx.Commit(); err != nil {
		logger.WithContext(ctx).Error().Err(err).Msg("error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}

	return
}
