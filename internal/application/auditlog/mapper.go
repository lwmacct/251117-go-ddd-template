package auditlog

import (
	domainAuditLog "github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
)

// ToAuditLogResponse 将领域实体转换为 DTO
func ToAuditLogResponse(log *domainAuditLog.AuditLog) *AuditLogResponse {
	if log == nil {
		return nil
	}

	return &AuditLogResponse{
		ID:        log.ID,
		UserID:    log.UserID,
		Action:    log.Action,
		Resource:  log.Resource,
		Details:   log.Details,
		IPAddress: log.IPAddress,
		UserAgent: log.UserAgent,
		Status:    log.Status,
		CreatedAt: log.CreatedAt,
	}
}
