package repository

import (
	"context"

	"be-dashboard-nba/pkg/entities"
)

type ReadListUserParams struct {
	SetSearch bool
	Search    string
	Order     string
	Limit     int
	Offset    int
}

func (q *Query) ReadListUserQuery(
	ctx context.Context,
	args ReadListUserParams,
) (data []entities.User, err error) {
	const stmt = `
		SELECT
			id, first_name, middle_name, last_name, 
			email, role, is_active,
			created_at, created_by, updated_at, updated_by
		FROM
			users
		WHERE
			(CASE WHEN $1::bool THEN(
				email ILIKE $2
				OR first_name ILIKE $2
				OR middle_name ILIKE $2
				OR last_name ILIKE $2
				OR role ILIKE $2
			) ELSE TRUE END)
		ORDER BY
			(CASE WHEN $3 = 'first_name ASC' THEN first_name END) ASC,
			(CASE WHEN $3 = 'first_name DESC' THEN first_name END) DESC,
			(CASE WHEN $3 = 'middle_name ASC' THEN middle_name END) ASC,
			(CASE WHEN $3 = 'middle_name DESC' THEN middle_name END) DESC,
			(CASE WHEN $3 = 'last_name ASC' THEN last_name END) ASC,
			(CASE WHEN $3 = 'last_name DESC' THEN last_name END) DESC,
			(CASE WHEN $3 = 'email ASC' THEN email END) ASC,
			(CASE WHEN $3 = 'email DESC' THEN email END) DESC,
			(CASE WHEN $3 = 'role ASC' THEN role END) ASC,
			(CASE WHEN $3 = 'role DESC' THEN role END) DESC,
			(CASE WHEN $3 = 'is_active ASC' THEN is_active END) ASC,
			(CASE WHEN $3 = 'is_active DESC' THEN is_active END) DESC,
			(CASE WHEN $3 = 'created_at ASC' THEN created_at END) ASC,
			(CASE WHEN $3 = 'created_at DESC' THEN created_at END) DESC,
			(CASE WHEN $3 = 'updated_at ASC' THEN updated_at END) ASC,
			(CASE WHEN $3 = 'updated_at DESC' THEN updated_at END) DESC
		LIMIT $4
		OFFSET $5
	`

	rows, err := q.db.QueryContext(ctx, stmt,
		args.SetSearch,
		args.Search,
		args.Order,
		args.Limit,
		args.Offset,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var u entities.User

		if err = rows.Scan(
			&u.ID,
			&u.FirstName,
			&u.MiddleName,
			&u.LastName,
			&u.Email,
			&u.Role,
			&u.IsActive,
			&u.CreatedAt,
			&u.CreatedBy,
			&u.UpdatedAt,
			&u.UpdatedBy,
		); err != nil {
			return
		}

		data = append(data, u)
	}

	return
}

func (q *Query) GetCountUserQuery(
	ctx context.Context,
	args ReadListUserParams,
) (count int64, err error) {
	const stmt = `
		SELECT
			COUNT(*)
		FROM
			users
		WHERE
			(CASE WHEN $1::bool THEN(
				email ILIKE $2
				OR first_name ILIKE $2
				OR middle_name ILIKE $2
				OR last_name ILIKE $2
				OR role ILIKE $2
			) ELSE TRUE END)
	`

	err = q.db.QueryRowContext(ctx, stmt,
		args.SetSearch,
		args.Search,
	).Scan(&count)

	return
}

func (q *Query) ReadDetailUserQuery(
	ctx context.Context,
	id int64,
) (data entities.User, err error) {
	const statement = `
		SELECT	
			id, first_name, middle_name, last_name, 
			email, role, is_active,
			created_at, created_by, updated_at, updated_by
		FROM
			users
		WHERE
			id = $1
	`

	err = q.db.QueryRowContext(ctx, statement, id).Scan(
		&data.ID,
		&data.FirstName,
		&data.MiddleName,
		&data.LastName,
		&data.Email,
		&data.Role,
		&data.IsActive,
		&data.CreatedAt,
		&data.CreatedBy,
		&data.UpdatedAt,
		&data.UpdatedBy,
	)

	return
}

func (q *Query) IsUserEmailExistsQuery(
	ctx context.Context,
	email string,
) (exists bool, err error) {
	const statement = `
		SELECT EXISTS (
			SELECT
				1
			FROM
				users
			WHERE
				email = $1
		)
	`

	err = q.db.QueryRowContext(ctx, statement, email).Scan(&exists)

	return
}

func (q *Query) IsUpdateUserEmailExistsQuery(
	ctx context.Context,
	email, id string,
) (exists bool, err error) {
	statement := `
		SELECT EXISTS (
			SELECT
				1
			FROM
				users
			WHERE
				email = $1
				AND id != $2
		)
	`

	err = q.db.QueryRowContext(ctx, statement, email, id).Scan(&exists)

	return
}
