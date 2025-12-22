package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
)

// mockMenuCommandRepo 是菜单命令仓储的 Mock 实现。
type mockMenuCommandRepo struct {
	menus       []*menu.Menu
	createError error
	idCounter   uint
}

func newMockMenuCommandRepo() *mockMenuCommandRepo {
	return &mockMenuCommandRepo{
		menus:     make([]*menu.Menu, 0),
		idCounter: 0,
	}
}

func (m *mockMenuCommandRepo) Create(ctx context.Context, menuItem *menu.Menu) error {
	if m.createError != nil {
		return m.createError
	}
	m.idCounter++
	menuItem.ID = m.idCounter
	m.menus = append(m.menus, menuItem)
	return nil
}

func (m *mockMenuCommandRepo) Update(ctx context.Context, menuItem *menu.Menu) error {
	return nil
}

func (m *mockMenuCommandRepo) Delete(ctx context.Context, id uint) error {
	return nil
}

func (m *mockMenuCommandRepo) UpdateOrder(ctx context.Context, menus []struct {
	ID       uint
	Order    int
	ParentID *uint
}) error {
	return nil
}

// mockMenuQueryRepo 是菜单查询仓储的 Mock 实现。
type mockMenuQueryRepo struct {
	menus     map[uint]*menu.Menu
	findError error
}

func newMockMenuQueryRepo() *mockMenuQueryRepo {
	return &mockMenuQueryRepo{
		menus: make(map[uint]*menu.Menu),
	}
}

func (m *mockMenuQueryRepo) FindByID(ctx context.Context, id uint) (*menu.Menu, error) {
	if m.findError != nil {
		return nil, m.findError
	}
	if menuItem, ok := m.menus[id]; ok {
		return menuItem, nil
	}
	return nil, errors.New("menu not found")
}

func (m *mockMenuQueryRepo) FindAll(ctx context.Context) ([]*menu.Menu, error) {
	result := make([]*menu.Menu, 0, len(m.menus))
	for _, menuItem := range m.menus {
		result = append(result, menuItem)
	}
	return result, nil
}

func (m *mockMenuQueryRepo) FindByParentID(ctx context.Context, parentID *uint) ([]*menu.Menu, error) {
	return nil, nil
}

type createMenuTestSetup struct {
	handler     *CreateMenuHandler
	commandRepo *mockMenuCommandRepo
	queryRepo   *mockMenuQueryRepo
}

func setupCreateMenuTest() *createMenuTestSetup {
	commandRepo := newMockMenuCommandRepo()
	queryRepo := newMockMenuQueryRepo()
	handler := NewCreateMenuHandler(commandRepo, queryRepo)

	return &createMenuTestSetup{
		handler:     handler,
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
	}
}

func TestCreateMenuHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功创建顶级菜单", func(t *testing.T) {
		setup := setupCreateMenuTest()

		cmd := CreateMenuCommand{
			Title:    "Dashboard",
			Path:     "/dashboard",
			Icon:     "mdi-home",
			ParentID: nil,
			Order:    1,
			Visible:  true,
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)

		// 验证菜单已创建
		require.Len(t, setup.commandRepo.menus, 1)
		assert.Equal(t, "Dashboard", setup.commandRepo.menus[0].Title)
		assert.Equal(t, "/dashboard", setup.commandRepo.menus[0].Path)
		assert.True(t, setup.commandRepo.menus[0].Visible)
	})

	t.Run("成功创建子菜单", func(t *testing.T) {
		setup := setupCreateMenuTest()

		// 先添加父菜单到 queryRepo
		parentID := uint(1)
		setup.queryRepo.menus[parentID] = &menu.Menu{
			ID:    parentID,
			Title: "Parent",
			Path:  "/parent",
		}

		cmd := CreateMenuCommand{
			Title:    "Child Menu",
			Path:     "/parent/child",
			Icon:     "mdi-folder",
			ParentID: &parentID,
			Order:    1,
			Visible:  true,
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)

		// 验证子菜单已创建且 ParentID 正确
		require.Len(t, setup.commandRepo.menus, 1)
		assert.Equal(t, &parentID, setup.commandRepo.menus[0].ParentID)
	})

	t.Run("父菜单不存在", func(t *testing.T) {
		setup := setupCreateMenuTest()

		nonExistentID := uint(999)
		cmd := CreateMenuCommand{
			Title:    "Orphan Menu",
			Path:     "/orphan",
			ParentID: &nonExistentID,
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "parent menu not found")
		assert.Nil(t, result)
	})

	t.Run("创建菜单失败", func(t *testing.T) {
		setup := setupCreateMenuTest()
		setup.commandRepo.createError = errors.New("database error")

		cmd := CreateMenuCommand{
			Title: "Failing Menu",
			Path:  "/fail",
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create menu")
		assert.Nil(t, result)
	})

	t.Run("隐藏菜单", func(t *testing.T) {
		setup := setupCreateMenuTest()

		cmd := CreateMenuCommand{
			Title:   "Hidden Menu",
			Path:    "/hidden",
			Visible: false,
		}

		result, err := setup.handler.Handle(ctx, cmd)

		require.NoError(t, err)
		require.NotNil(t, result)

		require.Len(t, setup.commandRepo.menus, 1)
		assert.False(t, setup.commandRepo.menus[0].Visible)
	})
}

func TestCreateMenuHandler_ValidationScenarios(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		cmd       CreateMenuCommand
		setupFunc func(*createMenuTestSetup)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "有效的顶级菜单",
			cmd: CreateMenuCommand{
				Title:   "Users",
				Path:    "/users",
				Icon:    "mdi-account",
				Order:   1,
				Visible: true,
			},
			wantErr: false,
		},
		{
			name: "有效的子菜单",
			cmd: CreateMenuCommand{
				Title:    "User List",
				Path:     "/users/list",
				ParentID: ptrUint(1),
			},
			setupFunc: func(s *createMenuTestSetup) {
				s.queryRepo.menus[1] = &menu.Menu{ID: 1, Title: "Users"}
			},
			wantErr: false,
		},
		{
			name: "父菜单不存在",
			cmd: CreateMenuCommand{
				Title:    "Invalid Child",
				Path:     "/invalid",
				ParentID: ptrUint(999),
			},
			wantErr: true,
			errMsg:  "parent menu not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupCreateMenuTest()
			if tt.setupFunc != nil {
				tt.setupFunc(setup)
			}

			result, err := setup.handler.Handle(ctx, tt.cmd)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestNewCreateMenuHandler(t *testing.T) {
	commandRepo := newMockMenuCommandRepo()
	queryRepo := newMockMenuQueryRepo()

	handler := NewCreateMenuHandler(commandRepo, queryRepo)

	assert.NotNil(t, handler)
}

// ptrUint 辅助函数，创建 uint 指针
func ptrUint(v uint) *uint {
	return &v
}
