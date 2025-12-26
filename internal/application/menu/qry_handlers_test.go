package menu

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domainMenu "github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
)

// ============================================================
// GetMenuHandler Tests
// ============================================================

func TestGetMenuHandler_Handle_Success(t *testing.T) {
	mockQryRepo := new(MockMenuQueryRepository)

	expectedMenu := newTestMenu(1, "Dashboard")
	mockQryRepo.On("FindByID", mock.Anything, uint(1)).Return(expectedMenu, nil)

	handler := NewGetMenuHandler(mockQryRepo)

	result, err := handler.Handle(context.Background(), GetMenuQuery{MenuID: 1})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedMenu.ID, result.ID)
	assert.Equal(t, expectedMenu.Title, result.Title)
	assert.Equal(t, expectedMenu.Path, result.Path)
}

func TestGetMenuHandler_Handle_NotFound(t *testing.T) {
	mockQryRepo := new(MockMenuQueryRepository)

	mockQryRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, errors.New("not found"))

	handler := NewGetMenuHandler(mockQryRepo)

	result, err := handler.Handle(context.Background(), GetMenuQuery{MenuID: 999})

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "menu not found")
}

// ============================================================
// ListMenusHandler Tests
// ============================================================

func TestListMenusHandler_Handle_Success(t *testing.T) {
	mockQryRepo := new(MockMenuQueryRepository)

	menus := []*domainMenu.Menu{
		newTestMenu(1, "Dashboard"),
		newTestMenu(2, "Users"),
		newTestMenu(3, "Settings"),
	}
	mockQryRepo.On("FindAll", mock.Anything).Return(menus, nil)

	handler := NewListMenusHandler(mockQryRepo)

	result, err := handler.Handle(context.Background(), ListMenusQuery{})

	require.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "Dashboard", result[0].Title)
	assert.Equal(t, "Users", result[1].Title)
}

func TestListMenusHandler_Handle_Empty(t *testing.T) {
	mockQryRepo := new(MockMenuQueryRepository)

	mockQryRepo.On("FindAll", mock.Anything).Return([]*domainMenu.Menu{}, nil)

	handler := NewListMenusHandler(mockQryRepo)

	result, err := handler.Handle(context.Background(), ListMenusQuery{})

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestListMenusHandler_Handle_Error(t *testing.T) {
	mockQryRepo := new(MockMenuQueryRepository)

	mockQryRepo.On("FindAll", mock.Anything).Return(nil, errors.New("database error"))

	handler := NewListMenusHandler(mockQryRepo)

	result, err := handler.Handle(context.Background(), ListMenusQuery{})

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch menus")
}

// ============================================================
// ReorderMenusHandler Tests
// ============================================================

func TestReorderMenusHandler_Handle_Success(t *testing.T) {
	mockCmdRepo := new(MockMenuCommandRepository)
	mockQryRepo := new(MockMenuQueryRepository)

	// Mock all menu existence checks
	mockQryRepo.On("FindByID", mock.Anything, uint(1)).Return(newTestMenu(1, "Menu1"), nil)
	mockQryRepo.On("FindByID", mock.Anything, uint(2)).Return(newTestMenu(2, "Menu2"), nil)
	mockCmdRepo.On("UpdateOrder", mock.Anything, mock.AnythingOfType("[]struct { ID uint; Order int; ParentID *uint }")).Return(nil)

	handler := NewReorderMenusHandler(mockCmdRepo, mockQryRepo)

	err := handler.Handle(context.Background(), ReorderMenusCommand{
		Menus: []MenuItemCommand{
			{ID: 1, Order: 2, ParentID: nil},
			{ID: 2, Order: 1, ParentID: nil},
		},
	})

	require.NoError(t, err)
	mockCmdRepo.AssertExpectations(t)
}

func TestReorderMenusHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		cmd        ReorderMenusCommand
		setupMocks func(*MockMenuCommandRepository, *MockMenuQueryRepository)
		wantErr    string
	}{
		{
			name: "菜单不存在",
			cmd: ReorderMenusCommand{
				Menus: []MenuItemCommand{{ID: 999, Order: 1}},
			},
			setupMocks: func(cmdRepo *MockMenuCommandRepository, qryRepo *MockMenuQueryRepository) {
				qryRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, errors.New("not found"))
			},
			wantErr: "menu 999 not found",
		},
		{
			name: "父菜单不存在",
			cmd: ReorderMenusCommand{
				Menus: []MenuItemCommand{{ID: 1, Order: 1, ParentID: ptrUint(999)}},
			},
			setupMocks: func(cmdRepo *MockMenuCommandRepository, qryRepo *MockMenuQueryRepository) {
				qryRepo.On("FindByID", mock.Anything, uint(1)).Return(newTestMenu(1, "Menu"), nil)
				qryRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, errors.New("not found"))
			},
			wantErr: "parent menu 999 not found",
		},
		{
			name: "更新排序失败",
			cmd: ReorderMenusCommand{
				Menus: []MenuItemCommand{{ID: 1, Order: 1}},
			},
			setupMocks: func(cmdRepo *MockMenuCommandRepository, qryRepo *MockMenuQueryRepository) {
				qryRepo.On("FindByID", mock.Anything, uint(1)).Return(newTestMenu(1, "Menu"), nil)
				cmdRepo.On("UpdateOrder", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			wantErr: "failed to update menu order",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCmdRepo := new(MockMenuCommandRepository)
			mockQryRepo := new(MockMenuQueryRepository)
			tt.setupMocks(mockCmdRepo, mockQryRepo)

			handler := NewReorderMenusHandler(mockCmdRepo, mockQryRepo)
			err := handler.Handle(context.Background(), tt.cmd)

			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}
