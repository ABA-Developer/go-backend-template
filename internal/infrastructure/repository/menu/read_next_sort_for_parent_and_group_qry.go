package repository

import (
	"context"
)

func (r *Query) ReadNextSortForParentAndGroupQry(ctx context.Context, parentID int32, group string) (int, error) {
	const stmt = `SELECT COALESCE(MAX(sort), -1) + 1 FROM app_menu WHERE parent_id = $1 AND "group" = $2`
	var nextSort int
	err := r.db.QueryRowContext(ctx, stmt, parentID, group).Scan(&nextSort)
	return nextSort, err
}
