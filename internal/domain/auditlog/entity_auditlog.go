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
	UserID     *uint
	Action     string
	Resource   string
	Status     string
	StartDate  *time.Time
	EndDate    *time.Time
	Page       int
	Limit      int
}
