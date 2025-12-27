package menu

import "time"

// CreateDTO 创建菜单请求 DTO
type CreateDTO struct {
	Title    string `json:"title" binding:"required,min=1,max=100" example:"系统管理"`
	Path     string `json:"path" binding:"required,max=255" example:"/system"`
	Icon     string `json:"icon" binding:"omitempty,max=100" example:"setting"`
	ParentID *uint  `json:"parent_id" example:"0"`
	Order    int    `json:"order" example:"1"`
	Visible  *bool  `json:"visible" example:"true"`
}

// UpdateDTO 更新菜单请求 DTO
type UpdateDTO struct {
	Title    *string `json:"title" binding:"omitempty,min=1,max=100"`
	Path     *string `json:"path" binding:"omitempty,max=255"`
	Icon     *string `json:"icon" binding:"omitempty,max=100"`
	ParentID *uint   `json:"parent_id"`
	Order    *int    `json:"order"`
	Visible  *bool   `json:"visible"`
}

// ReorderDTO 批量更新排序请求 DTO
type ReorderDTO struct {
	Menus []struct {
		ID       uint  `json:"id" binding:"required"`
		Order    int   `json:"order"`
		ParentID *uint `json:"parent_id"`
	} `json:"menus" binding:"required,dive"`
}

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

// CreateResultDTO 创建菜单结果 DTO
type CreateResultDTO struct {
	ID uint `json:"id"`
}
