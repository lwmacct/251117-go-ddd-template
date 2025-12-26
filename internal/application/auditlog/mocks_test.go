//nolint:forcetypeassert // Mock 返回值类型在测试中总是已知的
package auditlog

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	domainAuditLog "github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
)

// ============================================================
// MockAuditLogCommandRepository
// ============================================================

type MockAuditLogCommandRepository struct {
	mock.Mock
}

func (m *MockAuditLogCommandRepository) Create(ctx context.Context, log *domainAuditLog.AuditLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockAuditLogCommandRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAuditLogCommandRepository) DeleteOlderThan(ctx context.Context, days int) error {
	args := m.Called(ctx, days)
	return args.Error(0)
}

func (m *MockAuditLogCommandRepository) BatchCreate(ctx context.Context, logs []*domainAuditLog.AuditLog) error {
	args := m.Called(ctx, logs)
	return args.Error(0)
}

// ============================================================
// MockAuditLogQueryRepository
// ============================================================

type MockAuditLogQueryRepository struct {
	mock.Mock
}

func (m *MockAuditLogQueryRepository) FindByID(ctx context.Context, id uint) (*domainAuditLog.AuditLog, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainAuditLog.AuditLog), args.Error(1)
}

func (m *MockAuditLogQueryRepository) List(ctx context.Context, filter domainAuditLog.FilterOptions) ([]domainAuditLog.AuditLog, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domainAuditLog.AuditLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockAuditLogQueryRepository) ListByUser(ctx context.Context, userID uint, page, limit int) ([]domainAuditLog.AuditLog, int64, error) {
	args := m.Called(ctx, userID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domainAuditLog.AuditLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockAuditLogQueryRepository) ListByResource(ctx context.Context, resource string, page, limit int) ([]domainAuditLog.AuditLog, int64, error) {
	args := m.Called(ctx, resource, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domainAuditLog.AuditLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockAuditLogQueryRepository) ListByAction(ctx context.Context, action string, page, limit int) ([]domainAuditLog.AuditLog, int64, error) {
	args := m.Called(ctx, action, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domainAuditLog.AuditLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockAuditLogQueryRepository) Count(ctx context.Context, filter domainAuditLog.FilterOptions) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAuditLogQueryRepository) Search(ctx context.Context, keyword string, page, limit int) ([]domainAuditLog.AuditLog, int64, error) {
	args := m.Called(ctx, keyword, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domainAuditLog.AuditLog), args.Get(1).(int64), args.Error(2)
}

// ============================================================
// 测试辅助函数
// ============================================================

func newTestAuditLog(id uint) *domainAuditLog.AuditLog {
	return &domainAuditLog.AuditLog{
		ID:        id,
		UserID:    1,
		Username:  "testuser",
		Action:    "login",
		Resource:  "auth",
		Status:    "success",
		IPAddress: "127.0.0.1",
		UserAgent: "TestAgent/1.0",
		Details:   `{"event":"login_success"}`,
		CreatedAt: time.Now(),
	}
}
