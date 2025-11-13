package service

import (
	menuPresenter "be-dashboard-nba/api/presenter/menu"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/menu/repository"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *service) UpdateMenuOrderService(ctx context.Context, request menuPresenter.UpdateMenuOrderRequest, userID string) (err error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		s.log.Error().Err(err).Msg("error to begin transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				s.log.Error().Err(errRollback).AnErr("original_error", err).Msg("error to rollback transaction")
				err = errors.WithStack(constant.ErrUnknownSource)
				return
			}
		}
	}()

	q := repository.NewRepository(tx)

	paramsList := request.ToParamsList(userID)

	for _, params := range paramsList {
		err = q.UpdateMenuSortQuery(ctx, params)
		if err != nil {
			s.log.Error().Err(err).Int("menu_id", params.ID).Msg("Failed to update menu sort")
			err = errors.WithStack(constant.ErrUnknownSource)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		s.log.Error().Err(err).Msg("error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}

	return
}
