package payload

import (
	"time"

	"be-dashboard-nba/pkg/entities"
)

type ReadUserResponse struct {
	ID         int64      `json:"id"`
	FirstName  string     `json:"first_name"`
	MiddleName string     `json:"middle_name"`
	LastName   string     `json:"last_name"`
	Email      string     `json:"email"`
	Role       string     `json:"role"`
	IsActive   bool       `json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
	CreatedBy  *int64     `json:"created_by"`
	UpdatedAt  *time.Time `json:"updated_at"`
	UpdatedBy  *int64     `json:"updated_by"`
}

func ToReadUserResponse(entity entities.User) (response ReadUserResponse) {
	response.ID = entity.ID
	response.FirstName = entity.FirstName
	response.MiddleName = entity.MiddleName
	response.LastName = entity.LastName
	response.Email = entity.Email
	response.Role = entity.Role
	response.IsActive = entity.IsActive
	response.CreatedAt = entity.CreatedAt

	if entity.CreatedBy.Valid {
		response.CreatedBy = &entity.CreatedBy.Int64
	}

	if entity.UpdatedAt.Valid {
		response.UpdatedAt = &entity.UpdatedAt.Time
	}

	if entity.UpdatedBy.Valid {
		response.UpdatedBy = &entity.UpdatedBy.Int64
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
