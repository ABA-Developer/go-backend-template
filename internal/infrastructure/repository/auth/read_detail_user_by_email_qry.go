package repository

import (
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	"context"
)

func (r *Query) ReadDetailUserByEmailQry(
	ctx context.Context,
	email string,
) (data model.User, err error) {
	const statement = `
		SELECT
			id, name, full_name, password ,email, active, phone, img_path, img_name
		FROM
			app_user
		WHERE
			email = $1
	`
	err = r.db.QueryRowContext(ctx, statement, email).Scan(
		&data.ID,
		&data.Name,
		&data.FullName,
		&data.Password,
		&data.Email,
		&data.Active,
		&data.Phone,
		&data.ImgPath,
		&data.ImgName,
	)

	return
}
