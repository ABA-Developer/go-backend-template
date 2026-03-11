package presenter

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/utils"
	"be-dashboard-nba/repository/menu_permission"
)

type ReadMenuPermissionListRequest struct {
	utils.PaginationPayload
}

func (req *ReadMenuPermissionListRequest) ToParams(MenuID int) (params repository.ReadMenuPermissionService) {
	req.Init()

	if req.Limit <= 0 {
		req.Limit = constant.DefaultLimit
	}

	if req.Page <= 0 {
		req.Page = constant.DefaultPage
	}

	params = repository.ReadMenuPermissionService{
		SetSearch: req.SetSearch,
		Search:    req.Search,
		Order:     req.Order,
		Limit:     req.Limit,
		Offset:    req.Offset,
		MenuID:    MenuID,
	}

	return
}
