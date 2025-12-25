package auditlog

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
)

// CreateLogHandler 创建审计日志命令处理器
type CreateLogHandler struct {
	auditLogCommandRepo auditlog.CommandRepository
}

// NewCreateLogHandler 创建处理器实例
func NewCreateLogHandler(repo auditlog.CommandRepository) *CreateLogHandler {
	return &CreateLogHandler{
		auditLogCommandRepo: repo,
	}
}

// Handle 处理创建审计日志命令
func (h *CreateLogHandler) Handle(ctx context.Context, cmd CreateLogCommand) error {
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
