package repository

import (
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"context"
)

func (r *Query) ReadMenuPermissionCountQry(ctx context.Context,
	args repoParams.ReadMenuPermissionParams) (count int, err error) {
	const stmt = `
			SELECT COUNT(*)
			FROM
				app_menu_permission
			WHERE
				(CASE WHEN $1::bool THEN(
				action_name ILIKE $2
				) ELSE TRUE END)
				AND menu_id = $3
		`

	err = r.db.QueryRowContext(ctx, stmt,
		args.SetSearch,
		args.Search,
		args.MenuID,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return
}
