package usecase_test

import (
	"context"
	"testing"

	"be-dashboard-nba/internal/modules/menu/domain"
	"be-dashboard-nba/internal/modules/menu/usecase" // Sesuaikan path jika perlu

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestReadListMenuUsecase(t *testing.T) {
	ctx := context.Background()
	nopLogger := zerolog.Nop()

	// Helper untuk membuat pointer
	ptr := func(s string) *string { return &s }

	// Simulasi parameter filter dari frontend (misal mencari menu "Settings")
	mockFilter := domain.MenuFilter{
		Search: "Settings",
		// Tambahkan field lain seperti Limit, Offset jika ada di struct Anda
	}

	// Simulasi data kembalian dari database
	mockData := []domain.Menu{
		{
			ID:    1,
			Name:  "General Settings",
			URL:   ptr("/settings/general"),
			Group: "settings",
		},
		{
			ID:    2,
			Name:  "User Settings",
			URL:   ptr("/settings/users"),
			Group: "settings",
		},
	}

	errDb := errors.New("failed to execute read list query")

	tests := []struct {
		name         string
		params       domain.MenuFilter
		mockSetup    func(mockRepo *domain.MockMenuRepository)
		expectedData []domain.Menu
		wantErr      bool
		expectedErr  error
	}{
		{
			name:   "success - get list menu with filter",
			params: mockFilter,
			mockSetup: func(mockRepo *domain.MockMenuRepository) {
				// Pastikan mock menerima parameter mockFilter yang sesuai
				mockRepo.On("ReadListMenuQuery", ctx, mockFilter).Return(mockData, nil)
			},
			expectedData: mockData,
			wantErr:      false,
			expectedErr:  nil,
		},
		{
			name:   "success - empty result (no data matches filter)",
			params: mockFilter,
			mockSetup: func(mockRepo *domain.MockMenuRepository) {
				// Simulasi query sukses tapi tidak ada data yang cocok
				mockRepo.On("ReadListMenuQuery", ctx, mockFilter).Return([]domain.Menu(nil), nil)
			},
			expectedData: nil,
			wantErr:      false,
			expectedErr:  nil,
		},
		{
			name:   "error - repository failure",
			params: mockFilter,
			mockSetup: func(mockRepo *domain.MockMenuRepository) {
				mockRepo.On("ReadListMenuQuery", ctx, mockFilter).Return([]domain.Menu(nil), errDb)
			},
			expectedData: nil,
			wantErr:      true,
			// Karena Usecase langsung me-return raw error, kita validasi dengan error asli
			expectedErr: errDb,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(domain.MockMenuRepository)
			tc.mockSetup(mockRepo)

			// Inisialisasi Usecase dengan mock repo dan logger kosong
			uc := usecase.NewMenuUsecase(mockRepo, &nopLogger)

			// Eksekusi fungsi menggunakan parameter dari test case
			result, err := uc.ReadListMenuUsecase(ctx, tc.params)

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedData, result)
			}

			// Validasi bahwa mock dieksekusi sesuai ekspektasi
			mockRepo.AssertExpectations(t)
		})
	}
}
