package presenter

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/utils"
	"be-dashboard-nba/repository/user"
)

type ReadUserRequest struct {
	utils.PaginationPayload
}

func (req *ReadUserRequest) ToParams() (params repository.ReadListUserParams) {
	req.Init()

	if req.Limit <= 0 {
		req.Limit = constant.DefaultLimit
	}

	if req.Page <= 0 {
		req.Page = constant.DefaultPage
	}
	params = repository.ReadListUserParams{
		SetSearch: req.SetSearch,
		Search:    req.Search,
		Order:     req.Order,
		Offset:    req.Offset,
		Limit:     req.Limit,
	}

	return
}
