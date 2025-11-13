package presenter

import "be-dashboard-nba/pkg/entities"

type MenuParent struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Group string `json:"group"`
}

func ToReadMenuParentResponses(entities []entities.Menu) (data []MenuParent) {
	data = make([]MenuParent, len(entities))
	for i, e := range entities {

		parentData := MenuParent{
			ID:    e.ID,
			Name:  e.Name,
			Group: e.Group,
		}
		data[i] = parentData
	}

	return
}
