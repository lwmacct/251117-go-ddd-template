package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/role"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
)

func TestToUserDTO(t *testing.T) {
	t.Run("转换正常用户", func(t *testing.T) {
		now := time.Now()
		u := &user.User{
			ID:        1,
			Username:  "testuser",
			Email:     "test@example.com",
			FullName:  "Test User",
			Avatar:    "https://example.com/avatar.png",
			Bio:       "Test bio",
			Status:    "active",
			CreatedAt: now,
			UpdatedAt: now,
		}

		result := ToUserDTO(u)

		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "testuser", result.Username)
		assert.Equal(t, "test@example.com", result.Email)
		assert.Equal(t, "Test User", result.FullName)
		assert.Equal(t, "https://example.com/avatar.png", result.Avatar)
		assert.Equal(t, "Test bio", result.Bio)
		assert.Equal(t, "active", result.Status)
		assert.Equal(t, now, result.CreatedAt)
		assert.Equal(t, now, result.UpdatedAt)
	})

	t.Run("转换nil返回nil", func(t *testing.T) {
		result := ToUserDTO(nil)
		assert.Nil(t, result)
	})

	t.Run("转换空字段用户", func(t *testing.T) {
		u := &user.User{
			ID:       2,
			Username: "emptyuser",
			Email:    "empty@example.com",
		}

		result := ToUserDTO(u)

		assert.NotNil(t, result)
		assert.Equal(t, uint(2), result.ID)
		assert.Empty(t, result.FullName)
		assert.Empty(t, result.Avatar)
		assert.Empty(t, result.Bio)
	})
}

func TestToUserWithRolesDTO(t *testing.T) {
	t.Run("转换带角色的用户", func(t *testing.T) {
		now := time.Now()
		u := &user.User{
			ID:        1,
			Username:  "admin",
			Email:     "admin@example.com",
			FullName:  "Admin User",
			Avatar:    "https://example.com/admin.png",
			Bio:       "Admin bio",
			Status:    "active",
			CreatedAt: now,
			UpdatedAt: now,
			Roles: []role.Role{
				{
					ID:          1,
					Name:        "admin",
					DisplayName: "管理员",
					Description: "系统管理员",
				},
				{
					ID:          2,
					Name:        "user",
					DisplayName: "普通用户",
					Description: "普通用户角色",
				},
			},
		}

		result := ToUserWithRolesDTO(u)

		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "admin", result.Username)
		assert.Equal(t, "admin@example.com", result.Email)
		assert.Len(t, result.Roles, 2)
		assert.Equal(t, uint(1), result.Roles[0].ID)
		assert.Equal(t, "admin", result.Roles[0].Name)
		assert.Equal(t, "管理员", result.Roles[0].DisplayName)
		assert.Equal(t, uint(2), result.Roles[1].ID)
		assert.Equal(t, "user", result.Roles[1].Name)
	})

	t.Run("转换nil返回nil", func(t *testing.T) {
		result := ToUserWithRolesDTO(nil)
		assert.Nil(t, result)
	})

	t.Run("转换无角色用户", func(t *testing.T) {
		u := &user.User{
			ID:       1,
			Username: "noroles",
			Email:    "noroles@example.com",
			Status:   "active",
			Roles:    nil,
		}

		result := ToUserWithRolesDTO(u)

		assert.NotNil(t, result)
		assert.Equal(t, "noroles", result.Username)
		assert.Nil(t, result.Roles)
	})

	t.Run("转换空角色数组用户", func(t *testing.T) {
		u := &user.User{
			ID:       1,
			Username: "emptyroles",
			Email:    "empty@example.com",
			Status:   "active",
			Roles:    []role.Role{},
		}

		result := ToUserWithRolesDTO(u)

		assert.NotNil(t, result)
		assert.Nil(t, result.Roles)
	})
}
