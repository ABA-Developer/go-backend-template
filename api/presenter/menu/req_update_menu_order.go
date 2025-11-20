package presenter

import (
	"be-dashboard-nba/pkg/menu/repository"
	"database/sql"
)

type UpdateMenuOrderRequest struct {
	Group     string `json:"group" validate:"required"`
	ParentID  *int   `json:"parent_id" validate:"omitempty,min=1"`
	SortedIDs []int  `json:"sorted_ids" validate:"required,min=1,dive,min=1"`
}

func (req *UpdateMenuOrderRequest) ToParamsList(userID string) []repository.UpdateMenuSortParams {

	paramsList := make([]repository.UpdateMenuSortParams, len(req.SortedIDs))

	var parentID sql.NullInt32
	if req.ParentID != nil {
		parentID = sql.NullInt32{
			Int32: int32(*req.ParentID),
			Valid: true,
		}
	} else {
		parentID = sql.NullInt32{Valid: false}
	}

	for index, menuID := range req.SortedIDs {
		paramsList[index] = repository.UpdateMenuSortParams{
			ID:        menuID,
			Sort:      index,
			UpdatedBy: userID,
			Group:     req.Group,
			ParentID:  parentID,
		}
	}

	return paramsList
}
