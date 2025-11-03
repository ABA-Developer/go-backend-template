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

func (r *repository) ReadListUserQuery(
	ctx context.Context,
	args ReadListUserParams,
) (data []entities.User, err error) {
	const stmt = `
		SELECT
			id, name, full_name, email, password, active,
			phone, img_path, img_name,
			created_at, created_by, updated_at, updated_by
		FROM
			app_user
		WHERE
			(CASE WHEN $1::bool THEN(
				email ILIKE $2
				OR name ILIKE $2
				OR full_name ILIKE $2
				OR phone ILIKE $2
			) ELSE TRUE END)
		ORDER BY
			(CASE WHEN $3 = 'name ASC' THEN name END) ASC,
			(CASE WHEN $3 = 'name DESC' THEN name END) DESC,
			(CASE WHEN $3 = 'full_name ASC' THEN full_name END) ASC,
			(CASE WHEN $3 = 'full_name DESC' THEN full_name END) DESC,
			(CASE WHEN $3 = 'email ASC' THEN email END) ASC,
			(CASE WHEN $3 = 'email DESC' THEN email END) DESC,
			(CASE WHEN $3 = 'active ASC' THEN active END) ASC,
			(CASE WHEN $3 = 'active DESC' THEN active END) DESC,
			(CASE WHEN $3 = 'created_at ASC' THEN created_at END) ASC,
			(CASE WHEN $3 = 'created_at DESC' THEN created_at END) DESC,
			(CASE WHEN $3 = 'updated_at ASC' THEN updated_at END) ASC,
			(CASE WHEN $3 = 'updated_at DESC' THEN updated_at END) DESC
		LIMIT $4
		OFFSET $5
	`

	rows, err := r.db.QueryContext(ctx, stmt,
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
			&u.Name,
			&u.FullName,
			&u.Email,
			&u.Password,
			&u.Active,
			&u.Phone,
			&u.ImgPath,
			&u.ImgName,
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

func (r *repository) GetCountUserQuery(
	ctx context.Context,
	args ReadListUserParams,
) (count int64, err error) {
	const stmt = `
		SELECT
			COUNT(*)
		FROM
			app_user
		WHERE
			(CASE WHEN $1::bool THEN(
				email ILIKE $2
				OR name ILIKE $2
				OR full_name ILIKE $2
				OR phone ILIKE $2
			) ELSE TRUE END)
	`

	err = r.db.QueryRowContext(ctx, stmt,
		args.SetSearch,
		args.Search,
	).Scan(&count)

	return
}

func (r *repository) ReadDetailUserQuery(
	ctx context.Context,
	id string,
) (data entities.User, err error) {
	const statement = `
		SELECT 	
			id, name, full_name, email, active,
			phone, img_path, img_name
		FROM
			app_user
		WHERE
			id = $1
	`

	err = r.db.QueryRowContext(ctx, statement, id).Scan(
		&data.ID,
		&data.Name,
		&data.FullName,
		&data.Email,
		&data.Active,
		&data.Phone,
		&data.ImgPath,
		&data.ImgName,
	)

	return
}

func (r *repository) IsUserEmailExistsQuery(
	ctx context.Context,
	email string,
) (exists bool, err error) {
	const statement = `
		SELECT EXISTS (
			SELECT
				1
			FROM
				app_user
			WHERE
				email = $1
		)
	`

	err = r.db.QueryRowContext(ctx, statement, email).Scan(&exists)

	return
}

func (r *repository) IsUpdateUserEmailExistsQuery(
	ctx context.Context,
	email, id string,
) (exists bool, err error) {
	statement := `
		SELECT EXISTS (
			SELECT
				1
			FROM
				app_user
			WHERE
				email = $1
				AND id != $2
		)
	`

	err = r.db.QueryRowContext(ctx, statement, email, id).Scan(&exists)

	return
}
