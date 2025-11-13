package service

import (
	presenter "be-dashboard-nba/api/presenter/role"
	"be-dashboard-nba/constant"
	menuPermissionRepository "be-dashboard-nba/pkg/menu_permission/repository"
	roleRepository "be-dashboard-nba/pkg/role/repository"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *service) UpdateRoleAccessService(ctx context.Context, roleID int, request presenter.UpdateRoleAccessRequest) (err error) {
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

	rolerepo := roleRepository.NewRepository(tx)
	menuPermissionRepo := menuPermissionRepository.NewRepository(tx)
	for _, req := range request.AccessItem {
		params := req.ToParams(roleID)

		_, err = rolerepo.ReadRoleByIDQuery(ctx, params.RoleID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				s.log.Warn().Int("id", params.RoleID).Msg("role detail not found for update role menu permission")
				err = constant.ErrRoleIdNotFound
				return
			}
			s.log.Error().Err(err).Int("id", params.RoleID).Msg("error reading role detail query on update role menu permission")
			err = errors.WithStack(constant.ErrUnknownSource)
			return
		}

		_, err = menuPermissionRepo.ReadMenuPermissionByIdQuery(ctx, params.MenuPermissionID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				s.log.Warn().Int("id", params.RoleID).Msg("menu permission detail not found for update role menu permission")
				err = constant.ErrMenuPermissionIdNotFound
				return
			}
			s.log.Error().Err(err).Int("id", params.RoleID).Msg("error reading menu permission detail query on update role menu permission")
			err = errors.WithStack(constant.ErrUnknownSource)
			return
		}
		if *req.HasAccess {
			err = rolerepo.CreateRoleAccess(ctx, params)
			if err != nil {
				s.log.Error().Err(err).Interface("request_payload", request).Msg("error to create role menu permission")
				err = errors.WithStack(constant.ErrUnknownSource)
				return
			}
		} else {
			err = rolerepo.DeleteRoleAccess(ctx, params)
			if err != nil {
				s.log.Error().Err(err).Interface("request_payload", request).Msg("error to delete role menu permission")
				err = errors.WithStack(constant.ErrUnknownSource)
				return
			}
		}
	}

	if err = tx.Commit(); err != nil {
		s.log.Error().Err(err).Msg("error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}

	return

}
