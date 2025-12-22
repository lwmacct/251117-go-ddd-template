package role

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPermission_PermissionCode(t *testing.T) {
	tests := []struct {
		name string
		perm Permission
		want string
	}{
		{
			name: "标准三段式权限码",
			perm: Permission{
				ID:       1,
				Domain:   "user",
				Resource: "profile",
				Action:   "read",
				Code:     "user:profile:read",
			},
			want: "user:profile:read",
		},
		{
			name: "简单权限码",
			perm: Permission{
				ID:   2,
				Code: "admin",
			},
			want: "admin",
		},
		{
			name: "空权限码",
			perm: Permission{
				ID:   3,
				Code: "",
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.perm.PermissionCode()
			assert.Equal(t, tt.want, got, "Permission.PermissionCode()")
		})
	}
}

func TestRole_BasicFields(t *testing.T) {
	role := Role{
		ID:          1,
		Name:        "admin",
		DisplayName: "Administrator",
		Description: "System administrator with full access",
		IsSystem:    true,
		Permissions: []Permission{
			{ID: 1, Code: "user:read"},
			{ID: 2, Code: "user:write"},
		},
	}

	assert.Equal(t, uint(1), role.ID)
	assert.Equal(t, "admin", role.Name)
	assert.Equal(t, "Administrator", role.DisplayName)
	assert.Equal(t, "System administrator with full access", role.Description)
	assert.True(t, role.IsSystem)
	assert.Len(t, role.Permissions, 2)
}

func TestRole_EmptyPermissions(t *testing.T) {
	role := Role{
		Permissions: nil,
	}

	assert.Nil(t, role.Permissions)
	assert.Empty(t, role.Permissions)
}

func TestRole_SystemVsCustom(t *testing.T) {
	tests := []struct {
		name     string
		isSystem bool
	}{
		{
			name:     "系统角色",
			isSystem: true,
		},
		{
			name:     "自定义角色",
			isSystem: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role := Role{
				IsSystem: tt.isSystem,
			}
			assert.Equal(t, tt.isSystem, role.IsSystem)
		})
	}
}

func TestRole_IsSystemRole(t *testing.T) {
	t.Run("系统角色返回 true", func(t *testing.T) {
		role := Role{IsSystem: true}
		assert.True(t, role.IsSystemRole())
	})

	t.Run("非系统角色返回 false", func(t *testing.T) {
		role := Role{IsSystem: false}
		assert.False(t, role.IsSystemRole())
	})
}

func TestRole_CanBeDeleted(t *testing.T) {
	t.Run("系统角色不可删除", func(t *testing.T) {
		role := Role{IsSystem: true}
		assert.False(t, role.CanBeDeleted())
	})

	t.Run("非系统角色可删除", func(t *testing.T) {
		role := Role{IsSystem: false}
		assert.True(t, role.CanBeDeleted())
	})
}

func TestRole_HasPermission(t *testing.T) {
	role := Role{
		Permissions: []Permission{
			{Code: "user:read"},
			{Code: "user:write"},
			{Code: "role:read"},
		},
	}

	t.Run("拥有权限", func(t *testing.T) {
		assert.True(t, role.HasPermission("user:read"))
		assert.True(t, role.HasPermission("role:read"))
	})

	t.Run("没有权限", func(t *testing.T) {
		assert.False(t, role.HasPermission("user:delete"))
		assert.False(t, role.HasPermission("admin:*"))
	})

	t.Run("空权限列表", func(t *testing.T) {
		emptyRole := Role{}
		assert.False(t, emptyRole.HasPermission("any:perm"))
	})
}

func TestRole_HasAnyPermission(t *testing.T) {
	role := Role{
		Permissions: []Permission{
			{Code: "user:read"},
			{Code: "user:write"},
		},
	}

	t.Run("拥有其中一个权限", func(t *testing.T) {
		assert.True(t, role.HasAnyPermission("user:read", "admin:read"))
	})

	t.Run("没有任何权限", func(t *testing.T) {
		assert.False(t, role.HasAnyPermission("admin:read", "admin:write"))
	})

	t.Run("空参数", func(t *testing.T) {
		assert.False(t, role.HasAnyPermission())
	})
}

func TestRole_GetPermissionCodes(t *testing.T) {
	role := Role{
		Permissions: []Permission{
			{Code: "user:read"},
			{Code: "user:write"},
		},
	}

	codes := role.GetPermissionCodes()
	assert.Equal(t, []string{"user:read", "user:write"}, codes)
}

func TestRole_PermissionOperations(t *testing.T) {
	t.Run("AddPermission", func(t *testing.T) {
		role := Role{}
		role.AddPermission(Permission{Code: "user:read"})
		assert.Equal(t, 1, role.GetPermissionCount())

		// 重复添加不应增加
		role.AddPermission(Permission{Code: "user:read"})
		assert.Equal(t, 1, role.GetPermissionCount())
	})

	t.Run("RemovePermission", func(t *testing.T) {
		role := Role{
			Permissions: []Permission{
				{Code: "user:read"},
				{Code: "user:write"},
			},
		}
		assert.True(t, role.RemovePermission("user:read"))
		assert.Equal(t, 1, role.GetPermissionCount())
		assert.False(t, role.HasPermission("user:read"))

		// 移除不存在的权限
		assert.False(t, role.RemovePermission("nonexistent"))
	})

	t.Run("ClearPermissions", func(t *testing.T) {
		role := Role{
			Permissions: []Permission{
				{Code: "user:read"},
				{Code: "user:write"},
			},
		}
		role.ClearPermissions()
		assert.True(t, role.IsEmpty())
	})
}

func TestPermission_IsValid(t *testing.T) {
	tests := []struct {
		name string
		perm Permission
		want bool
	}{
		{
			name: "有效权限",
			perm: Permission{Domain: "user", Resource: "profile", Action: "read", Code: "user:profile:read"},
			want: true,
		},
		{
			name: "缺少 Domain",
			perm: Permission{Resource: "profile", Action: "read", Code: "user:profile:read"},
			want: false,
		},
		{
			name: "缺少 Code",
			perm: Permission{Domain: "user", Resource: "profile", Action: "read"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.perm.IsValid())
		})
	}
}

func TestPermission_Matches(t *testing.T) {
	perm := Permission{Domain: "user", Resource: "profile", Action: "read", Code: "user:profile:read"}

	tests := []struct {
		name    string
		pattern string
		want    bool
	}{
		{"精确匹配", "user:profile:read", true},
		{"全通配符", "*", true},
		{"三段通配符", "*:*:*", true},
		{"域通配符", "user:*:*", true},
		{"资源通配符", "user:profile:*", true},
		{"动作通配符", "*:*:read", true},
		{"不匹配的域", "admin:profile:read", false},
		{"不匹配的资源", "user:settings:read", false},
		{"不匹配的动作", "user:profile:write", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, perm.Matches(tt.pattern))
		})
	}
}

func TestPermission_GetComponents(t *testing.T) {
	perm := Permission{Domain: "user", Resource: "profile", Action: "read"}
	domain, resource, action := perm.GetComponents()

	assert.Equal(t, "user", domain)
	assert.Equal(t, "profile", resource)
	assert.Equal(t, "read", action)
}

func TestPermission_BuildCode(t *testing.T) {
	perm := Permission{Domain: "user", Resource: "profile", Action: "read"}
	assert.Equal(t, "user:profile:read", perm.BuildCode())
}
