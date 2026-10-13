package repository

import (
	"context"
)

func (r *Query) ReadSortForGroupQry(ctx context.Context, group string) (int, error) {
	var nextSort int
	const stmt = `
		SELECT COALESCE(MAX(sort), -1) + 1
		FROM app_menu
		WHERE parent_id IS NULL AND "group" = $1
	`
	err := r.db.QueryRowContext(ctx, stmt, group).Scan(&nextSort)
	return nextSort, err
}
