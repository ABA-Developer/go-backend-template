package user

import (
	presenter "be-dashboard-nba/api/presenter/user"
	"be-dashboard-nba/constant"
	"be-dashboard-nba/usecase/entities"
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUserService(t *testing.T) {
	phone := "08123456789"
	emptyPhone := ""
	activeTrue := true
	activeFalse := false

	tests := []struct {
		name         string
		userID       string
		updatedBy    string
		request      presenter.UpdateUserRequest
		existingUser entities.User
		existingRole entities.Role
		wantErr      bool
		expectedErr  error
		setupMock    func(mock sqlmock.Sqlmock, userID string, updatedBy string, request presenter.UpdateUserRequest, existingUser entities.User, existingRole entities.Role)
	}{
		{
			name:      "success update user with role change",
			userID:    "user-123",
			updatedBy: "admin-123",
			request: presenter.UpdateUserRequest{
				Name:     "updated_john",
				FullName: "Updated John Doe",
				Email:    "updated.john@example.com",
				RoleID:   2,
				Phone:    &phone,
				Active:   &activeTrue,
			},
			existingUser: entities.User{
				ID:       "user-123",
				Name:     "john_doe",
				FullName: "John Doe",
				Email:    "john@example.com",
				Phone:    sql.NullString{String: "08111111111", Valid: true},
				Active:   true,
				Role:     "Admin", // Different from new role
				RoleID:   1,
			},
			existingRole: entities.Role{
				ID:   2,
				Name: "Manager",
				Code: "MANAGER",
				Description: sql.NullString{
					String: "Manager Role",
					Valid:  true,
				},
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, userID string, updatedBy string, request presenter.UpdateUserRequest, existingUser entities.User, existingRole entities.Role) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock user existence check
				mock.ExpectQuery(`SELECT`).
					WithArgs(userID).
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

				// Mock role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(request.RoleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(existingRole.ID, existingRole.Name, existingRole.Code, existingRole.Description))

				// Mock update user role query (because role changed)
				mock.ExpectExec(`UPDATE app_user_role`).
					WithArgs(userID, request.RoleID).
					WillReturnResult(sqlmock.NewResult(0, 1))

				// Mock update user query
				mock.ExpectExec(`UPDATE app_user`).
					WithArgs(
						request.Name,
						request.FullName,
						request.Email,
						sql.NullString{String: phone, Valid: true},
						true,
						updatedBy,
						userID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))

				// Mock transaction commit
				mock.ExpectCommit()
			},
		},
		{
			name:      "success update user without role change",
			userID:    "user-123",
			updatedBy: "admin-123",
			request: presenter.UpdateUserRequest{
				Name:     "updated_john",
				FullName: "Updated John Doe",
				Email:    "updated.john@example.com",
				RoleID:   1, // Same role as existing
				Phone:    &phone,
				Active:   &activeTrue,
			},
			existingUser: entities.User{
				ID:       "user-123",
				Name:     "john_doe",
				FullName: "John Doe",
				Email:    "john@example.com",
				Phone:    sql.NullString{String: "08111111111", Valid: true},
				Active:   true,
				Role:     "Admin",
				RoleID:   1,
			},
			existingRole: entities.Role{
				ID:   1,
				Name: "Admin",
				Code: "ADMIN",
				Description: sql.NullString{
					String: "Admin Role",
					Valid:  true,
				},
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, userID string, updatedBy string, request presenter.UpdateUserRequest, existingUser entities.User, existingRole entities.Role) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock user existence check
				mock.ExpectQuery(`SELECT`).
					WithArgs(userID).
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

				// Mock role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(request.RoleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(existingRole.ID, existingRole.Name, existingRole.Code, existingRole.Description))

				// No user role update because role didn't change

				// Mock update user query
				mock.ExpectExec(`UPDATE app_user`).
					WithArgs(
						request.Name,
						request.FullName,
						request.Email,
						sql.NullString{String: phone, Valid: true},
						true,
						updatedBy,
						userID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))

				// Mock transaction commit
				mock.ExpectCommit()
			},
		},
		{
			name:      "success update user with empty phone",
			userID:    "user-123",
			updatedBy: "admin-123",
			request: presenter.UpdateUserRequest{
				Name:     "updated_john",
				FullName: "Updated John Doe",
				Email:    "updated.john@example.com",
				RoleID:   1,
				Phone:    &emptyPhone,
				Active:   &activeFalse,
			},
			existingUser: entities.User{
				ID:       "user-123",
				Name:     "john_doe",
				FullName: "John Doe",
				Email:    "john@example.com",
				Phone:    sql.NullString{String: "08111111111", Valid: true},
				Active:   true,
				Role:     "Admin",
				RoleID:   1,
			},
			existingRole: entities.Role{
				ID:   1,
				Name: "Admin",
				Code: "ADMIN",
				Description: sql.NullString{
					String: "Admin Role",
					Valid:  true,
				},
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, userID string, updatedBy string, request presenter.UpdateUserRequest, existingUser entities.User, existingRole entities.Role) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock user existence check
				mock.ExpectQuery(`SELECT`).
					WithArgs(userID).
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

				// Mock role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(request.RoleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(existingRole.ID, existingRole.Name, existingRole.Code, existingRole.Description))

				// Mock update user query - with empty phone (Valid: false)
				mock.ExpectExec(`UPDATE app_user`).
					WithArgs(
						request.Name,
						request.FullName,
						request.Email,
						sql.NullString{Valid: false}, // empty phone
						false,
						updatedBy,
						userID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))

				// Mock transaction commit
				mock.ExpectCommit()
			},
		},
		{
			name:      "success update user with nil phone and nil active",
			userID:    "user-123",
			updatedBy: "admin-123",
			request: presenter.UpdateUserRequest{
				Name:     "updated_john",
				FullName: "Updated John Doe",
				Email:    "updated.john@example.com",
				RoleID:   1,
				Phone:    nil,
				Active:   nil,
			},
			existingUser: entities.User{
				ID:       "user-123",
				Name:     "john_doe",
				FullName: "John Doe",
				Email:    "john@example.com",
				Phone:    sql.NullString{String: "08111111111", Valid: true},
				Active:   true,
				Role:     "Admin",
				RoleID:   1,
			},
			existingRole: entities.Role{
				ID:   1,
				Name: "Admin",
				Code: "ADMIN",
				Description: sql.NullString{
					String: "Admin Role",
					Valid:  true,
				},
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, userID string, updatedBy string, request presenter.UpdateUserRequest, existingUser entities.User, existingRole entities.Role) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock user existence check
				mock.ExpectQuery(`SELECT`).
					WithArgs(userID).
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

				// Mock role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(request.RoleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(existingRole.ID, existingRole.Name, existingRole.Code, existingRole.Description))

				// Mock update user query - with existing phone and active values
				mock.ExpectExec(`UPDATE app_user`).
					WithArgs(
						request.Name,
						request.FullName,
						request.Email,
						existingUser.Phone,  // keep existing phone
						existingUser.Active, // keep existing active
						updatedBy,
						userID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))

				// Mock transaction commit
				mock.ExpectCommit()
			},
		},
		{
			name:      "user not found",
			userID:    "user-999",
			updatedBy: "admin-123",
			request: presenter.UpdateUserRequest{
				Name:     "test_user",
				FullName: "Test User",
				Email:    "test@example.com",
				RoleID:   1,
				Phone:    &phone,
				Active:   &activeTrue,
			},
			existingUser: entities.User{},
			existingRole: entities.Role{},
			wantErr:      true,
			expectedErr:  constant.ErrUserIdNotFound,
			setupMock: func(mock sqlmock.Sqlmock, userID string, updatedBy string, request presenter.UpdateUserRequest, existingUser entities.User, existingRole entities.Role) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock user existence check returns no rows
				mock.ExpectQuery(`SELECT`).
					WithArgs(userID).
					WillReturnError(sql.ErrNoRows)

				// Mock transaction rollback
				mock.ExpectRollback()
			},
		},
		{
			name:      "role not found",
			userID:    "user-123",
			updatedBy: "admin-123",
			request: presenter.UpdateUserRequest{
				Name:     "test_user",
				FullName: "Test User",
				Email:    "test@example.com",
				RoleID:   999,
				Phone:    &phone,
				Active:   &activeTrue,
			},
			existingUser: entities.User{
				ID:       "user-123",
				Name:     "john_doe",
				FullName: "John Doe",
				Email:    "john@example.com",
				Phone:    sql.NullString{String: "08111111111", Valid: true},
				Active:   true,
				Role:     "Admin",
				RoleID:   1,
			},
			existingRole: entities.Role{},
			wantErr:      true,
			expectedErr:  constant.ErrRoleIdNotFound,
			setupMock: func(mock sqlmock.Sqlmock, userID string, updatedBy string, request presenter.UpdateUserRequest, existingUser entities.User, existingRole entities.Role) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock user existence check
				mock.ExpectQuery(`SELECT`).
					WithArgs(userID).
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

				// Mock role existence check returns no rows
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(request.RoleID).
					WillReturnError(sql.ErrNoRows)

				// Mock transaction rollback
				mock.ExpectRollback()
			},
		},
		{
			name:      "database error on user check",
			userID:    "user-123",
			updatedBy: "admin-123",
			request: presenter.UpdateUserRequest{
				Name:     "test_user",
				FullName: "Test User",
				Email:    "test@example.com",
				RoleID:   1,
				Phone:    &phone,
				Active:   &activeTrue,
			},
			existingUser: entities.User{},
			existingRole: entities.Role{},
			wantErr:      true,
			expectedErr:  constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, userID string, updatedBy string, request presenter.UpdateUserRequest, existingUser entities.User, existingRole entities.Role) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock user existence check returns database error
				mock.ExpectQuery(`SELECT`).
					WithArgs(userID).
					WillReturnError(sql.ErrConnDone)

				// Mock transaction rollback
				mock.ExpectRollback()
			},
		},
		{
			name:      "database error on role check",
			userID:    "user-123",
			updatedBy: "admin-123",
			request: presenter.UpdateUserRequest{
				Name:     "test_user",
				FullName: "Test User",
				Email:    "test@example.com",
				RoleID:   1,
				Phone:    &phone,
				Active:   &activeTrue,
			},
			existingUser: entities.User{
				ID:       "user-123",
				Name:     "john_doe",
				FullName: "John Doe",
				Email:    "john@example.com",
				Phone:    sql.NullString{String: "08111111111", Valid: true},
				Active:   true,
				Role:     "Admin",
				RoleID:   1,
			},
			existingRole: entities.Role{},
			wantErr:      true,
			expectedErr:  constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, userID string, updatedBy string, request presenter.UpdateUserRequest, existingUser entities.User, existingRole entities.Role) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock user existence check
				mock.ExpectQuery(`SELECT`).
					WithArgs(userID).
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

				// Mock role existence check returns database error
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(request.RoleID).
					WillReturnError(sql.ErrConnDone)

				// Mock transaction rollback
				mock.ExpectRollback()
			},
		},
		{
			name:      "database error on update user role",
			userID:    "user-123",
			updatedBy: "admin-123",
			request: presenter.UpdateUserRequest{
				Name:     "updated_john",
				FullName: "Updated John Doe",
				Email:    "updated.john@example.com",
				RoleID:   2, // Different role
				Phone:    &phone,
				Active:   &activeTrue,
			},
			existingUser: entities.User{
				ID:       "user-123",
				Name:     "john_doe",
				FullName: "John Doe",
				Email:    "john@example.com",
				Phone:    sql.NullString{String: "08111111111", Valid: true},
				Active:   true,
				Role:     "Admin", // Different from new role
				RoleID:   1,
			},
			existingRole: entities.Role{
				ID:   2,
				Name: "Manager",
				Code: "MANAGER",
				Description: sql.NullString{
					String: "Manager Role",
					Valid:  true,
				},
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, userID string, updatedBy string, request presenter.UpdateUserRequest, existingUser entities.User, existingRole entities.Role) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock user existence check
				mock.ExpectQuery(`SELECT`).
					WithArgs(userID).
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

				// Mock role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(request.RoleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(existingRole.ID, existingRole.Name, existingRole.Code, existingRole.Description))

				// Mock update user role query returns error
				mock.ExpectExec(`UPDATE app_user_role`).
					WithArgs(userID, request.RoleID).
					WillReturnError(sql.ErrConnDone)

				// Mock transaction rollback
				mock.ExpectRollback()
			},
		},
		{
			name:      "database error on update user",
			userID:    "user-123",
			updatedBy: "admin-123",
			request: presenter.UpdateUserRequest{
				Name:     "updated_john",
				FullName: "Updated John Doe",
				Email:    "updated.john@example.com",
				RoleID:   1,
				Phone:    &phone,
				Active:   &activeTrue,
			},
			existingUser: entities.User{
				ID:       "user-123",
				Name:     "john_doe",
				FullName: "John Doe",
				Email:    "john@example.com",
				Phone:    sql.NullString{String: "08111111111", Valid: true},
				Active:   true,
				Role:     "Admin",
				RoleID:   1,
			},
			existingRole: entities.Role{
				ID:   1,
				Name: "Admin",
				Code: "ADMIN",
				Description: sql.NullString{
					String: "Admin Role",
					Valid:  true,
				},
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, userID string, updatedBy string, request presenter.UpdateUserRequest, existingUser entities.User, existingRole entities.Role) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock user existence check
				mock.ExpectQuery(`SELECT`).
					WithArgs(userID).
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

				// Mock role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(request.RoleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(existingRole.ID, existingRole.Name, existingRole.Code, existingRole.Description))

				// Mock update user query returns error
				mock.ExpectExec(`UPDATE app_user`).
					WithArgs(
						request.Name,
						request.FullName,
						request.Email,
						sql.NullString{String: phone, Valid: true},
						true,
						updatedBy,
						userID,
					).
					WillReturnError(sql.ErrConnDone)

				// Mock transaction rollback
				mock.ExpectRollback()
			},
		},
		{
			name:      "database error on transaction begin",
			userID:    "user-123",
			updatedBy: "admin-123",
			request: presenter.UpdateUserRequest{
				Name:     "test_user",
				FullName: "Test User",
				Email:    "test@example.com",
				RoleID:   1,
				Phone:    &phone,
				Active:   &activeTrue,
			},
			existingUser: entities.User{},
			existingRole: entities.Role{},
			wantErr:      true,
			expectedErr:  constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, userID string, updatedBy string, request presenter.UpdateUserRequest, existingUser entities.User, existingRole entities.Role) {
				// Mock transaction begin returns error
				mock.ExpectBegin().WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:      "database error on transaction commit",
			userID:    "user-123",
			updatedBy: "admin-123",
			request: presenter.UpdateUserRequest{
				Name:     "updated_john",
				FullName: "Updated John Doe",
				Email:    "updated.john@example.com",
				RoleID:   1,
				Phone:    &phone,
				Active:   &activeTrue,
			},
			existingUser: entities.User{
				ID:       "user-123",
				Name:     "john_doe",
				FullName: "John Doe",
				Email:    "john@example.com",
				Phone:    sql.NullString{String: "08111111111", Valid: true},
				Active:   true,
				Role:     "Admin",
				RoleID:   1,
			},
			existingRole: entities.Role{
				ID:   1,
				Name: "Admin",
				Code: "ADMIN",
				Description: sql.NullString{
					String: "Admin Role",
					Valid:  true,
				},
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, userID string, updatedBy string, request presenter.UpdateUserRequest, existingUser entities.User, existingRole entities.Role) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock user existence check
				mock.ExpectQuery(`SELECT`).
					WithArgs(userID).
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

				// Mock role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id =`).
					WithArgs(request.RoleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(existingRole.ID, existingRole.Name, existingRole.Code, existingRole.Description))

				// Mock update user query
				mock.ExpectExec(`UPDATE app_user`).
					WithArgs(
						request.Name,
						request.FullName,
						request.Email,
						sql.NullString{String: phone, Valid: true},
						true,
						updatedBy,
						userID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))

				// Mock transaction commit returns error
				mock.ExpectCommit().WillReturnError(sql.ErrConnDone)

			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, mock, cleanup := newTestService(t)
			defer cleanup()

			ctx := context.Background()

			tc.setupMock(mock, tc.userID, tc.updatedBy, tc.request, tc.existingUser, tc.existingRole)

			err := svc.UpdateUserService(ctx, tc.request, tc.updatedBy, tc.userID)

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
