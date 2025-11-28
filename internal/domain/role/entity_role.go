// Package role 定义角色和权限领域模型。
//
// 本包是 RBAC (基于角色的访问控制) 系统的核心领域层，包含：
//   - Role: 角色实体，支持多对多关联权限
//   - Permission: 权限实体，采用 domain:resource:action 三段式格式
//   - Repository 接口: 定义角色和权限的数据访问契约
//
// 设计原则：
//   - 领域模型不依赖任何基础设施层（无 GORM Tag）
//   - 通过接口定义数据访问，实现依赖倒置
package role

import "time"

// Role represents a role in the RBAC system
type Role struct {
	ID          uint         `json:"id"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   *time.Time   `json:"-"`
	Name        string       `json:"name"`
	DisplayName string       `json:"display_name"`
	Description string       `json:"description"`
	IsSystem    bool         `json:"is_system"`
	Permissions []Permission `json:"permissions,omitempty"`
}
