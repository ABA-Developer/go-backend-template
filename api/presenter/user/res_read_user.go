package presenter

import (
	"be-dashboard-nba/pkg/entities"
)

type ReadUserResponse struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	Active   bool   `json:"active"`
}

func ToReadUserResponse(entity entities.User) (response ReadUserResponse) {
	response.ID = entity.ID
	response.FullName = entity.FullName
	response.Role = entity.Role
	response.Active = entity.Active

	return
}

func ToReadUserResponses(entities []entities.User) (response []ReadUserResponse) {
	response = make([]ReadUserResponse, len(entities))

	for i := range entities {
		response[i] = ToReadUserResponse(entities[i])
	}

	return
}
