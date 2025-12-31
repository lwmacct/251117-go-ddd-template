package persistence

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/product"
	"gorm.io/gorm"
)

type productCommandRepository struct {
	db *gorm.DB
}

// NewProductCommandRepository 创建产品写仓储实例
func NewProductCommandRepository(db *gorm.DB) product.CommandRepository {
	return &productCommandRepository{db: db}
}

func (r *productCommandRepository) Create(ctx context.Context, p *product.Product) error {
	model := newProductModelFromEntity(p)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	// 回写生成的 ID
	p.ID = model.ID
	p.CreatedAt = model.CreatedAt
	p.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *productCommandRepository) Update(ctx context.Context, p *product.Product) error {
	model := newProductModelFromEntity(p)
	result := r.db.WithContext(ctx).Model(&ProductModel{}).
		Where("id = ?", p.ID).
		Updates(map[string]any{
			"name":        model.Name,
			"description": model.Description,
			"price":       model.Price,
			"status":      model.Status,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return product.ErrProductNotFound
	}
	return nil
}

func (r *productCommandRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&ProductModel{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return product.ErrProductNotFound
	}
	return nil
}
