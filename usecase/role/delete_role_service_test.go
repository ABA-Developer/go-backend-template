package role

import (
	"be-dashboard-nba/constant"
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestDeleteRoleService(t *testing.T) {

	const (
		selectRoleQuery = `SELECT .* FROM app_role WHERE id = \$1`
		deleteRoleQuery = `DELETE FROM app_role WHERE id = \$1`
	)

	rowFields := []string{
		"id", "name", "code", "description",
	}

	createRoleRow := func() *sqlmock.Rows {
		row := sqlmock.NewRows(rowFields)

		row.AddRow(
			1,
			"role",
			"1000 1000 100",
			"descriptionrole",
		)

		return row
	}

	tests := []struct {
		name        string
		roleID      int
		wantErr     bool
		expectedErr error
		setupMock   func(mock sqlmock.Sqlmock, roleID int)
	}{
		{
			name:        "success delete role",
			roleID:      1,
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, roleID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(selectRoleQuery).
					WithArgs(roleID).
					WillReturnRows(createRoleRow())

				mock.ExpectExec(deleteRoleQuery).WithArgs(roleID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:        "role id not found",
			roleID:      999,
			wantErr:     true,
			expectedErr: constant.ErrRoleIdNotFound,
			setupMock: func(mock sqlmock.Sqlmock, roleID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(selectRoleQuery).
					WithArgs(roleID).
					WillReturnError(sql.ErrNoRows)

				mock.ExpectRollback()
			},
		},
		{
			name:        "sql query error when get role by id",
			roleID:      1,
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, roleID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(selectRoleQuery).
					WithArgs(roleID).
					WillReturnError(errors.New("sql error"))

				mock.ExpectRollback()
			},
		},
		{
			name:        "sql query error when delete role",
			roleID:      1,
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, roleID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(selectRoleQuery).
					WithArgs(roleID).
					WillReturnRows(createRoleRow())

				mock.ExpectExec(deleteRoleQuery).WithArgs(roleID).
					WillReturnError(errors.New("sql error"))

				mock.ExpectRollback()
			},
		},
		{
			name:        "error beginning transaction",
			roleID:      1,
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, roleID int) {
				mock.ExpectBegin().WillReturnError(errors.New("begin error"))
			},
		},
		{
			name:        "error commiting transaction",
			roleID:      1,
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, roleID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(selectRoleQuery).
					WithArgs(roleID).
					WillReturnRows(createRoleRow())

				mock.ExpectExec(deleteRoleQuery).WithArgs(roleID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, mock, cleanup := newTestService(t)
			defer cleanup()

			ctx := context.Background()

			tc.setupMock(mock, tc.roleID)

			err := svc.DeleteRoleService(ctx, tc.roleID)

			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})

	}

}
