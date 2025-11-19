// Package persistence 提供菜单查询仓储的 GORM 实现
package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
	"gorm.io/gorm"
)

// menuQueryRepository 菜单查询仓储的 GORM 实现
type menuQueryRepository struct {
	db *gorm.DB
}

// NewMenuQueryRepository 创建菜单查询仓储实例
func NewMenuQueryRepository(db *gorm.DB) menu.QueryRepository {
	return &menuQueryRepository{db: db}
}

// FindByID 根据 ID 查找菜单
func (r *menuQueryRepository) FindByID(ctx context.Context, id uint) (*menu.Menu, error) {
	var m menu.Menu
	err := r.db.WithContext(ctx).First(&m, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("menu not found with id: %d", id)
		}
		return nil, fmt.Errorf("failed to find menu by id: %w", err)
	}
	return &m, nil
}

// FindAll 查找所有菜单（树形结构）
func (r *menuQueryRepository) FindAll(ctx context.Context) ([]*menu.Menu, error) {
	var allMenus []*menu.Menu
	err := r.db.WithContext(ctx).
		Where("parent_id IS NULL").
		Order("`order` ASC, id ASC").
		Find(&allMenus).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find all menus: %w", err)
	}

	// 递归加载子菜单
	for _, m := range allMenus {
		if err := r.loadChildren(ctx, m); err != nil {
			return nil, err
		}
	}

	return allMenus, nil
}

// FindByParentID 根据父 ID 查找子菜单
func (r *menuQueryRepository) FindByParentID(ctx context.Context, parentID *uint) ([]*menu.Menu, error) {
	var menus []*menu.Menu
	query := r.db.WithContext(ctx)

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	err := query.Order("`order` ASC, id ASC").Find(&menus).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find menus by parent id: %w", err)
	}
	return menus, nil
}

// loadChildren 递归加载子菜单
func (r *menuQueryRepository) loadChildren(ctx context.Context, m *menu.Menu) error {
	var children []*menu.Menu
	err := r.db.WithContext(ctx).
		Where("parent_id = ?", m.ID).
		Order("`order` ASC, id ASC").
		Find(&children).Error
	if err != nil {
		return fmt.Errorf("failed to load children menus: %w", err)
	}

	m.Children = children

	for _, child := range children {
		if err := r.loadChildren(ctx, child); err != nil {
			return err
		}
	}

	return nil
}
