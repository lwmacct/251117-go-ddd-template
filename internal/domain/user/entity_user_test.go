package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
)

// newTestPermission 创建测试用权限。
func newTestPermission(id uint, code string) role.Permission {
	return role.Permission{
		ID:   id,
		Code: code,
	}
}

// newTestRole 创建测试用角色。
func newTestRole(id uint, name string, permissions ...role.Permission) role.Role {
	return role.Role{
		ID:          id,
		Name:        name,
		DisplayName: name,
		Permissions: permissions,
	}
}

// newTestUser 创建测试用用户。
func newTestUser(roles ...role.Role) *User {
	return &User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Status:   "active",
		Roles:    roles,
	}
}

func TestUser_HasRole(t *testing.T) {
	tests := []struct {
		name     string
		user     *User
		roleName string
		want     bool
	}{
		{
			name:     "用户拥有该角色",
			user:     newTestUser(newTestRole(1, "admin"), newTestRole(2, "editor")),
			roleName: "admin",
			want:     true,
		},
		{
			name:     "用户没有该角色",
			user:     newTestUser(newTestRole(1, "editor")),
			roleName: "admin",
			want:     false,
		},
		{
			name:     "用户没有任何角色",
			user:     newTestUser(),
			roleName: "admin",
			want:     false,
		},
		{
			name:     "角色名为空字符串",
			user:     newTestUser(newTestRole(1, "admin")),
			roleName: "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.user.HasRole(tt.roleName)
			assert.Equal(t, tt.want, got, "User.HasRole(%q)", tt.roleName)
		})
	}
}

func TestUser_HasAnyRole(t *testing.T) {
	tests := []struct {
		name      string
		user      *User
		roleNames []string
		want      bool
	}{
		{
			name:      "用户拥有其中一个角色",
			user:      newTestUser(newTestRole(1, "editor")),
			roleNames: []string{"admin", "editor", "viewer"},
			want:      true,
		},
		{
			name:      "用户拥有多个匹配角色",
			user:      newTestUser(newTestRole(1, "admin"), newTestRole(2, "editor")),
			roleNames: []string{"admin", "editor"},
			want:      true,
		},
		{
			name:      "用户没有任何匹配角色",
			user:      newTestUser(newTestRole(1, "viewer")),
			roleNames: []string{"admin", "editor"},
			want:      false,
		},
		{
			name:      "用户没有角色",
			user:      newTestUser(),
			roleNames: []string{"admin", "editor"},
			want:      false,
		},
		{
			name:      "空角色列表",
			user:      newTestUser(newTestRole(1, "admin")),
			roleNames: []string{},
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.user.HasAnyRole(tt.roleNames...)
			assert.Equal(t, tt.want, got, "User.HasAnyRole(%v)", tt.roleNames)
		})
	}
}

func TestUser_HasPermission(t *testing.T) {
	// 创建带权限的角色
	adminRole := newTestRole(1, "admin",
		newTestPermission(1, "user:read"),
		newTestPermission(2, "user:write"),
		newTestPermission(3, "user:delete"),
	)
	editorRole := newTestRole(2, "editor",
		newTestPermission(4, "post:read"),
		newTestPermission(5, "post:write"),
	)

	tests := []struct {
		name           string
		user           *User
		permissionCode string
		want           bool
	}{
		{
			name:           "用户拥有该权限",
			user:           newTestUser(adminRole),
			permissionCode: "user:read",
			want:           true,
		},
		{
			name:           "用户通过多个角色拥有权限",
			user:           newTestUser(adminRole, editorRole),
			permissionCode: "post:write",
			want:           true,
		},
		{
			name:           "用户没有该权限",
			user:           newTestUser(editorRole),
			permissionCode: "user:delete",
			want:           false,
		},
		{
			name:           "用户没有任何角色",
			user:           newTestUser(),
			permissionCode: "user:read",
			want:           false,
		},
		{
			name:           "角色没有任何权限",
			user:           newTestUser(newTestRole(1, "empty")),
			permissionCode: "user:read",
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.user.HasPermission(tt.permissionCode)
			assert.Equal(t, tt.want, got, "User.HasPermission(%q)", tt.permissionCode)
		})
	}
}

