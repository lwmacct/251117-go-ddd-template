package role

import "time"

// CreateRoleDTO 创建角色请求 DTO
type CreateRoleDTO struct {
	Name        string `json:"name" binding:"required,min=2,max=50" example:"developer"`
	DisplayName string `json:"display_name" binding:"required,max=100" example:"开发者"`
	Description string `json:"description" binding:"max=255" example:"系统开发人员角色"`
}

// UpdateRoleDTO 更新角色请求 DTO
type UpdateRoleDTO struct {
	DisplayName *string `json:"display_name,omitempty" binding:"omitempty,max=100"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=255"`
}

// SetPermissionsDTO 设置角色权限请求 DTO
type SetPermissionsDTO struct {
	PermissionIDs []uint `json:"permission_ids" binding:"required"`
}

// CreateRoleResultDTO 创建角色响应 DTO
type CreateRoleResultDTO struct {
	RoleID      uint   `json:"role_id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

// RoleDTO 角色响应 DTO
type RoleDTO struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	DisplayName string           `json:"display_name"`
	Description string           `json:"description"`
	IsSystem    bool             `json:"is_system"`
	Permissions []*PermissionDTO `json:"permissions,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// PermissionDTO 权限响应 DTO
type PermissionDTO struct {
	ID          uint      `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ListRolesDTO 角色列表响应 DTO
type ListRolesDTO struct {
	Roles []*RoleDTO `json:"roles"`
	Total int64      `json:"total"`
	Page  int        `json:"page"`
	Limit int        `json:"limit"`
}

// ListPermissionsDTO 权限列表响应 DTO
type ListPermissionsDTO struct {
	Permissions []*PermissionDTO `json:"permissions"`
	Total       int64            `json:"total"`
	Page        int              `json:"page"`
	Limit       int              `json:"limit"`
}
