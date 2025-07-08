package user

import (
	"be-dashboard-nba/pkg/entities"
	"context"
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	CreateUser(user *entities.User) (*entities.User, error)
	ReadUser(id int64) (*entities.User, error)
	UpdateUser(user *entities.User) (*entities.User, error)
	DeleteUser(id int64) error
	GetAllUser() ([]entities.User, error)
	GetUserById(id int64) (*entities.User, error)
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

func (r *repository) CreateUser(user *entities.User) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
		INSERT INTO users(first_name, middle_name, last_name, email, password, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW()) RETURNING id
	`
	// Hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	user.Password = string(hashedPassword)
	user.IsActive = true

	// Perform query
	row := r.Db.QueryRowContext(ctx, query, user.FirstName, user.MiddleName, user.LastName, user.Email, user.Password, user.Role, user.IsActive)
	err := row.Scan(&user.ID)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (r *repository) ReadUser(id int64) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
		SELECT id, first_name, middle_name, last_name, email, role, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	row := r.Db.QueryRowContext(ctx, query, id)

	var user entities.User
	err := row.Scan(&user.ID, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) UpdateUser(user *entities.User) (*entities.User, error) {
	existingUser, err := r.GetUserById(int64(user.ID))
	if err != nil {
		return nil, err
	}

	// Update only for non empty value
	if user.FirstName == "" {
		user.FirstName = existingUser.FirstName
	}
	if user.MiddleName == "" {
		user.MiddleName = existingUser.MiddleName
	}
	if user.LastName == "" {
		user.LastName = existingUser.LastName
	}
	if user.Email == "" {
		user.Email = existingUser.Email
	}
	if user.Role == "" {
		user.Role = existingUser.Role
	}
	if user.IsActive == existingUser.IsActive {
		user.IsActive = existingUser.IsActive
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
		UPDATE users SET first_name = $1, middle_name = $2, last_name = $3, email = $4, role = $5, is_active = $6, updated_at = NOW()
		WHERE id = $7
		RETURNING id, first_name, middle_name, last_name, email, role, is_active, created_at, updated_at;
	`
	row := r.Db.QueryRowContext(ctx, query, user.FirstName, user.MiddleName, user.LastName, user.Email, user.Role, user.IsActive, user.ID)

	var updatedUser entities.User
	err = row.Scan(
		&updatedUser.ID,
		&updatedUser.FirstName,
		&updatedUser.MiddleName,
		&updatedUser.LastName,
		&updatedUser.Email,
		&updatedUser.Role,
		&updatedUser.IsActive,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &updatedUser, nil
}

func (r *repository) DeleteUser(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
		DELETE FROM users WHERE id = $1
	`
	row, err := r.Db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	deletedRow, err := row.RowsAffected()
	if err != nil {
		return err
	}

	if deletedRow == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}
	return nil
}

func (r *repository) GetAllUser() ([]entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
		SELECT id, first_name, middle_name, last_name, email, role, is_active, created_at, updated_at
		FROM users
	`

	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var user entities.User
		err := rows.Scan(&user.ID, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *repository) GetUserById(id int64) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
		SELECT id, first_name, middle_name, last_name, email, role, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	row := r.Db.QueryRowContext(ctx, query, id)

	var user entities.User
	err := row.Scan(&user.ID, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetUserByEmail(email string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
		SELECT id, first_name, middle_name, last_name, email, role, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	row := r.Db.QueryRowContext(ctx, query, email)

	var user entities.User
	err := row.Scan(&user.ID, &user.FirstName, &user.MiddleName, &user.LastName, &user.Email, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
