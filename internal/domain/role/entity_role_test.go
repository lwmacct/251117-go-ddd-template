// Package role 提供角色和权限领域模型单元测试。
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
				ID:     1,
				Domain: "user",
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
