package repository

import (
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	"context"
)

func (r *Query) ReadUserProfileQry(
	ctx context.Context,
	id string,
) (data model.User, err error) {
	const statement = `
		SELECT
			au.id,
			au.name,
			au.full_name,
			au.email,
			au.phone,
			au.active,
			au.img_path,
			au.img_name,
			ar.name AS role,
			ar.id AS role_id
		FROM
			app_user au
		JOIN
			app_user_role aur ON au.id = aur.user_id
		JOIN
			app_role ar ON aur.role_id = ar.id
		WHERE
			au.id = $1
	`

	err = r.db.QueryRowContext(ctx, statement, id).Scan(
		&data.ID,
		&data.Name,
		&data.FullName,
		&data.Email,
		&data.Phone,
		&data.Active,
		&data.ImgPath,
		&data.ImgName,
		&data.Role,
		&data.RoleID,
	)

	return
}
