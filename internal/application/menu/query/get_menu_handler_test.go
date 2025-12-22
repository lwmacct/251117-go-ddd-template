package query

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
)

func TestGetMenuHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功获取菜单", func(t *testing.T) {
		mockRepo := &mockMenuQueryRepo{
			menu: &menu.Menu{
				ID:      1,
				Title:   "首页",
				Path:    "/home",
				Icon:    "home",
				Visible: true,
			},
		}

		handler := NewGetMenuHandler(mockRepo)
		result, err := handler.Handle(ctx, GetMenuQuery{MenuID: 1})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "首页", result.Title)
		assert.Equal(t, "/home", result.Path)
	})

	t.Run("菜单不存在", func(t *testing.T) {
		mockRepo := &mockMenuQueryRepo{
			findErr: errors.New("menu not found"),
		}

		handler := NewGetMenuHandler(mockRepo)
		result, err := handler.Handle(ctx, GetMenuQuery{MenuID: 999})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "menu not found")
		assert.Nil(t, result)
	})

	t.Run("获取带子菜单的菜单", func(t *testing.T) {
		parentID := uint(1)
		mockRepo := &mockMenuQueryRepo{
			menu: &menu.Menu{
				ID:      1,
				Title:   "系统管理",
				Path:    "/system",
				Visible: true,
				Children: []*menu.Menu{
					{ID: 2, Title: "用户管理", Path: "/system/users", ParentID: &parentID},
					{ID: 3, Title: "角色管理", Path: "/system/roles", ParentID: &parentID},
				},
			},
		}

		handler := NewGetMenuHandler(mockRepo)
		result, err := handler.Handle(ctx, GetMenuQuery{MenuID: 1})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "系统管理", result.Title)
	})

	t.Run("数据库错误", func(t *testing.T) {
		mockRepo := &mockMenuQueryRepo{
			findErr: errors.New("database connection failed"),
		}

		handler := NewGetMenuHandler(mockRepo)
		result, err := handler.Handle(ctx, GetMenuQuery{MenuID: 1})

		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestNewGetMenuHandler(t *testing.T) {
	mockRepo := &mockMenuQueryRepo{}
	handler := NewGetMenuHandler(mockRepo)

	assert.NotNil(t, handler)
}
