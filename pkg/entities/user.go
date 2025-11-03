package entities

import (
	"database/sql"
	"time"
)

type User struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`      // DIUBAH: dari FirstName
	FullName  string         `json:"full_name"` // DIUBAH: dari MiddleName/LastName
	Email     string         `json:"email"`     // DDL baru: NOT NULL
	Password  string         `json:"password"`
	Active    bool           `json:"active"`     // DIUBAH: dari IsActive
	Phone     sql.NullString `json:"phone"`      // BARU
	ImgPath   sql.NullString `json:"img_path"`   // BARU
	ImgName   sql.NullString `json:"img_name"`   // BARU
	CreatedBy string         `json:"created_by"` // DIUBAH: dari sql.NullInt64 ke string (NOT NULL)
	CreatedAt time.Time      `json:"created_at"` // DDL baru: NOT NULL
	UpdatedBy sql.NullString `json:"updated_by"` // DIUBAH: dari sql.NullInt64 ke sql.NullString
	UpdatedAt sql.NullTime   `json:"updated_at"`
}
