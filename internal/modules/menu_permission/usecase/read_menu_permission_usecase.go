package usecase

import (
	"context"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/core/utils"
	"be-dashboard-nba/internal/modules/menu_permission/domain"

	"github.com/pkg/errors"
)

func (s *menuPermissionUsecase) ReadListMenuPermissionUsecase(ctx context.Context, filter domain.MenuPermissionFilter) (data domain.MenuPermissionPaginationResponse, err error) {
	_, err = s.menuRepo.ReadMenuByIDQuery(ctx, filter.MenuID)
	if err != nil {
		if errors.Is(err, constant.ErrDataNotFound) {
			s.log.Warn().Int("menu_id", filter.MenuID).Msg("menu detail not found")
			err = constant.ErrMenuIdNotFound
			return
		}
		s.log.Error().Err(err).Int("menu_id", filter.MenuID).Msg("error reading menu detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	totalItems, err := s.menuPermissionRepo.ReadMenuPermissionCountQuery(ctx, filter)
	if err != nil {
		s.log.Error().Err(err).Msg("error query get menu permission count")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	pagination := utils.GeneratePagination(totalItems, filter.Page, filter.Limit)
	menuPermissionData, err := s.menuPermissionRepo.ReadMenuPermissionListQuery(ctx, filter)
	s.log.Print(filter)
	if err != nil {
		s.log.Error().Err(err).Msg("error query read list menu permission")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	var detailData []domain.MenuPermissionDetail
	for _, mp := range menuPermissionData {
		detailData = append(detailData, domain.MenuPermissionDetail{
			ID:         mp.ID,
			MenuID:     mp.MenuID,
			Code:       mp.Code,
			ActionName: mp.ActionName,
		})
	}

	if detailData == nil {
		detailData = []domain.MenuPermissionDetail{}
	}

	data = domain.MenuPermissionPaginationResponse{
		Data:       detailData,
		Pagination: pagination,
	}

	return
}
