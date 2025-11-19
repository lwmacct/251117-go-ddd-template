// Package menu 定义菜单模块的映射器
package menu

import (
	domainMenu "github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
)

// ToMenuResponse 将领域实体转换为 DTO
func ToMenuResponse(menu *domainMenu.Menu) *MenuResponse {
	if menu == nil {
		return nil
	}

	return &MenuResponse{
		ID:        menu.ID,
		Title:     menu.Title,
		Path:      menu.Path,
		Icon:      menu.Icon,
		ParentID:  menu.ParentID,
		Order:     menu.Order,
		Visible:   menu.Visible,
		CreatedAt: menu.CreatedAt,
		UpdatedAt: menu.UpdatedAt,
	}
}
