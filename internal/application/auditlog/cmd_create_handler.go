package auditlog

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
)

// CreateHandler 创建审计日志命令处理器
type CreateHandler struct {
	auditLogCommandRepo auditlog.CommandRepository
}

// NewCreateHandler 创建处理器实例
func NewCreateHandler(repo auditlog.CommandRepository) *CreateHandler {
	return &CreateHandler{
		auditLogCommandRepo: repo,
	}
}

// Handle 处理创建审计日志命令
func (h *CreateHandler) Handle(ctx context.Context, cmd CreateCommand) error {
	log := &auditlog.AuditLog{
		UserID:     cmd.UserID,
		Username:   cmd.Username,
		Action:     cmd.Action,
		Resource:   cmd.Resource,
		ResourceID: cmd.ResourceID,
		IPAddress:  cmd.IPAddress,
		UserAgent:  cmd.UserAgent,
		Details:    cmd.Details,
		Status:     cmd.Status,
	}

	return h.auditLogCommandRepo.Create(ctx, log)
}
