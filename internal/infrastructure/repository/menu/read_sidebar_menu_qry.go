package repository

import (
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	"context"
)

func (r *Query) ReadSidebarMenuQry(
	ctx context.Context,
	userID string,
) (data []model.Menu, err error) {

	const stmt = `
		SELECT DISTINCT
			m.id, m.parent_id, m.name, m.url, m.sort,
			m."group", m.icon, m.active, m.display
		FROM
			app_menu m
		JOIN
			app_menu_permission amp ON m.id = amp.menu_id
		JOIN
			app_role_access ara ON amp.id = ara.menu_permission_id
		JOIN
			app_user_role aur ON ara.role_id = aur.role_id
		WHERE
			aur.user_id = $1
			AND amp.code = 'R'
			AND m.display = true
			AND m.active = true
		ORDER BY
		m."group" ASC,
			m.sort ASC;
	`

	rows, err := r.db.QueryContext(ctx, stmt, userID)
	if err != nil {
		return
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var m model.Menu
		if err = rows.Scan(
			&m.ID, &m.ParentID, &m.Name, &m.URL, &m.Sort,
			&m.Group, &m.Icon, &m.Active, &m.Display,
		); err != nil {
			return
		}
		data = append(data, m)
	}
	return
}
