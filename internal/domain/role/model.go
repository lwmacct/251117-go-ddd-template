package role

import (
	"time"

	"gorm.io/gorm"
)

// Role represents a role in the RBAC system
type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"size:50;uniqueIndex;not null" json:"name"`
	DisplayName string         `gorm:"size:100;not null" json:"display_name"`
	Description string         `gorm:"size:255" json:"description"`
	IsSystem    bool           `gorm:"default:false;not null" json:"is_system"`
	Permissions []Permission   `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

// Permission represents a permission in the RBAC system
// Uses three-part format: domain:resource:action
type Permission struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Domain      string         `gorm:"size:50;not null;index" json:"domain"`     // admin, user, api
	Resource    string         `gorm:"size:50;not null;index" json:"resource"`   // users, roles, profiles
	Action      string         `gorm:"size:50;not null;index" json:"action"`     // create, read, update, delete
	Description string         `gorm:"size:255" json:"description"`
	Code        string         `gorm:"size:150;uniqueIndex;not null" json:"code"` // domain:resource:action
}

// TableName specifies the table name for Role model
func (Role) TableName() string {
	return "roles"
}

// TableName specifies the table name for Permission model
func (Permission) TableName() string {
	return "permissions"
}

// PermissionCode returns the full permission code in format "domain:resource:action"
func (p *Permission) PermissionCode() string {
	return p.Code
}
