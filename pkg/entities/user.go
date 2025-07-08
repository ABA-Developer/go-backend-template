package entities

import "time"

type User struct {
	ID         int       `json:"id"`
	FirstName  string    `json:"first_name,omitempty"`
	MiddleName string    `json:"middle_name,omitempty"`
	LastName   string    `json:"last_name,omitempty"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Role       string    `json:"role"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserUpdate struct {
	FirstName  string    `json:"first_name,omitempty"`
	MiddleName string    `json:"middle_name,omitempty"`
	LastName   string    `json:"last_name,omitempty"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Role       string    `json:"role"`
	IsActive   bool      `json:"is_active"`
	UpdatedAt  time.Time `json:"updated_at"`
}
