package repository

import (
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"context"
	"fmt"
)

func (r *Query) ReadRolesCountQry(ctx context.Context,
	args repoParams.ReadRolesParams) (count int, err error) {
	const stmt = `
			SELECT COUNT(*)
			FROM
				app_role
			WHERE
				(CASE WHEN $1::bool THEN(
				name ILIKE $2
				) ELSE TRUE END)
		`
	err = r.db.QueryRowContext(ctx, stmt,
		args.SetSearch,
		args.Search,
	).Scan(&count)
	if err != nil {
		fmt.Printf("Error in ReadRolesCountQry: %v\n", err)
	}
	return
}
