package repository

import (
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	"context"
)

func (r *Query) ReadMenuPermissionListQry(
	ctx context.Context,
	args repoParams.ReadMenuPermissionParams,
) (data []model.MenuPermission, err error) {

	const stmt = `
	SELECT
		id, menu_id, code, action_name
	FROM
		app_menu_permission
	WHERE
		(CASE WHEN $1::bool THEN(
			action_name ILIKE $2
		) ELSE TRUE END)
		 AND menu_id = $6
	ORDER BY
		(CASE WHEN $3 = 'action_name ASC' THEN action_name END)
		ASC,
		(CASE WHEN $3 = 'action_name DESC' THEN action_name END)
		DESC,
		(CASE WHEN $3 = 'code ASC' THEN code END)
		ASC,
		(CASE WHEN $3 = 'code DESC' THEN code END)
		DESC,
		(CASE WHEN $3 = 'id ASC' THEN id END)
		ASC,
		(CASE WHEN $3 = 'id DESC' THEN id END)
		DESC
	LIMIT $4
	OFFSET $5
	`
	rows, err := r.db.QueryContext(ctx, stmt,
		args.SetSearch,
		args.Search,
		args.Order,
		args.Limit,
		args.Offset,
		args.MenuID,
	)
	if err != nil {
		return
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var mp model.MenuPermission

		if err = rows.Scan(
			&mp.ID,
			&mp.MenuID,
			&mp.Code,
			&mp.ActionName,
		); err != nil {
			return
		}

		data = append(data, mp)
	}

	return
}
