package http

import (
	"be-dashboard-nba/internal/modules/role/domain"
)

type ReadRoleResponse struct {
	ID          int     `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

func ToReadRoleResponse(row domain.Role) (response ReadRoleResponse) {
	response.ID = row.ID
	response.Code = row.Code
	response.Name = row.Name
	response.Description = row.Description
	return
}

func ToReadRoleResponses(rows []domain.Role) (response []ReadRoleResponse) {
	response = make([]ReadRoleResponse, len(rows))
	for i := range rows {
		response[i] = ToReadRoleResponse(rows[i])
	}
	return
}

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

func ToReadRoleAccessResponse(rows []domain.RoleAccessResponse) []RoleAccessResponse {
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
