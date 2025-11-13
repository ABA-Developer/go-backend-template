package presenter

import (
	"be-dashboard-nba/pkg/menu/repository"
)

type UpdateMenuOrderRequest struct {
	ParentID  *int  `json:"parent_id" validate:"omitempty,min=1"`
	SortedIDs []int `json:"sorted_ids" validate:"required,min=1,dive,min=1"`
}

func (req *UpdateMenuOrderRequest) ToParamsList(userID string) []repository.UpdateMenuSortParams {

	paramsList := make([]repository.UpdateMenuSortParams, len(req.SortedIDs))

	for index, menuID := range req.SortedIDs {
		paramsList[index] = repository.UpdateMenuSortParams{
			ID:        menuID,
			Sort:      index,
			UpdatedBy: userID,
		}
	}

	return paramsList
}
