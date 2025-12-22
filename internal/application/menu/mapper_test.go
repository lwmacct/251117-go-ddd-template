package menu

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	domainMenu "github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
)

func TestToMenuResponse(t *testing.T) {
	t.Run("转换正常菜单", func(t *testing.T) {
		now := time.Now()
		m := &domainMenu.Menu{
			ID:        1,
			Title:     "首页",
			Path:      "/home",
			Icon:      "home",
			ParentID:  nil,
			Order:     1,
			Visible:   true,
			CreatedAt: now,
			UpdatedAt: now,
		}

		result := ToMenuResponse(m)

		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "首页", result.Title)
		assert.Equal(t, "/home", result.Path)
		assert.Equal(t, "home", result.Icon)
		assert.Nil(t, result.ParentID)
		assert.Equal(t, 1, result.Order)
		assert.True(t, result.Visible)
		assert.Equal(t, now, result.CreatedAt)
		assert.Equal(t, now, result.UpdatedAt)
	})

	t.Run("转换nil返回nil", func(t *testing.T) {
		result := ToMenuResponse(nil)
		assert.Nil(t, result)
	})

	t.Run("转换带父菜单ID的菜单", func(t *testing.T) {
		parentID := uint(1)
		m := &domainMenu.Menu{
			ID:       2,
			Title:    "用户管理",
			Path:     "/system/users",
			Icon:     "user",
			ParentID: &parentID,
			Order:    1,
			Visible:  true,
		}

		result := ToMenuResponse(m)

		assert.NotNil(t, result)
		assert.Equal(t, uint(2), result.ID)
		assert.Equal(t, "用户管理", result.Title)
		assert.NotNil(t, result.ParentID)
		assert.Equal(t, uint(1), *result.ParentID)
	})

	t.Run("转换不可见菜单", func(t *testing.T) {
		m := &domainMenu.Menu{
			ID:      3,
			Title:   "隐藏菜单",
			Path:    "/hidden",
			Visible: false,
		}

		result := ToMenuResponse(m)

		assert.NotNil(t, result)
		assert.Equal(t, "隐藏菜单", result.Title)
		assert.False(t, result.Visible)
	})

	t.Run("转换空字段菜单", func(t *testing.T) {
		m := &domainMenu.Menu{
			ID:    4,
			Title: "最小菜单",
		}

		result := ToMenuResponse(m)

		assert.NotNil(t, result)
		assert.Equal(t, uint(4), result.ID)
		assert.Equal(t, "最小菜单", result.Title)
		assert.Empty(t, result.Path)
		assert.Empty(t, result.Icon)
		assert.Equal(t, 0, result.Order)
	})
}
