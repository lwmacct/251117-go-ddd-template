package auditlog

import "context"

// CommandRepository 审计日志命令仓储接口（写操作）
type CommandRepository interface {
	// Create creates a new audit log entry
	Create(ctx context.Context, log *AuditLog) error

	// Delete deletes an audit log (soft delete, for data retention policy)
	Delete(ctx context.Context, id uint) error

	// DeleteOlderThan deletes audit logs older than the specified date
	DeleteOlderThan(ctx context.Context, days int) error

	// BatchCreate creates multiple audit log entries
	BatchCreate(ctx context.Context, logs []*AuditLog) error
}
