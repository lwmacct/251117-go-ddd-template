package persistence

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/menu"
	"gorm.io/gorm"
)

type menuRepository struct {
	db *gorm.DB
}

// NewMenuRepository 创建菜单仓储实例
func NewMenuRepository(db *gorm.DB) menu.Repository {
	return &menuRepository{db: db}
}

func (r *menuRepository) Create(ctx context.Context, m *menu.Menu) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *menuRepository) Update(ctx context.Context, m *menu.Menu) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *menuRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&menu.Menu{}, id).Error
}

func (r *menuRepository) FindByID(ctx context.Context, id uint) (*menu.Menu, error) {
	var m menu.Menu
	err := r.db.WithContext(ctx).First(&m, id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *menuRepository) FindAll(ctx context.Context) ([]*menu.Menu, error) {
	var allMenus []*menu.Menu
	err := r.db.WithContext(ctx).
		Where("parent_id IS NULL").
		Order("`order` ASC, id ASC").
		Find(&allMenus).Error
	if err != nil {
		return nil, err
	}

	// 递归加载子菜单
	for _, m := range allMenus {
		if err := r.loadChildren(ctx, m); err != nil {
			return nil, err
		}
	}

	return allMenus, nil
}

func (r *menuRepository) FindByParentID(ctx context.Context, parentID *uint) ([]*menu.Menu, error) {
	var menus []*menu.Menu
	query := r.db.WithContext(ctx)

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	err := query.Order("`order` ASC, id ASC").Find(&menus).Error
	return menus, err
}

func (r *menuRepository) UpdateOrder(ctx context.Context, menus []struct {
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
				return err
			}
		}
		return nil
	})
}

// loadChildren 递归加载子菜单
func (r *menuRepository) loadChildren(ctx context.Context, m *menu.Menu) error {
	var children []*menu.Menu
	err := r.db.WithContext(ctx).
		Where("parent_id = ?", m.ID).
		Order("`order` ASC, id ASC").
		Find(&children).Error
	if err != nil {
		return err
	}

	m.Children = children

	for _, child := range children {
		if err := r.loadChildren(ctx, child); err != nil {
			return err
		}
	}

	return nil
}
