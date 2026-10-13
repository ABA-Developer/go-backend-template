package mapper

import (
	"github.com/ABA-Developer/go-backend-template/internal/application/menu/dto"
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
)

func ToReadListMenuParams(req *dto.ReadMenuListPayload) repoParams.ReadListMenuParams {
	params := repoParams.ReadListMenuParams{
		SetSearch: req.SetSearch,
		Search:    req.Search,
	}

	if req.Search != "" {
		params.Search = "%" + req.Search + "%"
		params.SetSearch = true
	}

	return params
}
