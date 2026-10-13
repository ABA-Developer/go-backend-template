package repository

import (
	"github.com/ABA-Developer/go-backend-template/internal/domain/model"
	"context"
)

func (r *Query) ReadDetailUserByIdQry(
	ctx context.Context,
	id string,
) (data model.User, err error) {
	const statement = `
		SELECT
			id, name, full_name, password ,email, active, phone, img_path, img_name
		FROM
			app_user
		WHERE
			id = $1
	`
	err = r.db.QueryRowContext(ctx, statement, id).Scan(
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
