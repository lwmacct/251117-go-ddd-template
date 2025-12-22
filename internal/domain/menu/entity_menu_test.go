package menu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestMenu(id uint, title string, parentID *uint) *Menu {
	return &Menu{
		ID:       id,
		Title:    title,
		Path:     "/" + title,
		ParentID: parentID,
		Visible:  true,
	}
}

func TestMenu_IsRoot(t *testing.T) {
	tests := []struct {
		name     string
		parentID *uint
		want     bool
	}{
		{
			name:     "顶级菜单 - ParentID 为 nil",
			parentID: nil,
			want:     true,
		},
		{
			name:     "子菜单 - ParentID 不为 nil",
			parentID: ptrUint(1),
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menu := &Menu{ParentID: tt.parentID}
			assert.Equal(t, tt.want, menu.IsRoot())
		})
	}
}

func TestMenu_HasChildren(t *testing.T) {
	t.Run("无子菜单", func(t *testing.T) {
		menu := newTestMenu(1, "Dashboard", nil)
		assert.False(t, menu.HasChildren())
	})

	t.Run("有子菜单", func(t *testing.T) {
		menu := newTestMenu(1, "Dashboard", nil)
		menu.Children = []*Menu{newTestMenu(2, "Overview", ptrUint(1))}
		assert.True(t, menu.HasChildren())
	})
}

func TestMenu_Visibility(t *testing.T) {
	t.Run("默认可见", func(t *testing.T) {
		menu := newTestMenu(1, "Dashboard", nil)
		assert.True(t, menu.IsVisible())
	})

	t.Run("隐藏菜单", func(t *testing.T) {
		menu := newTestMenu(1, "Dashboard", nil)
		menu.Hide()
		assert.False(t, menu.IsVisible())
	})

	t.Run("显示菜单", func(t *testing.T) {
		menu := &Menu{Visible: false}
		menu.Show()
		assert.True(t, menu.IsVisible())
	})
}

func TestMenu_AddChild(t *testing.T) {
	t.Run("添加子菜单", func(t *testing.T) {
		parent := newTestMenu(1, "Parent", nil)
		child := newTestMenu(2, "Child", nil)

		parent.AddChild(child)

		assert.Len(t, parent.Children, 1)
		assert.Equal(t, uint(1), *child.ParentID)
		assert.Equal(t, "Child", parent.Children[0].Title)
	})

	t.Run("添加多个子菜单", func(t *testing.T) {
		parent := newTestMenu(1, "Parent", nil)

		parent.AddChild(newTestMenu(2, "Child1", nil))
		parent.AddChild(newTestMenu(3, "Child2", nil))
		parent.AddChild(newTestMenu(4, "Child3", nil))

		assert.Len(t, parent.Children, 3)
	})

	t.Run("添加 nil 子菜单", func(t *testing.T) {
		parent := newTestMenu(1, "Parent", nil)
		parent.AddChild(nil)
		assert.Empty(t, parent.Children)
	})
}

func TestMenu_RemoveChild(t *testing.T) {
	t.Run("成功移除子菜单", func(t *testing.T) {
		parent := newTestMenu(1, "Parent", nil)
		parent.Children = []*Menu{
			newTestMenu(2, "Child1", ptrUint(1)),
			newTestMenu(3, "Child2", ptrUint(1)),
		}

		removed := parent.RemoveChild(2)

		assert.True(t, removed)
		assert.Len(t, parent.Children, 1)
		assert.Equal(t, "Child2", parent.Children[0].Title)
	})

	t.Run("移除不存在的子菜单", func(t *testing.T) {
		parent := newTestMenu(1, "Parent", nil)
		parent.Children = []*Menu{newTestMenu(2, "Child", ptrUint(1))}

		removed := parent.RemoveChild(999)

		assert.False(t, removed)
		assert.Len(t, parent.Children, 1)
	})

	t.Run("从空列表移除", func(t *testing.T) {
		parent := newTestMenu(1, "Parent", nil)
		removed := parent.RemoveChild(1)
		assert.False(t, removed)
	})
}

func TestMenu_GetChildrenCount(t *testing.T) {
	parent := newTestMenu(1, "Parent", nil)
	assert.Equal(t, 0, parent.GetChildrenCount())

	parent.AddChild(newTestMenu(2, "Child1", nil))
	assert.Equal(t, 1, parent.GetChildrenCount())

	parent.AddChild(newTestMenu(3, "Child2", nil))
	assert.Equal(t, 2, parent.GetChildrenCount())
}

func TestMenu_GetDepth(t *testing.T) {
	t.Run("叶子节点深度为 0", func(t *testing.T) {
		menu := newTestMenu(1, "Leaf", nil)
		assert.Equal(t, 0, menu.GetDepth())
	})

	t.Run("一层子菜单深度为 1", func(t *testing.T) {
		parent := newTestMenu(1, "Parent", nil)
		parent.AddChild(newTestMenu(2, "Child", nil))
		assert.Equal(t, 1, parent.GetDepth())
	})

	t.Run("多层嵌套", func(t *testing.T) {
		root := newTestMenu(1, "Root", nil)
		level1 := newTestMenu(2, "Level1", nil)
		level2 := newTestMenu(3, "Level2", nil)

		level1.AddChild(level2)
		root.AddChild(level1)

		assert.Equal(t, 2, root.GetDepth())
		assert.Equal(t, 1, level1.GetDepth())
		assert.Equal(t, 0, level2.GetDepth())
	})
}

func TestMenu_FindChild(t *testing.T) {
	// 构建三层菜单树
	root := newTestMenu(1, "Root", nil)
	level1a := newTestMenu(2, "Level1A", nil)
	level1b := newTestMenu(3, "Level1B", nil)
	level2 := newTestMenu(4, "Level2", nil)

	level1a.AddChild(level2)
	root.AddChild(level1a)
	root.AddChild(level1b)

	t.Run("查找直接子菜单", func(t *testing.T) {
		found := root.FindChild(2)
		assert.NotNil(t, found)
		assert.Equal(t, "Level1A", found.Title)
	})

	t.Run("递归查找深层子菜单", func(t *testing.T) {
		found := root.FindChild(4)
		assert.NotNil(t, found)
		assert.Equal(t, "Level2", found.Title)
	})

	t.Run("查找不存在的菜单", func(t *testing.T) {
		found := root.FindChild(999)
		assert.Nil(t, found)
	})

	t.Run("从叶子节点查找", func(t *testing.T) {
		found := level2.FindChild(1)
		assert.Nil(t, found)
	})
}

// ptrUint 辅助函数，创建 uint 指针
func ptrUint(v uint) *uint {
	return &v
}
