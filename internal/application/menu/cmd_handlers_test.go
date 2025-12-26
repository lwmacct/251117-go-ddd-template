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
// CreateMenuHandler Tests
// ============================================================

func TestCreateMenuHandler_Handle_Success(t *testing.T) {
	tests := []struct {
		name string
		cmd  CreateMenuCommand
	}{
		{
			name: "创建顶级菜单",
			cmd: CreateMenuCommand{
				Title:   "Dashboard",
				Path:    "/dashboard",
				Icon:    "dashboard",
				Visible: true,
				Order:   1,
			},
		},
		{
			name: "创建子菜单",
			cmd: CreateMenuCommand{
				Title:    "Users",
				Path:     "/users",
				Icon:     "users",
				ParentID: ptrUint(1),
				Visible:  true,
				Order:    2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCmdRepo := new(MockMenuCommandRepository)
			mockQryRepo := new(MockMenuQueryRepository)

			if tt.cmd.ParentID != nil {
				mockQryRepo.On("FindByID", mock.Anything, *tt.cmd.ParentID).Return(newTestMenu(*tt.cmd.ParentID, "Parent"), nil)
			}
			mockCmdRepo.On("Create", mock.Anything, mock.AnythingOfType("*menu.Menu")).Return(nil)

			handler := NewCreateMenuHandler(mockCmdRepo, mockQryRepo)
			result, err := handler.Handle(context.Background(), tt.cmd)

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, uint(1), result.ID)
		})
	}
}

func TestCreateMenuHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		cmd        CreateMenuCommand
		setupMocks func(*MockMenuCommandRepository, *MockMenuQueryRepository)
		wantErr    string
	}{
		{
			name: "父菜单不存在",
			cmd:  CreateMenuCommand{Title: "Child", Path: "/child", ParentID: ptrUint(999)},
			setupMocks: func(cmdRepo *MockMenuCommandRepository, qryRepo *MockMenuQueryRepository) {
				qryRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, errors.New("not found"))
			},
			wantErr: "parent menu not found",
		},
		{
			name: "创建失败",
			cmd:  CreateMenuCommand{Title: "Menu", Path: "/menu"},
			setupMocks: func(cmdRepo *MockMenuCommandRepository, qryRepo *MockMenuQueryRepository) {
				cmdRepo.On("Create", mock.Anything, mock.AnythingOfType("*menu.Menu")).Return(errors.New("db error"))
			},
			wantErr: "failed to create menu",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCmdRepo := new(MockMenuCommandRepository)
			mockQryRepo := new(MockMenuQueryRepository)
			tt.setupMocks(mockCmdRepo, mockQryRepo)

			handler := NewCreateMenuHandler(mockCmdRepo, mockQryRepo)
			result, err := handler.Handle(context.Background(), tt.cmd)

			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

// ============================================================
// UpdateMenuHandler Tests
// ============================================================

func TestUpdateMenuHandler_Handle_Success(t *testing.T) {
	mockCmdRepo := new(MockMenuCommandRepository)
	mockQryRepo := new(MockMenuQueryRepository)

	existingMenu := newTestMenu(1, "OldTitle")
	mockQryRepo.On("FindByID", mock.Anything, uint(1)).Return(existingMenu, nil)
	mockCmdRepo.On("Update", mock.Anything, mock.AnythingOfType("*menu.Menu")).Return(nil)

	handler := NewUpdateMenuHandler(mockCmdRepo, mockQryRepo)

	result, err := handler.Handle(context.Background(), UpdateMenuCommand{
		MenuID: 1,
		Title:  ptrString("NewTitle"),
		Path:   ptrString("/new-path"),
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "NewTitle", existingMenu.Title)
	assert.Equal(t, "/new-path", existingMenu.Path)
}

func TestUpdateMenuHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		cmd        UpdateMenuCommand
		setupMocks func(*MockMenuCommandRepository, *MockMenuQueryRepository)
		wantErr    string
	}{
		{
			name: "菜单不存在",
			cmd:  UpdateMenuCommand{MenuID: 999},
			setupMocks: func(cmdRepo *MockMenuCommandRepository, qryRepo *MockMenuQueryRepository) {
				qryRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, errors.New("not found"))
			},
			wantErr: "menu not found",
		},
		{
			name: "不能设置自己为父菜单",
			cmd:  UpdateMenuCommand{MenuID: 1, ParentID: ptrUint(1)},
			setupMocks: func(cmdRepo *MockMenuCommandRepository, qryRepo *MockMenuQueryRepository) {
				qryRepo.On("FindByID", mock.Anything, uint(1)).Return(newTestMenu(1, "Menu"), nil)
			},
			wantErr: "cannot set menu as its own parent",
		},
		{
			name: "父菜单不存在",
			cmd:  UpdateMenuCommand{MenuID: 1, ParentID: ptrUint(999)},
			setupMocks: func(cmdRepo *MockMenuCommandRepository, qryRepo *MockMenuQueryRepository) {
				qryRepo.On("FindByID", mock.Anything, uint(1)).Return(newTestMenu(1, "Menu"), nil)
				qryRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, errors.New("not found"))
			},
			wantErr: "parent menu not found",
		},
		{
			name: "更新失败",
			cmd:  UpdateMenuCommand{MenuID: 1, Title: ptrString("New")},
			setupMocks: func(cmdRepo *MockMenuCommandRepository, qryRepo *MockMenuQueryRepository) {
				qryRepo.On("FindByID", mock.Anything, uint(1)).Return(newTestMenu(1, "Old"), nil)
				cmdRepo.On("Update", mock.Anything, mock.AnythingOfType("*menu.Menu")).Return(errors.New("db error"))
			},
			wantErr: "failed to update menu",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCmdRepo := new(MockMenuCommandRepository)
			mockQryRepo := new(MockMenuQueryRepository)
			tt.setupMocks(mockCmdRepo, mockQryRepo)

			handler := NewUpdateMenuHandler(mockCmdRepo, mockQryRepo)
			result, err := handler.Handle(context.Background(), tt.cmd)

			assert.Nil(t, result)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

// ============================================================
// DeleteMenuHandler Tests
// ============================================================

func TestDeleteMenuHandler_Handle_Success(t *testing.T) {
	mockCmdRepo := new(MockMenuCommandRepository)
	mockQryRepo := new(MockMenuQueryRepository)

	mockQryRepo.On("FindByID", mock.Anything, uint(1)).Return(newTestMenu(1, "Menu"), nil)
	mockQryRepo.On("FindByParentID", mock.Anything, ptrUint(1)).Return([]*domainMenu.Menu{}, nil)
	mockCmdRepo.On("Delete", mock.Anything, uint(1)).Return(nil)

	handler := NewDeleteMenuHandler(mockCmdRepo, mockQryRepo)

	err := handler.Handle(context.Background(), DeleteMenuCommand{MenuID: 1})

	require.NoError(t, err)
	mockCmdRepo.AssertExpectations(t)
}

func TestDeleteMenuHandler_Handle_Error(t *testing.T) {
	tests := []struct {
		name       string
		cmd        DeleteMenuCommand
		setupMocks func(*MockMenuCommandRepository, *MockMenuQueryRepository)
		wantErr    string
	}{
		{
			name: "菜单不存在",
			cmd:  DeleteMenuCommand{MenuID: 999},
			setupMocks: func(cmdRepo *MockMenuCommandRepository, qryRepo *MockMenuQueryRepository) {
				qryRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, errors.New("not found"))
			},
			wantErr: "menu not found",
		},
		{
			name: "检查子菜单失败",
			cmd:  DeleteMenuCommand{MenuID: 1},
			setupMocks: func(cmdRepo *MockMenuCommandRepository, qryRepo *MockMenuQueryRepository) {
				qryRepo.On("FindByID", mock.Anything, uint(1)).Return(newTestMenu(1, "Menu"), nil)
				qryRepo.On("FindByParentID", mock.Anything, ptrUint(1)).Return(nil, errors.New("db error"))
			},
			wantErr: "failed to check children",
		},
		{
			name: "有子菜单无法删除",
			cmd:  DeleteMenuCommand{MenuID: 1},
			setupMocks: func(cmdRepo *MockMenuCommandRepository, qryRepo *MockMenuQueryRepository) {
				qryRepo.On("FindByID", mock.Anything, uint(1)).Return(newTestMenu(1, "Parent"), nil)
				qryRepo.On("FindByParentID", mock.Anything, ptrUint(1)).Return([]*domainMenu.Menu{
					newTestMenuWithParent(2, "Child", 1),
				}, nil)
			},
			wantErr: "cannot delete menu with children",
		},
		{
			name: "删除失败",
			cmd:  DeleteMenuCommand{MenuID: 1},
			setupMocks: func(cmdRepo *MockMenuCommandRepository, qryRepo *MockMenuQueryRepository) {
				qryRepo.On("FindByID", mock.Anything, uint(1)).Return(newTestMenu(1, "Menu"), nil)
				qryRepo.On("FindByParentID", mock.Anything, ptrUint(1)).Return([]*domainMenu.Menu{}, nil)
				cmdRepo.On("Delete", mock.Anything, uint(1)).Return(errors.New("db error"))
			},
			wantErr: "failed to delete menu",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCmdRepo := new(MockMenuCommandRepository)
			mockQryRepo := new(MockMenuQueryRepository)
			tt.setupMocks(mockCmdRepo, mockQryRepo)

			handler := NewDeleteMenuHandler(mockCmdRepo, mockQryRepo)
			err := handler.Handle(context.Background(), tt.cmd)

			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

// ============================================================
// UpdateMenuHandler applyUpdates 覆盖测试
// ============================================================

func TestUpdateMenuHandler_ApplyUpdates_AllFields(t *testing.T) {
	mockCmdRepo := new(MockMenuCommandRepository)
	mockQryRepo := new(MockMenuQueryRepository)

	existingMenu := newTestMenu(1, "OldTitle")
	mockQryRepo.On("FindByID", mock.Anything, uint(1)).Return(existingMenu, nil)
	mockCmdRepo.On("Update", mock.Anything, mock.AnythingOfType("*menu.Menu")).Return(nil)

	handler := NewUpdateMenuHandler(mockCmdRepo, mockQryRepo)

	result, err := handler.Handle(context.Background(), UpdateMenuCommand{
		MenuID:   1,
		Title:    ptrString("NewTitle"),
		Path:     ptrString("/new-path"),
		Icon:     ptrString("new-icon"),
		ParentID: ptrUint(0),
		Order:    ptrInt(10),
		Visible:  ptrBool(false),
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "NewTitle", existingMenu.Title)
	assert.Equal(t, "/new-path", existingMenu.Path)
	assert.Equal(t, "new-icon", existingMenu.Icon)
	assert.Equal(t, 10, existingMenu.Order)
	assert.False(t, existingMenu.Visible)
}

func TestUpdateMenuHandler_ValidateParentID_ZeroValue(t *testing.T) {
	mockCmdRepo := new(MockMenuCommandRepository)
	mockQryRepo := new(MockMenuQueryRepository)

	existingMenu := newTestMenu(1, "Menu")
	mockQryRepo.On("FindByID", mock.Anything, uint(1)).Return(existingMenu, nil)
	mockCmdRepo.On("Update", mock.Anything, mock.AnythingOfType("*menu.Menu")).Return(nil)

	handler := NewUpdateMenuHandler(mockCmdRepo, mockQryRepo)

	// ParentID = 0 表示设为顶级菜单，应该成功
	result, err := handler.Handle(context.Background(), UpdateMenuCommand{
		MenuID:   1,
		ParentID: ptrUint(0),
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
}
