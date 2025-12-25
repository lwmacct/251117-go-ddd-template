package menu

import "time"

// MenuDTO 菜单响应 DTO
type MenuDTO struct {
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

// CreateMenuResultDTO 创建菜单结果 DTO
type CreateMenuResultDTO struct {
	ID uint `json:"id"`
}
