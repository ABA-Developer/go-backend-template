package presenter

import (
	"be-dashboard-nba/repository/menu"
)

type ReadMenuListRequest struct {
	SetSearch bool
	Search    string
}

func (r *ReadMenuListRequest) ToParams() repository.ReadListMenuParams {
	params := repository.ReadListMenuParams{
		SetSearch: r.SetSearch,
		Search:    r.Search,
	}

	if r.Search != "" {
		params.Search = "%" + r.Search + "%"
		params.SetSearch = true
	}

	return params
}
