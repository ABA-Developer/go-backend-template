package repository

import (
	"context"
	"database/sql"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/user/domain"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func isPostgresUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}

func isPostgresForeignKeyViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23503"
	}
	return false
}

func (r *userRepository) CreateUserWithRoleTx(
	ctx context.Context,
	user domain.User,
) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	defer tx.Rollback()

	var newUserID string
	const userStmt = `
		INSERT INTO app_user (
			name, full_name, email, password, active, phone, created_by, created_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, (now() at time zone 'UTC')::TIMESTAMP
		)
		RETURNING id
	`

	err = tx.QueryRowContext(ctx, userStmt,
		user.Name,
		user.FullName,
		user.Email,
		user.Password,
		user.Active,
		user.Phone,
		user.CreatedBy,
	).Scan(&newUserID)

	if err != nil {
		if isPostgresUniqueViolation(err) {
			return constant.ErrEmailAlreadyExists
		}
		return errors.Wrap(err, "failed to insert user")
	}

	const roleStmt = `
		INSERT INTO app_user_role(
			user_id, role_id
		)
		VALUES (
			$1, $2
		)
	`
	_, err = tx.ExecContext(ctx, roleStmt, newUserID, user.RoleID)
	if err != nil {
		if isPostgresForeignKeyViolation(err) {
			return constant.ErrRoleIdNotFound
		}
		return errors.Wrap(err, "failed to insert user role")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}
func (r *userRepository) UpdateUserWithRoleTx(
	ctx context.Context,
	user domain.User,
) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	// 2. Defer rollback: jika nanti return err != nil, otomatis dibatalkan
	defer tx.Rollback()

	const updateRoleStmt = `
		UPDATE app_user_role
		SET role_id = $2
		WHERE user_id = $1
	`
	// 3. ✨ GUNAKAN tx.ExecContext (bukan r.db.ExecContext) ✨
	_, err = tx.ExecContext(ctx, updateRoleStmt, user.ID, user.RoleID)
	if err != nil {
		if isPostgresForeignKeyViolation(err) {
			return constant.ErrRoleIdNotFound
		}
		return errors.Wrap(err, "failed to update user role")
	}

	const updateUserStmt = `
		UPDATE app_user 
		SET 
			name = $1,
			full_name = $2,
			email = $3,
			phone = $4,
			active = $5,
			updated_by = $6,
			updated_at = (now() at time zone 'UTC')::TIMESTAMP
		WHERE id = $7
	`
	// 4. ✨ GUNAKAN tx.ExecContext ✨
	_, err = tx.ExecContext(ctx, updateUserStmt,
		user.Name,
		user.FullName,
		user.Email,
		user.Phone,
		user.Active,
		user.UpdatedBy,
		user.ID,
	)
	if err != nil {
		if isPostgresUniqueViolation(err) {
			return constant.ErrEmailAlreadyExists
		}
		return errors.Wrap(err, "failed to update user details")
	}

	// 5. Commit transaksi jika semua query di atas berhasil
	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

func (r *userRepository) DeleteUserQuery(
	ctx context.Context,
	id string,
) (err error) {
	const statement = `
		DELETE FROM app_user
		WHERE
			id = $1
	`

	_, err = r.db.ExecContext(ctx, statement, id)

	return
}

func (r *userRepository) ReadUsersQuery(
	ctx context.Context,
	filter domain.UserFilter,
) (data []domain.UserListRow, err error) {
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
		filter.SetSearch,
		filter.Search,
		filter.Order,
		filter.Limit,
		filter.Offset,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var u domain.UserListRow

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

func (r *userRepository) ReadCountUserQuery(
	ctx context.Context,
	filter domain.UserFilter,
) (count int, err error) {
	const stmt = `
		SELECT
			COUNT(*)
		FROM
			app_user
		WHERE
			(CASE WHEN $1::bool THEN(
				full_name ILIKE $2
			) ELSE TRUE END)
	`

	err = r.db.QueryRowContext(ctx, stmt,
		filter.SetSearch,
		filter.Search,
	).Scan(&count)

	return
}

func (r *userRepository) ReadUserByIDQuery(
	ctx context.Context,
	id string,
) (data domain.UserDetailRow, err error) {
	const statement = `
		SELECT 	
			au.id, 
			au.name, 
			au.full_name, 
			au.email,
			au.phone,
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
			au.id = $1
	`

	err = r.db.QueryRowContext(ctx, statement, id).Scan(
		&data.ID,
		&data.Name,
		&data.FullName,
		&data.Email,
		&data.Phone,
		&data.Active,
		&data.Role,
		&data.RoleID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			err = constant.ErrDataNotFound
		}
		return
	}

	return
}

func (r *userRepository) ReadUserProfileQuery(
	ctx context.Context,
	id string,
) (data domain.UserProfileRow, err error) {
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
	if err != nil {
		if err == sql.ErrNoRows {
			err = constant.ErrDataNotFound
		}
		return
	}

	return
}

func (r *userRepository) IsUserEmailExistsQuery(
	ctx context.Context,
	email string,
) (exists bool, err error) {
	const statement = `
		SELECT EXISTS (
			SELECT
				1
			FROM
				app_user
			WHERE
				email = $1
		)
	`

	err = r.db.QueryRowContext(ctx, statement, email).Scan(&exists)

	return
}

func (r *userRepository) IsUpdateUserEmailExistsQuery(
	ctx context.Context,
	email, id string,
) (exists bool, err error) {
	statement := `
		SELECT EXISTS (
			SELECT
				1
			FROM
				app_user
			WHERE
				email = $1
				AND id != $2
		)
	`

	err = r.db.QueryRowContext(ctx, statement, email, id).Scan(&exists)

	return
}
