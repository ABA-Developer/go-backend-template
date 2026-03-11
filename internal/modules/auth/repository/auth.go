package repository

import (
	"context"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/auth/domain"

	"github.com/lib/pq"
)

func (r *repository) ReadDetailUserByEmailQuery(
	ctx context.Context,
	email string,
) (data domain.User, err error) {
	const statement = `
		SELECT
			id, name, full_name, password ,email, active, phone, img_path, img_name
		FROM
			app_user
		WHERE
			email = $1
	`
	err = r.db.QueryRowContext(ctx, statement, email).Scan(
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
) (data domain.User, err error) {
	const statement = `
		SELECT
			id, name, full_name, password ,email, active, phone, img_path, img_name
		FROM
			app_user
		WHERE
			id = $1
	`
	err = r.db.QueryRowContext(ctx, statement, id).Scan(
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
) (data domain.Session, err error) {
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

	err = r.db.QueryRowContext(ctx, statement, id).Scan(
		&data.ID,
		&data.UserID,
		&data.AccessToken,
		&data.AccessTokenExpiredAt,
		&data.RefreshToken,
		&data.RefreshTokenExpiredAt,
	)

	return
}

func (r *repository) CheckPermissionQuery(
	ctx context.Context,
	menuURL constant.MenuKey,
	userID string,
	permissionCodes []string,
) (bool, error) {

	const stmt = `
		WITH RECURSIVE menu_ancestry AS (
			SELECT id, parent_id, url
			FROM app_menu
			WHERE url = $1

			UNION ALL
	
			SELECT m.id, m.parent_id, m.url
			FROM app_menu m
			JOIN menu_ancestry ma ON m.id = ma.parent_id
		),

		child_menu AS (
			SELECT id FROM menu_ancestry WHERE url = $1
		),

		parent_menus AS (
			SELECT id FROM menu_ancestry WHERE url IS DISTINCT FROM $1
		),

		permitted_parents AS (
			SELECT DISTINCT ma.id
			FROM parent_menus ma
			JOIN app_menu_permission amp ON ma.id = amp.menu_id
			JOIN app_role_access ara ON amp.id = ara.menu_permission_id
			JOIN app_user_role aur ON ara.role_id = aur.role_id
			WHERE aur.user_id = $2
		),

		child_permission_check AS (
			SELECT 
			COUNT(DISTINCT amp.code) AS total_permissions
			FROM child_menu cm
			JOIN app_menu_permission amp ON cm.id = amp.menu_id
			JOIN app_role_access ara ON amp.id = ara.menu_permission_id
			JOIN app_user_role aur ON ara.role_id = aur.role_id
			WHERE 
				aur.user_id = $2
				AND amp.code = ANY($3) -- Syarat strict ini HANYA berlaku untuk child/target menu
		)

		SELECT (
			(SELECT COUNT(*) FROM parent_menus) = 0 
			OR 
			(SELECT COUNT(*) FROM permitted_parents) > 0
		)
		AND
		(
		(SELECT total_permissions FROM child_permission_check) = cardinality($3)
		 ) AS has_access;
	`

	var hasAccess bool
	err := r.db.QueryRowContext(
		ctx,
		stmt,
		menuURL.String(),
		userID,
		pq.Array(permissionCodes),
	).Scan(&hasAccess)

	if err != nil {
		return false, err
	}

	return hasAccess, nil
}
func (r *repository) CreateSessionQuery(
	ctx context.Context,
	args domain.Session,
) (err error) {
	const statement = `
		INSERT INTO sessions (
			id, user_id,
			access_token, access_token_expired_at,
			refresh_token, refresh_token_expired_at,
			ip_address, user_agent,
			created_at
		)
		VALUES (
			$1, $2,
			$3, $4,
			$5, $6,
			$7, $8,
			(now() at time zone 'UTC')::TIMESTAMP
		)
	`

	_, err = r.db.ExecContext(ctx, statement,
		args.ID,
		args.UserID,
		args.AccessToken,
		args.AccessTokenExpiredAt,
		args.RefreshToken,
		args.RefreshTokenExpiredAt,
		args.IPAddress,
		args.UserAgent,
	)

	return
}

func (r *repository) UpdateSessionQuery(
	ctx context.Context,
	args domain.Session,
) (err error) {
	const statement = `
		UPDATE
			sessions
		SET
			access_token = $2,
			access_token_expired_at = $3,
			refresh_token = $4,
			refresh_token_expired_at = $5,
			updated_at = (now() at time zone 'UTC')::TIMESTAMP
		WHERE
			id = $1
	`

	_, err = r.db.ExecContext(ctx, statement,
		args.ID,
		args.AccessToken,
		args.AccessTokenExpiredAt,
		args.RefreshToken,
		args.RefreshTokenExpiredAt,
	)

	return
}

func (r *repository) DeleteSessionQuery(
	ctx context.Context,
	id string,
) (err error) {

	const statement = `
		DELETE FROM sessions
		WHERE
			id = $1
	`

	_, err = r.db.ExecContext(ctx, statement, id)

	return
}

func (r *repository) CreateLoginAttemp(
	ctx context.Context,
	args domain.LoginAttemp,
) (err error) {

	const statement = `
		INSERT INTO app_login_attempt (
			password, ip_address, email
		)
		VALUES ($1, $2, $3)
	`
	_, err = r.db.ExecContext(ctx, statement,
		args.Password,
		args.IPAddress,
		args.Email,
	)
	return
}

func (r *repository) CreateLoginRecord(
	ctx context.Context,
	args domain.LoginRecord,
) (err error) {

	const statement = `
		INSERT INTO app_login (
			user_id, access_token, status, ip_address, type
		)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = r.db.ExecContext(ctx, statement,
		args.UserID,
		args.AccessToken,
		args.Status,
		args.IPAddress,
		args.Type,
	)
	return
}
