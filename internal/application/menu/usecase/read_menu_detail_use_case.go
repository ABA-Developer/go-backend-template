package usecase

import (
	"github.com/ABA-Developer/go-backend-template/constant"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	"github.com/ABA-Developer/go-backend-template/internal/infrastructure/logger"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *useCase) ReadMenuDetailUseCase(ctx context.Context, menuID int) (data model.Menu, err error) {

	r := s.newMenuRepo(s.db)
	data, err = r.ReadMenuByIDQry(ctx, menuID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.WithContext(ctx).Warn().Int("id", menuID).Msg("menu detail not found")
			err = constant.ErrMenuIdNotFound
			return
		}

		logger.WithContext(ctx).Error().Err(err).Int("id", menuID).Msg("error reading menu detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}
	return
}
