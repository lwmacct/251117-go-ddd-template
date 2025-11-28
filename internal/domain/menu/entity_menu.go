// Package menu 定义菜单领域模型。
//
// 本包管理系统导航菜单的领域逻辑，支持：
//   - 树形结构：通过 ParentID 实现无限级菜单嵌套
//   - 可见性控制：Visible 字段控制菜单显示/隐藏
//   - 排序支持：Order 字段定义同级菜单的显示顺序
//
// Menu 实体通常与 RBAC 权限系统配合使用，实现基于角色的菜单过滤。
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
