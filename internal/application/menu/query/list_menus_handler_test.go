package query

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
)

// mockMenuQueryRepo 菜单查询仓储 Mock
type mockMenuQueryRepo struct {
	menus   []*menu.Menu
	menu    *menu.Menu
	findErr error
}

func (m *mockMenuQueryRepo) FindByID(ctx context.Context, id uint) (*menu.Menu, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	return m.menu, nil
}

func (m *mockMenuQueryRepo) FindAll(ctx context.Context) ([]*menu.Menu, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	return m.menus, nil
}

func (m *mockMenuQueryRepo) FindByParentID(ctx context.Context, parentID *uint) ([]*menu.Menu, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	return m.menus, nil
}

func TestListMenusHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功获取菜单列表", func(t *testing.T) {
		mockRepo := &mockMenuQueryRepo{
			menus: []*menu.Menu{
				{ID: 1, Title: "首页", Path: "/home", Icon: "home", Visible: true},
				{ID: 2, Title: "用户管理", Path: "/users", Icon: "user", Visible: true},
				{ID: 3, Title: "系统设置", Path: "/settings", Icon: "setting", Visible: true},
			},
		}

		handler := NewListMenusHandler(mockRepo)
		result, err := handler.Handle(ctx, ListMenusQuery{})

		require.NoError(t, err)
		assert.Len(t, result, 3)
		assert.Equal(t, "首页", result[0].Title)
		assert.Equal(t, "/home", result[0].Path)
	})

	t.Run("获取空菜单列表", func(t *testing.T) {
		mockRepo := &mockMenuQueryRepo{
			menus: []*menu.Menu{},
		}

		handler := NewListMenusHandler(mockRepo)
		result, err := handler.Handle(ctx, ListMenusQuery{})

		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("获取带子菜单的树形结构", func(t *testing.T) {
		parentID := uint(1)
		mockRepo := &mockMenuQueryRepo{
			menus: []*menu.Menu{
				{
					ID:      1,
					Title:   "系统管理",
					Path:    "/system",
					Visible: true,
					Children: []*menu.Menu{
						{ID: 2, Title: "用户管理", Path: "/system/users", ParentID: &parentID},
						{ID: 3, Title: "角色管理", Path: "/system/roles", ParentID: &parentID},
					},
				},
			},
		}

		handler := NewListMenusHandler(mockRepo)
		result, err := handler.Handle(ctx, ListMenusQuery{})

		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "系统管理", result[0].Title)
	})

	t.Run("仓储返回错误", func(t *testing.T) {
		mockRepo := &mockMenuQueryRepo{
			findErr: errors.New("database error"),
		}

		handler := NewListMenusHandler(mockRepo)
		result, err := handler.Handle(ctx, ListMenusQuery{})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch menus")
		assert.Nil(t, result)
	})
}

func TestNewListMenusHandler(t *testing.T) {
	mockRepo := &mockMenuQueryRepo{}
	handler := NewListMenusHandler(mockRepo)

	assert.NotNil(t, handler)
}
