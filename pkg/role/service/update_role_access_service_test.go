package service

import (
	presenter "be-dashboard-nba/api/presenter/role"
	"be-dashboard-nba/constant"
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateRoleAccessService(t *testing.T) {
	trueVal := true
	falseVal := false

	tests := []struct {
		name        string
		roleID      int
		request     presenter.UpdateRoleAccessRequest
		wantErr     bool
		expectedErr error
		setupMock   func(mock sqlmock.Sqlmock, roleID int, request presenter.UpdateRoleAccessRequest)
	}{
		{
			name:   "success update multiple role accesses",
			roleID: 1,
			request: presenter.UpdateRoleAccessRequest{
				AccessItem: []presenter.UpdateRoleAccessItem{
					{
						AccessID:  1,
						HasAccess: &trueVal,
					},
					{
						AccessID:  2,
						HasAccess: &falseVal,
					},
					{
						AccessID:  3,
						HasAccess: &trueVal,
					},
				},
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, request presenter.UpdateRoleAccessRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// For each access item, mock the checks and operations
				for _, item := range request.AccessItem {
					// Mock role existence check
					mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
						WithArgs(roleID).
						WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
							AddRow(roleID, "Admin", "ADMIN", "Administrator Role"))

					// Mock menu permission existence check
					mock.ExpectQuery(`SELECT id, menu_id, code, action_name FROM app_menu_permission WHERE id =`).
						WithArgs(item.AccessID).
						WillReturnRows(sqlmock.NewRows([]string{"id", "menu_id", "code", "action_name"}).
							AddRow(item.AccessID, 1, "R", "Read"))

					// Mock create or delete based on HasAccess
					if *item.HasAccess {
						// Mock create role access
						mock.ExpectExec(`INSERT INTO app_role_access`).
							WithArgs(roleID, item.AccessID).
							WillReturnResult(sqlmock.NewResult(1, 1))
					} else {
						// Mock delete role access
						mock.ExpectExec(`DELETE FROM app_role_access`).
							WithArgs(roleID, item.AccessID).
							WillReturnResult(sqlmock.NewResult(0, 1))
					}
				}

				// Mock transaction commit
				mock.ExpectCommit()
			},
		},
		{
			name:   "success create single role access",
			roleID: 2,
			request: presenter.UpdateRoleAccessRequest{
				AccessItem: []presenter.UpdateRoleAccessItem{
					{
						AccessID:  5,
						HasAccess: &trueVal,
					},
				},
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, request presenter.UpdateRoleAccessRequest) {
				mock.ExpectBegin()

				// Mock role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(roleID, "User", "USER", "User Role"))

				// Mock menu permission existence check
				mock.ExpectQuery(`SELECT id, menu_id, code, action_name FROM app_menu_permission WHERE id =`).
					WithArgs(request.AccessItem[0].AccessID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "menu_id", "code", "action_name"}).
						AddRow(request.AccessItem[0].AccessID, 2, "R", "Read"))

				// Mock create role access
				mock.ExpectExec(`INSERT INTO app_role_access`).
					WithArgs(roleID, request.AccessItem[0].AccessID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Mock transaction commit
				mock.ExpectCommit()
			},
		},
		{
			name:   "success delete single role access",
			roleID: 3,
			request: presenter.UpdateRoleAccessRequest{
				AccessItem: []presenter.UpdateRoleAccessItem{
					{
						AccessID:  10,
						HasAccess: &falseVal,
					},
				},
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, request presenter.UpdateRoleAccessRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(roleID, "Manager", "MANAGER", "Manager Role"))

				// Mock menu permission existence check
				mock.ExpectQuery(`SELECT id, menu_id, code, action_name FROM app_menu_permission WHERE id =`).
					WithArgs(request.AccessItem[0].AccessID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "menu_id", "code", "action_name"}).
						AddRow(request.AccessItem[0].AccessID, 3, "W", "Write"))

				// Mock delete role access
				mock.ExpectExec(`DELETE FROM app_role_access`).
					WithArgs(roleID, request.AccessItem[0].AccessID).
					WillReturnResult(sqlmock.NewResult(0, 1))

				// Mock transaction commit
				mock.ExpectCommit()
			},
		},
		{
			name:   "role not found",
			roleID: 999,
			request: presenter.UpdateRoleAccessRequest{
				AccessItem: []presenter.UpdateRoleAccessItem{
					{
						AccessID:  1,
						HasAccess: &trueVal,
					},
				},
			},
			wantErr:     true,
			expectedErr: constant.ErrRoleIdNotFound,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, request presenter.UpdateRoleAccessRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock role existence check returns no rows
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnError(sql.ErrNoRows)

				// Mock transaction rollback
				mock.ExpectRollback()
			},
		},
		{
			name:   "menu permission not found",
			roleID: 1,
			request: presenter.UpdateRoleAccessRequest{
				AccessItem: []presenter.UpdateRoleAccessItem{
					{
						AccessID:  999,
						HasAccess: &trueVal,
					},
				},
			},
			wantErr:     true,
			expectedErr: constant.ErrMenuPermissionIdNotFound,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, request presenter.UpdateRoleAccessRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(roleID, "Admin", "ADMIN", "Administrator Role"))

				// Mock menu permission existence check returns no rows
				mock.ExpectQuery(`SELECT id, menu_id, code, action_name FROM app_menu_permission WHERE id =`).
					WithArgs(request.AccessItem[0].AccessID).
					WillReturnError(sql.ErrNoRows)

				// Mock transaction rollback
				mock.ExpectRollback()
			},
		},
		{
			name:   "database error on role check",
			roleID: 1,
			request: presenter.UpdateRoleAccessRequest{
				AccessItem: []presenter.UpdateRoleAccessItem{
					{
						AccessID:  1,
						HasAccess: &trueVal,
					},
				},
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, request presenter.UpdateRoleAccessRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock role existence check returns database error
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnError(sql.ErrConnDone)

				// Mock transaction rollback
				mock.ExpectRollback()
			},
		},
		{
			name:   "database error on menu permission check",
			roleID: 1,
			request: presenter.UpdateRoleAccessRequest{
				AccessItem: []presenter.UpdateRoleAccessItem{
					{
						AccessID:  1,
						HasAccess: &trueVal,
					},
				},
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, request presenter.UpdateRoleAccessRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(roleID, "Admin", "ADMIN", "Administrator Role"))

				// Mock menu permission existence check returns database error
				mock.ExpectQuery(`SELECT id, menu_id, code, action_name FROM app_menu_permission WHERE id =`).
					WithArgs(request.AccessItem[0].AccessID).
					WillReturnError(sql.ErrConnDone)

				// Mock transaction rollback
				mock.ExpectRollback()
			},
		},
		{
			name:   "database error on create role access",
			roleID: 1,
			request: presenter.UpdateRoleAccessRequest{
				AccessItem: []presenter.UpdateRoleAccessItem{
					{
						AccessID:  1,
						HasAccess: &trueVal,
					},
				},
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, request presenter.UpdateRoleAccessRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(roleID, "Admin", "ADMIN", "Administrator Role"))

				// Mock menu permission existence check
				mock.ExpectQuery(`SELECT id, menu_id, code, action_name FROM app_menu_permission WHERE id =`).
					WithArgs(request.AccessItem[0].AccessID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "menu_id", "code", "action_name"}).
						AddRow(request.AccessItem[0].AccessID, 1, "R", "Read"))

				// Mock create role access returns error
				mock.ExpectExec(`INSERT INTO app_role_access`).
					WithArgs(roleID, request.AccessItem[0].AccessID).
					WillReturnError(sql.ErrConnDone)

				// Mock transaction rollback
				mock.ExpectRollback()
			},
		},
		{
			name:   "database error on delete role access",
			roleID: 1,
			request: presenter.UpdateRoleAccessRequest{
				AccessItem: []presenter.UpdateRoleAccessItem{
					{
						AccessID:  1,
						HasAccess: &falseVal,
					},
				},
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, request presenter.UpdateRoleAccessRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(roleID, "Admin", "ADMIN", "Administrator Role"))

				// Mock menu permission existence check
				mock.ExpectQuery(`SELECT id, menu_id, code, action_name FROM app_menu_permission WHERE id =`).
					WithArgs(request.AccessItem[0].AccessID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "menu_id", "code", "action_name"}).
						AddRow(request.AccessItem[0].AccessID, 1, "R", "Read"))

				// Mock delete role access returns error
				mock.ExpectExec(`DELETE FROM app_role_access`).
					WithArgs(roleID, request.AccessItem[0].AccessID).
					WillReturnError(sql.ErrConnDone)

				// Mock transaction rollback
				mock.ExpectRollback()
			},
		},
		{
			name:   "database error on transaction begin",
			roleID: 1,
			request: presenter.UpdateRoleAccessRequest{
				AccessItem: []presenter.UpdateRoleAccessItem{
					{
						AccessID:  1,
						HasAccess: &trueVal,
					},
				},
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, request presenter.UpdateRoleAccessRequest) {
				// Mock transaction begin returns error
				mock.ExpectBegin().WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:   "database error on transaction commit",
			roleID: 1,
			request: presenter.UpdateRoleAccessRequest{
				AccessItem: []presenter.UpdateRoleAccessItem{
					{
						AccessID:  1,
						HasAccess: &trueVal,
					},
				},
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, request presenter.UpdateRoleAccessRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(roleID, "Admin", "ADMIN", "Administrator Role"))

				// Mock menu permission existence check
				mock.ExpectQuery(`SELECT id, menu_id, code, action_name FROM app_menu_permission WHERE id =`).
					WithArgs(request.AccessItem[0].AccessID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "menu_id", "code", "action_name"}).
						AddRow(request.AccessItem[0].AccessID, 1, "R", "Read"))

				// Mock create role access
				mock.ExpectExec(`INSERT INTO app_role_access`).
					WithArgs(roleID, request.AccessItem[0].AccessID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Mock transaction commit returns error
				mock.ExpectCommit().WillReturnError(sql.ErrConnDone)

			},
		},
		{
			name:   "error in second item of multiple updates",
			roleID: 1,
			request: presenter.UpdateRoleAccessRequest{
				AccessItem: []presenter.UpdateRoleAccessItem{
					{
						AccessID:  1,
						HasAccess: &trueVal,
					},
					{
						AccessID:  2,
						HasAccess: &trueVal,
					},
				},
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, request presenter.UpdateRoleAccessRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// First item - success
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(roleID, "Admin", "ADMIN", "Administrator Role"))

				mock.ExpectQuery(`SELECT id, menu_id, code, action_name FROM app_menu_permission WHERE id =`).
					WithArgs(request.AccessItem[0].AccessID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "menu_id", "code", "action_name"}).
						AddRow(request.AccessItem[0].AccessID, 1, "R", "Read"))

				mock.ExpectExec(`INSERT INTO app_role_access`).
					WithArgs(roleID, request.AccessItem[0].AccessID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Second item - error on role check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnError(sql.ErrConnDone)

				// Mock transaction rollback
				mock.ExpectRollback()
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, mock, cleanup := newTestService(t)
			defer cleanup()

			ctx := context.Background()

			tc.setupMock(mock, tc.roleID, tc.request)

			err := svc.UpdateRoleAccessService(ctx, tc.roleID, tc.request)

			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
			} else {
				assert.NoError(t, err)
			}

			// Check if all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
