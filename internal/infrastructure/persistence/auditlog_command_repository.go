package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
	"gorm.io/gorm"
)

// auditLogCommandRepository 审计日志命令仓储的 GORM 实现
// 嵌入 GenericCommandRepository 以复用 Create/Delete 操作
type auditLogCommandRepository struct {
	*GenericCommandRepository[auditlog.AuditLog, *AuditLogModel]
}

// NewAuditLogCommandRepository 创建审计日志命令仓储实例
func NewAuditLogCommandRepository(db *gorm.DB) auditlog.CommandRepository {
	return &auditLogCommandRepository{
		GenericCommandRepository: NewGenericCommandRepository(
			db, newAuditLogModelFromEntity,
		),
	}
}

// Create、Delete 方法由 GenericCommandRepository 提供

// DeleteOlderThan deletes audit logs older than the specified date
func (r *auditLogCommandRepository) DeleteOlderThan(ctx context.Context, days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	if err := r.DB().WithContext(ctx).
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
	if err := r.DB().WithContext(ctx).Create(models).Error; err != nil {
		return fmt.Errorf("failed to batch create audit logs: %w", err)
	}

	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			*logs[i] = *entity
		}
	}

	return nil
}
