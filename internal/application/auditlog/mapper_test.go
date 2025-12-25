package auditlog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
)

func TestToAuditLogDTO(t *testing.T) {
	t.Run("转换正常审计日志", func(t *testing.T) {
		now := time.Now()
		log := &auditlog.AuditLog{
			ID:        1,
			UserID:    100,
			Action:    auditlog.ActionLogin,
			Resource:  "auth",
			Details:   "用户登录成功",
			IPAddress: "192.168.1.1",
			UserAgent: "Mozilla/5.0",
			Status:    auditlog.StatusSuccess,
			CreatedAt: now,
		}

		result := ToAuditLogDTO(log)

		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, uint(100), result.UserID)
		assert.Equal(t, auditlog.ActionLogin, result.Action)
		assert.Equal(t, "auth", result.Resource)
		assert.Equal(t, "用户登录成功", result.Details)
		assert.Equal(t, "192.168.1.1", result.IPAddress)
		assert.Equal(t, "Mozilla/5.0", result.UserAgent)
		assert.Equal(t, auditlog.StatusSuccess, result.Status)
		assert.Equal(t, now, result.CreatedAt)
	})

	t.Run("转换nil返回nil", func(t *testing.T) {
		result := ToAuditLogDTO(nil)
		assert.Nil(t, result)
	})

	t.Run("转换失败状态日志", func(t *testing.T) {
		log := &auditlog.AuditLog{
			ID:       2,
			UserID:   101,
			Action:   auditlog.ActionLogin,
			Resource: "auth",
			Details:  "密码错误",
			Status:   auditlog.StatusFailed,
		}

		result := ToAuditLogDTO(log)

		assert.NotNil(t, result)
		assert.Equal(t, auditlog.StatusFailed, result.Status)
		assert.Equal(t, "密码错误", result.Details)
	})

	t.Run("转换空字段日志", func(t *testing.T) {
		log := &auditlog.AuditLog{
			ID:     3,
			UserID: 102,
			Action: auditlog.ActionCreate,
			Status: auditlog.StatusSuccess,
		}

		result := ToAuditLogDTO(log)

		assert.NotNil(t, result)
		assert.Equal(t, uint(3), result.ID)
		assert.Empty(t, result.Resource)
		assert.Empty(t, result.Details)
		assert.Empty(t, result.IPAddress)
		assert.Empty(t, result.UserAgent)
	})

	t.Run("转换各种操作类型", func(t *testing.T) {
		actions := []string{
			auditlog.ActionCreate,
			auditlog.ActionUpdate,
			auditlog.ActionDelete,
			auditlog.ActionLogin,
			auditlog.ActionLogout,
		}

		for _, action := range actions {
			log := &auditlog.AuditLog{
				ID:     1,
				Action: action,
			}

			result := ToAuditLogDTO(log)
			assert.NotNil(t, result)
			assert.Equal(t, action, result.Action)
		}
	})
}
