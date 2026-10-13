package repository

import (
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"context"
)

func (r *Query) ReadCountMenuQry(
	ctx context.Context,
	params repoParams.ReadListMenuParams,
) (count int64, err error) {
	const stmt = `
		SELECT
			COUNT(*)
		FROM
			app_menu
		WHERE
			(CASE WHEN $1::bool THEN(
				name ILIKE $2
				OR description ILIKE $2
				OR url ILIKE $2
			) ELSE TRUE END)
	`
	err = r.db.QueryRowContext(ctx, stmt,
		params.SetSearch,
		params.Search,
	).Scan(&count)
	return
}
