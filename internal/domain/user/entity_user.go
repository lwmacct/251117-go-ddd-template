package user

import (
	"slices"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
)

// User 用户实体
type User struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`

	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
	FullName string `json:"full_name"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
	Status   string `json:"status"`

	// RBAC: Many-to-Many relationship with roles
	Roles []role.Role `json:"roles,omitempty"`
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
	return slices.ContainsFunc(roleNames, u.HasRole)
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

// CanLogin 检查用户是否可以登录
func (u *User) CanLogin() bool {
	return u.Status == "active"
}

// IsBanned 检查用户是否被禁用
func (u *User) IsBanned() bool {
	return u.Status == "banned"
}

// IsInactive 检查用户是否未激活
func (u *User) IsInactive() bool {
	return u.Status == "inactive"
}

// Activate 激活用户
func (u *User) Activate() {
	u.Status = "active"
}

// Deactivate 停用用户
func (u *User) Deactivate() {
	u.Status = "inactive"
}

// Ban 禁用用户
func (u *User) Ban() {
	u.Status = "banned"
}

// AssignRole 分配角色（领域行为）
func (u *User) AssignRole(r role.Role) error {
	if u.HasRole(r.Name) {
		return ErrRoleAlreadyAssigned
	}
	u.Roles = append(u.Roles, r)
	return nil
}

// RemoveRole 移除角色（领域行为）
func (u *User) RemoveRole(roleName string) error {
	for i, r := range u.Roles {
		if r.Name == roleName {
			u.Roles = append(u.Roles[:i], u.Roles[i+1:]...)
			return nil
		}
	}
	return ErrRoleNotFound
}

// ClearRoles 清空所有角色
func (u *User) ClearRoles() {
	u.Roles = []role.Role{}
}

// UpdateProfile 更新用户资料（领域行为）
func (u *User) UpdateProfile(fullName, avatar, bio string) {
	if fullName != "" {
		u.FullName = fullName
	}
	if avatar != "" {
		u.Avatar = avatar
	}
	u.Bio = bio // Bio can be empty
}
