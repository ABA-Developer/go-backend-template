package usecase

import (
	"github.com/ABA-Developer/go-backend-template/constant"
	appdto "github.com/ABA-Developer/go-backend-template/internal/application/role/dto"
	"github.com/ABA-Developer/go-backend-template/internal/application/role/mapper"
	"context"

	"github.com/pkg/errors"
)

func (s *useCase) CreateRoleUseCase(ctx context.Context, paylaod appdto.CreateRolePayload, userID string) (err error) {
	r := s.newRoleRepo(s.db)

	err = r.CreateRoleQry(ctx, mapper.ToCreateRoleParams(&paylaod, userID))
	if err != nil {
		log(ctx).Error().Err(err).Interface("request_payload", paylaod).Msg("error to create role")
		err = errors.WithStack(constant.ErrUnknownSource)
	}
	return err
}
