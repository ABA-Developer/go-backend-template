package presenter

import (
	"be-dashboard-nba/pkg/user/repository"
	"database/sql" // Impor untuk sql.NullString
)

// DIUBAH: Disesuaikan dengan DDL 'app_user'
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	FullName string `json:"full_name" validate:"required"` // DDL baru NOT NULL
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Phone    string `json:"phone"`    // Opsional
	Active   *bool  `json:"active"`   // Opsional, DDL punya default
	ImgPath  string `json:"img_path"` // Opsional
	ImgName  string `json:"img_name"` // Opsional
	// DIHAPUS: FirstName, MiddleName, LastName, Role, IsActive
}

// DIUBAH: Disesuaikan dengan DDL 'app_user'
func (req *CreateUserRequest) ToParams(userID string, password string) (params repository.CreateUserParams) {
	// userID (CreatedBy) sekarang string

	params = repository.CreateUserParams{
		Name:      req.Name,
		FullName:  req.FullName,
		Email:     req.Email,
		Password:  password,
		CreatedBy: userID, // string -> string
	}

	// Menangani nilai default/opsional
	if req.Active != nil {
		params.Active = *req.Active
	} else {
		params.Active = true // Sesuai default DDL
	}

	// Menangani nullable string
	if req.Phone != "" {
		params.Phone = sql.NullString{String: req.Phone, Valid: true}
	}
	if req.ImgPath != "" {
		params.ImgPath = sql.NullString{String: req.ImgPath, Valid: true}
	}
	if req.ImgName != "" {
		params.ImgName = sql.NullString{String: req.ImgName, Valid: true}
	}

	return
}
