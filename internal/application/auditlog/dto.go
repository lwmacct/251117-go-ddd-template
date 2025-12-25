package auditlog

import "time"

// AuditLogDTO 审计日志响应 DTO
type AuditLogDTO struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	Details   string    `json:"details"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ListLogsDTO 审计日志列表响应 DTO
type ListLogsDTO struct {
	Logs  []*AuditLogDTO `json:"logs"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}
