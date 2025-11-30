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

// 操作状态常量
const (
	StatusSuccess = "success"
	StatusFailed  = "failed"
	StatusPending = "pending"
)

// 常见操作类型常量
const (
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionLogin  = "login"
	ActionLogout = "logout"
	ActionExport = "export"
	ActionImport = "import"
)

// IsSuccess 检查操作是否成功
func (a *AuditLog) IsSuccess() bool {
	return a.Status == StatusSuccess
}

// IsFailed 检查操作是否失败
func (a *AuditLog) IsFailed() bool {
	return a.Status == StatusFailed
}

// IsPending 检查操作是否待定
func (a *AuditLog) IsPending() bool {
	return a.Status == StatusPending
}

// GetResourceIdentifier 获取完整资源标识符
func (a *AuditLog) GetResourceIdentifier() string {
	if a.ResourceID == "" {
		return a.Resource
	}
	return a.Resource + ":" + a.ResourceID
}

// IsRecentlyCreated 检查是否在指定时间范围内创建
func (a *AuditLog) IsRecentlyCreated(duration time.Duration) bool {
	return time.Since(a.CreatedAt) <= duration
}

// HasDetails 检查是否有详情信息
func (a *AuditLog) HasDetails() bool {
	return a.Details != ""
}

// IsUserAction 检查是否为用户操作（有用户 ID）
func (a *AuditLog) IsUserAction() bool {
	return a.UserID > 0
}

// IsSystemAction 检查是否为系统操作（无用户 ID）
func (a *AuditLog) IsSystemAction() bool {
	return a.UserID == 0
}

// MatchesFilter 检查日志是否匹配过滤条件
func (a *AuditLog) MatchesFilter(filter FilterOptions) bool {
	if filter.UserID != nil && a.UserID != *filter.UserID {
		return false
	}
	if filter.Action != "" && a.Action != filter.Action {
		return false
	}
	if filter.Resource != "" && a.Resource != filter.Resource {
		return false
	}
	if filter.Status != "" && a.Status != filter.Status {
		return false
	}
	if filter.StartDate != nil && a.CreatedAt.Before(*filter.StartDate) {
		return false
	}
	if filter.EndDate != nil && a.CreatedAt.After(*filter.EndDate) {
		return false
	}
	return true
}

// IsValidFilter 检查过滤条件是否有效
func (f *FilterOptions) IsValidFilter() bool {
	if f.Page < 0 || f.Limit < 0 {
		return false
	}
	if f.StartDate != nil && f.EndDate != nil && f.StartDate.After(*f.EndDate) {
		return false
	}
	return true
}

// SetDefaults 设置默认分页值
func (f *FilterOptions) SetDefaults() {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 20
	}
	if f.Limit > 100 {
		f.Limit = 100
	}
}
