package usecase

import (
	appdto "github.com/ABA-Developer/go-backend-template/internal/application/menu/dto"
	"github.com/ABA-Developer/go-backend-template/internal/application/menu/mapper"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/logger"
	"context"
)

func (s *useCase) ReadListMenuUseCase(
	ctx context.Context,
	request appdto.ReadMenuListPayload,
) (data []model.Menu, err error) {
	r := s.newMenuRepo(s.db)
	data, err = r.ReadListMenuQry(ctx, mapper.ToReadListMenuParams(&request))
	if err != nil {
		logger.WithContext(ctx).Error().Err(err).Msg("error query read list menu")
		return
	}

	return
}
