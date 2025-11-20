package service

import (
	menuPermissionPresenter "be-dashboard-nba/api/presenter/menu_permission"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/menu_permission/repository"
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMenuPermissionService(t *testing.T) {
	tests := []struct {
		name             string
		payload          menuPermissionPresenter.UpdateMenuPermissionRequest
		userID           string
		menuPermissionID int
		wantErr          bool
		expectedErr      error
		setupMock        func(mock sqlmock.Sqlmock, payload menuPermissionPresenter.UpdateMenuPermissionRequest, userID string, menuPermissionID int)
	}{
		{
			name: "success update",
			payload: menuPermissionPresenter.UpdateMenuPermissionRequest{
				Code:       "U",
				ActionName: "Update",
			},
			userID:           "user-123",
			menuPermissionID: 1,
			wantErr:          false,
			expectedErr:      nil,
			setupMock: func(mock sqlmock.Sqlmock, payload menuPermissionPresenter.UpdateMenuPermissionRequest, userID string, menuPermissionID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(`(?s)SELECT.*FROM.*app_menu_permission.*WHERE.*id`).
					WithArgs(menuPermissionID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "menu_id", "code", "action_name",
					}).AddRow(
						menuPermissionID, 1, "R", "Read",
					))

				expectedParams := payload.ToParams(userID, menuPermissionID)
				mock.ExpectExec(`UPDATE app_menu_permission`).
					WithArgs(
						expectedParams.Code,
						expectedParams.ActionName,
						expectedParams.UpdatedBy,
						expectedParams.MenuPermissionID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "menu permission not found",
			payload: menuPermissionPresenter.UpdateMenuPermissionRequest{
				Code:       "U",
				ActionName: "Update",
			},
			userID:           "user-123",
			menuPermissionID: 999,
			wantErr:          true,
			expectedErr:      constant.ErrMenuPermissionIdNotFound,
			setupMock: func(mock sqlmock.Sqlmock, payload menuPermissionPresenter.UpdateMenuPermissionRequest, userID string, menuPermissionID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(`(?s)SELECT.*FROM.*app_menu_permission.*WHERE.*id`).
					WithArgs(menuPermissionID).
					WillReturnError(sql.ErrNoRows)

				mock.ExpectRollback()
			},
		},
		{
			name: "error reading menu permission",
			payload: menuPermissionPresenter.UpdateMenuPermissionRequest{
				Code:       "U",
				ActionName: "Update",
			},
			userID:           "user-123",
			menuPermissionID: 1,
			wantErr:          true,
			expectedErr:      constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, payload menuPermissionPresenter.UpdateMenuPermissionRequest, userID string, menuPermissionID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(`(?s)SELECT.*FROM.*app_menu_permission.*WHERE.*id`).
					WithArgs(menuPermissionID).
					WillReturnError(errors.New("database connection failed"))

				mock.ExpectRollback()
			},
		},
		{
			name: "error updating menu permission",
			payload: menuPermissionPresenter.UpdateMenuPermissionRequest{
				Code:       "U",
				ActionName: "Update",
			},
			userID:           "user-123",
			menuPermissionID: 1,
			wantErr:          true,
			expectedErr:      constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, payload menuPermissionPresenter.UpdateMenuPermissionRequest, userID string, menuPermissionID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(`(?s)SELECT.*FROM.*app_menu_permission.*WHERE.*id`).
					WithArgs(menuPermissionID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "menu_id", "code", "action_name",
					}).AddRow(
						menuPermissionID, 1, "R", "Read",
					))

				expectedParams := payload.ToParams(userID, menuPermissionID)
				mock.ExpectExec(`UPDATE app_menu_permission`).
					WithArgs(
						expectedParams.Code,
						expectedParams.ActionName,
						expectedParams.UpdatedBy,
						expectedParams.MenuPermissionID,
					).
					WillReturnError(errors.New("update failed"))

				mock.ExpectRollback()
			},
		},
		{
			name: "error beginning transaction",
			payload: menuPermissionPresenter.UpdateMenuPermissionRequest{
				Code:       "U",
				ActionName: "Update",
			},
			userID:           "user-123",
			menuPermissionID: 1,
			wantErr:          true,
			expectedErr:      constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, payload menuPermissionPresenter.UpdateMenuPermissionRequest, userID string, menuPermissionID int) {
				mock.ExpectBegin().WillReturnError(errors.New("begin transaction failed"))
			},
		},
		{
			name: "error committing transaction",
			payload: menuPermissionPresenter.UpdateMenuPermissionRequest{
				Code:       "U",
				ActionName: "Update",
			},
			userID:           "user-123",
			menuPermissionID: 1,
			wantErr:          true,
			expectedErr:      constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, payload menuPermissionPresenter.UpdateMenuPermissionRequest, userID string, menuPermissionID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(`(?s)SELECT.*FROM.*app_menu_permission.*WHERE.*id`).
					WithArgs(menuPermissionID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "menu_id", "code", "action_name",
					}).AddRow(
						menuPermissionID, 1, "R", "Read",
					))

				expectedParams := payload.ToParams(userID, menuPermissionID)
				mock.ExpectExec(`UPDATE app_menu_permission`).
					WithArgs(
						expectedParams.Code,
						expectedParams.ActionName,
						expectedParams.UpdatedBy,
						expectedParams.MenuPermissionID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit().WillReturnError(errors.New("commit failed"))
			},
		},
		{
			name: "error rolling back transaction",
			payload: menuPermissionPresenter.UpdateMenuPermissionRequest{
				Code:       "U",
				ActionName: "Update",
			},
			userID:           "user-123",
			menuPermissionID: 1,
			wantErr:          true,
			expectedErr:      constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, payload menuPermissionPresenter.UpdateMenuPermissionRequest, userID string, menuPermissionID int) {
				mock.ExpectBegin()

				mock.ExpectQuery(`(?s)SELECT.*FROM.*app_menu_permission.*WHERE.*id`).
					WithArgs(menuPermissionID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "menu_id", "code", "action_name",
					}).AddRow(
						menuPermissionID, 1, "R", "Read",
					))

				expectedParams := payload.ToParams(userID, menuPermissionID)
				mock.ExpectExec(`UPDATE app_menu_permission`).
					WithArgs(
						expectedParams.Code,
						expectedParams.ActionName,
						expectedParams.UpdatedBy,
						expectedParams.MenuPermissionID,
					).
					WillReturnError(errors.New("update failed"))

				mock.ExpectRollback().WillReturnError(errors.New("rollback failed"))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, mock, cleanup := newTestService(t)
			defer cleanup()

			ctx := context.Background()

			tc.setupMock(mock, tc.payload, tc.userID, tc.menuPermissionID)

			err := svc.UpdateMenuPermissionService(ctx, tc.payload, tc.userID, tc.menuPermissionID)

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

func TestUpdateMenuPermissionPayload_ToParams(t *testing.T) {
	req := menuPermissionPresenter.UpdateMenuPermissionRequest{
		Code:       "U",
		ActionName: "Update Permission",
	}

	userID := "user-123"
	menuPermissionID := 1

	params := req.ToParams(userID, menuPermissionID)

	expected := repository.UpdateMenuPermissionPayload{
		Code:             "U",
		ActionName:       "Update Permission",
		UpdatedBy:        "user-123",
		MenuPermissionID: 1,
	}

	assert.Equal(t, expected, params)
}
