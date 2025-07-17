package repository

import (
	"context"

	"be-dashboard-nba/pkg/entities"
)

func (q *Query) ReadDetailUserByEmailQuery(
	ctx context.Context,
	email string,
) (data entities.User, err error) {
	const statement = `
		SELECT
			id, first_name, middle_name, last_name, email, password, role, 
			created_at, created_by, updated_at, updated_by
		FROM
			users
		WHERE
			email = $1
	`

	err = q.db.QueryRowContext(ctx, statement, email).Scan(
		&data.ID,
		&data.FirstName,
		&data.MiddleName,
		&data.LastName,
		&data.Email,
		&data.Password,
		&data.Role,
		&data.CreatedAt,
		&data.CreatedBy,
		&data.UpdatedAt,
		&data.UpdatedBy,
	)

	return
}