func TestUser_GetRoleNames(t *testing.T) {
	tests := []struct {
		name string
		user *User
		want []string
	}{
		{
			name: "用户有多个角色",
			user: newTestUser(newTestRole(1, "admin"), newTestRole(2, "editor")),
			want: []string{"admin", "editor"},
		},
		{
			name: "用户有单个角色",
			user: newTestUser(newTestRole(1, "viewer")),
			want: []string{"viewer"},
		},
		{
			name: "用户没有角色",
			user: newTestUser(),
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.user.GetRoleNames()
			assert.Equal(t, tt.want, got, "User.GetRoleNames()")
		})
	}
}

func TestUser_GetPermissions(t *testing.T) {
	// 创建共享权限（测试去重）
	sharedPerm := newTestPermission(1, "shared:read")
	adminRole := newTestRole(1, "admin", sharedPerm, newTestPermission(2, "user:write"))
	editorRole := newTestRole(2, "editor", sharedPerm, newTestPermission(3, "post:write"))

	tests := []struct {
		name      string
		user      *User
		wantCount int
	}{
		{
			name:      "权限去重 - 多个角色共享相同权限",
			user:      newTestUser(adminRole, editorRole),
			wantCount: 3, // shared:read, user:write, post:write
		},
		{
			name:      "单个角色的权限",
			user:      newTestUser(adminRole),
			wantCount: 2,
		},
		{
			name:      "无角色无权限",
			user:      newTestUser(),
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.user.GetPermissions()
			assert.Len(t, got, tt.wantCount, "User.GetPermissions() count")
		})
	}
}

func TestUser_GetPermissionCodes(t *testing.T) {
	adminRole := newTestRole(1, "admin",
		newTestPermission(1, "user:read"),
		newTestPermission(2, "user:write"),
	)

	tests := []struct {
		name      string
		user      *User
		wantCount int
	}{
		{
			name:      "获取权限代码列表",
			user:      newTestUser(adminRole),
			wantCount: 2,
		},
		{
			name:      "无权限返回空列表",
			user:      newTestUser(),
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.user.GetPermissionCodes()
			assert.Len(t, got, tt.wantCount, "User.GetPermissionCodes() count")
		})
	}
}

func TestUser_IsAdmin(t *testing.T) {
	tests := []struct {
		name string
		user *User
		want bool
	}{
		{
			name: "用户是管理员",
			user: newTestUser(newTestRole(1, "admin")),
			want: true,
		},
		{
			name: "用户不是管理员",
			user: newTestUser(newTestRole(1, "editor")),
			want: false,
		},
		{
			name: "用户有admin角色但还有其他角色",
			user: newTestUser(newTestRole(1, "admin"), newTestRole(2, "editor")),
			want: true,
		},
		{
			name: "用户没有角色",
			user: newTestUser(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.user.IsAdmin()
			assert.Equal(t, tt.want, got, "User.IsAdmin()")
		})
	}
}

func TestUser_StatusChecks(t *testing.T) {
	tests := []struct {
		name       string
		status     string
		canLogin   bool
		isBanned   bool
		isInactive bool
	}{
		{
			name:       "活跃状态",
			status:     "active",
			canLogin:   true,
			isBanned:   false,
			isInactive: false,
		},
		{
			name:       "禁用状态",
			status:     "banned",
			canLogin:   false,
			isBanned:   true,
			isInactive: false,
		},
		{
			name:       "未激活状态",
			status:     "inactive",
			canLogin:   false,
			isBanned:   false,
			isInactive: true,
		},
		{
			name:       "未知状态",
			status:     "unknown",
			canLogin:   false,
			isBanned:   false,
			isInactive: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Status: tt.status}

			assert.Equal(t, tt.canLogin, user.CanLogin(), "User.CanLogin()")
			assert.Equal(t, tt.isBanned, user.IsBanned(), "User.IsBanned()")
			assert.Equal(t, tt.isInactive, user.IsInactive(), "User.IsInactive()")
		})
	}
}

func TestUser_StateTransitions(t *testing.T) {
	t.Run("Activate - 激活用户", func(t *testing.T) {
		user := &User{Status: "inactive"}
		user.Activate()
		assert.Equal(t, "active", user.Status, "User.Activate() status")
	})

	t.Run("Deactivate - 停用用户", func(t *testing.T) {
		user := &User{Status: "active"}
		user.Deactivate()
		assert.Equal(t, "inactive", user.Status, "User.Deactivate() status")
	})

	t.Run("Ban - 禁用用户", func(t *testing.T) {
		user := &User{Status: "active"}
		user.Ban()
		assert.Equal(t, "banned", user.Status, "User.Ban() status")
	})

	t.Run("状态转换链", func(t *testing.T) {
		user := &User{Status: "inactive"}

		// inactive -> active
		user.Activate()
		assert.True(t, user.CanLogin(), "激活后应该可以登录")

		// active -> banned
		user.Ban()
		assert.True(t, user.IsBanned(), "禁用后应该处于禁用状态")

		// banned -> active
		user.Activate()
		assert.True(t, user.CanLogin(), "再次激活后应该可以登录")

		// active -> inactive
		user.Deactivate()
		assert.True(t, user.IsInactive(), "停用后应该处于未激活状态")
	})
}

