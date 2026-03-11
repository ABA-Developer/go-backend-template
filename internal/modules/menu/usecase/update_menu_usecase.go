package usecase

import (
	"context"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/menu/domain"

	"github.com/pkg/errors"
)

func (s *menuUsecase) UpdateMenuUsecase(
	ctx context.Context,
	payload domain.MenuUpdatePayload,
) (err error) {
	existingMenu, err := s.menuRepo.ReadMenuByIDQuery(ctx, payload.ID)
	if err != nil {
		if errors.Is(err, constant.ErrDataNotFound) {
			s.log.Warn().Int("id", payload.ID).Msg("menu detail not found for update")
			return constant.ErrMenuIdNotFound
		}
		s.log.Error().Err(err).Msg("error reading menu detail query")
		return errors.WithStack(constant.ErrUnknownSource)
	}

	childCount, err := s.menuRepo.CountMenuChildren(ctx, payload.ID)
	if err != nil {
		s.log.Error().Err(err).Msg("error counting menu children")
		return errors.WithStack(constant.ErrUnknownSource)
	}

	if payload.Name == "" {
		payload.Name = existingMenu.Name
	}
	if payload.Description == nil {
		payload.Description = existingMenu.Description
	}
	if payload.URL == nil {
		payload.URL = existingMenu.URL
	}
	if payload.Icon == nil {
		payload.Icon = existingMenu.Icon
	}

	var parentChanged bool
	if existingMenu.ParentID == nil && payload.ParentID == nil {
		parentChanged = false
	} else if existingMenu.ParentID != nil && payload.ParentID != nil {
		parentChanged = *existingMenu.ParentID != *payload.ParentID
	} else {
		parentChanged = true
	}

	groupChanged := existingMenu.Group != payload.Group

	if childCount > 0 && parentChanged {
		return errors.WithStack(constant.ErrMenuHasChildren)
	}

	if payload.ParentID != nil {
		parentMenu, err := s.menuRepo.ReadMenuByIDQuery(ctx, int(*payload.ParentID))
		if err != nil {
			s.log.Error().Err(err).Msg("error reading parent menu")
			return errors.WithStack(constant.ErrUnknownSource)
		}
		payload.Group = parentMenu.Group
	}

	payload.Sort = existingMenu.Sort

	updateChildrenGroup := childCount > 0 && groupChanged && !parentChanged

	err = s.menuRepo.UpdateMenuQuery(ctx, payload, updateChildrenGroup)
	if err != nil {
		s.log.Error().Err(err).Msg("error to update menu transactionally")
		return errors.WithStack(constant.ErrUnknownSource)
	}

	return nil
}
