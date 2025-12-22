package persistence

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
	"gorm.io/gorm"
)

// menuCommandRepository 菜单命令仓储的 GORM 实现
// 嵌入 GenericCommandRepository 以复用基础 CRUD 操作
type menuCommandRepository struct {
	*GenericCommandRepository[menu.Menu, *MenuModel]
}

// NewMenuCommandRepository 创建菜单命令仓储实例
func NewMenuCommandRepository(db *gorm.DB) menu.CommandRepository {
	return &menuCommandRepository{
		GenericCommandRepository: NewGenericCommandRepository[menu.Menu, *MenuModel](
			db, newMenuModelFromEntity,
		),
	}
}

// Create、Update、Delete 方法由 GenericCommandRepository 提供

// UpdateOrder 批量更新菜单排序
func (r *menuCommandRepository) UpdateOrder(ctx context.Context, menus []struct {
	ID       uint
	Order    int
	ParentID *uint
}) error {
	return r.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, m := range menus {
			if err := tx.Model(&MenuModel{}).
				Where("id = ?", m.ID).
				Updates(map[string]any{
					"order":     m.Order,
					"parent_id": m.ParentID,
				}).Error; err != nil {
				return fmt.Errorf("failed to update menu order: %w", err)
			}
		}
		return nil
	})
}
