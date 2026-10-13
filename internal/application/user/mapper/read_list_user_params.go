package mapper

import (
	"github.com/ABA-Developer/go-backend-template/constant"
	"github.com/ABA-Developer/go-backend-template/internal/application/user/dto"
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
)

func ToReadListUserParams(req *dto.ReadUserPayload) repoParams.ReadListUserParams {
	req.Init()

	if req.Limit <= 0 {
		req.Limit = constant.DefaultLimit
	}

	if req.Page <= 0 {
		req.Page = constant.DefaultPage
	}

	// Preserve existing behavior: Offset comes from PaginationPayload.Init().
	return repoParams.ReadListUserParams{
		SetSearch: req.SetSearch,
		Search:    req.Search,
		Order:     req.Order,
		Offset:    req.Offset,
		Limit:     req.Limit,
	}
}
