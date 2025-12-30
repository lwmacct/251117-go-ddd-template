package user

import "time"

// CreateDTO 创建用户 DTO
type CreateDTO struct {
	Username string  `json:"username" binding:"required,min=3,max=50"`
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	FullName string  `json:"full_name" binding:"max=100"`
	Status   *string `json:"status" binding:"omitempty,oneof=active inactive"`
	RoleIDs  []uint  `json:"role_ids" binding:"omitempty,dive,gt=0"`
}

// UpdateDTO 更新用户 DTO
type UpdateDTO struct {
	Username *string `json:"username" binding:"omitempty,min=3,max=50"`
	Email    *string `json:"email" binding:"omitempty,email"`
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

// UserDTO 用户响应 DTO (不包含敏感信息)
type UserDTO struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Avatar    string    `json:"avatar"`
	Bio       string    `json:"bio"`
	Status    string    `json:"status"`
	Type      string    `json:"type"`
	IsSystem  bool      `json:"is_system"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserWithRolesDTO 用户响应 DTO（包含角色信息）
type UserWithRolesDTO struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Avatar    string    `json:"avatar"`
	Bio       string    `json:"bio"`
	Status    string    `json:"status"`
	Type      string    `json:"type"`
	IsSystem  bool      `json:"is_system"`
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

// UserListDTO 用户列表响应 DTO
type UserListDTO struct {
	Users []*UserDTO `json:"users"`
	Total int64      `json:"total"`
}

// BatchCreateDTO 批量创建用户请求 DTO
type BatchCreateDTO struct {
	Users []BatchItemDTO `json:"users" binding:"required,min=1,max=100,dive"`
}

// BatchItemDTO 批量创建中的单个用户 DTO
type BatchItemDTO struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"max=100"`
	Status   string `json:"status" binding:"omitempty,oneof=active inactive"`
	RoleIDs  []uint `json:"role_ids" binding:"omitempty,dive,gt=0"`
}

// BatchCreateResultDTO 批量创建用户响应 DTO
type BatchCreateResultDTO struct {
	Total   int                   `json:"total"`
	Success int                   `json:"success"`
	Failed  int                   `json:"failed"`
	Errors  []BatchCreateErrorDTO `json:"errors,omitempty"`
}

// BatchCreateErrorDTO 批量创建错误详情 DTO
type BatchCreateErrorDTO struct {
	Index    int    `json:"index"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Error    string `json:"error"`
}

// CreateResultDTO 创建用户结果 DTO（Handler 返回类型）
type CreateResultDTO struct {
	UserID   uint
	Username string
	Email    string
}

// UpdateResultDTO 更新用户结果 DTO（Handler 返回类型）
type UpdateResultDTO struct {
	UserID uint
}
