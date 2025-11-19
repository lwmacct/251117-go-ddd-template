package auditlog

import "context"

// Repository defines the interface for audit log data operations
type Repository interface {
	// Create creates a new audit log entry
	Create(ctx context.Context, log *AuditLog) error

	// FindByID finds an audit log by ID
	FindByID(ctx context.Context, id uint) (*AuditLog, error)

	// List returns audit logs with filtering and pagination
	List(ctx context.Context, filter FilterOptions) ([]AuditLog, int64, error)

	// ListByUser returns audit logs for a specific user
	ListByUser(ctx context.Context, userID uint, page, limit int) ([]AuditLog, int64, error)

	// ListByResource returns audit logs for a specific resource
	ListByResource(ctx context.Context, resource string, page, limit int) ([]AuditLog, int64, error)

	// Delete deletes an audit log (soft delete, for data retention policy)
	Delete(ctx context.Context, id uint) error

	// DeleteOlderThan deletes audit logs older than the specified date
	DeleteOlderThan(ctx context.Context, days int) error
}
