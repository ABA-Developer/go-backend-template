package usecase

import (
	"context"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/menu/domain"

	"github.com/pkg/errors"
)

func (s *menuUsecase) ReadMenuDetailUsecase(ctx context.Context, menuID int) (data domain.MenuDetail, err error) {
	menuData, err := s.menuRepo.ReadMenuByIDQuery(ctx, menuID)

	if err != nil {
		if errors.Is(err, constant.ErrDataNotFound) {
			s.log.Warn().Int("id", menuID).Msg("menu detail not found")
			return data, constant.ErrMenuIdNotFound
		}

		s.log.Error().Err(err).Int("id", menuID).Msg("error reading menu detail query")
		return data, errors.WithStack(constant.ErrUnknownSource)
	}

	data = domain.MenuDetail{
		ID:          menuData.ID,
		ParentID:    menuData.ParentID,
		Name:        menuData.Name,
		Description: menuData.Description,
		URL:         menuData.URL,
		Sort:        menuData.Sort,
		Group:       menuData.Group,
		Icon:        menuData.Icon,
		Active:      menuData.Active,
		Display:     menuData.Display,
	}

	return data, nil
}
