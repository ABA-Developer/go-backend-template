package presenter

import "be-dashboard-nba/pkg/entities"

type Access struct {
	AccessID   int    `json:"access_id"`
	AccessName string `json:"access_name"`
	HasAccess  bool   `json:"has_access"`
}

type RoleAccessResponse struct {
	MenuID   int      `json:"menu_id"`
	MenuName string   `json:"menu_name"`
	Accesses []Access `json:"accesses"`
}

func ToReadRoleAccessResponse(rows []entities.RoleAccessResponse) []RoleAccessResponse {
	menuMap := make(map[int]*RoleAccessResponse)

	for _, r := range rows {
		if _, exists := menuMap[r.MenuID]; !exists {
			menuMap[r.MenuID] = &RoleAccessResponse{
				MenuID:   r.MenuID,
				MenuName: r.MenuName,
				Accesses: []Access{},
			}
		}

		menuMap[r.MenuID].Accesses = append(menuMap[r.MenuID].Accesses, Access{
			AccessID:   r.PermissionID,
			AccessName: r.PermissionName,
			HasAccess:  r.HasAccess,
		})
	}

	result := make([]RoleAccessResponse, 0, len(menuMap))
	for _, m := range menuMap {
		result = append(result, *m)
	}

	return result
}
