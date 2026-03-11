package user

import (
	userPresenter "be-dashboard-nba/api/presenter/user"
	"be-dashboard-nba/constant"
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserService(t *testing.T) {
	// Set environment variable for bcrypt
	os.Setenv("AUTH_BCRYPT_COST", "10")

	phone := "08123456789"
	activeTrue := true
	activeFalse := false

	tests := []struct {
		name        string
		userID      string
		request     userPresenter.CreateUserRequest
		wantErr     bool
		expectedErr error
		setupMock   func(mock sqlmock.Sqlmock, userID string, request userPresenter.CreateUserRequest)
	}{
		{
			name:   "success create user with phone and active true",
			userID: "admin-123",
			request: userPresenter.CreateUserRequest{
				Name:     "john_doe",
				FullName: "John Doe",
				Email:    "john.doe@example.com",
				Password: "password123",
				RoleID:   1,
				Phone:    &phone,
				Active:   &activeTrue,
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, userID string, request userPresenter.CreateUserRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock role existence check - use exact SQL matching
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id = \$1`).
					WithArgs(request.RoleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(1, "Admin", "ADMIN", "Administrator Role"))

				// Mock create user query
				mock.ExpectQuery(`INSERT INTO app_user`).
					WithArgs(
						request.Name,
						request.FullName,
						request.Email,
						sqlmock.AnyArg(), // hashed password
						true,             // active
						sql.NullString{String: phone, Valid: true},
						userID,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("user-123"))

				// Mock create user role query
				mock.ExpectExec(`INSERT INTO app_user_role`).
					WithArgs("user-123", request.RoleID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Mock transaction commit
				mock.ExpectCommit()
			},
		},
		{
			name:   "success create user without phone and active false",
			userID: "admin-456",
			request: userPresenter.CreateUserRequest{
				Name:     "jane_smith",
				FullName: "Jane Smith",
				Email:    "jane.smith@example.com",
				Password: "password456",
				RoleID:   2,
				Phone:    nil,
				Active:   &activeFalse,
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, userID string, request userPresenter.CreateUserRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id = \$1`).
					WithArgs(request.RoleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(2, "User", "USER", "User Role"))

				// Mock create user query - without phone
				mock.ExpectQuery(`INSERT INTO app_user`).
					WithArgs(
						request.Name,
						request.FullName,
						request.Email,
						sqlmock.AnyArg(),             // hashed password
						false,                        // active
						sql.NullString{Valid: false}, // null phone
						userID,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("user-456"))

				// Mock create user role query
				mock.ExpectExec(`INSERT INTO app_user_role`).
					WithArgs("user-456", request.RoleID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Mock transaction commit
				mock.ExpectCommit()
			},
		},
		{
			name:   "success create user with default active true",
			userID: "admin-789",
			request: userPresenter.CreateUserRequest{
				Name:     "bob_wilson",
				FullName: "Bob Wilson",
				Email:    "bob.wilson@example.com",
				Password: "password789",
				RoleID:   3,
				Phone:    &phone,
				Active:   nil, // Should default to true
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, userID string, request userPresenter.CreateUserRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id = \$1`).
					WithArgs(request.RoleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(3, "Manager", "MANAGER", "Manager Role"))

				// Mock create user query - default active true
				mock.ExpectQuery(`INSERT INTO app_user`).
					WithArgs(
						request.Name,
						request.FullName,
						request.Email,
						sqlmock.AnyArg(), // hashed password
						true,             // default active true
						sql.NullString{String: phone, Valid: true},
						userID,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("user-789"))

				// Mock create user role query
				mock.ExpectExec(`INSERT INTO app_user_role`).
					WithArgs("user-789", request.RoleID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Mock transaction commit
				mock.ExpectCommit()
			},
		},
		{
			name:   "role not found",
			userID: "admin-123",
			request: userPresenter.CreateUserRequest{
				Name:     "test_user",
				FullName: "Test User",
				Email:    "test@example.com",
				Password: "password123",
				RoleID:   999,
				Phone:    &phone,
				Active:   &activeTrue,
			},
			wantErr:     true,
			expectedErr: constant.ErrRoleIdNotFound,
			setupMock: func(mock sqlmock.Sqlmock, userID string, request userPresenter.CreateUserRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock role existence check returns no rows
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id = \$1`).
					WithArgs(request.RoleID).
					WillReturnError(sql.ErrNoRows)

				// Mock transaction rollback
				mock.ExpectRollback()
			},
		},
		{
			name:   "database error on role check",
			userID: "admin-123",
			request: userPresenter.CreateUserRequest{
				Name:     "test_user",
				FullName: "Test User",
				Email:    "test@example.com",
				Password: "password123",
				RoleID:   1,
				Phone:    &phone,
				Active:   &activeTrue,
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, userID string, request userPresenter.CreateUserRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock role existence check returns database error
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id = \$1`).
					WithArgs(request.RoleID).
					WillReturnError(sql.ErrConnDone)

				// Mock transaction rollback
				mock.ExpectRollback()
			},
		},
		{
			name:   "database error on create user",
			userID: "admin-123",
			request: userPresenter.CreateUserRequest{
				Name:     "test_user",
				FullName: "Test User",
				Email:    "test@example.com",
				Password: "password123",
				RoleID:   1,
				Phone:    &phone,
				Active:   &activeTrue,
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, userID string, request userPresenter.CreateUserRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id = \$1`).
					WithArgs(request.RoleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(1, "Admin", "ADMIN", "Administrator Role"))

				// Mock create user query returns error
				mock.ExpectQuery(`INSERT INTO app_user`).
					WithArgs(
						request.Name,
						request.FullName,
						request.Email,
						sqlmock.AnyArg(),
						true,
						sql.NullString{String: phone, Valid: true},
						userID,
					).
					WillReturnError(sql.ErrConnDone)

				// Mock transaction rollback
				mock.ExpectRollback()
			},
		},
		{
			name:   "database error on create user role",
			userID: "admin-123",
			request: userPresenter.CreateUserRequest{
				Name:     "test_user",
				FullName: "Test User",
				Email:    "test@example.com",
				Password: "password123",
				RoleID:   1,
				Phone:    &phone,
				Active:   &activeTrue,
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, userID string, request userPresenter.CreateUserRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id = \$1`).
					WithArgs(request.RoleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(1, "Admin", "ADMIN", "Administrator Role"))

				// Mock create user query success
				mock.ExpectQuery(`INSERT INTO app_user`).
					WithArgs(
						request.Name,
						request.FullName,
						request.Email,
						sqlmock.AnyArg(),
						true,
						sql.NullString{String: phone, Valid: true},
						userID,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("user-123"))

				// Mock create user role query returns error
				mock.ExpectExec(`INSERT INTO app_user_role`).
					WithArgs("user-123", request.RoleID).
					WillReturnError(sql.ErrConnDone)

				// Mock transaction rollback
				mock.ExpectRollback()
			},
		},
		{
			name:   "database error on transaction begin",
			userID: "admin-123",
			request: userPresenter.CreateUserRequest{
				Name:     "test_user",
				FullName: "Test User",
				Email:    "test@example.com",
				Password: "password123",
				RoleID:   1,
				Phone:    &phone,
				Active:   &activeTrue,
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, userID string, request userPresenter.CreateUserRequest) {
				// Mock transaction begin returns error
				mock.ExpectBegin().WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:   "database error on transaction commit",
			userID: "admin-123",
			request: userPresenter.CreateUserRequest{
				Name:     "test_user",
				FullName: "Test User",
				Email:    "test@example.com",
				Password: "password123",
				RoleID:   1,
				Phone:    &phone,
				Active:   &activeTrue,
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, userID string, request userPresenter.CreateUserRequest) {
				// Mock transaction begin
				mock.ExpectBegin()

				// Mock role existence check
				mock.ExpectQuery(`SELECT id, name, code, description FROM app_role WHERE id = \$1`).
					WithArgs(request.RoleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "description"}).
						AddRow(1, "Admin", "ADMIN", "Administrator Role"))

				// Mock create user query
				mock.ExpectQuery(`INSERT INTO app_user`).
					WithArgs(
						request.Name,
						request.FullName,
						request.Email,
						sqlmock.AnyArg(),
						true,
						sql.NullString{String: phone, Valid: true},
						userID,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("user-123"))

				// Mock create user role query
				mock.ExpectExec(`INSERT INTO app_user_role`).
					WithArgs("user-123", request.RoleID).
					WillReturnResult(sqlmock.NewResult(1, 1))

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

			tc.setupMock(mock, tc.userID, tc.request)

			err := svc.CreateUserService(ctx, tc.request, tc.userID)

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

// Test untuk error generating hash password secara terpisah
func TestCreateUserService_HashPasswordError(t *testing.T) {
	// Set environment variable yang tidak valid untuk menyebabkan error
	os.Setenv("AUTH_BCRYPT_COST", "invalid")

	phone := "08123456789"
	activeTrue := true

	request := userPresenter.CreateUserRequest{
		Name:     "test_user",
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "password123",
		RoleID:   1,
		Phone:    &phone,
		Active:   &activeTrue,
	}

	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	ctx := context.Background()

	// Mock transaction begin dan rollback karena akan error di hash password
	mock.ExpectBegin()
	mock.ExpectRollback()

	err := svc.CreateUserService(ctx, request, "admin-123")

	assert.Error(t, err)
	assert.ErrorIs(t, err, constant.ErrUnknownSource)

	// Check if all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
