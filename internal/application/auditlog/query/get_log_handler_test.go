package query

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
)

// mockAuditLogQueryRepoForGet 审计日志查询仓储 Mock
type mockAuditLogQueryRepoForGet struct {
	log     *auditlog.AuditLog
	logs    []auditlog.AuditLog
	total   int64
	findErr error
}

func (m *mockAuditLogQueryRepoForGet) FindByID(ctx context.Context, id uint) (*auditlog.AuditLog, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	return m.log, nil
}

func (m *mockAuditLogQueryRepoForGet) List(ctx context.Context, filter auditlog.FilterOptions) ([]auditlog.AuditLog, int64, error) {
	return m.logs, m.total, m.findErr
}

func (m *mockAuditLogQueryRepoForGet) ListByUser(ctx context.Context, userID uint, page, limit int) ([]auditlog.AuditLog, int64, error) {
	return m.logs, m.total, m.findErr
}

func (m *mockAuditLogQueryRepoForGet) ListByResource(ctx context.Context, resource string, page, limit int) ([]auditlog.AuditLog, int64, error) {
	return m.logs, m.total, m.findErr
}

func (m *mockAuditLogQueryRepoForGet) ListByAction(ctx context.Context, action string, page, limit int) ([]auditlog.AuditLog, int64, error) {
	return m.logs, m.total, m.findErr
}

func (m *mockAuditLogQueryRepoForGet) Count(ctx context.Context, filter auditlog.FilterOptions) (int64, error) {
	return m.total, m.findErr
}

func (m *mockAuditLogQueryRepoForGet) Search(ctx context.Context, keyword string, page, limit int) ([]auditlog.AuditLog, int64, error) {
	return m.logs, m.total, m.findErr
}

func TestGetLogHandler_Handle(t *testing.T) {
	ctx := context.Background()

	t.Run("成功获取审计日志", func(t *testing.T) {
		now := time.Now()
		mockRepo := &mockAuditLogQueryRepoForGet{
			log: &auditlog.AuditLog{
				ID:        1,
				UserID:    100,
				Username:  "testuser",
				Action:    auditlog.ActionLogin,
				Resource:  "auth",
				Status:    auditlog.StatusSuccess,
				IPAddress: "192.168.1.1",
				UserAgent: "Mozilla/5.0",
				Details:   "登录成功",
				CreatedAt: now,
			},
		}

		handler := NewGetLogHandler(mockRepo)
		result, err := handler.Handle(ctx, GetLogQuery{LogID: 1})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, uint(100), result.UserID)
		assert.Equal(t, auditlog.ActionLogin, result.Action)
		assert.Equal(t, auditlog.StatusSuccess, result.Status)
		assert.Equal(t, "192.168.1.1", result.IPAddress)
	})

	t.Run("审计日志不存在", func(t *testing.T) {
		mockRepo := &mockAuditLogQueryRepoForGet{
			log:     nil,
			findErr: errors.New("not found"),
		}

		handler := NewGetLogHandler(mockRepo)
		result, err := handler.Handle(ctx, GetLogQuery{LogID: 999})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit log not found")
		assert.Nil(t, result)
	})

	t.Run("仓储返回nil但无错误", func(t *testing.T) {
		mockRepo := &mockAuditLogQueryRepoForGet{
			log:     nil,
			findErr: nil,
		}

		handler := NewGetLogHandler(mockRepo)
		result, err := handler.Handle(ctx, GetLogQuery{LogID: 1})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit log not found")
		assert.Nil(t, result)
	})

	t.Run("获取带详情的审计日志", func(t *testing.T) {
		mockRepo := &mockAuditLogQueryRepoForGet{
			log: &auditlog.AuditLog{
				ID:         2,
				UserID:     100,
				Username:   "admin",
				Action:     auditlog.ActionUpdate,
				Resource:   "user",
				ResourceID: "123",
				Status:     auditlog.StatusSuccess,
				Details:    `{"old": {"name": "old"}, "new": {"name": "new"}}`,
			},
		}

		handler := NewGetLogHandler(mockRepo)
		result, err := handler.Handle(ctx, GetLogQuery{LogID: 2})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "user", result.Resource)
		assert.Contains(t, result.Details, "old")
	})
}

func TestNewGetLogHandler(t *testing.T) {
	mockRepo := &mockAuditLogQueryRepoForGet{}
	handler := NewGetLogHandler(mockRepo)

	assert.NotNil(t, handler)
}
