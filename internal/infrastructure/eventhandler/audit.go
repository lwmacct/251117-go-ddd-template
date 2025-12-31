package eventhandler

import (
	"context"
	"log/slog"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/audit"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/event"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/event/events"
)

// AuditEventHandler 审计日志事件处理器
// 订阅业务事件并创建审计日志记录
type AuditEventHandler struct {
	auditRepo audit.CommandRepository
	logger    *slog.Logger
}

// NewAuditEventHandler 创建审计日志处理器
func NewAuditEventHandler(auditRepo audit.CommandRepository) *AuditEventHandler {
	return &AuditEventHandler{
		auditRepo: auditRepo,
		logger:    slog.Default(),
	}
}

// Handle 处理事件
func (h *AuditEventHandler) Handle(ctx context.Context, e event.Event) error {
	switch evt := e.(type) {
	case *events.CommandExecutedEvent:
		return h.handleCommandExecuted(ctx, evt)
	case *events.LoginSucceededEvent:
		return h.handleLoginSucceeded(ctx, evt)
	case *events.LoginFailedEvent:
		return h.handleLoginFailed(ctx, evt)
	case *events.UserCreatedEvent:
		return h.handleUserCreated(ctx, evt)
	case *events.UserDeletedEvent:
		return h.handleUserDeleted(ctx, evt)
	case *events.UserRoleAssignedEvent:
		return h.handleUserRoleAssigned(ctx, evt)
	case *events.RolePermissionsChangedEvent:
		return h.handleRolePermissionsChanged(ctx, evt)
	default:
		// 忽略不处理的事件
		return nil
	}
}

// handleCommandExecuted 处理命令执行事件
func (h *AuditEventHandler) handleCommandExecuted(ctx context.Context, evt *events.CommandExecutedEvent) error {
	status := "success"
	if !evt.Success {
		status = "failure"
	}

	log := &audit.Audit{
		UserID:     evt.UserID,
		Username:   evt.Username,
		Action:     string(evt.Action),
		Resource:   evt.Resource,
		ResourceID: evt.ResourceID,
		IPAddress:  evt.IPAddress,
		UserAgent:  evt.UserAgent,
		Details:    evt.Details,
		Status:     status,
	}

	return h.createAuditLog(ctx, log, "command_executed")
}

// handleLoginSucceeded 处理登录成功事件
func (h *AuditEventHandler) handleLoginSucceeded(ctx context.Context, evt *events.LoginSucceededEvent) error {
	log := &audit.Audit{
		UserID:    evt.UserID,
		Username:  evt.Username,
		Action:    "login",
		Resource:  "session",
		IPAddress: evt.IPAddress,
		UserAgent: evt.UserAgent,
		Status:    "success",
	}

	return h.createAuditLog(ctx, log, "login_succeeded")
}

// handleLoginFailed 处理登录失败事件
func (h *AuditEventHandler) handleLoginFailed(ctx context.Context, evt *events.LoginFailedEvent) error {
	log := &audit.Audit{
		UserID:    0, // 登录失败时可能没有用户ID
		Username:  evt.Username,
		Action:    "login",
		Resource:  "session",
		IPAddress: evt.IPAddress,
		Details:   evt.Reason,
		Status:    "failure",
	}

	return h.createAuditLog(ctx, log, "login_failed")
}

// handleUserCreated 处理用户创建事件
func (h *AuditEventHandler) handleUserCreated(ctx context.Context, evt *events.UserCreatedEvent) error {
	log := &audit.Audit{
		UserID:     evt.UserID,
		Username:   evt.Username,
		Action:     "create",
		Resource:   "user",
		ResourceID: evt.AggregateID(),
		Status:     "success",
	}

	return h.createAuditLog(ctx, log, "user_created")
}

// handleUserDeleted 处理用户删除事件
func (h *AuditEventHandler) handleUserDeleted(ctx context.Context, evt *events.UserDeletedEvent) error {
	log := &audit.Audit{
		Action:     "delete",
		Resource:   "user",
		ResourceID: evt.AggregateID(),
		Status:     "success",
	}

	return h.createAuditLog(ctx, log, "user_deleted")
}

// handleUserRoleAssigned 处理用户角色分配事件
func (h *AuditEventHandler) handleUserRoleAssigned(ctx context.Context, evt *events.UserRoleAssignedEvent) error {
	log := &audit.Audit{
		Action:     "assign_roles",
		Resource:   "user",
		ResourceID: evt.AggregateID(),
		Status:     "success",
	}

	return h.createAuditLog(ctx, log, "user_role_assigned")
}

// handleRolePermissionsChanged 处理角色权限变更事件
func (h *AuditEventHandler) handleRolePermissionsChanged(ctx context.Context, evt *events.RolePermissionsChangedEvent) error {
	log := &audit.Audit{
		Action:     "set_permissions",
		Resource:   "role",
		ResourceID: evt.AggregateID(),
		Status:     "success",
	}

	return h.createAuditLog(ctx, log, "role_permissions_changed")
}

// createAuditLog 创建审计日志（带错误处理）
func (h *AuditEventHandler) createAuditLog(ctx context.Context, log *audit.Audit, eventType string) error {
	if err := h.auditRepo.Create(ctx, log); err != nil {
		h.logger.Error("failed to create audit log",
			"event_type", eventType,
			"error", err,
		)
		// 审计日志写入失败不应阻塞业务流程
		return nil
	}

	h.logger.Debug("audit log created",
		"event_type", eventType,
		"action", log.Action,
		"resource", log.Resource,
	)

	return nil
}

// Ensure interface is implemented
var _ event.EventHandler = (*AuditEventHandler)(nil)
