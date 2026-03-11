package usecase_test

import (
	"context"
	"testing"

	"be-dashboard-nba/constant"
	menuDomain "be-dashboard-nba/internal/modules/menu/domain"
	"be-dashboard-nba/internal/modules/menu_permission/domain"  // Sesuaikan path jika menggunakan /modules/
	"be-dashboard-nba/internal/modules/menu_permission/usecase" // Sesuaikan path

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMenuPermissionUsecase(t *testing.T) {
	ctx := context.Background()
	nopLogger := zerolog.Nop()

	errDb := errors.New("database connection error")

	payload := domain.MenuPermissionUpdatePayload{
		ID:         77,
		Code:       "UPDATE_ROLE",
		ActionName: "Update Role Detail",
		UpdatedBy:  "user-uuid-1234",
	}

	tests := []struct {
		name        string
		payload     domain.MenuPermissionUpdatePayload
		mockSetup   func(mockPermRepo *domain.MockMenuPermissionRepository)
		wantErr     bool
		expectedErr error
	}{
		{
			name:    "success - update menu permission",
			payload: payload,
			mockSetup: func(mockPermRepo *domain.MockMenuPermissionRepository) {
				// 1. Ekspektasi pengecekan ID berhasil (data ditemukan)
				mockPermRepo.On("ReadMenuPermissionByIDQuery", ctx, payload.ID).
					Return(domain.MenuPermissionDetail{ID: payload.ID}, nil).Once()

				// 2. Ekspektasi update query berhasil
				mockPermRepo.On("UpdateMenuPermissionQuery", ctx, payload).
					Return(nil).Once()
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:    "failed - menu permission id not found",
			payload: payload,
			mockSetup: func(mockPermRepo *domain.MockMenuPermissionRepository) {
				// Simulasi data tidak ditemukan
				mockPermRepo.On("ReadMenuPermissionByIDQuery", ctx, payload.ID).
					Return(domain.MenuPermissionDetail{}, constant.ErrDataNotFound).Once()

				// UpdateMenuPermissionQuery TIDAK BOLEH dipanggil
			},
			wantErr:     true,
			expectedErr: constant.ErrMenuPermissionIdNotFound,
		},
		{
			name:    "failed - error reading menu permission detail",
			payload: payload,
			mockSetup: func(mockPermRepo *domain.MockMenuPermissionRepository) {
				// Simulasi error saat membaca DB
				mockPermRepo.On("ReadMenuPermissionByIDQuery", ctx, payload.ID).
					Return(domain.MenuPermissionDetail{}, errDb).Once()
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
		},
		{
			name:    "failed - error executing update query",
			payload: payload,
			mockSetup: func(mockPermRepo *domain.MockMenuPermissionRepository) {
				// Pengecekan ID berhasil
				mockPermRepo.On("ReadMenuPermissionByIDQuery", ctx, payload.ID).
					Return(domain.MenuPermissionDetail{ID: payload.ID}, nil).Once()

				// Simulasi error saat melakukan update
				mockPermRepo.On("UpdateMenuPermissionQuery", ctx, payload).
					Return(errDb).Once()
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
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
			err := uc.UpdateMenuPermissionUsecase(ctx, tc.payload)

			// Validasi Error
			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
			} else {
				assert.NoError(t, err)
			}

			// Validasi pemanggilan mock
			mockPermRepo.AssertExpectations(t)
		})
	}
}
