package http

import (
	"be-dashboard-nba/internal/modules/auth/domain"
)

type UserResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	FullName string  `json:"full_name"`
	Email    string  `json:"email"`
	Active   bool    `json:"active"`
	Phone    *string `json:"phone"`
	ImgPath  *string `json:"img_path"`
	ImgName  *string `json:"img_name"`
}

type SessionResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

func ToReadAuthMeResponse(entity domain.User) UserResponse {
	return UserResponse{
		ID:       entity.ID,
		Name:     entity.Name,
		FullName: entity.FullName,
		Email:    entity.Email,
		Active:   entity.Active,
		Phone:    entity.Phone,
		ImgPath:  entity.ImgPath,
		ImgName:  entity.ImgName,
	}
}

func ToSessionResponse(entity domain.Session, user domain.User) SessionResponse {
	return SessionResponse{
		AccessToken:  entity.AccessToken,
		RefreshToken: entity.RefreshToken,
		User:         ToReadAuthMeResponse(user),
	}
}
