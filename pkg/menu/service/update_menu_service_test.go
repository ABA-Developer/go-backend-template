package service

import (
	menuPresenter "be-dashboard-nba/api/presenter/menu"
	"be-dashboard-nba/constant"
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMenuService(t *testing.T) {
	// Helper functions untuk pointer
	getPtrBool := func(b bool) *bool { return &b }
	getPtrInt := func(i int) *int { return &i }
	getPtrString := func(s string) *string { return &s }

	// Test data constants
	const (
		userID       = "user123"
		menuID       = 1
		parentMenuID = 2
		newParentID  = 5
	)

	// Menu groups
	const (
		groupMain      = "main"
		groupAdmin     = "admin"
		groupNew       = "new_group"
		groupOld       = "old_group"
		parentGroup    = "parent_group"
		groupSystem    = "system"
		groupRequested = "requested_group"
	)

	// Menu names
	const (
		menuUpdated     = "Updated Menu"
		menuRoot        = "Root Menu"
		menuParent      = "Parent Menu"
		menuChild       = "Child Menu"
		menuWithDetails = "Menu with Details"
		menuSame        = "Same Menu"
		menuNonExistent = "Non-existent"
		menuOld         = "Old Menu"
		menuOldName     = "Old Name"
		menuNewParent   = "New Parent"
	)

	// Common descriptions
	const (
		descriptionHere = "Description here"
		urlPath         = "/url"
		iconName        = "icon-name"
	)

	// Database field names
	menuFields := []string{
		"id", "parent_id", "name", "description", "url", "sort", "group", "icon",
		"active", "display",
	}

	// SQL query constants
	const (
		selectMenuQuery          = `SELECT`
		countChildrenQuery       = `SELECT COUNT\(\*\) FROM app_menu WHERE parent_id = \$1`
		selectNextSortRootQuery  = `SELECT COALESCE\(MAX\(sort\), -1\) \+ 1 FROM app_menu WHERE parent_id IS NULL AND "group" = \$1`
		selectNextSortChildQuery = `SELECT COALESCE\(MAX\(sort\), -1\) \+ 1 FROM app_menu WHERE parent_id = \$1 AND "group" = \$2`
		updateMenuQuery          = `UPDATE app_menu`
		updateChildrenGroupQuery = `UPDATE app_menu SET "group" = \$1, updated_at = NOW\(\) WHERE parent_id = \$2`
	)

	// Helper functions untuk test data
	createMenuRow := func(id int64, parentID interface{}, name string, sort int32, groupName string, active, display bool) *sqlmock.Rows {
		rows := sqlmock.NewRows(menuFields)
		var parentIDValue interface{}

		switch v := parentID.(type) {
		case int:
			parentIDValue = int32(v)
		case int32:
			parentIDValue = v
		case nil:
			parentIDValue = nil
		default:
			parentIDValue = parentID
		}

		return rows.AddRow(
			id,
			parentIDValue,
			name,
			nil,
			nil,
			sort,
			groupName,
			nil,
			active,
			display,
		)
	}

	createNullString := func(s string) sql.NullString {
		if s == "" {
			return sql.NullString{Valid: false}
		}
		return sql.NullString{String: s, Valid: true}
	}

	createNullInt32 := func(i int) sql.NullInt32 {
		if i == 0 {
			return sql.NullInt32{Valid: false}
		}
		return sql.NullInt32{Int32: int32(i), Valid: true}
	}

	// Test cases
	tests := []struct {
		name        string
		request     menuPresenter.UpdateMenuRequest
		userID      string
		menuID      int
		wantErr     bool
		expectedErr error
		setupMock   func(mock sqlmock.Sqlmock, request menuPresenter.UpdateMenuRequest, userID string, menuID int)
	}{
		// ==================== SUCCESS CASES ====================
		{
			name: "success_update_basic_fields",
			request: menuPresenter.UpdateMenuRequest{
				Name:    menuUpdated,
				Group:   groupMain,
				Active:  getPtrBool(true),
				Display: getPtrBool(true),
			},
			userID:      userID,
			menuID:      menuID,
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.UpdateMenuRequest, userID string, menuID int) {
				mock.ExpectBegin()

				// Mock existing menu (root without children)
				mock.ExpectQuery(selectMenuQuery).
					WithArgs(menuID).
					WillReturnRows(createMenuRow(1, nil, menuOldName, 0, groupMain, true, true))

				// Mock count children - no children
				mock.ExpectQuery(countChildrenQuery).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				// Mock update menu
				mock.ExpectExec(updateMenuQuery).
					WithArgs(
						menuID,
						menuUpdated,
						sql.NullInt32{Valid: false},
						sql.NullString{Valid: false},
						sql.NullString{Valid: false},
						0,
						groupMain,
						sql.NullString{Valid: false},
						true,
						true,
						createNullString(userID),
					).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "success_child_menu_become_root",
			request: menuPresenter.UpdateMenuRequest{
				Name:     menuChild,
				Group:    groupMain,
				ParentID: nil, // Explicitly null - become root
				Active:   getPtrBool(true),
				Display:  getPtrBool(true),
			},
			userID:      userID,
			menuID:      menuID,
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.UpdateMenuRequest, userID string, menuID int) {
				mock.ExpectBegin()

				// Mock existing menu (child with parent)
				mock.ExpectQuery(selectMenuQuery).
					WithArgs(menuID).
					WillReturnRows(createMenuRow(1, parentMenuID, menuChild, 1, groupMain, true, true))

				// Mock count children - no children
				mock.ExpectQuery(countChildrenQuery).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				// Mock get next sort for root group (karena menjadi root)
				mock.ExpectQuery(selectNextSortRootQuery).
					WithArgs(groupMain).
					WillReturnRows(sqlmock.NewRows([]string{"sort"}).AddRow(2))

				// Mock update menu - parent_id menjadi NULL
				mock.ExpectExec(updateMenuQuery).
					WithArgs(
						menuID,
						menuChild,
						sql.NullInt32{Valid: false}, // ParentID menjadi NULL
						sql.NullString{Valid: false},
						sql.NullString{Valid: false},
						2, // New sort dari query di atas
						groupMain,
						sql.NullString{Valid: false},
						true,
						true,
						createNullString(userID),
					).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "success_root_menu_change_group_without_children",
			request: menuPresenter.UpdateMenuRequest{
				Name:    menuRoot,
				Group:   groupNew,
				Active:  getPtrBool(true),
				Display: getPtrBool(true),
			},
			userID:      userID,
			menuID:      menuID,
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.UpdateMenuRequest, userID string, menuID int) {
				mock.ExpectBegin()

				// Mock existing menu (root without children)
				mock.ExpectQuery(selectMenuQuery).
					WithArgs(menuID).
					WillReturnRows(createMenuRow(1, nil, menuRoot, 0, groupOld, true, true))

				// Mock count children - no children
				mock.ExpectQuery(countChildrenQuery).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				// Mock get next sort for new group
				mock.ExpectQuery(selectNextSortRootQuery).
					WithArgs(groupNew).
					WillReturnRows(sqlmock.NewRows([]string{"sort"}).AddRow(2))

				// Mock update menu
				mock.ExpectExec(updateMenuQuery).
					WithArgs(
						menuID,
						menuRoot,
						sql.NullInt32{Valid: false},
						sql.NullString{Valid: false},
						sql.NullString{Valid: false},
						2,
						groupNew,
						sql.NullString{Valid: false},
						true,
						true,
						createNullString(userID),
					).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "success_root_menu_change_group_with_children",
			request: menuPresenter.UpdateMenuRequest{
				Name:    menuParent,
				Group:   groupNew,
				Active:  getPtrBool(true),
				Display: getPtrBool(true),
			},
			userID:      userID,
			menuID:      menuID,
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.UpdateMenuRequest, userID string, menuID int) {
				mock.ExpectBegin()

				// Mock existing menu (root with children)
				mock.ExpectQuery(selectMenuQuery).
					WithArgs(menuID).
					WillReturnRows(createMenuRow(1, nil, menuParent, 0, groupOld, true, true))

				// Mock count children - has 2 children
				mock.ExpectQuery(countChildrenQuery).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

				// Mock get next sort for new group
				mock.ExpectQuery(selectNextSortRootQuery).
					WithArgs(groupNew).
					WillReturnRows(sqlmock.NewRows([]string{"sort"}).AddRow(3))

				// Mock update menu parent
				mock.ExpectExec(updateMenuQuery).
					WithArgs(
						menuID,
						menuParent,
						sql.NullInt32{Valid: false},
						sql.NullString{Valid: false},
						sql.NullString{Valid: false},
						3,
						groupNew,
						sql.NullString{Valid: false},
						true,
						true,
						createNullString(userID),
					).
					WillReturnResult(sqlmock.NewResult(0, 1))

				// Mock update children group
				mock.ExpectExec(updateChildrenGroupQuery).
					WithArgs(groupNew, menuID).
					WillReturnResult(sqlmock.NewResult(0, 2))

				mock.ExpectCommit()
			},
		},
		{
			name: "success_child_menu_change_to_new_parent",
			request: menuPresenter.UpdateMenuRequest{
				Name:     menuChild,
				Group:    groupMain,
				ParentID: getPtrInt(newParentID),
				Active:   getPtrBool(true),
				Display:  getPtrBool(true),
			},
			userID:      userID,
			menuID:      menuID,
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.UpdateMenuRequest, userID string, menuID int) {
				mock.ExpectBegin()

				// Mock existing menu (child without children)
				mock.ExpectQuery(selectMenuQuery).
					WithArgs(menuID).
					WillReturnRows(createMenuRow(1, parentMenuID, menuChild, 1, groupMain, true, true))

				// Mock count children - no children
				mock.ExpectQuery(countChildrenQuery).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				// Mock read new parent menu
				mock.ExpectQuery(selectMenuQuery).
					WithArgs(newParentID).
					WillReturnRows(createMenuRow(5, nil, menuNewParent, 0, groupAdmin, true, true))

				// Mock get next sort for new parent and group
				mock.ExpectQuery(selectNextSortChildQuery).
					WithArgs(int32(newParentID), groupAdmin).
					WillReturnRows(sqlmock.NewRows([]string{"sort"}).AddRow(2))

				// Mock update menu
				mock.ExpectExec(updateMenuQuery).
					WithArgs(
						menuID,
						menuChild,
						createNullInt32(newParentID),
						sql.NullString{Valid: false},
						sql.NullString{Valid: false},
						2,
						groupAdmin,
						sql.NullString{Valid: false},
						true,
						true,
						createNullString(userID),
					).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "success_update_with_description_url_icon",
			request: menuPresenter.UpdateMenuRequest{
				Name:        menuWithDetails,
				Group:       groupMain,
				Description: getPtrString(descriptionHere),
				URL:         getPtrString(urlPath),
				Icon:        getPtrString(iconName),
				Active:      getPtrBool(true),
				Display:     getPtrBool(true),
			},
			userID:      userID,
			menuID:      menuID,
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.UpdateMenuRequest, userID string, menuID int) {
				mock.ExpectBegin()

				// Mock existing menu
				mock.ExpectQuery(selectMenuQuery).
					WithArgs(menuID).
					WillReturnRows(createMenuRow(1, nil, menuOld, 0, groupMain, true, true))

				// Mock count children - no children
				mock.ExpectQuery(countChildrenQuery).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				// Mock update menu with all fields
				mock.ExpectExec(updateMenuQuery).
					WithArgs(
						menuID,
						menuWithDetails,
						sql.NullInt32{Valid: false},
						createNullString(descriptionHere),
						createNullString(urlPath),
						0,
						groupMain,
						createNullString(iconName),
						true,
						true,
						createNullString(userID),
					).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "success_no_changes_keep_existing",
			request: menuPresenter.UpdateMenuRequest{
				Name:    menuSame,
				Group:   groupMain,
				Active:  getPtrBool(true),
				Display: getPtrBool(true),
			},
			userID:      userID,
			menuID:      menuID,
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.UpdateMenuRequest, userID string, menuID int) {
				mock.ExpectBegin()

				// Mock existing menu
				mock.ExpectQuery(selectMenuQuery).
					WithArgs(menuID).
					WillReturnRows(createMenuRow(1, nil, menuSame, 2, groupMain, true, true))

				// Mock count children - no children
				mock.ExpectQuery(countChildrenQuery).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				// Mock update menu (keep existing sort)
				mock.ExpectExec(updateMenuQuery).
					WithArgs(
						menuID,
						menuSame,
						sql.NullInt32{Valid: false},
						sql.NullString{Valid: false},
						sql.NullString{Valid: false},
						2,
						groupMain,
						sql.NullString{Valid: false},
						true,
						true,
						createNullString(userID),
					).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},

		// ==================== ERROR CASES ====================
		{
			name: "error_menu_not_found",
			request: menuPresenter.UpdateMenuRequest{
				Name:    menuNonExistent,
				Group:   groupMain,
				Active:  getPtrBool(true),
				Display: getPtrBool(true),
			},
			userID:      userID,
			menuID:      999,
			wantErr:     true,
			expectedErr: constant.ErrMenuIdNotFound,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.UpdateMenuRequest, userID string, menuID int) {
				mock.ExpectBegin()
				mock.ExpectQuery(selectMenuQuery).
					WithArgs(menuID).
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
		},
		{
			name: "error_menu_with_children_cannot_change_parent",
			request: menuPresenter.UpdateMenuRequest{
				Name:     menuParent,
				Group:    groupMain,
				ParentID: getPtrInt(newParentID), // Try to change parent
				Active:   getPtrBool(true),
				Display:  getPtrBool(true),
			},
			userID:      userID,
			menuID:      menuID,
			wantErr:     true,
			expectedErr: constant.ErrMenuHasChildren,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.UpdateMenuRequest, userID string, menuID int) {
				mock.ExpectBegin()

				// Mock existing menu (with children)
				mock.ExpectQuery(selectMenuQuery).
					WithArgs(menuID).
					WillReturnRows(createMenuRow(1, nil, menuParent, 0, groupMain, true, true))

				// Mock count children - has children
				mock.ExpectQuery(countChildrenQuery).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

				mock.ExpectRollback()
			},
		},
		{
			name: "error_get_sort_for_new_parent_fails",
			request: menuPresenter.UpdateMenuRequest{
				Name:     menuChild,
				Group:    groupMain,
				ParentID: getPtrInt(newParentID),
				Active:   getPtrBool(true),
				Display:  getPtrBool(true),
			},
			userID:      userID,
			menuID:      menuID,
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.UpdateMenuRequest, userID string, menuID int) {
				mock.ExpectBegin()

				// Mock existing menu
				mock.ExpectQuery(selectMenuQuery).
					WithArgs(menuID).
					WillReturnRows(createMenuRow(1, parentMenuID, menuChild, 1, groupMain, true, true))

				// Mock count children - no children
				mock.ExpectQuery(countChildrenQuery).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				// Mock read new parent menu
				mock.ExpectQuery(selectMenuQuery).
					WithArgs(newParentID).
					WillReturnRows(createMenuRow(5, nil, menuNewParent, 0, groupAdmin, true, true))

				// Mock get next sort for new parent - fails
				mock.ExpectQuery(selectNextSortChildQuery).
					WithArgs(int32(newParentID), groupAdmin).
					WillReturnError(errors.New("database error"))

				mock.ExpectRollback()
			},
		},
		{
			name: "error_get_sort_for_root_fails",
			request: menuPresenter.UpdateMenuRequest{
				Name:     menuChild,
				Group:    groupMain,
				ParentID: getPtrInt(0), // Become root
				Active:   getPtrBool(true),
				Display:  getPtrBool(true),
			},
			userID:      userID,
			menuID:      menuID,
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.UpdateMenuRequest, userID string, menuID int) {
				mock.ExpectBegin()

				// Mock existing menu
				mock.ExpectQuery(selectMenuQuery).
					WithArgs(menuID).
					WillReturnRows(createMenuRow(1, parentMenuID, menuChild, 1, groupMain, true, true))

				// Mock count children - no children
				mock.ExpectQuery(countChildrenQuery).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				// Mock get next sort for root - fails
				mock.ExpectQuery(selectNextSortRootQuery).
					WithArgs(groupMain).
					WillReturnError(errors.New("database error"))

				mock.ExpectRollback()
			},
		},
		{
			name: "error_update_children_group_fails",
			request: menuPresenter.UpdateMenuRequest{
				Name:    menuParent,
				Group:   groupNew,
				Active:  getPtrBool(true),
				Display: getPtrBool(true),
			},
			userID:      userID,
			menuID:      menuID,
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.UpdateMenuRequest, userID string, menuID int) {
				mock.ExpectBegin()

				// Mock existing menu (root with children)
				mock.ExpectQuery(selectMenuQuery).
					WithArgs(menuID).
					WillReturnRows(createMenuRow(1, nil, menuParent, 0, groupOld, true, true))

				// Mock count children - has children
				mock.ExpectQuery(countChildrenQuery).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

				// Mock get next sort for new group
				mock.ExpectQuery(selectNextSortRootQuery).
					WithArgs(groupNew).
					WillReturnRows(sqlmock.NewRows([]string{"sort"}).AddRow(3))

				// Mock update menu parent
				mock.ExpectExec(updateMenuQuery).
					WithArgs(
						menuID,
						menuParent,
						sql.NullInt32{Valid: false},
						sql.NullString{Valid: false},
						sql.NullString{Valid: false},
						3,
						groupNew,
						sql.NullString{Valid: false},
						true,
						true,
						createNullString(userID),
					).
					WillReturnResult(sqlmock.NewResult(0, 1))

				// Mock update children group - fails
				mock.ExpectExec(updateChildrenGroupQuery).
					WithArgs(groupNew, menuID).
					WillReturnError(errors.New("update children error"))

				mock.ExpectRollback()
			},
		},
		{
			name: "error_update_menu_query_fails",
			request: menuPresenter.UpdateMenuRequest{
				Name:    menuChild,
				Group:   groupMain,
				Active:  getPtrBool(true),
				Display: getPtrBool(true),
			},
			userID:      userID,
			menuID:      menuID,
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.UpdateMenuRequest, userID string, menuID int) {
				mock.ExpectBegin()

				// Mock existing menu
				mock.ExpectQuery(selectMenuQuery).
					WithArgs(menuID).
					WillReturnRows(createMenuRow(1, nil, menuChild, 0, groupMain, true, true))

				// Mock count children - no children
				mock.ExpectQuery(countChildrenQuery).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				// Mock update menu - fails
				mock.ExpectExec(updateMenuQuery).
					WillReturnError(errors.New("update error"))

				mock.ExpectRollback()
			},
		},
		{
			name: "error_commit_fails",
			request: menuPresenter.UpdateMenuRequest{
				Name:    menuChild,
				Group:   groupMain,
				Active:  getPtrBool(true),
				Display: getPtrBool(true),
			},
			userID:      userID,
			menuID:      menuID,
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.UpdateMenuRequest, userID string, menuID int) {
				mock.ExpectBegin()

				// Mock existing menu
				mock.ExpectQuery(selectMenuQuery).
					WithArgs(menuID).
					WillReturnRows(createMenuRow(1, nil, menuChild, 0, groupMain, true, true))

				// Mock count children - no children
				mock.ExpectQuery(countChildrenQuery).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				// Mock update menu
				mock.ExpectExec(updateMenuQuery).
					WillReturnResult(sqlmock.NewResult(0, 1))

				// Mock commit - fails
				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
		},
		{
			name: "error_begin_transaction_fails",
			request: menuPresenter.UpdateMenuRequest{
				Name:    menuChild,
				Group:   groupMain,
				Active:  getPtrBool(true),
				Display: getPtrBool(true),
			},
			userID:      userID,
			menuID:      menuID,
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.UpdateMenuRequest, userID string, menuID int) {
				mock.ExpectBegin().WillReturnError(errors.New("begin error"))
			},
		},
		{
			name: "error_read_parent_menu_fails",
			request: menuPresenter.UpdateMenuRequest{
				Name:     menuChild,
				Group:    groupMain,
				ParentID: getPtrInt(newParentID),
				Active:   getPtrBool(true),
				Display:  getPtrBool(true),
			},
			userID:      userID,
			menuID:      menuID,
			wantErr:     true,
			expectedErr: constant.ErrUnknownSource,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.UpdateMenuRequest, userID string, menuID int) {
				mock.ExpectBegin()

				// Mock existing menu
				mock.ExpectQuery(selectMenuQuery).
					WithArgs(menuID).
					WillReturnRows(createMenuRow(1, parentMenuID, menuChild, 1, groupMain, true, true))

				// Mock count children - no children
				mock.ExpectQuery(countChildrenQuery).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				// Mock read new parent menu - fails
				mock.ExpectQuery(selectMenuQuery).
					WithArgs(newParentID).
					WillReturnError(errors.New("parent not found"))

				mock.ExpectRollback()
			},
		},

		// ==================== EDGE CASES ====================
		{
			name: "edge_root_menu_already_root",
			request: menuPresenter.UpdateMenuRequest{
				Name:     menuRoot,
				Group:    groupMain,
				ParentID: nil, // Already root
				Active:   getPtrBool(true),
				Display:  getPtrBool(true),
			},
			userID:      userID,
			menuID:      menuID,
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.UpdateMenuRequest, userID string, menuID int) {
				mock.ExpectBegin()

				// Mock existing menu (already root)
				mock.ExpectQuery(selectMenuQuery).
					WithArgs(menuID).
					WillReturnRows(createMenuRow(1, nil, menuRoot, 0, groupMain, true, true))

				// Mock count children - no children
				mock.ExpectQuery(countChildrenQuery).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				// Mock update menu (no parent change)
				mock.ExpectExec(updateMenuQuery).
					WithArgs(
						menuID,
						menuRoot,
						sql.NullInt32{Valid: false}, // Still root
						sql.NullString{Valid: false},
						sql.NullString{Valid: false},
						0, // Keep existing sort
						groupMain,
						sql.NullString{Valid: false},
						true,
						true,
						createNullString(userID),
					).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "edge_same_parent_no_change",
			request: menuPresenter.UpdateMenuRequest{
				Name:     menuChild,
				Group:    groupMain,
				ParentID: getPtrInt(parentMenuID), // Same as existing
				Active:   getPtrBool(true),
				Display:  getPtrBool(true),
			},
			userID:      userID,
			menuID:      menuID,
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.UpdateMenuRequest, userID string, menuID int) {
				mock.ExpectBegin()

				// Mock existing menu
				mock.ExpectQuery(selectMenuQuery).
					WithArgs(menuID).
					WillReturnRows(createMenuRow(1, parentMenuID, menuChild, 1, groupMain, true, true))

				// Mock count children - no children
				mock.ExpectQuery(countChildrenQuery).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				// Mock read parent menu
				mock.ExpectQuery(selectMenuQuery).
					WithArgs(parentMenuID).
					WillReturnRows(createMenuRow(int64(parentMenuID), nil, "Parent Menu", 0, groupMain, true, true))

				// Mock update menu (no parent change, no group change)
				mock.ExpectExec(updateMenuQuery).
					WithArgs(
						menuID,
						menuChild,
						createNullInt32(parentMenuID),
						sql.NullString{Valid: false},
						sql.NullString{Valid: false},
						1,         // Keep existing sort
						groupMain, // Same group (from parent)
						sql.NullString{Valid: false},
						true,
						true,
						createNullString(userID),
					).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "edge_menu_with_children_update_other_fields",
			request: menuPresenter.UpdateMenuRequest{
				Name:    menuParent,
				Group:   groupMain,         // Same group
				Active:  getPtrBool(false), // Change active
				Display: getPtrBool(true),
			},
			userID:      userID,
			menuID:      menuID,
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.UpdateMenuRequest, userID string, menuID int) {
				mock.ExpectBegin()

				// Mock existing menu (with children)
				mock.ExpectQuery(selectMenuQuery).
					WithArgs(menuID).
					WillReturnRows(createMenuRow(1, nil, menuParent, 0, groupMain, true, true))

				// Mock count children - has children
				mock.ExpectQuery(countChildrenQuery).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

				// Mock update menu (only active changed, no children update needed)
				mock.ExpectExec(updateMenuQuery).
					WithArgs(
						menuID,
						menuParent,
						sql.NullInt32{Valid: false},
						sql.NullString{Valid: false},
						sql.NullString{Valid: false},
						0,         // Keep existing sort
						groupMain, // Same group
						sql.NullString{Valid: false},
						false, // Active changed to false
						true,
						createNullString(userID),
					).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},

		// ==================== PROTECTION CASES ====================
		{
			name: "protection_child_menu_follows_new_parent_group",
			request: menuPresenter.UpdateMenuRequest{
				Name:     menuChild,
				Group:    groupRequested, // Will be overridden
				ParentID: getPtrInt(newParentID),
				Active:   getPtrBool(true),
				Display:  getPtrBool(true),
			},
			userID:      userID,
			menuID:      menuID,
			wantErr:     false,
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, request menuPresenter.UpdateMenuRequest, userID string, menuID int) {
				mock.ExpectBegin()

				// Mock existing menu
				mock.ExpectQuery(selectMenuQuery).
					WithArgs(menuID).
					WillReturnRows(createMenuRow(1, parentMenuID, menuChild, 1, groupOld, true, true))

				// Mock count children - no children
				mock.ExpectQuery(countChildrenQuery).
					WithArgs(menuID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				// Mock read new parent menu
				mock.ExpectQuery(selectMenuQuery).
					WithArgs(newParentID).
					WillReturnRows(createMenuRow(5, nil, menuNewParent, 0, parentGroup, true, true))

				// Mock get next sort for new parent and parent's group
				mock.ExpectQuery(selectNextSortChildQuery).
					WithArgs(int32(newParentID), parentGroup).
					WillReturnRows(sqlmock.NewRows([]string{"sort"}).AddRow(2))

				// Mock update menu (group should be "parent_group", not "requested_group")
				mock.ExpectExec(updateMenuQuery).
					WithArgs(
						menuID,
						menuChild,
						createNullInt32(newParentID),
						sql.NullString{Valid: false},
						sql.NullString{Valid: false},
						2,
						parentGroup, // Follows parent group, not requested group
						sql.NullString{Valid: false},
						true,
						true,
						createNullString(userID),
					).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, mock, cleanup := newTestService(t)
			defer cleanup()

			ctx := context.Background()
			tc.setupMock(mock, tc.request, tc.userID, tc.menuID)

			err := svc.UpdateMenuService(ctx, tc.request, tc.userID, tc.menuID)

			if tc.wantErr {
				assert.Error(t, err)
				if tc.expectedErr != nil {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
