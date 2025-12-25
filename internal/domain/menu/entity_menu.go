package menu

import (
	"time"
)

// Menu 菜单实体
type Menu struct {
	ID        uint       `json:"id"`
	Title     string     `json:"title"`
	Path      string     `json:"path"`
	Icon      string     `json:"icon"`
	ParentID  *uint      `json:"parent_id"`
	Order     int        `json:"order"`
	Visible   bool       `json:"visible"`
	Children  []*Menu    `json:"children,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

// IsRoot 检查是否为顶级菜单
func (m *Menu) IsRoot() bool {
	return m.ParentID == nil
}

// HasChildren 检查是否有子菜单
func (m *Menu) HasChildren() bool {
	return len(m.Children) > 0
}

// IsVisible 检查菜单是否可见
func (m *Menu) IsVisible() bool {
	return m.Visible
}

// Hide 隐藏菜单
func (m *Menu) Hide() {
	m.Visible = false
}

// Show 显示菜单
func (m *Menu) Show() {
	m.Visible = true
}

// AddChild 添加子菜单
func (m *Menu) AddChild(child *Menu) {
	if child == nil {
		return
	}
	child.ParentID = &m.ID
	m.Children = append(m.Children, child)
}

// RemoveChild 移除指定 ID 的子菜单
func (m *Menu) RemoveChild(childID uint) bool {
	for i, child := range m.Children {
		if child.ID == childID {
			m.Children = append(m.Children[:i], m.Children[i+1:]...)
			return true
		}
	}
	return false
}

// GetChildrenCount 获取子菜单数量
func (m *Menu) GetChildrenCount() int {
	return len(m.Children)
}

// GetDepth 获取菜单在树中的深度（根节点深度为 0）
// 注意：此方法需要完整的树结构，仅计算 Children 链
func (m *Menu) GetDepth() int {
	if !m.HasChildren() {
		return 0
	}
	maxChildDepth := 0
	for _, child := range m.Children {
		childDepth := child.GetDepth()
		if childDepth > maxChildDepth {
			maxChildDepth = childDepth
		}
	}
	return maxChildDepth + 1
}

// FindChild 在子菜单中查找指定 ID 的菜单（递归）
func (m *Menu) FindChild(id uint) *Menu {
	for _, child := range m.Children {
		if child.ID == id {
			return child
		}
		if found := child.FindChild(id); found != nil {
			return found
		}
	}
	return nil
}
