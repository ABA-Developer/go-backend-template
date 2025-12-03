package service

import (
	"be-dashboard-nba/constant"
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestDeleteMenuService(t *testing.T) {
	tests := []struct {
		name        string
		menuID      int
		wantErr     bool
		expectedErr error
		setupMock   func(mock sqlmock.Sqlmock, menuID int)
	}{
		{
			name:        "success delete menu",
			menuID:      1,
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, menuID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(`SELECT .* FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "description", "url", "sort", "group", "icon", "active", "display",
					}).AddRow(
						menuID, nil, "Dashboard", nil, "/dashboard", 0, "main", nil, true, true,
					))

				mock.ExpectExec(`DELETE FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:        "menu not found",
			menuID:      99,
			wantErr:     true,
			expectedErr: constant.ErrMenuIdNotFound,
			setupMock: func(mock sqlmock.Sqlmock, menuID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(`SELECT .* FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnError(sql.ErrNoRows)

				mock.ExpectRollback()
			},
		},
		{
			name:        "error beginning transaction",
			menuID:      1,
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, menuID int) {
				mock.ExpectBegin().WillReturnError(errors.New("begin tx error"))
			},
		},
		{
			name:        "error reading menu",
			menuID:      1,
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, menuID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(`SELECT .* FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnError(errors.New("database error"))

				mock.ExpectRollback()
			},
		},
		{
			name:        "error deleting menu",
			menuID:      1,
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, menuID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(`SELECT .* FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "description", "url", "sort", "group", "icon", "active", "display",
					}).AddRow(
						menuID, nil, "Dashboard", nil, "/dashboard", 0, "main", nil, true, true,
					))

				mock.ExpectExec(`DELETE FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnError(errors.New("delete failed"))

				mock.ExpectRollback()
			},
		},
		{
			name:        "error committing transaction",
			menuID:      1,
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, menuID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(`SELECT .* FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "description", "url", "sort", "group", "icon", "active", "display",
					}).AddRow(
						menuID, nil, "Dashboard", nil, "/dashboard", 0, "main", nil, true, true,
					))

				mock.ExpectExec(`DELETE FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
		},
		{
			name:        "error rolling back transaction",
			menuID:      1,
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, menuID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(`SELECT .* FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnError(errors.New("read error"))

				// Rollback also fails
				mock.ExpectRollback().WillReturnError(errors.New("rollback failed"))
			},
		},
		{
			name:        "delete menu with children",
			menuID:      2,
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, menuID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(`SELECT .* FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "description", "url", "sort", "group", "icon", "active", "display",
					}).AddRow(
						menuID, nil, "Settings", stringPtr("Application settings"), "/settings", 1, "main", stringPtr("icon-settings"), true, true,
					))

				mock.ExpectExec(`DELETE FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},
		{
			name:        "delete menu with parent ID",
			menuID:      3,
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, menuID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(`SELECT .* FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "parent_id", "name", "description", "url", "sort", "group", "icon", "active", "display",
					}).AddRow(
						menuID, int32Ptr(1), "Submenu", nil, "/submenu", 0, "main", nil, true, true,
					))

				mock.ExpectExec(`DELETE FROM app_menu WHERE id = \$1`).
					WithArgs(menuID).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, mock, cleanup := newTestService(t)
			defer cleanup()

			ctx := context.Background()

			// Setup mock expectations
			tc.setupMock(mock, tc.menuID)

			// Execute service
			err := svc.DeleteMenuService(ctx, tc.menuID)

			// Assertions
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

// Helper functions untuk test
func stringPtr(s string) *string { return &s }
func int32Ptr(i int32) *int32    { return &i }
