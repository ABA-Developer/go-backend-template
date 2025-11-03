package presenter

import (
	"be-dashboard-nba/pkg/entities"
)

type ReadUserResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	FullName string  `json:"full_name"`
	Email    string  `json:"email"`
	Phone    *string `json:"phone"`
	Active   bool    `json:"active"`
	ImgPath  *string `json:"img_path"`
	ImgName  *string `json:"img_name"`
}

func ToReadUserResponse(entity entities.User) (response ReadUserResponse) {
	response.ID = entity.ID
	response.Name = entity.Name
	response.FullName = entity.FullName
	response.Email = entity.Email
	response.Active = entity.Active

	if entity.Phone.Valid {
		response.Phone = &entity.Phone.String
	}
	if entity.ImgPath.Valid {
		response.ImgPath = &entity.ImgPath.String
	}
	if entity.ImgName.Valid {
		response.ImgName = &entity.ImgName.String
	}

	return
}

func ToReadUserResponses(entities []entities.User) (response []ReadUserResponse) {
	response = make([]ReadUserResponse, len(entities))

	for i := range entities {
		response[i] = ToReadUserResponse(entities[i])
	}

	return
}
