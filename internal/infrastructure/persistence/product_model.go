package persistence

import (
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/product"
	"gorm.io/gorm"
)

// ProductModel 定义产品的 GORM 持久化模型
//
//nolint:recvcheck // TableName uses value receiver per GORM convention
type ProductModel struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"uniqueIndex;size:100;not null"`
	Description string  `gorm:"type:text"`
	Price       float64 `gorm:"type:decimal(10,2);not null;default:0"`
	Status      string  `gorm:"size:20;default:'active';not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName 指定产品表名
func (ProductModel) TableName() string {
	return "products"
}

func newProductModelFromEntity(entity *product.Product) *ProductModel {
	if entity == nil {
		return nil
	}

	return &ProductModel{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		Price:       entity.Price,
		Status:      string(entity.Status),
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

// ToEntity 将 GORM Model 转换为 Domain Entity
func (m *ProductModel) ToEntity() *product.Product {
	if m == nil {
		return nil
	}

	return &product.Product{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Price:       m.Price,
		Status:      product.Status(m.Status),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func mapProductModelsToEntities(models []ProductModel) []*product.Product {
	if len(models) == 0 {
		return nil
	}

	products := make([]*product.Product, 0, len(models))
	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			products = append(products, entity)
		}
	}
	return products
}
