package repository

import (
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"context"
)

func (r *Query) CreateLoginAttempQry(
	ctx context.Context,
	args repoParams.LoginAttempParams,
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
