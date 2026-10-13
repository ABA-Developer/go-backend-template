package repository

import (
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	"context"
)

func (r *Query) ReadUsersQry(
	ctx context.Context,
	args repoParams.ReadListUserParams,
) (data []model.User, err error) {
	const stmt = `
		SELECT
			au.id,
			au.full_name,
			au.active,
			ar.name AS role,
			ar.id AS role_id
		FROM
			app_user au
		JOIN
			app_user_role aur ON au.id = aur.user_id
		JOIN
			app_role ar ON aur.role_id = ar.id
		WHERE
			(CASE WHEN $1::bool THEN(
				full_name ILIKE $2
			) ELSE TRUE END)
		ORDER BY
			(CASE WHEN $3 = 'full_name ASC' THEN au.full_name END) ASC,
			(CASE WHEN $3 = 'full_name DESC' THEN au.full_name END) DESC,
			(CASE WHEN $3 = 'role ASC' THEN ar.name END) ASC,
			(CASE WHEN $3 = 'role DESC' THEN ar.name END) DESC,
			(CASE WHEN $3 = 'active ASC' THEN au.active END) ASC,
			(CASE WHEN $3 = 'active DESC' THEN au.active END) DESC,
			(CASE WHEN $3 = 'created_at ASC' THEN au.created_at END) ASC,
			(CASE WHEN $3 = 'created_at DESC' THEN au.created_at END) DESC
		LIMIT $4
		OFFSET $5
	`

	rows, err := r.db.QueryContext(ctx, stmt,
		args.SetSearch,
		args.Search,
		args.Order,
		args.Limit,
		args.Offset,
	)
	if err != nil {
		return
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var u model.User

		if err = rows.Scan(
			&u.ID,
			&u.FullName,
			&u.Active,
			&u.Role,
			&u.RoleID,
		); err != nil {
			return
		}

		data = append(data, u)
	}

	return
}
