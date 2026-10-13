package repository

import (
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	"context"
)

func (r *Query) ReadRolesQry(
	ctx context.Context,
	args repoParams.ReadRolesParams,
) (data []model.Role, err error) {
	const stmt = `
		SELECT
			id, name, code ,description
		FROM
			app_role
		WHERE
			(CASE WHEN $1::bool THEN(
				name ILIKE $2
			) ELSE TRUE END)
		ORDER BY
			(CASE WHEN $3 = 'name ASC' THEN name END)
			ASC,
			(CASE WHEN $3 = 'name DESC' THEN name END)
			DESC,
			(CASE WHEN $3 = 'description ASC' THEN description END)
			ASC,
			(CASE WHEN $3 = 'description DESC' THEN description END)
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
	)
	if err != nil {
		return
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var r model.Role

		if err = rows.Scan(
			&r.ID,
			&r.Name,
			&r.Code,
			&r.Description,
		); err != nil {
			return
		}

		data = append(data, r)
	}

	return
}
