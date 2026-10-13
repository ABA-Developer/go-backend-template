package mapper

import (
	"database/sql"

	"github.com/ABA-Developer/go-backend-template/internal/application/menu/dto"
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
)

func ToUpdateMenuSortParamsList(req *dto.UpdateMenuOrderPayload, updatedBy string) []repoParams.UpdateMenuSortParams {
	paramsList := make([]repoParams.UpdateMenuSortParams, len(req.SortedIDs))

	var parentID sql.NullInt32
	if req.ParentID != nil {
		parentID = sql.NullInt32{Int32: int32(*req.ParentID), Valid: true}
	} else {
		parentID = sql.NullInt32{Valid: false}
	}

	for index, menuID := range req.SortedIDs {
		paramsList[index] = repoParams.UpdateMenuSortParams{
			ID:        menuID,
			Sort:      index,
			UpdatedBy: updatedBy,
			Group:     req.Group,
			ParentID:  parentID,
		}
	}

	return paramsList
}
