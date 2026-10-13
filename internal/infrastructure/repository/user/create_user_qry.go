package repository

import (
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"context"
)

func (r *Query) CreateUserQry(
	ctx context.Context,
	args repoParams.CreateUserParams,
) (userID string, err error) {
	const statement = `
		INSERT INTO app_user (
			name, full_name, email, password, active, phone, created_by, created_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, (now() at time zone 'UTC')::TIMESTAMP
		)
		RETURNING id
	`

	err = r.db.QueryRowContext(ctx, statement,
		args.Name,
		args.FullName,
		args.Email,
		args.Password,
		args.Active,
		args.Phone,
		args.CreatedBy,
	).Scan(&userID)

	return
}
