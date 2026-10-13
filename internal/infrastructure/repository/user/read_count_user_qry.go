package repository

import (
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"context"
)

func (r *Query) ReadCountUserQry(
	ctx context.Context,
	args repoParams.ReadListUserParams,
) (count int, err error) {
	const stmt = `
		SELECT
			COUNT(*)
		FROM
			app_user
		WHERE
			(CASE WHEN $1::bool THEN(
				full_name ILIKE $2
			) ELSE TRUE END)
	`

	err = r.db.QueryRowContext(ctx, stmt,
		args.SetSearch,
		args.Search,
	).Scan(&count)

	return
}
