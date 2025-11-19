// Package user 定义用户领域模型
package user

import (
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 软删除

	Username string `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"` // 不在 JSON 中返回密码
	FullName string `gorm:"size:100" json:"full_name"`
	Avatar   string `gorm:"size:255" json:"avatar"`
	Bio      string `gorm:"type:text" json:"bio"`
	Status   string `gorm:"size:20;default:'active'" json:"status"` // active, inactive, banned

	// RBAC: Many-to-Many relationship with roles
	Roles []role.Role `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate GORM 钩子：创建前
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 这里可以添加创建前的逻辑，例如密码加密
	return nil
}

// UserCreateDTO 创建用户 DTO
type UserCreateDTO struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"max=100"`
}

// UserUpdateDTO 更新用户 DTO
type UserUpdateDTO struct {
	FullName *string `json:"full_name" binding:"omitempty,max=100"`
	Avatar   *string `json:"avatar" binding:"omitempty,max=255"`
	Bio      *string `json:"bio"`
	Status   *string `json:"status" binding:"omitempty,oneof=active inactive banned"`
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

// ToResponse 转换为响应 DTO
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FullName:  u.FullName,
		Avatar:    u.Avatar,
		Bio:       u.Bio,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// HasRole checks if the user has a specific role
func (u *User) HasRole(roleName string) bool {
	for _, r := range u.Roles {
		if r.Name == roleName {
			return true
		}
	}
	return false
}

// HasAnyRole checks if the user has any of the specified roles
func (u *User) HasAnyRole(roleNames ...string) bool {
	for _, roleName := range roleNames {
		if u.HasRole(roleName) {
			return true
		}
	}
	return false
}

// HasPermission checks if the user has a specific permission
func (u *User) HasPermission(permissionCode string) bool {
	for _, r := range u.Roles {
		for _, p := range r.Permissions {
			if p.Code == permissionCode {
				return true
			}
		}
	}
	return false
}

// GetRoleNames returns all role names of the user
func (u *User) GetRoleNames() []string {
	names := make([]string, 0, len(u.Roles))
	for _, r := range u.Roles {
		names = append(names, r.Name)
	}
	return names
}

// GetPermissions returns all unique permissions of the user
func (u *User) GetPermissions() []role.Permission {
	permissionMap := make(map[uint]role.Permission)
	for _, r := range u.Roles {
		for _, p := range r.Permissions {
			permissionMap[p.ID] = p
		}
	}

	permissions := make([]role.Permission, 0, len(permissionMap))
	for _, p := range permissionMap {
		permissions = append(permissions, p)
	}
	return permissions
}

// GetPermissionCodes returns all unique permission codes of the user
func (u *User) GetPermissionCodes() []string {
	permissions := u.GetPermissions()
	codes := make([]string, 0, len(permissions))
	for _, p := range permissions {
		codes = append(codes, p.Code)
	}
	return codes
}

// IsAdmin checks if the user has admin role
func (u *User) IsAdmin() bool {
	return u.HasRole("admin")
}
