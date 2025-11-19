// Package persistence 提供菜单命令仓储的 GORM 实现
package persistence

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
	"gorm.io/gorm"
)

// menuCommandRepository 菜单命令仓储的 GORM 实现
type menuCommandRepository struct {
	db *gorm.DB
}

// NewMenuCommandRepository 创建菜单命令仓储实例
func NewMenuCommandRepository(db *gorm.DB) menu.CommandRepository {
	return &menuCommandRepository{db: db}
}

// Create 创建菜单
func (r *menuCommandRepository) Create(ctx context.Context, m *menu.Menu) error {
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("failed to create menu: %w", err)
	}
	return nil
}

// Update 更新菜单
func (r *menuCommandRepository) Update(ctx context.Context, m *menu.Menu) error {
	if err := r.db.WithContext(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("failed to update menu: %w", err)
	}
	return nil
}

// Delete 删除菜单
func (r *menuCommandRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&menu.Menu{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete menu: %w", err)
	}
	return nil
}

// UpdateOrder 批量更新菜单排序
func (r *menuCommandRepository) UpdateOrder(ctx context.Context, menus []struct {
	ID       uint
	Order    int
	ParentID *uint
}) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, m := range menus {
			if err := tx.Model(&menu.Menu{}).
				Where("id = ?", m.ID).
				Updates(map[string]interface{}{
					"order":     m.Order,
					"parent_id": m.ParentID,
				}).Error; err != nil {
				return fmt.Errorf("failed to update menu order: %w", err)
			}
		}
		return nil
	})
}
