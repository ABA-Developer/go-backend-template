package usecase_test

import (
	"context"
	"testing"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/menu/domain"
	"be-dashboard-nba/internal/modules/menu/usecase" // Sesuaikan path

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestDeleteMenuUsecase(t *testing.T) {
	ctx := context.Background()
	menuID := 15 // Contoh ID menu yang akan dihapus
	nopLogger := zerolog.Nop()

	tests := []struct {
		name        string
		menuID      int
		mockSetup   func(mockRepo *domain.MockMenuRepository)
		wantErr     bool
		expectedErr error
	}{
		{
			name:   "success - delete menu",
			menuID: menuID,
			mockSetup: func(mockRepo *domain.MockMenuRepository) {
				// Simulasi Repo berhasil mengeksekusi Delete
				mockRepo.On("DeleteMenuQuery", ctx, menuID).Return(nil)
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:   "error - menu not found (rows affected 0)",
			menuID: menuID,
			mockSetup: func(mockRepo *domain.MockMenuRepository) {
				// Simulasi Repo tidak menemukan data (mengembalikan ErrDataNotFound)
				mockRepo.On("DeleteMenuQuery", ctx, menuID).Return(constant.ErrDataNotFound)
			},
			wantErr: true,
			// Ekspektasi: Usecase harus menerjemahkannya menjadi ErrMenuIdNotFound
			expectedErr: constant.ErrMenuIdNotFound,
		},
		{
			name:   "error - repository internal error",
			menuID: menuID,
			mockSetup: func(mockRepo *domain.MockMenuRepository) {
				// Simulasi Repo mengalami error database yang tidak terduga
				mockRepo.On("DeleteMenuQuery", ctx, menuID).Return(errors.New("connection lost"))
			},
			wantErr: true,
			// Ekspektasi: Usecase harus membungkusnya menjadi ErrUnknownSource
			expectedErr: constant.ErrUnknownSource,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(domain.MockMenuRepository)
			tc.mockSetup(mockRepo)

			// Inisialisasi Usecase dengan mock repo dan logger kosong
			uc := usecase.NewMenuUsecase(mockRepo, &nopLogger)

			err := uc.DeleteMenuUsecase(ctx, tc.menuID)

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			// Validasi bahwa method di mockRepo dipanggil sesuai skenario
			mockRepo.AssertExpectations(t)
		})
	}
}
