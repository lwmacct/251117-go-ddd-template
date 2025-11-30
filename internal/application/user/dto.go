// Package user 定义用户应用层的 DTO
package user

import "time"

// CreateUserDTO 创建用户 DTO
type CreateUserDTO struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"max=100"`
	RoleIDs  []uint `json:"role_ids" binding:"omitempty,dive,gt=0"`
}

// UpdateUserDTO 更新用户 DTO
type UpdateUserDTO struct {
	FullName *string `json:"full_name" binding:"omitempty,max=100"`
	Avatar   *string `json:"avatar" binding:"omitempty,max=255"`
	Bio      *string `json:"bio"`
	Status   *string `json:"status" binding:"omitempty,oneof=active inactive banned"`
}

// ChangePasswordDTO 修改密码 DTO
type ChangePasswordDTO struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// AssignRolesDTO 分配角色 DTO
type AssignRolesDTO struct {
	RoleIDs []uint `json:"role_ids" binding:"required"`
}

// UserResponse 用户响应 DTO (不包含敏感信息)
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Avatar    string    `json:"avatar"`
	Bio       string    `json:"bio"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserWithRolesResponse 用户响应 DTO（包含角色信息）
type UserWithRolesResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Avatar    string    `json:"avatar"`
	Bio       string    `json:"bio"`
	Status    string    `json:"status"`
	Roles     []RoleDTO `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RoleDTO 角色 DTO（嵌套在用户响应中）
type RoleDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

// UserListResponse 用户列表响应 DTO
type UserListResponse struct {
	Users []*UserResponse `json:"users"`
	Total int64           `json:"total"`
}

// BatchCreateUserDTO 批量创建用户请求 DTO
type BatchCreateUserDTO struct {
	Users []BatchUserItemDTO `json:"users" binding:"required,min=1,max=100,dive"`
}

// BatchUserItemDTO 批量创建中的单个用户 DTO
type BatchUserItemDTO struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"max=100"`
	Status   string `json:"status" binding:"omitempty,oneof=active inactive"`
	RoleIDs  []uint `json:"role_ids" binding:"omitempty,dive,gt=0"`
}

// BatchCreateUserResponse 批量创建用户响应 DTO
type BatchCreateUserResponse struct {
	Total   int                      `json:"total"`
	Success int                      `json:"success"`
	Failed  int                      `json:"failed"`
	Errors  []BatchCreateErrorDetail `json:"errors,omitempty"`
}

// BatchCreateErrorDetail 批量创建错误详情
type BatchCreateErrorDetail struct {
	Index    int    `json:"index"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Error    string `json:"error"`
}
