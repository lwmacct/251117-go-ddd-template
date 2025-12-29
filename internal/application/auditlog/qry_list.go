package auditlog

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
)

// ListHandler 获取审计日志列表查询处理器
type ListHandler struct {
	auditLogQueryRepo auditlog.QueryRepository
}

// NewListHandler 创建 ListHandler 实例
func NewListHandler(auditLogQueryRepo auditlog.QueryRepository) *ListHandler {
	return &ListHandler{
		auditLogQueryRepo: auditLogQueryRepo,
	}
}

// Handle 处理获取审计日志列表查询
func (h *ListHandler) Handle(ctx context.Context, query ListQuery) (*ListLogsDTO, error) {
	filter := auditlog.FilterOptions{
		Page:      query.Page,
		Limit:     query.Limit,
		UserID:    query.UserID,
		Action:    query.Action,
		Resource:  query.Resource,
		Status:    query.Status,
		StartDate: query.StartDate,
		EndDate:   query.EndDate,
	}

	logs, total, err := h.auditLogQueryRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list audit logs: %w", err)
	}

	// 转换为 DTO
	logResponses := make([]*AuditLogDTO, 0, len(logs))
	for i := range logs {
		logResponses = append(logResponses, ToAuditLogDTO(&logs[i]))
	}

	return &ListLogsDTO{
		Logs:  logResponses,
		Total: total,
		Page:  query.Page,
		Limit: query.Limit,
	}, nil
}
