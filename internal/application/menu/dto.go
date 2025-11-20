// Package menu 定义菜单模块的 DTO
package menu

import "time"

// MenuResponse 菜单响应 DTO
type MenuResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Path      string    `json:"path"`
	Icon      string    `json:"icon"`
	ParentID  *uint     `json:"parent_id"`
	Order     int       `json:"order"`
	Visible   bool      `json:"visible"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
