package presenter

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/utils"
	"be-dashboard-nba/repository/role"
)

type ReadRoleAccessesRequest struct {
	utils.PaginationPayload
}

func (req *ReadRoleAccessesRequest) ToParams(roleID int) (params repository.ReadRoleAccessParams) {
	req.Init()
	if req.Limit <= 0 {
		req.Limit = constant.DefaultLimit
	}

	if req.Page <= 0 {
		req.Page = constant.DefaultPage
	}

	params = repository.ReadRoleAccessParams{
		SetSearch: req.SetSearch,
		Search:    req.Search,
		Order:     req.Order,
		Limit:     req.Limit,
		Offset:    req.Offset,
		RoleID:    roleID,
	}
	return
}
