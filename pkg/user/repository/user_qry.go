package repository

import (
	"context"
)

type CreateUserParams struct {
	FirstName  string
	MiddleName string
	LastName   string
	Email      string
	Password   string
	Role       string
	IsActive   bool
	CreatedBy  int64
}

func (q *Query) CreateUserQuery(
	ctx context.Context,
	args CreateUserParams,
) (err error) {
	const statement = `
		INSERT INTO users (
			first_name, middle_name, last_name, 
			email, password, role, is_active,
			created_by, created_at
		)
		VALUES (
			$1, $2, $3, 
			$4, $5, $6, $7,
			$8, (now() at time zone 'UTC')::TIMESTAMP
		)
	`

	_, err = q.db.ExecContext(ctx, statement,
		args.FirstName,
		args.MiddleName,
		args.LastName,
		args.Email,
		args.Password,
		args.Role,
		args.IsActive,
		args.CreatedBy,
	)

	return
}

type UpdateUserParams struct {
	ID         int64
	FirstName  string
	MiddleName string
	LastName   string
	Email      string
	Password   string
	Role       string
	IsActive   bool
	UpdatedBy  int64
}

func (q *Query) UpdateUserQuery(
	ctx context.Context,
	args UpdateUserParams,
) (err error) {
	const statement = `
		UPDATE
			users
		SET
			first_name = CASE WHEN $2 <> '' THEN $2 ELSE first_name END,
			middle_name = CASE WHEN $3 <> '' THEN $3 ELSE middle_name END,
			last_name = CASE WHEN $4 <> '' THEN $4 ELSE last_name END,
			email = CASE WHEN $5 <> '' THEN $5 ELSE email END,
			password = CASE WHEN $6 <> '' THEN $6 ELSE password END,
			role = CASE WHEN $7 <> '' THEN $7 ELSE role END,
			is_active = $8,
			updated_by = $9,
			updated_at = (now() at time zone 'UTC')::TIMESTAMP
		WHERE
			id = $1
	`

	_, err = q.db.ExecContext(ctx, statement,
		args.ID,
		args.FirstName,
		args.MiddleName,
		args.LastName,
		args.Email,
		args.Password,
		args.Role,
		args.IsActive,
		args.UpdatedBy,
	)

	return
}

func (q *Query) DeleteUserQuery(
	ctx context.Context,
	id int64,
) (err error) {
	const statement = `
		DELETE FROM users
		WHERE
			id = $1
	`

	_, err = q.db.ExecContext(ctx, statement, id)

	return
}