func TestUser_AssignRole(t *testing.T) {
	t.Run("成功分配新角色", func(t *testing.T) {
		user := newTestUser()
		newRole := newTestRole(1, "editor")

		err := user.AssignRole(newRole)

		require.NoError(t, err, "User.AssignRole() 应该成功")
		assert.True(t, user.HasRole("editor"), "分配后应该拥有该角色")
		assert.Len(t, user.Roles, 1, "角色数量应该是 1")
	})

	t.Run("分配已有角色返回错误", func(t *testing.T) {
		user := newTestUser(newTestRole(1, "admin"))

		err := user.AssignRole(newTestRole(2, "admin"))

		assert.ErrorIs(t, err, ErrRoleAlreadyAssigned, "User.AssignRole() 应该返回 ErrRoleAlreadyAssigned")
	})

	t.Run("分配多个不同角色", func(t *testing.T) {
		user := newTestUser()

		_ = user.AssignRole(newTestRole(1, "admin"))
		_ = user.AssignRole(newTestRole(2, "editor"))
		_ = user.AssignRole(newTestRole(3, "viewer"))

		assert.Len(t, user.Roles, 3, "角色数量应该是 3")
	})
}

func TestUser_RemoveRole(t *testing.T) {
	t.Run("成功移除已有角色", func(t *testing.T) {
		user := newTestUser(newTestRole(1, "admin"), newTestRole(2, "editor"))

		err := user.RemoveRole("admin")

		require.NoError(t, err, "User.RemoveRole() 应该成功")
		assert.False(t, user.HasRole("admin"), "移除后不应该拥有该角色")
		assert.Len(t, user.Roles, 1, "角色数量应该是 1")
	})

	t.Run("移除不存在的角色返回错误", func(t *testing.T) {
		user := newTestUser(newTestRole(1, "admin"))

		err := user.RemoveRole("editor")

		assert.ErrorIs(t, err, ErrRoleNotFound, "User.RemoveRole() 应该返回 ErrRoleNotFound")
	})

	t.Run("从空角色列表移除返回错误", func(t *testing.T) {
		user := newTestUser()

		err := user.RemoveRole("admin")

		assert.ErrorIs(t, err, ErrRoleNotFound, "User.RemoveRole() 应该返回 ErrRoleNotFound")
	})

	t.Run("移除中间的角色", func(t *testing.T) {
		user := newTestUser(
			newTestRole(1, "admin"),
			newTestRole(2, "editor"),
			newTestRole(3, "viewer"),
		)

		err := user.RemoveRole("editor")

		require.NoError(t, err, "User.RemoveRole() 应该成功")
		assert.True(t, user.HasRole("admin"), "admin 角色不应该被影响")
		assert.True(t, user.HasRole("viewer"), "viewer 角色不应该被影响")
		assert.False(t, user.HasRole("editor"), "被移除的角色不应该存在")
	})
}

func TestUser_ClearRoles(t *testing.T) {
	t.Run("清空所有角色", func(t *testing.T) {
		user := newTestUser(
			newTestRole(1, "admin"),
			newTestRole(2, "editor"),
			newTestRole(3, "viewer"),
		)

		user.ClearRoles()

		assert.Empty(t, user.Roles, "ClearRoles() 后角色应该为空")
	})

	t.Run("清空空角色列表", func(t *testing.T) {
		user := newTestUser()

		user.ClearRoles()

		assert.Empty(t, user.Roles, "ClearRoles() 后角色应该为空")
	})
}

