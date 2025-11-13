package entities

import (
	"database/sql"
	"time"
)

type User struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	FullName  string         `json:"full_name"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	Active    bool           `json:"active"`
	Phone     sql.NullString `json:"phone"`
	ImgPath   sql.NullString `json:"img_path"`
	ImgName   sql.NullString `json:"img_name"`
	CreatedBy string         `json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedBy sql.NullString `json:"updated_by"`
	UpdatedAt sql.NullTime   `json:"updated_at"`
}
