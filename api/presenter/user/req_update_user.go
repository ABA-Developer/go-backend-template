package presenter

import "be-dashboard-nba/pkg/user/repository"

type UpdateUserRequest struct {
	Name     string  `json:"name"`
	FullName string  `json:"full_name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Phone    string  `json:"phone"`
	Active   *bool   `json:"active"`
	ImgPath  *string `json:"img_path"`
	ImgName  *string `json:"img_name"`
}

func (req *UpdateUserRequest) ToParams(userID string, password string) (params repository.UpdateUserParams) {

	params = repository.UpdateUserParams{
		ID:        userID,
		Name:      req.Name,
		FullName:  req.FullName,
		Email:     req.Email,
		Password:  password,
		Phone:     req.Phone,
		UpdatedBy: userID,
	}

	if req.Active != nil {
		params.Active = *req.Active
	}

	if req.Active == nil {
		params.Active = true
	}
	if req.ImgPath != nil {
		params.ImgPath = *req.ImgPath
	} else {
		params.ImgPath = ""
	}

	if req.ImgName != nil {
		params.ImgName = *req.ImgName
	} else {
		params.ImgName = ""
	}

	return
}
