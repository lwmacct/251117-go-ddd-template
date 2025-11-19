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
	model := newAuditLogModelFromEntity(log)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}
	if saved := model.toEntity(); saved != nil {
		*log = *saved
	}
	return nil
}

// Delete deletes an audit log (soft delete, for data retention policy)
func (r *auditLogCommandRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&AuditLogModel{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete audit log: %w", err)
	}
	return nil
}

// DeleteOlderThan deletes audit logs older than the specified date
func (r *auditLogCommandRepository) DeleteOlderThan(ctx context.Context, days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	if err := r.db.WithContext(ctx).
		Where("created_at < ?", cutoffDate).
		Delete(&AuditLogModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete old audit logs: %w", err)
	}
	return nil
}

// BatchCreate creates multiple audit log entries
func (r *auditLogCommandRepository) BatchCreate(ctx context.Context, logs []*auditlog.AuditLog) error {
	if len(logs) == 0 {
		return nil
	}
	models := make([]*AuditLogModel, 0, len(logs))
	for _, log := range logs {
		if model := newAuditLogModelFromEntity(log); model != nil {
			models = append(models, model)
		}
	}
	if err := r.db.WithContext(ctx).Create(models).Error; err != nil {
		return fmt.Errorf("failed to batch create audit logs: %w", err)
	}

	for i := range models {
		if entity := models[i].toEntity(); entity != nil {
			*logs[i] = *entity
		}
	}

	return nil
}
