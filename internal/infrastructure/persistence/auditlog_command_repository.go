// Package persistence 提供审计日志命令仓储的 GORM 实现
package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
	"gorm.io/gorm"
)

// auditLogCommandRepository 审计日志命令仓储的 GORM 实现
type auditLogCommandRepository struct {
	db *gorm.DB
}

// NewAuditLogCommandRepository 创建审计日志命令仓储实例
func NewAuditLogCommandRepository(db *gorm.DB) auditlog.CommandRepository {
	return &auditLogCommandRepository{db: db}
}

// Create creates a new audit log entry
func (r *auditLogCommandRepository) Create(ctx context.Context, log *auditlog.AuditLog) error {
	if err := r.db.WithContext(ctx).Create(log).Error; err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}
	return nil
}

// Delete deletes an audit log (soft delete, for data retention policy)
func (r *auditLogCommandRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&auditlog.AuditLog{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete audit log: %w", err)
	}
	return nil
}

// DeleteOlderThan deletes audit logs older than the specified date
func (r *auditLogCommandRepository) DeleteOlderThan(ctx context.Context, days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	if err := r.db.WithContext(ctx).
		Where("created_at < ?", cutoffDate).
		Delete(&auditlog.AuditLog{}).Error; err != nil {
		return fmt.Errorf("failed to delete old audit logs: %w", err)
	}
	return nil
}

// BatchCreate creates multiple audit log entries
func (r *auditLogCommandRepository) BatchCreate(ctx context.Context, logs []*auditlog.AuditLog) error {
	if len(logs) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).Create(logs).Error; err != nil {
		return fmt.Errorf("failed to batch create audit logs: %w", err)
	}
	return nil
}
