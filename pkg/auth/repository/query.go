package repository

import (
	"context"
	"time"
)

type SessionParams struct {
	ID                    string
	UserID                string
	AccessToken           string
	AccessTokenExpiredAt  time.Time
	RefreshToken          string
	RefreshTokenExpiredAt time.Time
	IPAddress             string
	UserAgent             string
}

func (r *repository) CreateSessionQuery(
	ctx context.Context,
	args SessionParams,
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

	_, err = r.DB.ExecContext(ctx, statement,
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
	args SessionParams,
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

	_, err = r.DB.ExecContext(ctx, statement,
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

	_, err = r.DB.ExecContext(ctx, statement, id)

	return
}

type LoginAttempParams struct {
	Email     string
	Password  string
	IPAddress string
}

func (r *repository) CreateLoginAttemp(
	ctx context.Context,
	args LoginAttempParams,
) (err error) {

	const statement = `
		INSERT INTO app_login_attempt (
			password, ip_address, email
		)
		VALUES ($1, $2, $3)
	`
	_, err = r.DB.ExecContext(ctx, statement,
		args.Password,
		args.IPAddress,
		args.Email,
	)
	return
}

type LoginRecord struct {
	UserID      string
	AccessToken string
	Status      string
	IPAddress   string
	Type        string
}

func (r *repository) CreateLoginRecord(
	ctx context.Context,
	args LoginRecord,
) (err error) {

	const statement = `
		INSERT INTO app_login (
			user_id, access_token, status, ip_address, type
		)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = r.DB.ExecContext(ctx, statement,
		args.UserID,
		args.AccessToken,
		args.Status,
		args.IPAddress,
		args.Type,
	)
	return
}
