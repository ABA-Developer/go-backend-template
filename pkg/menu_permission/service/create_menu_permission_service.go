package service

import (
	menuPermissionPresenter "be-dashboard-nba/api/presenter/menu_permission"
	"be-dashboard-nba/constant"
	menuRepo "be-dashboard-nba/pkg/menu/repository"
	menuPermissionRepo "be-dashboard-nba/pkg/menu_permission/repository"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *service) CreateMenuPermissionService(ctx context.Context, payload menuPermissionPresenter.CreateMenuPermissionRequest, userID string, menuID int) (err error) {

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

	rmp := menuPermissionRepo.NewRepository(tx)
	rm := menuRepo.NewRepository(tx)

	_, err = rm.ReadMenuByIDQuery(ctx, menuID)

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

	err = rmp.CreateMenuPermissionQuery(ctx, payload.ToParams(userID, menuID))
	if err != nil {
		s.log.Error().Err(err).Interface("request_payload", payload).Msg("error to create menu permission")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	if err = tx.Commit(); err != nil {
		s.log.Error().Err(err).Msg("error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}
	return
}
