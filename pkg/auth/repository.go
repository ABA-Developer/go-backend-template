package auth

import (
	"be-dashboard-nba/pkg/entities"
	"context"
	"database/sql"
	"time"
)

type Repository interface {
	GetUserByEmail(email string) (*entities.User, error)
}

type repository struct {
	Db *sql.DB
}

func NewRepo(db *sql.DB) Repository {
	return &repository{
		Db: db,
	}
}

func (r *repository) GetUserByEmail(email string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
		SELECT id, first_name, middle_name, last_name, email, password, role, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	row := r.Db.QueryRowContext(ctx, query, email)

	var user entities.User
	err := row.Scan(&user.ID, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email, &user.Password, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
