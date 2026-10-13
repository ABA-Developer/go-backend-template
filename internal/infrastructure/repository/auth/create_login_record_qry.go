package repository

import (
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"context"
)

func (r *Query) CreateLoginRecordQry(
	ctx context.Context,
	args repoParams.LoginRecord,
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
