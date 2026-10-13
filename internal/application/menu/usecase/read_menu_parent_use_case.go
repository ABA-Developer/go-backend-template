package usecase

import (
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/logger"
	"context"
)

func (s *useCase) ReadMenuParentUseCase(ctx context.Context) (data []model.Menu, err error) {
	r := s.newMenuRepo(s.db)
	data, err = r.ReadParentMenuQry(ctx)
	if err != nil {
		logger.WithContext(ctx).Error().Err(err).Msg("error query read  menu parent")
		return
	}

	return
}
