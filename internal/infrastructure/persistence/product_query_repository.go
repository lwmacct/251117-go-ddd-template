package persistence

import (
	"context"
	"errors"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/product"
	"gorm.io/gorm"
)

type productQueryRepository struct {
	db *gorm.DB
}

// NewProductQueryRepository 创建产品读仓储实例
func NewProductQueryRepository(db *gorm.DB) product.QueryRepository {
	return &productQueryRepository{db: db}
}

func (r *productQueryRepository) GetByID(ctx context.Context, id uint) (*product.Product, error) {
	var model ProductModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, product.ErrProductNotFound
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *productQueryRepository) GetByName(ctx context.Context, name string) (*product.Product, error) {
	var model ProductModel
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, product.ErrProductNotFound
		}
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *productQueryRepository) List(ctx context.Context, offset, limit int) ([]*product.Product, error) {
	var models []ProductModel
	if err := r.db.WithContext(ctx).
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}
	return mapProductModelsToEntities(models), nil
}

func (r *productQueryRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&ProductModel{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *productQueryRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&ProductModel{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
