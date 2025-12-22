// Package query 定义统计查询处理器
package query

import (
	"context"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/stats"
)

// GetStatsQuery 获取统计信息的查询
type GetStatsQuery struct {
	RecentLogsLimit int // 最近审计日志条数，默认 5
}

// GetStatsResult 统计查询结果
type GetStatsResult struct {
	TotalUsers       int64                   `json:"total_users"`
	ActiveUsers      int64                   `json:"active_users"`
	InactiveUsers    int64                   `json:"inactive_users"`
	BannedUsers      int64                   `json:"banned_users"`
	TotalRoles       int64                   `json:"total_roles"`
	TotalPermissions int64                   `json:"total_permissions"`
	TotalMenus       int64                   `json:"total_menus"`
	RecentAuditLogs  []AuditLogSummaryResult `json:"recent_audit_logs,omitempty"`
}

// AuditLogSummaryResult 审计日志摘要结果
type AuditLogSummaryResult struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// GetStatsHandler 获取统计信息的处理器
type GetStatsHandler struct {
	statsQueryRepo stats.QueryRepository
}

// NewGetStatsHandler 创建 GetStatsHandler 实例
func NewGetStatsHandler(statsQueryRepo stats.QueryRepository) *GetStatsHandler {
	return &GetStatsHandler{
		statsQueryRepo: statsQueryRepo,
	}
}

// Handle 处理获取统计信息的查询
func (h *GetStatsHandler) Handle(_ context.Context, query GetStatsQuery) (*GetStatsResult, error) {
	limit := query.RecentLogsLimit
	if limit <= 0 {
		limit = 5
	}

	systemStats, err := h.statsQueryRepo.GetSystemStats(limit)
	if err != nil {
		return nil, err
	}

	return toGetStatsResult(systemStats), nil
}

// toGetStatsResult 将 domain 统计信息转换为结果 DTO
func toGetStatsResult(s *stats.SystemStats) *GetStatsResult {
	result := &GetStatsResult{
		TotalUsers:       s.TotalUsers,
		ActiveUsers:      s.ActiveUsers,
		InactiveUsers:    s.InactiveUsers,
		BannedUsers:      s.BannedUsers,
		TotalRoles:       s.TotalRoles,
		TotalPermissions: s.TotalPermissions,
		TotalMenus:       s.TotalMenus,
	}

	if len(s.RecentAuditLogs) > 0 {
		result.RecentAuditLogs = make([]AuditLogSummaryResult, len(s.RecentAuditLogs))
		for i, log := range s.RecentAuditLogs {
			result.RecentAuditLogs[i] = AuditLogSummaryResult{
				ID:        log.ID,
				UserID:    log.UserID,
				Username:  log.Username,
				Action:    log.Action,
				Resource:  log.Resource,
				Status:    log.Status,
				CreatedAt: log.CreatedAt,
			}
		}
	}

	return result
}
