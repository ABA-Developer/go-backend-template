package usecase

import (
	"context"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/menu/domain"

	"github.com/pkg/errors"
)

func (s *menuUsecase) CreateMenuUsecase(
	ctx context.Context,
	payload domain.MenuCreatePayload,
	userID string,
) (err error) {

	// 2. Lempar seluruh logika Transaksi & Pencarian Sort ke Repository
	err = s.menuRepo.CreateMenuQuery(ctx, payload, userID)
	if err != nil {
		s.log.Error().Err(err).Interface("request_payload", payload).Msg("error to create menu transactionally")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	return nil
}
