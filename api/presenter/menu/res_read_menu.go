package presenter

import (
	"be-dashboard-nba/pkg/entities"
)

type MenuListItem struct {
	ID       int            `json:"id"`
	Name     string         `json:"name"`
	Icon     *string        `json:"icon"`
	Url      *string        `json:"url"`
	Children []MenuListItem `json:"children"`
	Sort     int            `json:"sort"`
}

type ReadMenuListResponse struct {
	GroupName   string         `json:"group_name"`
	GroupChilds []MenuListItem `json:"group_childs"`
}

func ToReadMenuListResponse(menuEntities []entities.Menu) []ReadMenuListResponse {
	childrenMap := make(map[int32][]entities.Menu)

	groupMap := make(map[string][]entities.Menu)

	for _, e := range menuEntities {
		if !e.ParentID.Valid {
			groupMap[e.Group] = append(groupMap[e.Group], e)
		} else {
			childrenMap[e.ParentID.Int32] = append(childrenMap[e.ParentID.Int32], e)
		}
	}

	var response []ReadMenuListResponse
	for groupName, roots := range groupMap {
		groupItem := ReadMenuListResponse{
			GroupName:   groupName,
			GroupChilds: []MenuListItem{},
		}

		for _, root := range roots {
			groupItem.GroupChilds = append(groupItem.GroupChilds, buildTreeRecursiveListItem(root, childrenMap))
		}
		response = append(response, groupItem)
	}
	return response
}

func buildTreeRecursiveListItem(current entities.Menu, childrenMap map[int32][]entities.Menu) MenuListItem {

	item := MenuListItem{
		ID:       current.ID,
		Name:     current.Name,
		Sort:     current.Sort,
		Children: []MenuListItem{},
	}

	if current.Icon.Valid {
		item.Icon = &current.Icon.String
	}
	if current.URL.Valid {
		item.Url = &current.URL.String
	}

	if children, ok := childrenMap[int32(current.ID)]; ok {
		for _, childEntity := range children {
			item.Children = append(item.Children, buildTreeRecursiveListItem(childEntity, childrenMap))
		}
	}

	return item
}
