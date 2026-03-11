package usecase_test

import (
	"context"
	"testing"

	"be-dashboard-nba/constant"
	menuDomain "be-dashboard-nba/internal/modules/menu/domain"
	"be-dashboard-nba/internal/modules/menu_permission/domain"
	"be-dashboard-nba/internal/modules/menu_permission/usecase"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestReadMenuPermissionDetailUsecase(t *testing.T) {
	ctx := context.Background()
	nopLogger := zerolog.Nop()

	menuPermissionID := 5
	errDb := errors.New("query execution error")

	// Mock data yang diharapkan kembali dari database
	expectedDetail := domain.MenuPermissionDetail{
		ID:         menuPermissionID,
		MenuID:     10,
		Code:       "READ_DETAIL",
		ActionName: "Read Detail Data",
	}

	tests := []struct {
		name         string
		id           int
		mockSetup    func(mockPermRepo *domain.MockMenuPermissionRepository)
		expectedData domain.MenuPermissionDetail
		wantErr      bool
		expectedErr  error
	}{
		{
			name: "success - get menu permission detail",
			id:   menuPermissionID,
			mockSetup: func(mockPermRepo *domain.MockMenuPermissionRepository) {
				mockPermRepo.On("ReadMenuPermissionByIDQuery", ctx, menuPermissionID).
					Return(expectedDetail, nil).Once()
			},
			expectedData: expectedDetail,
			wantErr:      false,
			expectedErr:  nil,
		},
		{
			name: "failed - menu permission id not found",
			id:   menuPermissionID,
			mockSetup: func(mockPermRepo *domain.MockMenuPermissionRepository) {
				mockPermRepo.On("ReadMenuPermissionByIDQuery", ctx, menuPermissionID).
					Return(domain.MenuPermissionDetail{}, constant.ErrDataNotFound).Once()
			},
			expectedData: domain.MenuPermissionDetail{},
			wantErr:      true,
			expectedErr:  constant.ErrMenuPermissionIdNotFound, // Pastikan translasi error berhasil
		},
		{
			name: "failed - error reading from database",
			id:   menuPermissionID,
			mockSetup: func(mockPermRepo *domain.MockMenuPermissionRepository) {
				mockPermRepo.On("ReadMenuPermissionByIDQuery", ctx, menuPermissionID).
					Return(domain.MenuPermissionDetail{}, errDb).Once()
			},
			expectedData: domain.MenuPermissionDetail{},
			wantErr:      true,
			expectedErr:  constant.ErrUnknownSource, // Pastikan wrapping error berhasil
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Inisialisasi Mocks
			mockPermRepo := new(domain.MockMenuPermissionRepository)
			mockMenuRepo := new(menuDomain.MockMenuRepository) // Dummy mock untuk kebutuhan constructor

			tc.mockSetup(mockPermRepo)

			// Inisialisasi Usecase
			uc := usecase.NewMenuPermissionUsecase(mockPermRepo, mockMenuRepo, &nopLogger)

			// Eksekusi fungsi Usecase
			result, err := uc.ReadMenuPermissionDetailUsecase(ctx, tc.id)

			// Validasi hasil
			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedData, result)
			}

			// Validasi pemanggilan mock
			mockPermRepo.AssertExpectations(t)
		})
	}
}
