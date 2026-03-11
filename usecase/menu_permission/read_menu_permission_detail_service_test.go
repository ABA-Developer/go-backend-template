package menu_permission

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/usecase/entities"
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestReadMenuPermissionDetail(t *testing.T) {
	tests := []struct {
		name             string
		menuPermissionID int
		expectedData     entities.MenuPermission
		wantErr          bool
		expectedErr      error
		setupMock        func(mock sqlmock.Sqlmock, menuPermissionID int)
	}{
		{
			name:             "success",
			menuPermissionID: 1,
			expectedData: entities.MenuPermission{
				ID:         1,
				MenuID:     1,
				Code:       "R",
				ActionName: "Read",
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, menuPermissionID int) {
				mock.ExpectQuery(`(?s)SELECT.*FROM.*app_menu_permission.*WHERE.*id`).
					WithArgs(menuPermissionID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "menu_id", "code", "action_name",
					}).AddRow(
						menuPermissionID, 1, "R", "Read",
					))
			},
		},
		{
			name:             "menu permission not found",
			menuPermissionID: 999,
			expectedData:     entities.MenuPermission{},
			wantErr:          true,
			expectedErr:      constant.ErrMenuPermissionIdNotFound,
			setupMock: func(mock sqlmock.Sqlmock, menuPermissionID int) {
				mock.ExpectQuery(`(?s)SELECT.*FROM.*app_menu_permission.*WHERE.*id`).
					WithArgs(menuPermissionID).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:             "database error",
			menuPermissionID: 1,
			expectedData:     entities.MenuPermission{},
			wantErr:          true,
			expectedErr:      constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, menuPermissionID int) {
				mock.ExpectQuery(`(?s)SELECT.*FROM.*app_menu_permission.*WHERE.*id`).
					WithArgs(menuPermissionID).
					WillReturnError(errors.New("database connection failed"))
			},
		},
		{
			name:             "success with different data",
			menuPermissionID: 2,
			expectedData: entities.MenuPermission{
				ID:         2,
				MenuID:     1,
				Code:       "U",
				ActionName: "Update",
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, menuPermissionID int) {
				mock.ExpectQuery(`(?s)SELECT.*FROM.*app_menu_permission.*WHERE.*id`).
					WithArgs(menuPermissionID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "menu_id", "code", "action_name",
					}).AddRow(
						menuPermissionID, 1, "U", "Update",
					))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, mock, cleanup := newTestService(t)
			defer cleanup()

			ctx := context.Background()

			// Setup mock expectations
			tc.setupMock(mock, tc.menuPermissionID)

			// Execute service
			data, err := svc.ReadMenuPermissionDetail(ctx, tc.menuPermissionID)

			// Assertions
			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
				assert.Equal(t, tc.expectedData, data) // Should be empty on error
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedData, data)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
