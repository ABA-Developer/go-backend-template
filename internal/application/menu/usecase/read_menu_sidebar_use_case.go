package usecase

import (
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/logger"
	"context"
)

func (s *useCase) ReadSidebarMenuUseCase(
	ctx context.Context,
	userID string,
) (data []model.Menu, err error) {
	r := s.newMenuRepo(s.db)
	data, err = r.ReadSidebarMenuQry(ctx, userID)
	if err != nil {
		logger.WithContext(ctx).Error().Err(err).Msg("error query read list menu")
		return
	}

	return
}
