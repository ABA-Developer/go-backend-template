package service

import (
	presenter "be-dashboard-nba/api/presenter/user"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/db"
	"be-dashboard-nba/internal/utils"
	"be-dashboard-nba/pkg/entities"
	"be-dashboard-nba/pkg/user/repository"
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Service interface {
	CreateUserService(ctx context.Context, request presenter.CreateUserRequest, userID string) (err error)
	DeleteUserService(ctx context.Context, id string) (err error)
	ReadDetailUserService(ctx context.Context, id string) (data entities.User, err error)
	UpdateUserService(ctx context.Context, request presenter.UpdateUserRequest, userID string) (err error)
}

type service struct {
	db  db.DB
	log *zerolog.Logger
}

func NewService(db db.DB, log *zerolog.Logger) Service {
	return &service{
		db:  db,
		log: log,
	}
}

func (s *service) CreateUserService(
	ctx context.Context,
	request presenter.CreateUserRequest,
	userID string,
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
			if !errors.Is(err, constant.ErrUnknownSource) {
				err = errors.WithStack(constant.ErrUnknownSource)
			}
		}
	}()

	q := repository.NewRepository(tx)

	password, err := utils.GenerateHashPassword(request.Password)
	if err != nil {
		s.log.Error().Err(err).Msg("error generate hash password")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	err = q.CreateUserQuery(ctx, request.ToParams(userID, password))
	if err != nil {
		s.log.Error().Err(err).Interface("request_payload", request).Msg("error to create user")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	if err = tx.Commit(); err != nil {
		s.log.Error().Err(err).Msg("error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}

	return
}

func (s *service) DeleteUserService(
	ctx context.Context,
	id string,
) (err error) {
	r := repository.NewRepository(s.db)

	err = r.DeleteUserQuery(ctx, id)
	if err != nil {
		s.log.Error().Err(err).Str("id", id).Msg("error delete user query")

		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	return nil
}

func (s *service) ReadDetailUserService(
	ctx context.Context,
	id string,
) (data entities.User, err error) {
	r := repository.NewRepository(s.db)

	data, err = r.ReadDetailUserQuery(ctx, id)
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

func (s *service) UpdateUserService(
	ctx context.Context,
	request presenter.UpdateUserRequest,
	userID string,
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

			if !errors.Is(err, constant.ErrUnknownSource) {
				err = errors.WithStack(constant.ErrUnknownSource)
			}
		}
	}()

	r := repository.NewRepository(tx)

	var password string

	if request.Password != "" {
		password, err = utils.GenerateHashPassword(request.Password)
		if err != nil {
			s.log.Error().Err(err).Msg("error generate hash password")
			err = errors.WithStack(constant.ErrUnknownSource)
			return
		}
	}

	err = r.UpdateUserQuery(ctx, request.ToParams(userID, password))
	if err != nil {
		s.log.Error().Err(err).Interface("request_payload", request).Msg("error to update user")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	if err = tx.Commit(); err != nil {
		s.log.Error().Err(err).Msg("error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}

	return
}
