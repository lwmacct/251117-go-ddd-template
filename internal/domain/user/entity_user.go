package user

import (
	"slices"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/permission"
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

	// Type 用户类型：human（人类用户）或 service（服务账户）。
	// 服务账户无密码，仅通过 PAT 认证。
	Type UserType `json:"type"`

	// IsSystem 系统预置用户标记。
	// 系统用户（如 root、admin）不可删除，部分字段不可修改。
	IsSystem bool `json:"is_system"`

	// RBAC: Many-to-Many relationship with roles
	Roles []role.Role `json:"roles,omitempty"`
}

// HasRole 检查用户是否拥有指定角色
func (u *User) HasRole(roleName string) bool {
	for _, r := range u.Roles {
		if r.Name == roleName {
			return true
		}
	}
	return false
}

// HasAnyRole 检查用户是否拥有任一指定角色
func (u *User) HasAnyRole(roleNames ...string) bool {
	return slices.ContainsFunc(roleNames, u.HasRole)
}

// HasPermission 检查用户是否有指定操作对指定资源的权限。
// 遍历用户所有角色的权限进行模式匹配。
func (u *User) HasPermission(op permission.Operation, res permission.Resource) bool {
	for _, r := range u.Roles {
		if r.HasPermission(op, res) {
			return true
		}
	}
	return false
}

// HasOperationPermission 检查用户是否有指定操作的权限（资源为 *）。
func (u *User) HasOperationPermission(op permission.Operation) bool {
	return u.HasPermission(op, permission.ResourceAll)
}

// GetRoleNames 获取用户所有角色名称
func (u *User) GetRoleNames() []string {
	names := make([]string, 0, len(u.Roles))
	for _, r := range u.Roles {
		names = append(names, r.Name)
	}
	return names
}

// GetPermissions 获取用户所有去重后的权限。
// 基于 OperationPattern + ResourcePattern 去重。
func (u *User) GetPermissions() []role.Permission {
	seen := make(map[string]bool)
	var permissions []role.Permission

	for _, r := range u.Roles {
		for _, p := range r.Permissions {
			key := p.OperationPattern + ":" + p.ResourcePattern
			if !seen[key] {
				seen[key] = true
				permissions = append(permissions, p)
			}
		}
	}

	return permissions
}

// IsAdmin 检查用户是否拥有管理员角色
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

// ============================================================================
// 类型判断方法
// ============================================================================

// IsHuman 报告用户是否为人类用户。
func (u *User) IsHuman() bool {
	return u.Type == UserTypeHuman
}

// IsServiceAccount 报告用户是否为服务账户。
func (u *User) IsServiceAccount() bool {
	return u.Type == UserTypeService
}

// IsSystemUser 报告用户是否为系统预置用户。
// 系统用户不可删除，部分字段不可修改。
func (u *User) IsSystemUser() bool {
	return u.IsSystem
}

// IsRoot 报告用户是否为 root 超级管理员。
func (u *User) IsRoot() bool {
	return u.Username == RootUsername
}

// ============================================================================
// 保护策略方法
// ============================================================================

// CanBeDeleted 报告用户是否可以被删除。
// 系统用户不可删除。
func (u *User) CanBeDeleted() bool {
	return !u.IsSystem
}

// CanModifyUsername 报告用户名是否可以被修改。
// 系统用户的用户名不可修改。
func (u *User) CanModifyUsername() bool {
	return !u.IsSystem
}

// CanModifyStatus 报告用户状态是否可以被修改。
// 仅 root 用户状态不可修改。
func (u *User) CanModifyStatus() bool {
	return u.Username != RootUsername
}

// CanModifyRoles 报告用户角色是否可以被修改。
// root 用户角色不可修改（始终拥有 *:*:* 权限）。
func (u *User) CanModifyRoles() bool {
	return u.Username != RootUsername
}

// ============================================================================
// 认证相关方法
// ============================================================================

// CanPasswordLogin 报告用户是否可以使用密码登录。
// 服务账户不支持密码登录。
func (u *User) CanPasswordLogin() bool {
	return u.IsHuman() && u.CanLogin()
}

// RequiresPAT 报告用户是否必须使用 PAT 认证。
// 服务账户仅支持 PAT 认证。
func (u *User) RequiresPAT() bool {
	return u.IsServiceAccount()
}
