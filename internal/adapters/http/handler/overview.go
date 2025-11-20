package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	"gorm.io/gorm"
)

type OverviewHandler struct {
	db *gorm.DB
}

func NewOverviewHandler(db *gorm.DB) *OverviewHandler {
	return &OverviewHandler{db: db}
}

type SystemStatsResponse struct {
	TotalUsers       int64             `json:"total_users"`
	ActiveUsers      int64             `json:"active_users"`
	InactiveUsers    int64             `json:"inactive_users"`
	BannedUsers      int64             `json:"banned_users"`
	TotalRoles       int64             `json:"total_roles"`
	TotalPermissions int64             `json:"total_permissions"`
	TotalMenus       int64             `json:"total_menus"`
	RecentAuditLogs  []AuditLogSummary `json:"recent_audit_logs,omitempty"`
}

type AuditLogSummary struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *OverviewHandler) GetStats(c *gin.Context) {
	stats := SystemStatsResponse{}

	// 统计用户数（排除软删除）
	if err := h.db.Table("users").Where("deleted_at IS NULL").Count(&stats.TotalUsers).Error; err != nil {
		response.InternalError(c, "failed to count users", err.Error())
		return
	}
	if err := h.db.Table("users").Where("deleted_at IS NULL AND status = ?", "active").Count(&stats.ActiveUsers).Error; err != nil {
		response.InternalError(c, "failed to count active users", err.Error())
		return
	}
	if err := h.db.Table("users").Where("deleted_at IS NULL AND status = ?", "inactive").Count(&stats.InactiveUsers).Error; err != nil {
		response.InternalError(c, "failed to count inactive users", err.Error())
		return
	}
	if err := h.db.Table("users").Where("deleted_at IS NULL AND status = ?", "banned").Count(&stats.BannedUsers).Error; err != nil {
		response.InternalError(c, "failed to count banned users", err.Error())
		return
	}

	// 统计角色和权限
	if err := h.db.Table("roles").Where("deleted_at IS NULL").Count(&stats.TotalRoles).Error; err != nil {
		response.InternalError(c, "failed to count roles", err.Error())
		return
	}
	if err := h.db.Table("permissions").Where("deleted_at IS NULL").Count(&stats.TotalPermissions).Error; err != nil {
		response.InternalError(c, "failed to count permissions", err.Error())
		return
	}

	// 统计菜单（使用表名）
	if err := h.db.Table("menus").Count(&stats.TotalMenus).Error; err != nil {
		response.InternalError(c, "failed to count menus", err.Error())
		return
	}

	// 获取最近的审计日志
	var recentLogs []AuditLogSummary
	if err := h.db.Table("audit_logs").
		Select("id, user_id, username, action, resource, status, created_at").
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(5).
		Scan(&recentLogs).Error; err != nil {
		response.InternalError(c, "failed to fetch recent audit logs", err.Error())
		return
	}
	stats.RecentAuditLogs = recentLogs

	response.OK(c, "success", stats)
}
