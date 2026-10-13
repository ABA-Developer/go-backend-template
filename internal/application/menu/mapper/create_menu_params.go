package mapper

import (
	"database/sql"
	"strings"

	"github.com/ABA-Developer/go-backend-template/internal/application/menu/dto"
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
)

func ToCreateMenuParams(req *dto.CreateMenuPayload, createdBy string) repoParams.CreateMenuParams {
	params := repoParams.CreateMenuParams{
		Name:      req.Name,
		CreatedBy: createdBy,
		Group:     req.Group,
	}

	if req.Active != nil {
		params.Active = *req.Active
	}

	if req.Display != nil {
		params.Display = *req.Display
	}

	if req.ParentID != nil {
		params.ParentID = sql.NullInt32{Int32: int32(*req.ParentID), Valid: true}
	}

	if req.Description != nil {
		params.Description = sql.NullString{String: *req.Description, Valid: true}
	}

	if req.URL != nil {
		urlToStore := *req.URL
		if !strings.HasPrefix(urlToStore, "/") {
			urlToStore = "/" + urlToStore
		}
		params.URL = sql.NullString{String: urlToStore, Valid: true}
	}

	if req.Icon != nil {
		params.Icon = sql.NullString{String: *req.Icon, Valid: true}
	}

	return params
}
