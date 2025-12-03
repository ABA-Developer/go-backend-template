package service

import (
	rolePresenter "be-dashboard-nba/api/presenter/role"
	"be-dashboard-nba/constant"
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateRoleService(t *testing.T) {
	description := "Updated role description"
	emptyDescription := ""

	tests := []struct {
		name        string
		roleID      int
		userID      string
		payload     rolePresenter.UpdateRoleRequest
		wantErr     bool
		expectedErr error
		setupMock   func(mock sqlmock.Sqlmock, roleID int, userID string, payload rolePresenter.UpdateRoleRequest)
	}{
		{
			name:   "success update role with description",
			roleID: 1,
			userID: "user-123",
			payload: rolePresenter.UpdateRoleRequest{
				Name:        "Updated Admin",
				Code:        "UPDATED_ADMIN",
				Description: &description,
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, userID string, payload rolePresenter.UpdateRoleRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(1, "Admin", "ADMIN", "Original description"))

				// Mock update query
				mock.ExpectExec(`UPDATE app_role`).
					WithArgs(
						payload.Code,
						payload.Name,
						sql.NullString{String: description, Valid: true},
						userID,
						roleID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))

				// Mock transaction commit
				mock.ExpectCommit()
			},
		},
		{
			name:   "success update role without description",
			roleID: 2,
			userID: "user-456",
			payload: rolePresenter.UpdateRoleRequest{
				Name:        "Updated User",
				Code:        "UPDATED_USER",
				Description: nil,
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, userID string, payload rolePresenter.UpdateRoleRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(2, "User", "USER", "User role"))

				// Mock update query - description will use COALESCE to keep existing value
				mock.ExpectExec(`UPDATE app_role`).
					WithArgs(
						payload.Code,
						payload.Name,
						sql.NullString{Valid: false}, // nil description
						userID,
						roleID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))

				// Mock transaction commit
				mock.ExpectCommit()
			},
		},
		{
			name:   "success update role with empty description",
			roleID: 3,
			userID: "user-789",
			payload: rolePresenter.UpdateRoleRequest{
				Name:        "Manager",
				Code:        "MANAGER",
				Description: &emptyDescription,
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, userID string, payload rolePresenter.UpdateRoleRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(3, "Old Manager", "OLD_MANAGER", "Old description"))

				// Mock update query - empty description
				mock.ExpectExec(`UPDATE app_role`).
					WithArgs(
						payload.Code,
						payload.Name,
						sql.NullString{String: "", Valid: true},
						userID,
						roleID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))

				// Mock transaction commit
				mock.ExpectCommit()
			},
		},
		{
			name:   "role not found",
			roleID: 999,
			userID: "user-123",
			payload: rolePresenter.UpdateRoleRequest{
				Name:        "Non-existent",
				Code:        "NON_EXISTENT",
				Description: &description,
			},
			wantErr:     true,
			expectedErr: constant.ErrRoleIdNotFound,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, userID string, payload rolePresenter.UpdateRoleRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock role existence check returns no rows
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnError(sql.ErrNoRows)

				// Mock transaction rollback due to error
				mock.ExpectRollback()
			},
		},
		{
			name:   "database error on role check",
			roleID: 1,
			userID: "user-123",
			payload: rolePresenter.UpdateRoleRequest{
				Name:        "Admin",
				Code:        "ADMIN",
				Description: &description,
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, userID string, payload rolePresenter.UpdateRoleRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock role existence check returns database error
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnError(sql.ErrConnDone)

				// Mock transaction rollback due to error
				mock.ExpectRollback()
			},
		},
		{
			name:   "database error on update query",
			roleID: 1,
			userID: "user-123",
			payload: rolePresenter.UpdateRoleRequest{
				Name:        "Admin",
				Code:        "ADMIN",
				Description: &description,
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, userID string, payload rolePresenter.UpdateRoleRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(1, "Admin", "ADMIN", "Original description"))

				// Mock update query returns error
				mock.ExpectExec(`UPDATE app_role`).
					WithArgs(
						payload.Code,
						payload.Name,
						sql.NullString{String: description, Valid: true},
						userID,
						roleID,
					).
					WillReturnError(sql.ErrConnDone)

				// Mock transaction rollback due to error
				mock.ExpectRollback()
			},
		},
		{
			name:   "database error on transaction begin",
			roleID: 1,
			userID: "user-123",
			payload: rolePresenter.UpdateRoleRequest{
				Name:        "Admin",
				Code:        "ADMIN",
				Description: &description,
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, userID string, payload rolePresenter.UpdateRoleRequest) {
				// Mock transaction begin returns error
				mock.ExpectBegin().WillReturnError(sql.ErrConnDone)
				// No rollback expected because transaction never started
			},
		},
		{
			name:   "database error on transaction commit",
			roleID: 1,
			userID: "user-123",
			payload: rolePresenter.UpdateRoleRequest{
				Name:        "Admin",
				Code:        "ADMIN",
				Description: &description,
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, userID string, payload rolePresenter.UpdateRoleRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(1, "Admin", "ADMIN", "Original description"))

				// Mock update query
				mock.ExpectExec(`UPDATE app_role`).
					WithArgs(
						payload.Code,
						payload.Name,
						sql.NullString{String: description, Valid: true},
						userID,
						roleID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))

				// Mock transaction commit returns error
				mock.ExpectCommit().WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:   "rollback error after update error",
			roleID: 1,
			userID: "user-123",
			payload: rolePresenter.UpdateRoleRequest{
				Name:        "Admin",
				Code:        "ADMIN",
				Description: &description,
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, roleID int, userID string, payload rolePresenter.UpdateRoleRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(roleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(1, "Admin", "ADMIN", "Original description"))

				// Mock update query returns error
				mock.ExpectExec(`UPDATE app_role`).
					WithArgs(
						payload.Code,
						payload.Name,
						sql.NullString{String: description, Valid: true},
						userID,
						roleID,
					).
					WillReturnError(sql.ErrConnDone)

				// Mock transaction rollback also returns error
				mock.ExpectRollback().WillReturnError(sql.ErrTxDone)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, mock, cleanup := newTestService(t)
			defer cleanup()

			ctx := context.Background()

			tc.setupMock(mock, tc.roleID, tc.userID, tc.payload)

			err := svc.UpdateRoleService(ctx, tc.payload, tc.userID, tc.roleID)

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
