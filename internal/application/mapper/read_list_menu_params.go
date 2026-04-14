package mapper

import (
	"be-dashboard-nba/internal/application/dto"
	menuPresenter "be-dashboard-nba/internal/presentation/presenter/menu"
)

func ToReadListMenuParams(req *menuPresenter.ReadMenuListRequest) dto.ReadListMenuParams {
	params := dto.ReadListMenuParams{
		SetSearch: req.SetSearch,
		Search:    req.Search,
	}

	if req.Search != "" {
		params.Search = "%" + req.Search + "%"
		params.SetSearch = true
	}

	return params
}
