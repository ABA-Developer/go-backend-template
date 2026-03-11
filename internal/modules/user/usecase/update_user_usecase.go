package usecase

import (
	"context"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/user/domain"

	"github.com/pkg/errors"
)

func (s *userUsecase) UpdateUserUsecase(ctx context.Context, payload domain.UpdateUserPayload) (err error) {
	existingUser, err := s.userRepo.ReadUserByIDQuery(ctx, payload.ID)
	if err != nil {
		if errors.Is(err, constant.ErrDataNotFound) {
			s.log.Warn().Str("id", payload.ID).Msg("user detail not found")
			err = constant.ErrUserIdNotFound
			return
		}
		s.log.Error().Err(err).Str("id", payload.ID).Msg("error reading user detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	userDomain := domain.User{
		ID:        payload.ID,
		Name:      payload.Name,
		FullName:  payload.FullName,
		Email:     payload.Email,
		UpdatedBy: payload.UpdatedBy,
		RoleID:    payload.RoleID,
		Active:    existingUser.Active,
		Phone:     existingUser.Phone,
	}

	if payload.Active != nil {
		userDomain.Active = *payload.Active
	}

	if payload.Phone != nil {
		userDomain.Phone = payload.Phone
	}

	err = s.userRepo.UpdateUserWithRoleTx(ctx, userDomain)
	if err != nil {
		if errors.Is(err, constant.ErrEmailAlreadyExists) || errors.Is(err, constant.ErrRoleIdNotFound) {
			s.log.Warn().Err(err).Msg("validation error on update user")
			return err
		}
		s.log.Error().Err(err).Interface("request_payload", payload).Msg("error to update user transactionally")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	return nil
}
