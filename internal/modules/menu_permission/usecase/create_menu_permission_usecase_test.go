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

func TestCreateMenuPermissionUsecase(t *testing.T) {
	ctx := context.Background()
	nopLogger := zerolog.Nop()

	payload := domain.MenuPermissionCreatePayload{
		MenuID:     10,
		Code:       "CREATE_USER",
		ActionName: "Create User",
	}

	errDb := errors.New("database connection error")

	tests := []struct {
		name        string
		payload     domain.MenuPermissionCreatePayload
		mockSetup   func(mockMenuRepo *menuDomain.MockMenuRepository, mockPermRepo *domain.MockMenuPermissionRepository)
		wantErr     bool
		expectedErr error
	}{
		{
			name:    "success - create menu permission",
			payload: payload,
			mockSetup: func(mockMenuRepo *menuDomain.MockMenuRepository, mockPermRepo *domain.MockMenuPermissionRepository) {
				// 1. Ekspektasi pengecekan MenuID berhasil (menu ditemukan)
				mockMenuRepo.On("ReadMenuByIDQuery", ctx, payload.MenuID).
					Return(menuDomain.Menu{ID: payload.MenuID}, nil).Once()

				// 2. Ekspektasi insert query berhasil
				mockPermRepo.On("CreateMenuPermissionQuery", ctx, payload).
					Return(nil).Once()
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:    "failed - menu id not found",
			payload: payload,
			mockSetup: func(mockMenuRepo *menuDomain.MockMenuRepository, mockPermRepo *domain.MockMenuPermissionRepository) {
				// Simulasi error data not found dari Menu Repository
				mockMenuRepo.On("ReadMenuByIDQuery", ctx, payload.MenuID).
					Return(menuDomain.Menu{}, constant.ErrDataNotFound).Once()

				// CreateMenuPermissionQuery TIDAK BOLEH dipanggil karena proses harus terhenti
			},
			wantErr:     true,
			expectedErr: constant.ErrMenuIdNotFound,
		},
		{
			name:    "failed - error reading menu detail",
			payload: payload,
			mockSetup: func(mockMenuRepo *menuDomain.MockMenuRepository, mockPermRepo *domain.MockMenuPermissionRepository) {
				// Simulasi error database umum saat mengecek menu
				mockMenuRepo.On("ReadMenuByIDQuery", ctx, payload.MenuID).
					Return(menuDomain.Menu{}, errDb).Once()
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
		},
		{
			name:    "failed - error inserting menu permission",
			payload: payload,
			mockSetup: func(mockMenuRepo *menuDomain.MockMenuRepository, mockPermRepo *domain.MockMenuPermissionRepository) {
				// Pengecekan menu berhasil
				mockMenuRepo.On("ReadMenuByIDQuery", ctx, payload.MenuID).
					Return(menuDomain.Menu{ID: payload.MenuID}, nil).Once()

				// Simulasi error saat melakukan insert permission
				mockPermRepo.On("CreateMenuPermissionQuery", ctx, payload).
					Return(errDb).Once()
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockMenuRepo := new(menuDomain.MockMenuRepository)
			mockPermRepo := new(domain.MockMenuPermissionRepository)

			tc.mockSetup(mockMenuRepo, mockPermRepo)

			// Inisialisasi usecase dengan kedua mock repository
			uc := usecase.NewMenuPermissionUsecase(mockPermRepo, mockMenuRepo, &nopLogger)

			err := uc.CreateMenuPermissionUsecase(ctx, tc.payload)

			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
			} else {
				assert.NoError(t, err)
			}

			// Validasi bahwa semua ekspektasi mock terpenuhi
			mockMenuRepo.AssertExpectations(t)
			mockPermRepo.AssertExpectations(t)
		})
	}
}
