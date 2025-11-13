package service

import (
	menuPresenter "be-dashboard-nba/api/presenter/menu"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/menu/repository"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *service) UpdateMenuService(
	ctx context.Context,
	request menuPresenter.UpdateMenuRequest,
	userID string,
	menuID int,
) (err error) {

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

	_, err = q.ReadMenuByIDQuery(ctx, menuID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warn().Int("id", menuID).Msg("menu detail not found for update")
			err = constant.ErrMenuIdNotFound
			return
		}
		s.log.Error().Err(err).Int("id", menuID).Msg("error reading menu detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	params := request.ToUpdateParams(userID, menuID)

	err = q.UpdateMenuQuery(ctx, params)
	if err != nil {
		s.log.Error().Err(err).Interface("request_payload", request).Msg("error to update menu")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	if err = tx.Commit(); err != nil {
		s.log.Error().Err(err).Msg("error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}

	return
}
