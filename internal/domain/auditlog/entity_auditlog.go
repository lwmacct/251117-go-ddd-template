// Package auditlog 定义审计日志领域模型。
//
// 审计日志用于记录系统中的关键操作，支持：
//   - 用户行为追踪：记录 who (UserID/Username) 做了 what (Action)
//   - 资源变更审计：记录对哪个资源 (Resource/ResourceID) 的操作
//   - 安全分析：记录 IP 地址和 User-Agent 用于安全审计
//   - 合规需求：满足 SOC2、GDPR 等合规性审计要求
//
// FilterOptions 提供灵活的日志查询过滤能力。
package auditlog

import "time"

// AuditLog represents an audit log entry for tracking user actions
type AuditLog struct {
	ID         uint       `json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"-"`
	UserID     uint       `json:"user_id"`
	Username   string     `json:"username"`
	Action     string     `json:"action"`
	Resource   string     `json:"resource"`
	ResourceID string     `json:"resource_id,omitempty"`
	IPAddress  string     `json:"ip_address,omitempty"`
	UserAgent  string     `json:"user_agent,omitempty"`
	Details    string     `json:"details,omitempty"`
	Status     string     `json:"status"`
}

// FilterOptions represents options for filtering audit logs
type FilterOptions struct {
	UserID    *uint
	Action    string
	Resource  string
	Status    string
	StartDate *time.Time
	EndDate   *time.Time
	Page      int
	Limit     int
}
