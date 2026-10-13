package repository

import (
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	"context"
)

func (r *Query) ReadParentMenuQry(
	ctx context.Context,
) (data []model.Menu, err error) {
	const stmt = `
		SELECT
			id, name, "group"
		FROM
			app_menu
		WHERE
			parent_id IS NULL
		ORDER BY
			sort ASC;
	`
	rows, err := r.db.QueryContext(ctx, stmt)
	if err != nil {
		return
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var m model.Menu
		if err = rows.Scan(
			&m.ID, &m.Name, &m.Group,
		); err != nil {
			return
		}
		data = append(data, m)
	}
	return
}
