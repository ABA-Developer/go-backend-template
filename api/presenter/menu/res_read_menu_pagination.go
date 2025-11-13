package presenter

import (
	"be-dashboard-nba/internal/utils"
	"be-dashboard-nba/pkg/menu/repository"
)

type ReadMenuListRequest struct {
	utils.PaginationPayload
}

func (req *ReadMenuListRequest) ToParams() (params repository.ReadListMenuParams) {
	req.Init()

	params = repository.ReadListMenuParams{
		SetSearch: req.SetSearch,
		Search:    req.Search,
	}

	return
}
