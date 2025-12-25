package auditlog

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auditlog"
)

// ToAuditLogDTO 将领域实体转换为 DTO
func ToAuditLogDTO(log *auditlog.AuditLog) *AuditLogDTO {
	if log == nil {
		return nil
	}

	return &AuditLogDTO{
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
