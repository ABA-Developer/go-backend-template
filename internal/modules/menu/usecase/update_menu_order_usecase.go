package usecase

import (
	"context"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/menu/domain"

	"github.com/pkg/errors"
)

func (s *menuUsecase) UpdateMenuOrderUsecase(
	ctx context.Context,
	payload domain.MenuUpdateSortPayload,
) (err error) {

	var repoPayloads []domain.MenuUpdateSortItemPayload

	var parentID32 *int32
	if payload.ParentID != nil {
		val := int32(*payload.ParentID)
		parentID32 = &val
	}

	for index, id := range payload.SortedIDs {
		repoPayloads = append(repoPayloads, domain.MenuUpdateSortItemPayload{
			ID:        id,
			Sort:      index + 1,
			UpdatedBy: payload.UpdatedBy,
			ParentID:  parentID32,
			Group:     payload.Group,
		})
	}

	err = s.menuRepo.UpdateMenuOrderQuery(ctx, repoPayloads)
	if err != nil {
		s.log.Error().Err(err).Interface("payload", payload).Msg("error to update menu order transactionally")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	return nil
}
