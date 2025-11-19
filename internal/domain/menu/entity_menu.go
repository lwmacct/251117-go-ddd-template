package menu

import (
	"time"
)

// Menu 菜单实体
type Menu struct {
	ID        uint       `json:"id"`
	Title     string     `json:"title"`
	Path      string     `json:"path"`
	Icon      string     `json:"icon"`
	ParentID  *uint      `json:"parent_id"`
	Order     int        `json:"order"`
	Visible   bool       `json:"visible"`
	Children  []*Menu    `json:"children,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}
