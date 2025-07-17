package payload

import (
	"time"

	"be-dashboard-nba/pkg/entities"
)

type UserResponse struct {
	ID         int64      `json:"id"`
	FirstName  string     `json:"first_name,omitempty"`
	MiddleName string     `json:"middle_name,omitempty"`
	LastName   string     `json:"last_name,omitempty"`
	Email      string     `json:"email"`
	Role       string     `json:"role"`
	CreatedAt  time.Time  `json:"created_at"`
	CreatedBy  *int64     `json:"created_by"`
	UpdatedAt  *time.Time `json:"updated_at"`
	UpdatedBy  *int64     `json:"updated_by"`
}

type SessionResponse struct {
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiredAt  time.Time    `json:"access_token_expired_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiredAt time.Time    `json:"refresh_token_expired_at"`
	User                  UserResponse `json:"user"`
}

func ToSessionResponse(entity entities.Session, user entities.User) (response SessionResponse) {
	response.AccessToken = entity.AccessToken
	response.AccessTokenExpiredAt = entity.AccessTokenExpiredAt
	response.RefreshToken = entity.RefreshToken
	response.RefreshTokenExpiredAt = entity.RefreshTokenExpiredAt
	response.User.ID = user.ID
	response.User.FirstName = user.FirstName
	response.User.MiddleName = user.MiddleName
	response.User.LastName = user.LastName
	response.User.Email = user.Email
	response.User.Role = user.Role
	response.User.CreatedAt = user.CreatedAt

	if user.CreatedBy.Valid {
		response.User.CreatedBy = &user.CreatedBy.Int64
	}

	if user.UpdatedAt.Valid {
		response.User.UpdatedAt = &user.UpdatedAt.Time
	}

	if user.UpdatedBy.Valid {
		response.User.UpdatedBy = &user.UpdatedBy.Int64
	}

	return
}
