package usecase_test

import (
	"context"
	"testing"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/core/utils"
	menuDomain "be-dashboard-nba/internal/modules/menu/domain"
	"be-dashboard-nba/internal/modules/menu_permission/domain"
	"be-dashboard-nba/internal/modules/menu_permission/usecase"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestReadListMenuPermissionUsecase(t *testing.T) {
	ctx := context.Background()
	nopLogger := zerolog.Nop()

	errDb := errors.New("database error")

	// Setup Filter Input
	filter := domain.MenuPermissionFilter{
		MenuID: 10,
	}
	// Asumsi struct Anda mewarisi PaginationPayload
	filter.Page = 1
	filter.Limit = 10

	// Mock Data Balikan dari DB
	mockDBData := []domain.MenuPermissionDetail{
		{ID: 1, MenuID: 10, Code: "CREATE", ActionName: "Create Data"},
		{ID: 2, MenuID: 10, Code: "UPDATE", ActionName: "Update Data"},
	}

	// Ekspektasi Data setelah di-mapping oleh Usecase
	expectedDetailData := []domain.MenuPermissionDetail{
		{ID: 1, MenuID: 10, Code: "CREATE", ActionName: "Create Data"},
		{ID: 2, MenuID: 10, Code: "UPDATE", ActionName: "Update Data"},
	}

	// Ekspektasi Paginasi (karena fungsi ini deterministik, kita bisa langsung generate)
	expectedPagination := utils.GeneratePagination(2, filter.Page, filter.Limit)
	emptyPagination := utils.GeneratePagination(0, filter.Page, filter.Limit)

	tests := []struct {
		name         string
		filter       domain.MenuPermissionFilter
		mockSetup    func(mockMenuRepo *menuDomain.MockMenuRepository, mockPermRepo *domain.MockMenuPermissionRepository)
		expectedData domain.MenuPermissionPaginationResponse
		wantErr      bool
		expectedErr  error
	}{
		{
			name:   "success - read list with data",
			filter: filter,
			mockSetup: func(mockMenuRepo *menuDomain.MockMenuRepository, mockPermRepo *domain.MockMenuPermissionRepository) {
				// 1. Validasi Menu
				mockMenuRepo.On("ReadMenuByIDQuery", ctx, filter.MenuID).Return(menuDomain.Menu{ID: filter.MenuID}, nil).Once()
				// 2. Count Data
				mockPermRepo.On("ReadMenuPermissionCountQuery", ctx, filter).Return(2, nil).Once()
				// 3. Get Data List
				mockPermRepo.On("ReadMenuPermissionListQuery", ctx, filter).Return(mockDBData, nil).Once()
			},
			expectedData: domain.MenuPermissionPaginationResponse{
				Data:       expectedDetailData,
				Pagination: expectedPagination,
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:   "success - read list with empty data (nil safety check)",
			filter: filter,
			mockSetup: func(mockMenuRepo *menuDomain.MockMenuRepository, mockPermRepo *domain.MockMenuPermissionRepository) {
				mockMenuRepo.On("ReadMenuByIDQuery", ctx, filter.MenuID).Return(menuDomain.Menu{ID: filter.MenuID}, nil).Once()
				mockPermRepo.On("ReadMenuPermissionCountQuery", ctx, filter).Return(0, nil).Once()

				mockPermRepo.On("ReadMenuPermissionListQuery", ctx, filter).Return([]domain.MenuPermissionDetail(nil), nil).Once()
			},
			expectedData: domain.MenuPermissionPaginationResponse{
				Data:       []domain.MenuPermissionDetail{}, // Ekspektasi array kosong, bukan nil
				Pagination: emptyPagination,
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:   "failed - menu id not found",
			filter: filter,
			mockSetup: func(mockMenuRepo *menuDomain.MockMenuRepository, mockPermRepo *domain.MockMenuPermissionRepository) {
				mockMenuRepo.On("ReadMenuByIDQuery", ctx, filter.MenuID).Return(menuDomain.Menu{}, constant.ErrDataNotFound).Once()
				// Berhenti di sini, fungsi count & list tidak dipanggil
			},
			expectedData: domain.MenuPermissionPaginationResponse{},
			wantErr:      true,
			expectedErr:  constant.ErrMenuIdNotFound,
		},
		{
			name:   "failed - error validating menu",
			filter: filter,
			mockSetup: func(mockMenuRepo *menuDomain.MockMenuRepository, mockPermRepo *domain.MockMenuPermissionRepository) {
				mockMenuRepo.On("ReadMenuByIDQuery", ctx, filter.MenuID).Return(menuDomain.Menu{}, errDb).Once()
			},
			expectedData: domain.MenuPermissionPaginationResponse{},
			wantErr:      true,
			expectedErr:  constant.ErrUnknownSource,
		},
		{
			name:   "failed - error counting data",
			filter: filter,
			mockSetup: func(mockMenuRepo *menuDomain.MockMenuRepository, mockPermRepo *domain.MockMenuPermissionRepository) {
				mockMenuRepo.On("ReadMenuByIDQuery", ctx, filter.MenuID).Return(menuDomain.Menu{ID: filter.MenuID}, nil).Once()
				mockPermRepo.On("ReadMenuPermissionCountQuery", ctx, filter).Return(0, errDb).Once()
			},
			expectedData: domain.MenuPermissionPaginationResponse{},
			wantErr:      true,
			expectedErr:  constant.ErrUnknownSource,
		},
		{
			name:   "failed - error fetching list data",
			filter: filter,
			mockSetup: func(mockMenuRepo *menuDomain.MockMenuRepository, mockPermRepo *domain.MockMenuPermissionRepository) {
				mockMenuRepo.On("ReadMenuByIDQuery", ctx, filter.MenuID).Return(menuDomain.Menu{ID: filter.MenuID}, nil).Once()
				mockPermRepo.On("ReadMenuPermissionCountQuery", ctx, filter).Return(2, nil).Once()

				mockPermRepo.On("ReadMenuPermissionListQuery", ctx, filter).Return([]domain.MenuPermissionDetail(nil), errDb).Once()
			},
			expectedData: domain.MenuPermissionPaginationResponse{},
			wantErr:      true,
			expectedErr:  constant.ErrUnknownSource,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockMenuRepo := new(menuDomain.MockMenuRepository)
			mockPermRepo := new(domain.MockMenuPermissionRepository)

			tc.mockSetup(mockMenuRepo, mockPermRepo)

			uc := usecase.NewMenuPermissionUsecase(mockPermRepo, mockMenuRepo, &nopLogger)

			result, err := uc.ReadListMenuPermissionUsecase(ctx, tc.filter)

			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedData, result)
			}

			mockMenuRepo.AssertExpectations(t)
			mockPermRepo.AssertExpectations(t)
		})
	}
}
