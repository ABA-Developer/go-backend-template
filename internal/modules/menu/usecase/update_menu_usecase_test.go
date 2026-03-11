package usecase_test

import (
	"context"
	"testing"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/modules/menu/domain"
	"be-dashboard-nba/internal/modules/menu/usecase"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateMenuUsecase(t *testing.T) {
	ctx := context.Background()
	nopLogger := zerolog.Nop()

	ptrStr := func(s string) *string { return &s }
	ptrInt32 := func(i int32) *int32 { return &i }

	menuID := 1
	updatedBy := ptrStr("user-123")

	// Data Existing di Database
	existingMenu := domain.Menu{
		ID:          menuID,
		ParentID:    nil, // Awalnya root menu
		Name:        "Old Name",
		Description: ptrStr("Old Desc"),
		URL:         ptrStr("/old"),
		Icon:        ptrStr("old-icon"),
		Sort:        5,
		Group:       "admin",
	}

	tests := []struct {
		name        string
		payload     domain.MenuUpdatePayload
		mockSetup   func(mockRepo *domain.MockMenuRepository)
		wantErr     bool
		expectedErr error
	}{
		{
			name: "success - simple update without changing parent/group (fallback test)",
			payload: domain.MenuUpdatePayload{
				ID:        menuID,
				Name:      "New Name", // Hanya ubah nama
				Group:     "admin",
				UpdatedBy: updatedBy,
			},
			mockSetup: func(mockRepo *domain.MockMenuRepository) {
				mockRepo.On("ReadMenuByIDQuery", ctx, menuID).Return(existingMenu, nil)
				mockRepo.On("CountMenuChildren", ctx, menuID).Return(0, nil)

				// Ekspektasi payload akhir setelah melewati logika fallback Usecase
				expectedPayload := domain.MenuUpdatePayload{
					ID:          menuID,
					ParentID:    nil,
					Name:        "New Name",
					Description: existingMenu.Description, // Fallback ke data lama
					URL:         existingMenu.URL,         // Fallback ke data lama
					Sort:        5,                        // Disalin dari existingMenu.Sort
					Group:       "admin",
					Icon:        existingMenu.Icon, // Fallback ke data lama
					UpdatedBy:   updatedBy,
				}
				// updateChildrenGroup = false karena tidak ada anak dan tidak ubah grup
				mockRepo.On("UpdateMenuQuery", ctx, expectedPayload, false).Return(nil)
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "success - update group and cascade to children",
			payload: domain.MenuUpdatePayload{
				ID:        menuID,
				Name:      "Old Name",
				Group:     "superadmin", // Ubah grup
				UpdatedBy: updatedBy,
			},
			mockSetup: func(mockRepo *domain.MockMenuRepository) {
				mockRepo.On("ReadMenuByIDQuery", ctx, menuID).Return(existingMenu, nil)
				// Skenario punya anak
				mockRepo.On("CountMenuChildren", ctx, menuID).Return(3, nil)

				expectedPayload := domain.MenuUpdatePayload{
					ID:          menuID,
					Name:        "Old Name",
					Description: existingMenu.Description,
					URL:         existingMenu.URL,
					Icon:        existingMenu.Icon,
					Sort:        5,
					Group:       "superadmin",
					UpdatedBy:   updatedBy,
				}
				// ✨ updateChildrenGroup = true karena group berubah dan punya anak
				mockRepo.On("UpdateMenuQuery", ctx, expectedPayload, true).Return(nil)
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "success - move to new parent and inherit group",
			payload: domain.MenuUpdatePayload{
				ID:        menuID,
				ParentID:  ptrInt32(2), // Pindah parent
				Group:     "admin",     // Payload awal admin
				UpdatedBy: updatedBy,
			},
			mockSetup: func(mockRepo *domain.MockMenuRepository) {
				mockRepo.On("ReadMenuByIDQuery", ctx, menuID).Return(existingMenu, nil)
				mockRepo.On("CountMenuChildren", ctx, menuID).Return(0, nil) // Syarat pindah parent: tidak boleh punya anak

				// Mock saat Usecase membaca data parent baru
				parentMenu := domain.Menu{ID: 2, Group: "finance"}
				mockRepo.On("ReadMenuByIDQuery", ctx, 2).Return(parentMenu, nil)

				// Payload akhir harus MEWARISI grup dari parent
				expectedPayload := domain.MenuUpdatePayload{
					ID:          menuID,
					ParentID:    ptrInt32(2),
					Name:        existingMenu.Name, // Fallback
					Description: existingMenu.Description,
					URL:         existingMenu.URL,
					Icon:        existingMenu.Icon,
					Sort:        5,
					Group:       "finance", // ✨ Ter-override oleh parentMenu.Group
					UpdatedBy:   updatedBy,
				}
				mockRepo.On("UpdateMenuQuery", ctx, expectedPayload, false).Return(nil)
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "error - cannot change parent if menu has children",
			payload: domain.MenuUpdatePayload{
				ID:       menuID,
				ParentID: ptrInt32(99), // Mencoba pindah parent
			},
			mockSetup: func(mockRepo *domain.MockMenuRepository) {
				mockRepo.On("ReadMenuByIDQuery", ctx, menuID).Return(existingMenu, nil)
				// ✨ Jika punya anak (> 0) dan parent berubah, Usecase harus melempar error
				mockRepo.On("CountMenuChildren", ctx, menuID).Return(2, nil)
			},
			wantErr:     true,
			expectedErr: constant.ErrMenuHasChildren,
		},
		{
			name:    "error - menu not found",
			payload: domain.MenuUpdatePayload{ID: 999},
			mockSetup: func(mockRepo *domain.MockMenuRepository) {
				mockRepo.On("ReadMenuByIDQuery", ctx, 999).Return(domain.Menu{}, constant.ErrDataNotFound)
			},
			wantErr:     true,
			expectedErr: constant.ErrMenuIdNotFound,
		},
		{
			name:    "error - update query failed",
			payload: domain.MenuUpdatePayload{ID: menuID},
			mockSetup: func(mockRepo *domain.MockMenuRepository) {
				mockRepo.On("ReadMenuByIDQuery", ctx, menuID).Return(existingMenu, nil)
				mockRepo.On("CountMenuChildren", ctx, menuID).Return(0, nil)
				// Menggunakan mock.AnythingOfType agar fleksibel jika hanya ingin ngetest kegagalan query
				mockRepo.On("UpdateMenuQuery", ctx, mock.AnythingOfType("domain.MenuUpdatePayload"), false).
					Return(errors.New("db error"))
			},
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(domain.MockMenuRepository)
			tc.mockSetup(mockRepo)

			uc := usecase.NewMenuUsecase(mockRepo, &nopLogger)

			err := uc.UpdateMenuUsecase(ctx, tc.payload)

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
