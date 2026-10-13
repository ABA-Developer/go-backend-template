package repository

import (
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"context"
)

func (r *Query) UpdateUserQry(
	ctx context.Context,
	args repoParams.UpdateUserParams,
) (err error) {
	const statement = `
		UPDATE app_user
		SET
		name = $1,
		full_name = $2,
		email = $3,
		phone = $4,
		active = $5,
		updated_by = $6,
		updated_at = (now() at time zone 'UTC')::TIMESTAMP
		WHERE id = $7
	`

	_, err = r.db.ExecContext(ctx, statement,
		args.Name,
		args.FullName,
		args.Email,
		args.Phone,
		args.Active,
		args.UpdatedBy,
		args.ID,
	)

	return
}
