package auditlog

import (
	"context"
	"errors"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
)

// GetHandler 获取审计日志查询处理器
type GetHandler struct {
	auditLogQueryRepo auditlog.QueryRepository
}

// NewGetHandler 创建 GetHandler 实例
func NewGetHandler(auditLogQueryRepo auditlog.QueryRepository) *GetHandler {
	return &GetHandler{
		auditLogQueryRepo: auditLogQueryRepo,
	}
}

// Handle 处理获取审计日志查询
func (h *GetHandler) Handle(ctx context.Context, query GetQuery) (*AuditLogDTO, error) {
	log, err := h.auditLogQueryRepo.FindByID(ctx, query.LogID)
	if err != nil || log == nil {
		return nil, errors.New("audit log not found")
	}

	return ToAuditLogDTO(log), nil
}
