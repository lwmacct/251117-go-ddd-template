package auditlog

import (
	"context"
	"errors"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
)

// GetLogHandler 获取审计日志查询处理器
type GetLogHandler struct {
	auditLogQueryRepo auditlog.QueryRepository
}

// NewGetLogHandler 创建 GetLogHandler 实例
func NewGetLogHandler(auditLogQueryRepo auditlog.QueryRepository) *GetLogHandler {
	return &GetLogHandler{
		auditLogQueryRepo: auditLogQueryRepo,
	}
}

// Handle 处理获取审计日志查询
func (h *GetLogHandler) Handle(ctx context.Context, query GetLogQuery) (*AuditLogDTO, error) {
	log, err := h.auditLogQueryRepo.FindByID(ctx, query.LogID)
	if err != nil || log == nil {
		return nil, errors.New("audit log not found")
	}

	return ToAuditLogDTO(log), nil
}
