package auditlog

import "context"

// QueryRepository 审计日志查询仓储接口（读操作）
type QueryRepository interface {
	// FindByID finds an audit log by ID
	FindByID(ctx context.Context, id uint) (*AuditLog, error)

	// List returns audit logs with filtering and pagination
	List(ctx context.Context, filter FilterOptions) ([]AuditLog, int64, error)

	// ListByUser returns audit logs for a specific user
	ListByUser(ctx context.Context, userID uint, page, limit int) ([]AuditLog, int64, error)

	// ListByResource returns audit logs for a specific resource
	ListByResource(ctx context.Context, resource string, page, limit int) ([]AuditLog, int64, error)

	// ListByAction returns audit logs for a specific action
	ListByAction(ctx context.Context, action string, page, limit int) ([]AuditLog, int64, error)

	// Count returns the total number of audit logs matching the filter
	Count(ctx context.Context, filter FilterOptions) (int64, error)

	// Search searches audit logs by keyword
	Search(ctx context.Context, keyword string, page, limit int) ([]AuditLog, int64, error)
}
