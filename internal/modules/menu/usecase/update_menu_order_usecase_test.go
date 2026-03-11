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

func TestUpdateMenuOrderUsecase(t *testing.T) {
	ctx := context.Background()
	nopLogger := zerolog.Nop()

	// Helper untuk membuat pointer tipe int dan int32
	ptrInt := func(i int) *int { return &i }
	ptrInt32 := func(i int32) *int32 { return &i }

	// Skenario 1: Parameter dengan ParentID
	payloadWithParent := domain.MenuUpdateSortPayload{
		Group:     "admin",
		ParentID:  ptrInt(5),
		SortedIDs: []int{10, 20, 30}, // 10 urutan 1, 20 urutan 2, 30 urutan 3
		UpdatedBy: "user-123",
	}
	expectedRepoPayloadsWithParent := []domain.MenuUpdateSortItemPayload{
		{ID: 10, Sort: 1, UpdatedBy: "user-123", ParentID: ptrInt32(5), Group: "admin"},
		{ID: 20, Sort: 2, UpdatedBy: "user-123", ParentID: ptrInt32(5), Group: "admin"},
		{ID: 30, Sort: 3, UpdatedBy: "user-123", ParentID: ptrInt32(5), Group: "admin"},
	}

	// Skenario 2: Parameter tanpa ParentID (Menu Utama/Root)
	payloadWithoutParent := domain.MenuUpdateSortPayload{
		Group:     "main",
		ParentID:  nil,
		SortedIDs: []int{15, 25}, // 15 urutan 1, 25 urutan 2
		UpdatedBy: "user-456",
	}
	expectedRepoPayloadsWithoutParent := []domain.MenuUpdateSortItemPayload{
		{ID: 15, Sort: 1, UpdatedBy: "user-456", ParentID: nil, Group: "main"},
		{ID: 25, Sort: 2, UpdatedBy: "user-456", ParentID: nil, Group: "main"},
	}

	tests := []struct {
		name        string
		payload     domain.MenuUpdateSortPayload
		mockSetup   func(mockRepo *domain.MockMenuRepository)
		wantErr     bool
		expectedErr error
	}{
		{
			name:    "success - update order with parent_id",
			payload: payloadWithParent,
			mockSetup: func(mockRepo *domain.MockMenuRepository) {
				// ✨ Validasi krusial: Pastikan mockRepo MENGHARAPKAN array yang sudah di-mapping dengan benar
				mockRepo.On("UpdateMenuOrderQuery", ctx, expectedRepoPayloadsWithParent).Return(nil)
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:    "success - update order without parent_id (nil pointer)",
			payload: payloadWithoutParent,
			mockSetup: func(mockRepo *domain.MockMenuRepository) {
				// Validasi untuk ParentID bernilai nil
				mockRepo.On("UpdateMenuOrderQuery", ctx, expectedRepoPayloadsWithoutParent).Return(nil)
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:    "error - repository transaction failed",
			payload: payloadWithParent,
			mockSetup: func(mockRepo *domain.MockMenuRepository) {
				mockRepo.On("UpdateMenuOrderQuery", ctx, expectedRepoPayloadsWithParent).
					Return(errors.New("deadlock detected"))
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

			err := uc.UpdateMenuOrderUsecase(ctx, tc.payload)

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			// Validasi bahwa method di mock dieksekusi dengan argumen yang tepat persis
			mockRepo.AssertExpectations(t)
		})
	}
}
