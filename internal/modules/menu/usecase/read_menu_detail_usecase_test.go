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

func TestReadMenuDetailUsecase(t *testing.T) {
	ctx := context.Background()
	menuID := 10
	nopLogger := zerolog.Nop()

	// Siapkan pointer helper untuk mock data
	valParentID := int32(2)
	valDesc := "Menu untuk pengaturan pengguna"
	valURL := "/settings/users"
	valIcon := "icon-user"

	// Mock data kembalian dari layer Repository (entitas asli database)
	mockMenuFromRepo := domain.Menu{
		ID:          menuID,
		ParentID:    &valParentID,
		Name:        "User Settings",
		Description: &valDesc,
		URL:         &valURL,
		Sort:        1,
		Group:       "Admin",
		Icon:        &valIcon,
		Active:      true,
		Display:     true,
	}

	// Mock ekspektasi hasil mapping di layer Usecase (entitas View/Response)
	expectedMenuDetail := domain.MenuDetail{
		ID:          mockMenuFromRepo.ID,
		ParentID:    mockMenuFromRepo.ParentID,
		Name:        mockMenuFromRepo.Name,
		Description: mockMenuFromRepo.Description,
		URL:         mockMenuFromRepo.URL,
		Sort:        mockMenuFromRepo.Sort,
		Group:       mockMenuFromRepo.Group,
		Icon:        mockMenuFromRepo.Icon,
		Active:      mockMenuFromRepo.Active,
		Display:     mockMenuFromRepo.Display,
	}

	tests := []struct {
		name         string
		menuID       int
		mockSetup    func(mockRepo *domain.MockMenuRepository)
		wantErr      bool
		expectedErr  error
		expectedData domain.MenuDetail
	}{
		{
			name:   "success - read and map menu detail",
			menuID: menuID,
			mockSetup: func(mockRepo *domain.MockMenuRepository) {
				// Simulasi Repo menemukan data dan mengembalikan object domain.Menu
				mockRepo.On("ReadMenuByIDQuery", ctx, menuID).Return(mockMenuFromRepo, nil)
			},
			wantErr:      false,
			expectedErr:  nil,
			expectedData: expectedMenuDetail,
		},
		{
			name:   "error - menu not found",
			menuID: menuID,
			mockSetup: func(mockRepo *domain.MockMenuRepository) {
				// Simulasi Repo tidak menemukan data
				mockRepo.On("ReadMenuByIDQuery", ctx, menuID).Return(domain.Menu{}, constant.ErrDataNotFound)
			},
			wantErr:      true,
			expectedErr:  constant.ErrMenuIdNotFound,
			expectedData: domain.MenuDetail{}, // Jika error, kembalian struct harus kosong
		},
		{
			name:   "error - repository internal error",
			menuID: menuID,
			mockSetup: func(mockRepo *domain.MockMenuRepository) {
				// Simulasi error database yang tidak terduga
				mockRepo.On("ReadMenuByIDQuery", ctx, menuID).Return(domain.Menu{}, errors.New("db connection timeout"))
			},
			wantErr:      true,
			expectedErr:  constant.ErrUnknownSource,
			expectedData: domain.MenuDetail{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(domain.MockMenuRepository)
			tc.mockSetup(mockRepo)

			// Inisialisasi Usecase dengan mock repo dan logger kosong
			uc := usecase.NewMenuUsecase(mockRepo, &nopLogger)

			// Eksekusi fungsi
			resultData, err := uc.ReadMenuDetailUsecase(ctx, tc.menuID)

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			// ✨ Validasi krusial: pastikan hasil mapping struct sesuai dengan ekspektasi
			assert.Equal(t, tc.expectedData, resultData)

			// Validasi bahwa method di mock dipanggil sesuai skenario
			mockRepo.AssertExpectations(t)
		})
	}
}
