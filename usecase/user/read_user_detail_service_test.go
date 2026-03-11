package user

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/usecase/entities"
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestReadDetailUserService(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		expected    entities.User
		wantErr     bool
		expectedErr error
		setupMock   func(mock sqlmock.Sqlmock, userID string, expected entities.User)
	}{
		{
			name:   "success get user detail",
			userID: "user-123",
			expected: entities.User{
				ID:       "user-123",
				Name:     "john_doe",
				FullName: "John Doe",
				Email:    "john.doe@example.com",
				Phone:    sql.NullString{String: "08123456789", Valid: true},
				Active:   true,
				Role:     "Admin",
				RoleID:   1,
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, userID string, expected entities.User) {
				// Mock read detail user query
				mock.ExpectQuery(`SELECT`).
					WithArgs(userID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "name", "full_name", "email", "phone", "active", "role", "role_id",
					}).
						AddRow(
							expected.ID,
							expected.Name,
							expected.FullName,
							expected.Email,
							expected.Phone,
							expected.Active,
							expected.Role,
							expected.RoleID,
						))
			},
		},
		{
			name:   "success get user detail without phone",
			userID: "user-456",
			expected: entities.User{
				ID:       "user-456",
				Name:     "jane_smith",
				FullName: "Jane Smith",
				Email:    "jane.smith@example.com",
				Phone:    sql.NullString{Valid: false},
				Active:   false,
				Role:     "User",
				RoleID:   1,
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, userID string, expected entities.User) {
				// Mock read detail user query
				mock.ExpectQuery(`SELECT`).
					WithArgs(userID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "name", "full_name", "email", "phone", "active", "role", "role_id",
					}).
						AddRow(
							expected.ID,
							expected.Name,
							expected.FullName,
							expected.Email,
							nil, // null phone
							expected.Active,
							expected.Role,
							expected.RoleID,
						))
			},
		},
		{
			name:        "user not found",
			userID:      "user-999",
			expected:    entities.User{},
			wantErr:     true,
			expectedErr: constant.ErrUserIdNotFound,
			setupMock: func(mock sqlmock.Sqlmock, userID string, expected entities.User) {
				// Mock read detail user query returns no rows
				mock.ExpectQuery(`SELECT`).
					WithArgs(userID).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:        "database error",
			userID:      "user-123",
			expected:    entities.User{},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, userID string, expected entities.User) {
				// Mock read detail user query returns database error
				mock.ExpectQuery(`SELECT`).
					WithArgs(userID).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, mock, cleanup := newTestService(t)
			defer cleanup()

			ctx := context.Background()

			tc.setupMock(mock, tc.userID, tc.expected)

			result, err := svc.ReadDetailUserService(ctx, tc.userID)

			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					if tc.expectedErr == sql.ErrNoRows {
						// Untuk ErrNoRows, service mengembalikan error asli tanpa wrapping
						assert.ErrorIs(t, err, sql.ErrNoRows)
					} else {
						assert.ErrorIs(t, err, tc.expectedErr)
					}
				}
				assert.Equal(t, entities.User{}, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.ID, result.ID)
				assert.Equal(t, tc.expected.Name, result.Name)
				assert.Equal(t, tc.expected.FullName, result.FullName)
				assert.Equal(t, tc.expected.Email, result.Email)
				assert.Equal(t, tc.expected.Phone, result.Phone)
				assert.Equal(t, tc.expected.Active, result.Active)
				assert.Equal(t, tc.expected.Role, result.Role)
			}

			// Check if all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
