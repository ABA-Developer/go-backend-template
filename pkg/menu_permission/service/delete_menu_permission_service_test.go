package service

import (
	"be-dashboard-nba/constant"
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDeleteMenuPermissionService(t *testing.T) {
	tests := []struct {
		name             string
		menuPermissionID int
		wantErr          bool
		expectedErr      error
		setupMock        func(mock sqlmock.Sqlmock, menuPermissionID int)
	}{
		{
			name:             "success delete menu permission",
			menuPermissionID: 1,
			wantErr:          false,
			expectedErr:      nil,
			setupMock: func(mock sqlmock.Sqlmock, menuPermissionID int) {
				mock.ExpectBegin()

				// Mock menu permission exists
				mock.ExpectQuery(`(?s)SELECT.*FROM.*app_menu_permission.*WHERE.*id`).
					WithArgs(menuPermissionID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "menu_id", "code", "action_name",
					}).AddRow(
						menuPermissionID, 1, "R", "read",
					))

				// Mock delete query
				mock.ExpectExec(`DELETE FROM app_menu_permission`).
					WithArgs(menuPermissionID).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:             "menu permission not found",
			menuPermissionID: 999,
			wantErr:          true,
			expectedErr:      constant.ErrMenuPermissionIdNotFound,
			setupMock: func(mock sqlmock.Sqlmock, menuPermissionID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(`(?s)SELECT.*FROM.*app_menu_permission.*WHERE.*id`).
					WithArgs(menuPermissionID).
					WillReturnError(sql.ErrNoRows)

				mock.ExpectRollback()
			},
		},
		{
			name:             "error reading menu permission",
			menuPermissionID: 1,
			wantErr:          true,
			expectedErr:      constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, menuPermissionID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(`(?s)SELECT.*FROM.*app_menu_permission.*WHERE.*id`).
					WithArgs(menuPermissionID).
					WillReturnError(errors.New("select error"))

				mock.ExpectRollback()
			},
		},
		{
			name:             "error deleting menu permission",
			menuPermissionID: 1,
			wantErr:          true,
			expectedErr:      constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, menuPermissionID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(`(?s)SELECT.*FROM.*app_menu_permission.*WHERE.*id`).
					WithArgs(menuPermissionID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "menu_id", "code", "action_name",
					}).AddRow(
						menuPermissionID, 1, "R", "read",
					))

				mock.ExpectExec(`DELETE FROM app_menu_permission`).
					WithArgs(menuPermissionID).
					WillReturnError(errors.New("delete error"))

				mock.ExpectRollback()
			},
		},
		{
			name:             "error beginning transaction",
			menuPermissionID: 1,
			wantErr:          true,
			expectedErr:      constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, menuPermissionID int) {
				mock.ExpectBegin().WillReturnError(errors.New("begin error"))
			},
		},
		{
			name:             "error committing transaction",
			menuPermissionID: 1,
			wantErr:          true,
			expectedErr:      constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, menuPermissionID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(`(?s)SELECT.*FROM.*app_menu_permission.*WHERE.*id`).
					WithArgs(menuPermissionID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "menu_id", "code", "action_name",
					}).AddRow(
						menuPermissionID, 1, "R", "read",
					))

				mock.ExpectExec(`DELETE FROM app_menu_permission`).
					WithArgs(menuPermissionID).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, mock, cleanup := newTestService(t)
			defer cleanup()

			ctx := context.Background()

			tc.setupMock(mock, tc.menuPermissionID)

			err := svc.DeleteMenuPermissionService(ctx, tc.menuPermissionID)

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
