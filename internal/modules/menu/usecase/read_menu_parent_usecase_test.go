package usecase_test

import (
	"context"
	"testing"

	"be-dashboard-nba/internal/modules/menu/domain"
	"be-dashboard-nba/internal/modules/menu/usecase" // Sesuaikan path

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestReadMenuParentUsecase(t *testing.T) {
	ctx := context.Background()
	nopLogger := zerolog.Nop()

	// Data dummy simulasi dari database
	mockParentData := []domain.MenuParent{
		{
			ID:    1,
			Name:  "Dashboard Utama",
			Group: "admin",
		},
		{
			ID:    5,
			Name:  "Pengaturan Sistem",
			Group: "admin",
		},
	}

	// Buat variabel error spesifik agar bisa divalidasi oleh assert.ErrorIs
	errMockDB := errors.New("database timeout")

	tests := []struct {
		name         string
		mockSetup    func(mockRepo *domain.MockMenuRepository)
		wantErr      bool
		expectedErr  error
		expectedData []domain.MenuParent
	}{
		{
			name: "success - get parent menu list",
			mockSetup: func(mockRepo *domain.MockMenuRepository) {
				// Simulasi Repo berhasil mengambil array data
				mockRepo.On("ReadParentMenuQuery", ctx).Return(mockParentData, nil)
			},
			wantErr:      false,
			expectedErr:  nil,
			expectedData: mockParentData,
		},
		{
			name: "error - repository failed",
			mockSetup: func(mockRepo *domain.MockMenuRepository) {
				// Simulasi Repo gagal, kembalikan array kosong dan error mock
				mockRepo.On("ReadParentMenuQuery", ctx).Return([]domain.MenuParent(nil), errMockDB)
			},
			wantErr:      true,
			expectedErr:  errMockDB, // Usecase saat ini mengembalikan error mentah ini
			expectedData: nil,       // Jika error, data kembalian harus nil/kosong
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(domain.MockMenuRepository)
			tc.mockSetup(mockRepo)

			uc := usecase.NewMenuUsecase(mockRepo, &nopLogger)

			resultData, err := uc.ReadMenuParentUsecase(ctx)

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			// Validasi array data yang dikembalikan sama persis
			assert.Equal(t, tc.expectedData, resultData)

			mockRepo.AssertExpectations(t)
		})
	}
}
