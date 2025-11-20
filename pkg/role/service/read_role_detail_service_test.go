package service

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/pkg/entities"
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestReadRoleDetailService(t *testing.T) {
	tests := []struct {
		name         string
		roleID       int
		expectedData entities.Role
		wantErr      bool
		expectedErr  error
		setupMock    func(mock sqlmock.Sqlmock, roleID int)
	}{
		{
			name:   "success",
			roleID: 1,
			expectedData: entities.Role{
				ID:   1,
				Name: "Admin",
				Code: "ADMIN",
				Description: sql.NullString{
					String: "Administrator Role",
					Valid:  true,
				},
			},
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, roleID int) {
				mock.ExpectQuery(`(?s)SELECT.*FROM.*app_role.*WHERE.*id`).
					WithArgs(roleID).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "name", "code", "description",
					}).AddRow(
						roleID, "Admin", "ADMIN", "Administrator Role",
					))
			},
		},
		{
			name:         "role id not found",
			roleID:       999,
			expectedData: entities.Role{},
			wantErr:      true,
			expectedErr:  constant.ErrRoleIdNotFound,
			setupMock: func(mock sqlmock.Sqlmock, roleID int) {
				mock.ExpectQuery(`(?s)SELECT.*FROM.*app_role.*WHERE.*id`).
					WithArgs(roleID).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:         "database error",
			roleID:       1,
			expectedData: entities.Role{},
			wantErr:      true,
			expectedErr:  constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, roleID int) {
				mock.ExpectQuery(`(?s)SELECT.*FROM.*app_role.*WHERE.*id`).
					WithArgs(roleID).
					WillReturnError(errors.New("sql error"))
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, mock, cleanup := newTestService(t)
			defer cleanup()

			ctx := context.Background()

			tc.setupMock(mock, tc.roleID)

			data, err := svc.ReadRoleDetail(ctx, tc.roleID)

			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
				assert.Equal(t, tc.expectedData, data)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedData, data)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
