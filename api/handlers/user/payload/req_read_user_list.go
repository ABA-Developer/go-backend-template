package payload

import (
	"be-dashboard-nba/internal/utils"
	"be-dashboard-nba/pkg/user/repository"
)

type ReadUserListRequest struct {
	utils.PaginationPayload
}

func (req *ReadUserListRequest) ToParams() (params repository.ReadListUserParams) {
	req.Init()

	params = repository.ReadListUserParams{
		SetSearch: req.SetSearch,
		Search:    req.Search,
		Order:     req.Order,
		Offset:    req.Offset,
		Limit:     req.Limit,
	}

	return
}
