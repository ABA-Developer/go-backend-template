package repository

import (
	"be-dashboard-nba/pkg/entities"
	"context"
)

func (r *repository) ReadDetailUserByEmailQuery(
	ctx context.Context,
	email string,
) (data entities.User, err error) {
	const statement = `
		SELECT
			id, name, full_name, password ,email, active, phone, img_path, img_name
		FROM
			app_user
		WHERE
			email = $1
	`
	err = r.DB.QueryRowContext(ctx, statement, email).Scan(
		&data.ID,
		&data.Name,
		&data.FullName,
		&data.Password,
		&data.Email,
		&data.Active,
		&data.Phone,
		&data.ImgPath,
		&data.ImgName,
	)

	return
}

func (r *repository) ReadDetailUserByIdQuery(
	ctx context.Context,
	id string,
) (data entities.User, err error) {
	const statement = `
		SELECT
			id, name, full_name, password ,email, active, phone, img_path, img_name
		FROM
			app_user
		WHERE
			id = $1
	`
	err = r.DB.QueryRowContext(ctx, statement, id).Scan(
		&data.ID,
		&data.Name,
		&data.FullName,
		&data.Password,
		&data.Email,
		&data.Active,
		&data.Phone,
		&data.ImgPath,
		&data.ImgName,
	)

	return
}

func (r *repository) ReadDetailSessionQuery(
	ctx context.Context,
	id string,
) (data entities.Session, err error) {
	const statement = `
		SELECT
			id, user_id,
			access_token, access_token_expired_at,
			refresh_token, refresh_token_expired_at
		FROM
			sessions
		WHERE
			id = $1
	`

	err = r.DB.QueryRowContext(ctx, statement, id).Scan(
		&data.ID,
		&data.UserID,
		&data.AccessToken,
		&data.AccessTokenExpiredAt,
		&data.RefreshToken,
		&data.RefreshTokenExpiredAt,
	)

	return
}
