package persistence

import (
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/stats"
	"gorm.io/gorm"
)

// statsQueryRepository 统计查询仓储的 GORM 实现
type statsQueryRepository struct {
	db *gorm.DB
}

// NewStatsQueryRepository 创建统计查询仓储实例
func NewStatsQueryRepository(db *gorm.DB) stats.QueryRepository {
	return &statsQueryRepository{db: db}
}

// GetSystemStats 获取系统统计信息
func (r *statsQueryRepository) GetSystemStats(recentLogsLimit int) (*stats.SystemStats, error) {
	s := &stats.SystemStats{}

	// 统计用户数
	total, err := r.GetTotalUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to get total users: %w", err)
	}
	s.TotalUsers = total

	active, err := r.GetUserCountByStatus("active")
	if err != nil {
		return nil, fmt.Errorf("failed to get active users: %w", err)
	}
	s.ActiveUsers = active

	inactive, err := r.GetUserCountByStatus("inactive")
	if err != nil {
		return nil, fmt.Errorf("failed to get inactive users: %w", err)
	}
	s.InactiveUsers = inactive

	banned, err := r.GetUserCountByStatus("banned")
	if err != nil {
		return nil, fmt.Errorf("failed to get banned users: %w", err)
	}
	s.BannedUsers = banned

	// 统计角色
	roles, err := r.GetTotalRoles()
	if err != nil {
		return nil, fmt.Errorf("failed to get total roles: %w", err)
	}
	s.TotalRoles = roles

	// 统计权限
	permissions, err := r.GetTotalPermissions()
	if err != nil {
		return nil, fmt.Errorf("failed to get total permissions: %w", err)
	}
	s.TotalPermissions = permissions

	// 统计菜单
	menus, err := r.GetTotalMenus()
	if err != nil {
		return nil, fmt.Errorf("failed to get total menus: %w", err)
	}
	s.TotalMenus = menus

	// 获取最近审计日志
	logs, err := r.GetRecentAuditLogs(recentLogsLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent audit logs: %w", err)
	}
	s.RecentAuditLogs = logs

	return s, nil
}

// GetUserCountByStatus 按状态统计用户数量
func (r *statsQueryRepository) GetUserCountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Table("users").
		Where("deleted_at IS NULL AND status = ?", status).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetTotalUsers 获取用户总数
func (r *statsQueryRepository) GetTotalUsers() (int64, error) {
	var count int64
	err := r.db.Table("users").
		Where("deleted_at IS NULL").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetTotalRoles 获取角色总数
func (r *statsQueryRepository) GetTotalRoles() (int64, error) {
	var count int64
	err := r.db.Table("roles").
		Where("deleted_at IS NULL").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetTotalPermissions 获取权限总数
func (r *statsQueryRepository) GetTotalPermissions() (int64, error) {
	var count int64
	err := r.db.Table("permissions").
		Where("deleted_at IS NULL").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetTotalMenus 获取菜单总数
func (r *statsQueryRepository) GetTotalMenus() (int64, error) {
	var count int64
	err := r.db.Table("menus").
		Where("deleted_at IS NULL").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetRecentAuditLogs 获取最近的审计日志
func (r *statsQueryRepository) GetRecentAuditLogs(limit int) ([]stats.AuditLogSummary, error) {
	var logs []stats.AuditLogSummary
	err := r.db.Table("audit_logs").
		Select("id, user_id, username, action, resource, status, created_at").
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(limit).
		Scan(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}
