package usecase_test

import (
	"context"
	"testing"

	"be-dashboard-nba/internal/modules/menu/domain"
	"be-dashboard-nba/internal/modules/menu/usecase"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestReadSidebarMenuUsecase(t *testing.T) {
	ctx := context.Background()
	userID := "user-uuid-123"
	nopLogger := zerolog.Nop()

	// Helper untuk pointer string
	ptr := func(s string) *string { return &s }

	// Mock data sidebar yang biasanya bersifat rekursif/nested (tergantung implementasi repo Anda)
	mockSidebar := []domain.Menu{
		{
			ID:    1,
			Name:  "Dashboard",
			URL:   ptr("/dashboard"),
			Icon:  ptr("home-icon"),
			Group: "main",
		},
		{
			ID:    2,
			Name:  "User Management",
			URL:   ptr("/users"),
			Icon:  ptr("user-icon"),
			Group: "settings",
		},
	}

	errDb := errors.New("query execution error")

	tests := []struct {
		name         string
		userID       string
		mockSetup    func(mockRepo *domain.MockMenuRepository)
		expectedData []domain.Menu
		wantErr      bool
		expectedErr  error
	}{
		{
			name:   "success - get sidebar for authorized user",
			userID: userID,
			mockSetup: func(mockRepo *domain.MockMenuRepository) {
				mockRepo.On("ReadSidebarMenuQuery", ctx, userID).Return(mockSidebar, nil)
			},
			expectedData: mockSidebar,
			wantErr:      false,
			expectedErr:  nil,
		},
		{
			name:   "success - empty sidebar for user with no access",
			userID: "user-newbie",
			mockSetup: func(mockRepo *domain.MockMenuRepository) {
				// User ada tapi tidak punya akses menu apa pun
				mockRepo.On("ReadSidebarMenuQuery", ctx, "user-newbie").Return([]domain.Menu(nil), nil)
			},
			expectedData: nil,
			wantErr:      false,
			expectedErr:  nil,
		},
		{
			name:   "error - repository failure",
			userID: userID,
			mockSetup: func(mockRepo *domain.MockMenuRepository) {
				mockRepo.On("ReadSidebarMenuQuery", ctx, userID).Return([]domain.Menu(nil), errDb)
			},
			expectedData: nil,
			wantErr:      true,
			expectedErr:  errDb,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(domain.MockMenuRepository)
			tc.mockSetup(mockRepo)

			uc := usecase.NewMenuUsecase(mockRepo, &nopLogger)

			result, err := uc.ReadSidebarMenuUsecase(ctx, tc.userID)

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedData, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
