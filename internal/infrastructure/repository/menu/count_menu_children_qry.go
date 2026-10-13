package repository

import (
	"context"
)

func (r *Query) CountMenuChildrenQry(ctx context.Context, menuID int) (int, error) {
	const stmt = `SELECT COUNT(*) FROM app_menu WHERE parent_id = $1`
	var count int
	err := r.db.QueryRowContext(ctx, stmt, menuID).Scan(&count)
	return count, err
}
