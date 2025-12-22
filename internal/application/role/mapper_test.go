package role

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	domainRole "github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
)

func TestToRoleResponse(t *testing.T) {
	t.Run("转换正常角色", func(t *testing.T) {
		now := time.Now()
		r := &domainRole.Role{
			ID:          1,
			Name:        "admin",
			DisplayName: "管理员",
			Description: "系统管理员角色",
			IsSystem:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		result := ToRoleResponse(r)

		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "admin", result.Name)
		assert.Equal(t, "管理员", result.DisplayName)
		assert.Equal(t, "系统管理员角色", result.Description)
		assert.True(t, result.IsSystem)
		assert.Equal(t, now, result.CreatedAt)
		assert.Equal(t, now, result.UpdatedAt)
		assert.Nil(t, result.Permissions)
	})

	t.Run("转换nil返回nil", func(t *testing.T) {
		result := ToRoleResponse(nil)
		assert.Nil(t, result)
	})

	t.Run("转换带权限的角色", func(t *testing.T) {
		now := time.Now()
		r := &domainRole.Role{
			ID:          2,
			Name:        "editor",
			DisplayName: "编辑者",
			Description: "内容编辑角色",
			IsSystem:    false,
			Permissions: []domainRole.Permission{
				{
					ID:          1,
					Code:        "user:read",
					Description: "读取用户",
					Resource:    "user",
					Action:      "read",
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:          2,
					Code:        "user:write",
					Description: "写入用户",
					Resource:    "user",
					Action:      "write",
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			CreatedAt: now,
			UpdatedAt: now,
		}

		result := ToRoleResponse(r)

		assert.NotNil(t, result)
		assert.Equal(t, "editor", result.Name)
		assert.False(t, result.IsSystem)
		assert.Len(t, result.Permissions, 2)
		assert.Equal(t, "user:read", result.Permissions[0].Code)
		assert.Equal(t, "user:write", result.Permissions[1].Code)
	})

	t.Run("转换空权限角色", func(t *testing.T) {
		r := &domainRole.Role{
			ID:          3,
			Name:        "guest",
			DisplayName: "访客",
			Permissions: []domainRole.Permission{},
		}

		result := ToRoleResponse(r)

		assert.NotNil(t, result)
		assert.Equal(t, "guest", result.Name)
		assert.Nil(t, result.Permissions)
	})
}

func TestToPermissionResponse(t *testing.T) {
	t.Run("转换正常权限", func(t *testing.T) {
		now := time.Now()
		p := &domainRole.Permission{
			ID:          1,
			Code:        "user:create",
			Description: "创建用户权限",
			Resource:    "user",
			Action:      "create",
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		result := ToPermissionResponse(p)

		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "user:create", result.Code)
		assert.Equal(t, "user:create", result.Name) // Name 使用 Code
		assert.Equal(t, "创建用户权限", result.Description)
		assert.Equal(t, "user", result.Resource)
		assert.Equal(t, "create", result.Action)
		assert.Equal(t, now, result.CreatedAt)
		assert.Equal(t, now, result.UpdatedAt)
	})

	t.Run("转换nil返回nil", func(t *testing.T) {
		result := ToPermissionResponse(nil)
		assert.Nil(t, result)
	})

	t.Run("转换空字段权限", func(t *testing.T) {
		p := &domainRole.Permission{
			ID:   2,
			Code: "test:action",
		}

		result := ToPermissionResponse(p)

		assert.NotNil(t, result)
		assert.Equal(t, uint(2), result.ID)
		assert.Equal(t, "test:action", result.Code)
		assert.Empty(t, result.Description)
		assert.Empty(t, result.Resource)
		assert.Empty(t, result.Action)
	})
}
