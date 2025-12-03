package service

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/entities"
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUserService(t *testing.T) {
	tests := []struct {
		name           string
		userIDToDelete string
		deletedBy      string
		existingUser   entities.User
		wantErr        bool
		expectedErr    error
		setupMock      func(mock sqlmock.Sqlmock, userIDToDelete string, existingUser entities.User)
	}{
		{
			name:           "success delete user",
			userIDToDelete: "user-123",
			deletedBy:      "user-456",
			existingUser: entities.User{
				ID:       "user-123", // Different from userIDToDelete to avoid self-delete
				Name:     "John Doe",
				FullName: "John Doe",
				Email:    "john@example.com",
				Phone:    sql.NullString{String: "08123456789", Valid: true},
				Active:   true,
				Role:     "Admin",
				RoleID:   1,
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, userIDToDelete string, existingUser entities.User) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock user existence check - sesuai dengan query sebenarnya
				mock.ExpectQuery(`SELECT`).
					WithArgs(userIDToDelete).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "name", "full_name", "email", "phone", "active", "role", "role_id",
					}).
						AddRow(
							existingUser.ID,
							existingUser.Name,
							existingUser.FullName,
							existingUser.Email,
							existingUser.Phone,
							existingUser.Active,
							existingUser.Role,
							existingUser.RoleID,
						))

				// Mock delete user query
				mock.ExpectExec(`DELETE FROM app_user`).
					WithArgs(userIDToDelete).
					WillReturnResult(sqlmock.NewResult(0, 1))

				// Mock transaction commit
				mock.ExpectCommit()
			},
		},
		{
			name:           "user not found",
			userIDToDelete: "user-999",
			deletedBy:      "user-456",
			existingUser:   entities.User{},
			wantErr:        true,
			expectedErr:    constant.ErrUserIdNotFound,
			setupMock: func(mock sqlmock.Sqlmock, userIDToDelete string, existingUser entities.User) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock user existence check returns no rows
				mock.ExpectQuery(`SELECT`).
					WithArgs(userIDToDelete).
					WillReturnError(sql.ErrNoRows)

				// Mock transaction rollback
				mock.ExpectRollback()
			},
		},
		{
			name:           "forbidden self delete",
			userIDToDelete: "user-456",
			deletedBy:      "user-456",
			existingUser: entities.User{
				ID:       "user-456", // Same as userIDToDelete - self delete
				Name:     "John Doe",
				FullName: "John Doe",
				Email:    "john@example.com",
				Phone:    sql.NullString{String: "08123456789", Valid: true},
				Active:   true,
				Role:     "Admin",
				RoleID:   1,
			},
			wantErr:     true,
			expectedErr: constant.ErrForbiddenSelfDelete,
			setupMock: func(mock sqlmock.Sqlmock, userIDToDelete string, existingUser entities.User) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock user existence check
				mock.ExpectQuery(`SELECT`).
					WithArgs(userIDToDelete).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "name", "full_name", "email", "phone", "active", "role", "role_id",
					}).
						AddRow(
							existingUser.ID,
							existingUser.Name,
							existingUser.FullName,
							existingUser.Email,
							existingUser.Phone,
							existingUser.Active,
							existingUser.Role,
							existingUser.RoleID,
						))

				// Mock transaction rollback (due to self delete error)
				mock.ExpectRollback()
			},
		},
		{
			name:           "database error on user check",
			userIDToDelete: "user-123",
			deletedBy:      "user-456",
			existingUser:   entities.User{},
			wantErr:        true,
			expectedErr:    constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, userIDToDelete string, existingUser entities.User) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock user existence check returns database error
				mock.ExpectQuery(`SELECT`).
					WithArgs(userIDToDelete).
					WillReturnError(sql.ErrConnDone)

				// Mock transaction rollback
				mock.ExpectRollback()
			},
		},
		{
			name:           "database error on delete user",
			userIDToDelete: "user-123",
			deletedBy:      "user-456",
			existingUser: entities.User{
				ID:       "user-123",
				Name:     "John Doe",
				FullName: "John Doe",
				Email:    "john@example.com",
				Phone:    sql.NullString{String: "08123456789", Valid: true},
				Active:   true,
				Role:     "Admin",
				RoleID:   1,
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, userIDToDelete string, existingUser entities.User) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock user existence check
				mock.ExpectQuery(`SELECT`).
					WithArgs(userIDToDelete).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "name", "full_name", "email", "phone", "active", "role", "role_id",
					}).
						AddRow(
							existingUser.ID,
							existingUser.Name,
							existingUser.FullName,
							existingUser.Email,
							existingUser.Phone,
							existingUser.Active,
							existingUser.Role,
							existingUser.RoleID,
						))

				// Mock delete user query returns error
				mock.ExpectExec(`DELETE FROM app_user`).
					WithArgs(userIDToDelete).
					WillReturnError(sql.ErrConnDone)

				// Mock transaction rollback
				mock.ExpectRollback()
			},
		},
		{
			name:           "database error on transaction begin",
			userIDToDelete: "user-123",
			deletedBy:      "user-456",
			existingUser:   entities.User{},
			wantErr:        true,
			expectedErr:    constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, userIDToDelete string, existingUser entities.User) {
				// Mock transaction begin returns error
				mock.ExpectBegin().WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:           "database error on transaction commit",
			userIDToDelete: "user-123",
			deletedBy:      "user-456",
			existingUser: entities.User{
				ID:       "user-123",
				Name:     "John Doe",
				FullName: "John Doe",
				Email:    "john@example.com",
				Phone:    sql.NullString{String: "08123456789", Valid: true},
				Active:   true,
				Role:     "Admin",
				RoleID:   1,
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, userIDToDelete string, existingUser entities.User) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock user existence check
				mock.ExpectQuery(`SELECT`).
					WithArgs(userIDToDelete).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "name", "full_name", "email", "phone", "active", "role", "role_id",
					}).
						AddRow(
							existingUser.ID,
							existingUser.Name,
							existingUser.FullName,
							existingUser.Email,
							existingUser.Phone,
							existingUser.Active,
							existingUser.Role,
							existingUser.RoleID,
						))

				// Mock delete user query
				mock.ExpectExec(`DELETE FROM app_user`).
					WithArgs(userIDToDelete).
					WillReturnResult(sqlmock.NewResult(0, 1))

				// Mock transaction commit returns error
				mock.ExpectCommit().WillReturnError(sql.ErrConnDone)

			},
		},
		{
			name:           "rollback error after delete error",
			userIDToDelete: "user-123",
			deletedBy:      "user-456",
			existingUser: entities.User{
				ID:       "user-123",
				Name:     "John Doe",
				FullName: "John Doe",
				Email:    "john@example.com",
				Phone:    sql.NullString{String: "08123456789", Valid: true},
				Active:   true,
				Role:     "Admin",
				RoleID:   1,
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, userIDToDelete string, existingUser entities.User) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock user existence check
				mock.ExpectQuery(`SELECT`).
					WithArgs(userIDToDelete).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "name", "full_name", "email", "phone", "active", "role", "role_id",
					}).
						AddRow(
							existingUser.ID,
							existingUser.Name,
							existingUser.FullName,
							existingUser.Email,
							existingUser.Phone,
							existingUser.Active,
							existingUser.Role,
							existingUser.RoleID,
						))

				// Mock delete user query returns error
				mock.ExpectExec(`DELETE FROM app_user`).
					WithArgs(userIDToDelete).
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

			tc.setupMock(mock, tc.userIDToDelete, tc.existingUser)

			err := svc.DeleteUserService(ctx, tc.userIDToDelete, tc.deletedBy)

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
