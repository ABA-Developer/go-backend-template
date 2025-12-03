package presenter

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/utils"
	"be-dashboard-nba/pkg/role/repository"
)

type ReadRolesRequest struct {
	utils.PaginationPayload
}

func (req *ReadRolesRequest) ToParams() (params repository.ReadRolesParams) {
	req.Init()
	if req.Limit <= 0 {
		req.Limit = constant.DefaultLimit
	}

	if req.Page <= 0 {
		req.Page = constant.DefaultPage
	}

	params = repository.ReadRolesParams{
		SetSearch: req.SetSearch,
		Search:    req.Search,
		Order:     req.Order,
		Limit:     req.Limit,
		Offset:    req.Offset,
	}
	return
}
