package auditlog

import (
	"time"

	"gorm.io/gorm"
)

// AuditLog represents an audit log entry for tracking user actions
type AuditLog struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	Username   string         `gorm:"size:100;not null" json:"username"`
	Action     string         `gorm:"size:100;not null" json:"action"`
	Resource   string         `gorm:"size:100;not null" json:"resource"`
	ResourceID string         `gorm:"size:100" json:"resource_id,omitempty"`
	IPAddress  string         `gorm:"size:45" json:"ip_address,omitempty"`
	UserAgent  string         `gorm:"size:255" json:"user_agent,omitempty"`
	Details    string         `gorm:"type:text" json:"details,omitempty"`
	Status     string         `gorm:"size:20;default:'success'" json:"status"`
}

// TableName specifies the table name for AuditLog model
func (AuditLog) TableName() string {
	return "audit_logs"
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