func TestUser_UpdateProfile(t *testing.T) {
	t.Run("更新所有字段", func(t *testing.T) {
		user := &User{
			FullName: "Old Name",
			Avatar:   "old-avatar.png",
			Bio:      "Old bio",
		}

		user.UpdateProfile("New Name", "new-avatar.png", "New bio")

		assert.Equal(t, "New Name", user.FullName, "FullName 应该被更新")
		assert.Equal(t, "new-avatar.png", user.Avatar, "Avatar 应该被更新")
		assert.Equal(t, "New bio", user.Bio, "Bio 应该被更新")
	})

	t.Run("只更新FullName", func(t *testing.T) {
		user := &User{
			FullName: "Old Name",
			Avatar:   "old-avatar.png",
			Bio:      "Old bio",
		}

		user.UpdateProfile("New Name", "", "")

		assert.Equal(t, "New Name", user.FullName, "FullName 应该被更新")
		assert.Equal(t, "old-avatar.png", user.Avatar, "Avatar 不应该被清空")
		assert.Empty(t, user.Bio, "Bio 可以被清空")
	})

	t.Run("Bio可以设置为空", func(t *testing.T) {
		user := &User{Bio: "Some bio"}

		user.UpdateProfile("", "", "")

		assert.Empty(t, user.Bio, "Bio 应该被清空")
	})

	t.Run("空字符串不覆盖FullName和Avatar", func(t *testing.T) {
		user := &User{
			FullName: "Keep This",
			Avatar:   "keep-this.png",
		}

		user.UpdateProfile("", "", "new bio")

		assert.Equal(t, "Keep This", user.FullName, "FullName 不应该被清空")
		assert.Equal(t, "keep-this.png", user.Avatar, "Avatar 不应该被清空")
	})
}

func TestUser_EdgeCases(t *testing.T) {
	t.Run("nil Roles slice", func(t *testing.T) {
		user := &User{Roles: nil}

		// 这些操作不应该 panic
		assert.NotPanics(t, func() {
			_ = user.HasRole("admin")
			_ = user.HasAnyRole("admin", "editor")
			_ = user.HasPermission("user:read")
			_ = user.GetRoleNames()
			_ = user.GetPermissions()
			_ = user.GetPermissionCodes()
			_ = user.IsAdmin()
		}, "nil Roles slice 不应该导致 panic")
	})

	t.Run("角色中 nil Permissions slice", func(t *testing.T) {
		roleWithNilPerms := role.Role{
			ID:          1,
			Name:        "test",
			Permissions: nil,
		}
		user := &User{Roles: []role.Role{roleWithNilPerms}}

		// 不应该 panic
		assert.NotPanics(t, func() {
			_ = user.HasPermission("test:read")
			_ = user.GetPermissions()
			_ = user.GetPermissionCodes()
		}, "nil Permissions slice 不应该导致 panic")
	})
}

func BenchmarkUser_HasRole(b *testing.B) {
	// 创建有10个角色的用户
	roles := make([]role.Role, 10)
	for i := range 10 {
		roles[i] = newTestRole(uint(i+1), "role"+string(rune('0'+i))) //nolint:gosec // benchmark with small indices
	}
	user := &User{Roles: roles}

	for b.Loop() {
		user.HasRole("role5")
	}
}

func BenchmarkUser_HasPermission(b *testing.B) {
	// 创建有10个角色，每个角色10个权限的用户
	roles := make([]role.Role, 10)
	for i := range 10 {
		perms := make([]role.Permission, 10)
		for j := range 10 {
			perms[j] = newTestPermission(uint(i*10+j+1), "perm:"+string(rune('0'+i))+":"+string(rune('0'+j))) //nolint:gosec // benchmark with small indices
		}
		roles[i] = newTestRole(uint(i+1), "role"+string(rune('0'+i)), perms...) //nolint:gosec // benchmark with small indices
	}
	user := &User{Roles: roles}

	for b.Loop() {
		user.HasPermission("perm:5:5")
	}
}

func BenchmarkUser_GetPermissions_Deduplication(b *testing.B) {
	// 创建有重叠权限的角色
	sharedPerm := newTestPermission(1, "shared:read")
	roles := make([]role.Role, 5)
	for i := range 5 {
		perms := []role.Permission{
			sharedPerm, // 共享权限
			newTestPermission(uint(i*10+2), "unique:"+string(rune('0'+i))), //nolint:gosec // benchmark with small indices
		}
		roles[i] = newTestRole(uint(i+1), "role"+string(rune('0'+i)), perms...) //nolint:gosec // benchmark with small indices
	}
	user := &User{Roles: roles}

	for b.Loop() {
		_ = user.GetPermissions()
	}
}
