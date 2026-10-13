package repository

import (
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	"context"
)

func (r *Query) ReadRoleByIDQry(ctx context.Context, roleID int) (data model.Role, err error) {
	const stmt = `
		SELECT
			id, name, code, description
		FROM
			app_role
		WHERE
			id = $1
	`
	err = r.db.QueryRowContext(ctx, stmt, roleID).Scan(
		&data.ID,
		&data.Name,
		&data.Code,
		&data.Description,
	)

	return
}
