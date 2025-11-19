// Package persistence 提供审计日志查询仓储的 GORM 实现
package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
	"gorm.io/gorm"
)

// auditLogQueryRepository 审计日志查询仓储的 GORM 实现
type auditLogQueryRepository struct {
	db *gorm.DB
}

// NewAuditLogQueryRepository 创建审计日志查询仓储实例
func NewAuditLogQueryRepository(db *gorm.DB) auditlog.QueryRepository {
	return &auditLogQueryRepository{db: db}
}

// FindByID finds an audit log by ID
func (r *auditLogQueryRepository) FindByID(ctx context.Context, id uint) (*auditlog.AuditLog, error) {
	var log auditlog.AuditLog
	err := r.db.WithContext(ctx).First(&log, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("audit log not found")
		}
		return nil, fmt.Errorf("failed to find audit log: %w", err)
	}
	return &log, nil
}

// List returns audit logs with filtering and pagination
func (r *auditLogQueryRepository) List(ctx context.Context, filter auditlog.FilterOptions) ([]auditlog.AuditLog, int64, error) {
	var logs []auditlog.AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&auditlog.AuditLog{})

	// Apply filters
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Resource != "" {
		query = query.Where("resource = ?", filter.Resource)
	}
	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	// Apply pagination
	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}

	// Order by created_at desc
	query = query.Order("created_at DESC")

	if err := query.Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list audit logs: %w", err)
	}

	return logs, total, nil
}

// ListByUser returns audit logs for a specific user
func (r *auditLogQueryRepository) ListByUser(ctx context.Context, userID uint, page, limit int) ([]auditlog.AuditLog, int64, error) {
	var logs []auditlog.AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&auditlog.AuditLog{}).
		Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs by user: %w", err)
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list audit logs by user: %w", err)
	}

	return logs, total, nil
}

// ListByResource returns audit logs for a specific resource
func (r *auditLogQueryRepository) ListByResource(ctx context.Context, resource string, page, limit int) ([]auditlog.AuditLog, int64, error) {
	var logs []auditlog.AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&auditlog.AuditLog{}).
		Where("resource = ?", resource)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs by resource: %w", err)
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list audit logs by resource: %w", err)
	}

	return logs, total, nil
}

// ListByAction returns audit logs for a specific action
func (r *auditLogQueryRepository) ListByAction(ctx context.Context, action string, page, limit int) ([]auditlog.AuditLog, int64, error) {
	var logs []auditlog.AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&auditlog.AuditLog{}).
		Where("action = ?", action)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs by action: %w", err)
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list audit logs by action: %w", err)
	}

	return logs, total, nil
}

// Count returns the total number of audit logs matching the filter
func (r *auditLogQueryRepository) Count(ctx context.Context, filter auditlog.FilterOptions) (int64, error) {
	var total int64

	query := r.db.WithContext(ctx).Model(&auditlog.AuditLog{})

	// Apply filters
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Resource != "" {
		query = query.Where("resource = ?", filter.Resource)
	}
	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	return total, nil
}

// Search searches audit logs by keyword
func (r *auditLogQueryRepository) Search(ctx context.Context, keyword string, page, limit int) ([]auditlog.AuditLog, int64, error) {
	var logs []auditlog.AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&auditlog.AuditLog{}).
		Where("resource LIKE ? OR action LIKE ? OR details LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search audit logs: %w", err)
	}

	return logs, total, nil
}
