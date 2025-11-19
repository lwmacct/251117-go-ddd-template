package role

import "time"

// Permission represents a permission in the RBAC system
// Uses three-part format: domain:resource:action
type Permission struct {
	ID          uint       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"-"`
	Domain      string     `json:"domain"`
	Resource    string     `json:"resource"`
	Action      string     `json:"action"`
	Description string     `json:"description"`
	Code        string     `json:"code"`
}

// PermissionCode returns the full permission code in format "domain:resource:action"
func (p *Permission) PermissionCode() string {
	return p.Code
}
