package repository

import (
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	"context"
)

func (r *Query) ReadMenuByIDQry(ctx context.Context, menuID int) (data model.Menu, err error) {
	const stmt = `
		SELECT
			id, parent_id, name, description, url, sort, "group", icon,
			active, display
		FROM
			app_menu
		WHERE
			id = $1
	`
	err = r.db.QueryRowContext(ctx, stmt, menuID).Scan(
		&data.ID, &data.ParentID, &data.Name, &data.Description, &data.URL, &data.Sort,
		&data.Group, &data.Icon, &data.Active, &data.Display,
	)
	return
}
