package entities

import (
	"database/sql"
	"time"
)

type User struct {
	ID         int64         `json:"id"`
	FirstName  string        `json:"first_name"`
	MiddleName string        `json:"middle_name"`
	LastName   string        `json:"last_name"`
	Email      string        `json:"email"`
	Password   string        `json:"password"`
	Role       string        `json:"role"`
	IsActive   bool          `json:"is_active"`
	CreatedAt  time.Time     `json:"created_at"`
	CreatedBy  sql.NullInt64 `json:"created_by"`
	UpdatedAt  sql.NullTime  `json:"updated_at"`
	UpdatedBy  sql.NullInt64 `json:"updated_by"`
}

type UserUpdate struct {
	FirstName  string    `json:"first_name"`
	MiddleName string    `json:"middle_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Role       string    `json:"role"`
	IsActive   bool      `json:"is_active"`
	UpdatedAt  time.Time `json:"updated_at"`
}
