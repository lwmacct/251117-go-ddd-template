package menu

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
)

// ToMenuDTO 将领域实体转换为 DTO
func ToMenuDTO(menu *menu.Menu) *MenuDTO {
	if menu == nil {
		return nil
	}

	return &MenuDTO{
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
