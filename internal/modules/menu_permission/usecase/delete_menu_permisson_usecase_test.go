package usecase_test

import (
	"context"
	"testing"

	"be-dashboard-nba/constant"
	menuDomain "be-dashboard-nba/internal/modules/menu/domain"  // Untuk mockMenuRepo
	"be-dashboard-nba/internal/modules/menu_permission/domain"  // Sesuaikan path jika menggunakan /modules/
	"be-dashboard-nba/internal/modules/menu_permission/usecase" // Sesuaikan path

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestDeleteMenuPermissionUsecase(t *testing.T) {
	ctx := context.Background()
	nopLogger := zerolog.Nop()

	menuPermissionID := 99
	errDb := errors.New("database connection error")

	tests := []struct {
		name        string
		id          int
		mockSetup   func(mockPermRepo *domain.MockMenuPermissionRepository)
		wantErr     bool
		expectedErr error
	}{
		{
			name: "success - delete menu permission",
			id:   menuPermissionID,
			mockSetup: func(mockPermRepo *domain.MockMenuPermissionRepository) {
				// 1. Ekspektasi pengecekan ID berhasil (data ditemukan)
				mockPermRepo.On("ReadMenuPermissionByIDQuery", ctx, menuPermissionID).
					Return(domain.MenuPermissionDetail{ID: menuPermissionID}, nil).Once()

				// 2. Ekspektasi delete query berhasil
				mockPermRepo.On("DeleteMenuPermissionQuery", ctx, menuPermissionID).
					Return(nil).Once()
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "failed - menu permission id not found",
			id:   menuPermissionID,
			mockSetup: func(mockPermRepo *domain.MockMenuPermissionRepository) {
				// Simulasi data tidak ditemukan
				mockPermRepo.On("ReadMenuPermissionByIDQuery", ctx, menuPermissionID).
					Return(domain.MenuPermissionDetail{}, constant.ErrDataNotFound).Once()

				// DeleteMenuPermissionQuery TIDAK BOLEH dipanggil
			},
			wantErr:     true,
			expectedErr: constant.ErrMenuPermissionIdNotFound, // Error yang ditranslasikan oleh Usecase
		},
		{
			name: "failed - error reading menu permission detail",
			id:   menuPermissionID,
			mockSetup: func(mockPermRepo *domain.MockMenuPermissionRepository) {
				// Simulasi error database umum
				mockPermRepo.On("ReadMenuPermissionByIDQuery", ctx, menuPermissionID).
					Return(domain.MenuPermissionDetail{}, errDb).Once()
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
		},
		{
			name: "failed - error executing delete query",
			id:   menuPermissionID,
			mockSetup: func(mockPermRepo *domain.MockMenuPermissionRepository) {
				// Pengecekan ID berhasil
				mockPermRepo.On("ReadMenuPermissionByIDQuery", ctx, menuPermissionID).
					Return(domain.MenuPermissionDetail{ID: menuPermissionID}, nil).Once()

				// Simulasi error saat melakukan delete
				mockPermRepo.On("DeleteMenuPermissionQuery", ctx, menuPermissionID).
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
			mockMenuRepo := new(menuDomain.MockMenuRepository) // Diinisialisasi karena dibutuhkan oleh constructor Usecase

			tc.mockSetup(mockPermRepo)

			// Inisialisasi Usecase (Pastikan urutan parameter sesuai dengan constructor Anda)
			uc := usecase.NewMenuPermissionUsecase(mockPermRepo, mockMenuRepo, &nopLogger)

			// Eksekusi
			err := uc.DeleteMenuPermissionUsecase(ctx, tc.id)

			// Validasi Error
			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
			} else {
				assert.NoError(t, err)
			}

			// Validasi bahwa mock dieksekusi sesuai ekspektasi
			mockPermRepo.AssertExpectations(t)
		})
	}
}
