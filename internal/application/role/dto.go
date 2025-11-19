// Package role 定义角色模块的 DTO
package role

import "time"

// RoleResponse 角色响应 DTO
type RoleResponse struct {
	ID          uint                  `json:"id"`
	Name        string                `json:"name"`
	DisplayName string                `json:"display_name"`
	Description string                `json:"description"`
	IsSystem    bool                  `json:"is_system"`
	Permissions []*PermissionResponse `json:"permissions,omitempty"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

// PermissionResponse 权限响应 DTO
type PermissionResponse struct {
	ID          uint      `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ListRolesResponse 角色列表响应 DTO
type ListRolesResponse struct {
	Roles []*RoleResponse `json:"roles"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
}

// ListPermissionsResponse 权限列表响应 DTO
type ListPermissionsResponse struct {
	Permissions []*PermissionResponse `json:"permissions"`
	Total       int64                 `json:"total"`
	Page        int                   `json:"page"`
	Limit       int                   `json:"limit"`
}
