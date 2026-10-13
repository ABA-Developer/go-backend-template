package usecase

import (
	"github.com/ABA-Developer/go-backend-template/constant"
	appdto "github.com/ABA-Developer/go-backend-template/internal/application/user/dto"
	"github.com/ABA-Developer/go-backend-template/internal/application/user/mapper"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *useCase) UpdateUserUseCase(
	ctx context.Context,
	request appdto.UpdateUserPayload,
	updatedBy string,
	userID string,
) (err error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		log(ctx).Error().Err(err).Msg("error to begin transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				log(ctx).Error().Err(errRollback).AnErr("original_error", err).Msg("error to rollback transaction")
				err = errors.WithStack(constant.ErrUnknownSource)
				return
			}
		}
	}()

	ru := s.newUserRepo(tx)
	rr := s.newRoleRepo(tx)

	existingUser, err := ru.ReadUserByIDQry(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log(ctx).Warn().Str("id", userID).Msg("user detail not found")
			err = constant.ErrUserIdNotFound
			return
		}
		log(ctx).Error().Err(err).Str("id", userID).Msg("error reading user detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	existingRoleUser, err := rr.ReadRoleByIDQry(ctx, request.RoleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log(ctx).Warn().Int("id", request.RoleID).Msg("role detail not found")
			err = constant.ErrRoleIdNotFound
			return
		}
		log(ctx).Error().Err(err).Int("id", request.RoleID).Msg("error reading role detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	params := mapper.ToUpdateUserParams(&request, userID, updatedBy, existingUser)

	if existingUser.Role != existingRoleUser.Name {
		err = ru.UpdateUserRoleQry(ctx, params.RoleID, params.ID)
		if err != nil {
			log(ctx).Error().Err(err).Interface("request_payload", request).Msg("error to update user role")
			err = errors.WithStack(constant.ErrUnknownSource)
			return
		}
	}

	err = ru.UpdateUserQry(ctx, params)
	if err != nil {
		log(ctx).Error().Err(err).Interface("request_payload", request).Msg("error to update user")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	if err = tx.Commit(); err != nil {
		log(ctx).Error().Err(err).Msg("error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}

	return
}
