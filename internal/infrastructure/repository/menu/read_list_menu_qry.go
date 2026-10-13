package repository

import (
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	"context"
)

func (r *Query) ReadListMenuQry(
	ctx context.Context,
	params repoParams.ReadListMenuParams,
) (data []model.Menu, err error) {
	const stmt = `
		WITH RECURSIVE menu_with_parents AS (
			SELECT
				id, parent_id, name, description, url, sort, "group", icon,
				active, display, created_by, created_at, updated_by, updated_at
			FROM
				app_menu
			WHERE
				(CASE WHEN $1::bool THEN
					name ILIKE $2
				ELSE TRUE
				END)

			UNION

			SELECT
				m.id, m.parent_id, m.name, m.description, m.url, m.sort, m."group", m.icon,
				m.active, m.display, m.created_by, m.created_at, m.updated_by, m.updated_at
			FROM
				app_menu m
			JOIN
				menu_with_parents mwp ON m.id = mwp.parent_id
		)
		SELECT DISTINCT * FROM menu_with_parents
		ORDER BY
			sort ASC;
	`

	rows, err := r.db.QueryContext(ctx, stmt,
		params.SetSearch,
		params.Search,
	)
	if err != nil {
		return
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var m model.Menu
		if err = rows.Scan(
			&m.ID, &m.ParentID, &m.Name, &m.Description, &m.URL, &m.Sort,
			&m.Group, &m.Icon, &m.Active, &m.Display,
			&m.CreatedBy, &m.CreatedAt, &m.UpdatedBy, &m.UpdatedAt,
		); err != nil {
			return
		}
		data = append(data, m)
	}

	return
}
