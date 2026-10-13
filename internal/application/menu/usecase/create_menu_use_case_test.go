package usecase

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ABA-Developer/go-backend-template/constant"
	"github.com/ABA-Developer/go-backend-template/internal/application/menu/dto"
	repoParams "github.com/ABA-Developer/go-backend-template/internal/domain/contract/repository/params"
	mockrepo "github.com/ABA-Developer/go-backend-template/internal/test/mock/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func TestCreateMenuUseCase(t *testing.T) {
	parentID := 10

	tests := []struct {
		name    string
		req     dto.CreateMenuPayload
		setup   func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock)
		wantErr error
	}{
		{
			name: "positive/with-parent",
			req: dto.CreateMenuPayload{
				Name:     "Dashboard",
				Group:    "main",
				ParentID: &parentID,
				Active:   ptrBool(true),
				Display:  ptrBool(true),
			},
			setup: func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock) {
				m.ExpectBegin()
				m.ExpectCommit()

				menuRepo.On("ReadNextSortForParentQry", testifymock.Anything, int32(parentID)).Return(5, nil).Once()
				menuRepo.
					On("CreateMenuQry", testifymock.Anything, testifymock.Anything).
					Run(func(args testifymock.Arguments) {
						params := args.Get(1).(repoParams.CreateMenuParams)
						assert.Equal(t, 5, params.Sort)
						assert.True(t, params.ParentID.Valid)
						assert.Equal(t, int32(parentID), params.ParentID.Int32)
						assert.Equal(t, "main", params.Group)
					}).
					Return(nil).
					Once()
			},
		},
		{
			name: "edge/without-parent-uses-group-sort",
			req: dto.CreateMenuPayload{
				Name:    "Home",
				Group:   "main",
				Active:  ptrBool(true),
				Display: ptrBool(true),
			},
			setup: func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock) {
				m.ExpectBegin()
				m.ExpectCommit()

				menuRepo.On("ReadSortForGroupQry", testifymock.Anything, "main").Return(2, nil).Once()
				menuRepo.
					On("CreateMenuQry", testifymock.Anything, testifymock.Anything).
					Run(func(args testifymock.Arguments) {
						params := args.Get(1).(repoParams.CreateMenuParams)
						assert.Equal(t, 2, params.Sort)
						assert.False(t, params.ParentID.Valid)
						assert.Equal(t, "main", params.Group)
					}).
					Return(nil).
					Once()
			},
		},
		{
			name: "negative/sort-query-error",
			req: dto.CreateMenuPayload{
				Name:    "Home",
				Group:   "main",
				Active:  ptrBool(true),
				Display: ptrBool(true),
			},
			setup: func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				menuRepo.On("ReadSortForGroupQry", testifymock.Anything, "main").Return(0, errors.New("db down")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name: "negative/create-error",
			req: dto.CreateMenuPayload{
				Name:    "Home",
				Group:   "main",
				Active:  ptrBool(true),
				Display: ptrBool(true),
			},
			setup: func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				menuRepo.On("ReadSortForGroupQry", testifymock.Anything, "main").Return(2, nil).Once()
				menuRepo.On("CreateMenuQry", testifymock.Anything, testifymock.Anything).Return(errors.New("insert failed")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name: "negative/begin-tx-error",
			req:  dto.CreateMenuPayload{Name: "X", Group: "main"},
			setup: func(m sqlmock.Sqlmock, _ *mockrepo.MenuRepositoryMock) {
				m.ExpectBegin().WillReturnError(errors.New("begin failed"))
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name: "edge/rollback-error",
			req:  dto.CreateMenuPayload{Name: "X", Group: "main"},
			setup: func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback().WillReturnError(sql.ErrConnDone)

				menuRepo.On("ReadSortForGroupQry", testifymock.Anything, "main").Return(0, errors.New("sort failed")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menuRepo := mockrepo.NewMenuRepository()
			uc, mock, cleanup := newTestUseCaseWithMenuRepo(t, menuRepo)
			defer cleanup()

			tt.setup(mock, menuRepo)

			err := uc.CreateMenuUseCase(context.Background(), tt.req, "user-1")
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func ptrBool(v bool) *bool { return &v }
