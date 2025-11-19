package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
	"gorm.io/gorm"
)

// AuditLogRepositoryImpl implements auditlog.Repository interface
type AuditLogRepositoryImpl struct {
	db *gorm.DB
}

// NewAuditLogRepository creates a new instance of AuditLogRepositoryImpl
func NewAuditLogRepository(db *gorm.DB) auditlog.Repository {
	return &AuditLogRepositoryImpl{db: db}
}

// Create creates a new audit log entry
func (a *AuditLogRepositoryImpl) Create(ctx context.Context, log *auditlog.AuditLog) error {
	return a.db.WithContext(ctx).Create(log).Error
}

// FindByID finds an audit log by ID
func (a *AuditLogRepositoryImpl) FindByID(ctx context.Context, id uint) (*auditlog.AuditLog, error) {
	var log auditlog.AuditLog
	err := a.db.WithContext(ctx).First(&log, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &log, nil
}

// List returns audit logs with filtering and pagination
func (a *AuditLogRepositoryImpl) List(ctx context.Context, filter auditlog.FilterOptions) ([]auditlog.AuditLog, int64, error) {
	var logs []auditlog.AuditLog
	var total int64

	query := a.db.WithContext(ctx).Model(&auditlog.AuditLog{})

	// Apply filters
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}
	if filter.Resource != "" {
		query = query.Where("resource = ?", filter.Resource)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}

	offset := (filter.Page - 1) * filter.Limit
	err := query.Offset(offset).Limit(filter.Limit).
		Order("created_at DESC").
		Find(&logs).Error

	return logs, total, err
}

// ListByUser returns audit logs for a specific user
func (a *AuditLogRepositoryImpl) ListByUser(ctx context.Context, userID uint, page, limit int) ([]auditlog.AuditLog, int64, error) {
	var logs []auditlog.AuditLog
	var total int64

	query := a.db.WithContext(ctx).Model(&auditlog.AuditLog{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&logs).Error

	return logs, total, err
}

// ListByResource returns audit logs for a specific resource
func (a *AuditLogRepositoryImpl) ListByResource(ctx context.Context, resource string, page, limit int) ([]auditlog.AuditLog, int64, error) {
	var logs []auditlog.AuditLog
	var total int64

	query := a.db.WithContext(ctx).Model(&auditlog.AuditLog{}).Where("resource = ?", resource)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&logs).Error

	return logs, total, err
}

// Delete deletes an audit log (soft delete, for data retention policy)
func (a *AuditLogRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return a.db.WithContext(ctx).Delete(&auditlog.AuditLog{}, id).Error
}

// DeleteOlderThan deletes audit logs older than the specified days
func (a *AuditLogRepositoryImpl) DeleteOlderThan(ctx context.Context, days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	return a.db.WithContext(ctx).
		Where("created_at < ?", cutoffDate).
		Delete(&auditlog.AuditLog{}).Error
}
