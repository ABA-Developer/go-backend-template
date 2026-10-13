package repository

import (
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"context"
)

func (r *Query) UpdateMenuQry(
	ctx context.Context,
	params repoParams.UpdateMenuParams,
) (err error) {
	const stmt = `
        UPDATE app_menu
        SET
            name = $2,
            parent_id = $3,
            description = $4,
            url = $5,
            sort = $6,
            "group" = $7,
            icon = $8,
            active = $9,
            display = $10,
            updated_by = $11,
            updated_at = NOW()
        WHERE
            id = $1;
    `

	_, err = r.db.ExecContext(ctx, stmt,
		params.ID,
		params.Name,
		params.ParentID,
		params.Description,
		params.URL,
		params.Sort,
		params.Group,
		params.Icon,
		params.Active,
		params.Display,
		params.UpdatedBy,
	)
	return
}
