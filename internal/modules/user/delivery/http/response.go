package http

import (
	"be-dashboard-nba/internal/modules/user/domain"
)

type ReadUserResponse struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	Active   bool   `json:"active"`
}

func ToReadUserResponse(row domain.UserListRow) (response ReadUserResponse) {
	response.ID = row.ID
	response.FullName = row.FullName
	response.Role = row.Role
	response.Active = row.Active

	return
}

func ToReadUserResponses(rows []domain.UserListRow) (response []ReadUserResponse) {
	response = make([]ReadUserResponse, len(rows))

	for i := range rows {
		response[i] = ToReadUserResponse(rows[i])
	}

	return
}

type ReadUserDetailResponse struct {
	ID       string  `json:"id"`
	FullName string  `json:"full_name"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Role     string  `json:"role"`
	RoleID   int     `json:"role_id"`
	Phone    *string `json:"phone"`
	Active   bool    `json:"active"`
}

func ToReadUserDetailResponse(row domain.UserDetailRow) (response ReadUserDetailResponse) {
	response.ID = row.ID
	response.FullName = row.FullName
	response.Name = row.Name
	response.Email = row.Email
	response.Role = row.Role
	response.RoleID = row.RoleID
	response.Active = row.Active

	response.Phone = row.Phone

	return
}

type ReadUserProfileResponse struct {
	ID       string  `json:"id"`
	FullName string  `json:"full_name"`
	Name     string  `json:"name"`
	ImgPath  *string `json:"img_path"`
	ImgName  *string `json:"img_name"`
	Email    string  `json:"email"`
	Role     string  `json:"role"`
	RoleID   int     `json:"role_id"`
	Phone    *string `json:"phone"`
	Active   bool    `json:"active"`
}

func ToReadUserProfileResponse(row domain.UserProfileRow) (response ReadUserProfileResponse) {
	response.ID = row.ID
	response.FullName = row.FullName
	response.Name = row.Name
	response.Email = row.Email
	response.Role = row.Role
	response.RoleID = row.RoleID
	response.Active = row.Active

	response.Phone = row.Phone

	response.ImgName = row.ImgName
	response.ImgPath = row.ImgPath

	return
}
