package repository

import (
	"context"
	"database/sql"
)

type CreateUserParams struct {
	Name      string
	FullName  string
	Email     string
	Password  string
	Active    bool
	Phone     sql.NullString
	ImgPath   sql.NullString
	ImgName   sql.NullString
	CreatedBy string
}

func (r *repository) CreateUserQuery(
	ctx context.Context,
	args CreateUserParams,
) (err error) {
	const statement = `
		INSERT INTO app_user (
			name, full_name, email, password, active,
			phone, img_path, img_name,
			created_by, created_at
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8,
			$9, (now() at time zone 'UTC')::TIMESTAMP
		)
	`

	// DIUBAH: Urutan argumen
	_, err = r.db.ExecContext(ctx, statement,
		args.Name,
		args.FullName,
		args.Email,
		args.Password,
		args.Active,
		args.Phone,
		args.ImgPath,
		args.ImgName,
		args.CreatedBy,
	)

	return
}

type UpdateUserParams struct {
	ID        string
	Name      string
	FullName  string
	Email     string
	Password  string
	Phone     string
	Active    bool
	ImgPath   string
	ImgName   string
	UpdatedBy string
}

func (r *repository) UpdateUserQuery(
	ctx context.Context,
	args UpdateUserParams,
) (err error) {
	const statement = `
		UPDATE
			app_user
		SET
			name = CASE WHEN $2 <> '' THEN $2 ELSE name END,
			full_name = CASE WHEN $3 <> '' THEN $3 ELSE full_name END,
			email = CASE WHEN $4 <> '' THEN $4 ELSE email END,
			password = CASE WHEN $5 <> '' THEN $5 ELSE password END,
			phone = CASE WHEN $6 <> '' THEN $6 ELSE phone END,
			active = $7,
			img_path = CASE WHEN $8 <> '' THEN $8 ELSE img_path END,
			img_name = CASE WHEN $9 <> '' THEN $9 ELSE img_name END,
			updated_by = $10,
			updated_at = (now() at time zone 'UTC')::TIMESTAMP
		WHERE
			id = $1
	`

	_, err = r.db.ExecContext(ctx, statement,
		args.ID,
		args.Name,
		args.FullName,
		args.Email,
		args.Password,
		args.Phone,
		args.Active,
		args.ImgPath,
		args.ImgName,
		args.UpdatedBy,
	)

	return
}

func (r *repository) DeleteUserQuery(
	ctx context.Context,
	id string,
) (err error) {
	const statement = `
		DELETE FROM app_user
		WHERE
			id = $1
	`

	_, err = r.db.ExecContext(ctx, statement, id)

	return
}
