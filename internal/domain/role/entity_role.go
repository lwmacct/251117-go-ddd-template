package role

import "time"

// Role 角色实体，RBAC 系统的核心组件
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

// IsSystemRole 检查是否为系统角色
func (r *Role) IsSystemRole() bool {
	return r.IsSystem
}

// CanBeDeleted 检查角色是否可以被删除（系统角色不可删除）
func (r *Role) CanBeDeleted() bool {
	return !r.IsSystem
}

// CanBeModified 检查角色是否可以被修改（系统角色不可修改）
func (r *Role) CanBeModified() bool {
	return !r.IsSystem
}

// HasPermission 检查角色是否拥有指定权限
func (r *Role) HasPermission(code string) bool {
	for _, p := range r.Permissions {
		if p.Code == code {
			return true
		}
	}
	return false
}

// HasAnyPermission 检查角色是否拥有任一指定权限
func (r *Role) HasAnyPermission(codes ...string) bool {
	codeSet := make(map[string]struct{}, len(codes))
	for _, code := range codes {
		codeSet[code] = struct{}{}
	}
	for _, p := range r.Permissions {
		if _, ok := codeSet[p.Code]; ok {
			return true
		}
	}
	return false
}

// GetPermissionCodes 获取所有权限代码
func (r *Role) GetPermissionCodes() []string {
	codes := make([]string, len(r.Permissions))
	for i, p := range r.Permissions {
		codes[i] = p.Code
	}
	return codes
}

// GetPermissionCount 获取权限数量
func (r *Role) GetPermissionCount() int {
	return len(r.Permissions)
}

// IsEmpty 检查角色是否没有权限
func (r *Role) IsEmpty() bool {
	return len(r.Permissions) == 0
}

// AddPermission 添加权限到角色
func (r *Role) AddPermission(p Permission) {
	if !r.HasPermission(p.Code) {
		r.Permissions = append(r.Permissions, p)
	}
}

// RemovePermission 从角色移除权限
func (r *Role) RemovePermission(code string) bool {
	for i, p := range r.Permissions {
		if p.Code == code {
			r.Permissions = append(r.Permissions[:i], r.Permissions[i+1:]...)
			return true
		}
	}
	return false
}

// ClearPermissions 清空所有权限
func (r *Role) ClearPermissions() {
	r.Permissions = nil
}
