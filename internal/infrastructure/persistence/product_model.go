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
	Code        string  `gorm:"uniqueIndex;size:50;not null"` // 产品代码
	Name        string  `gorm:"size:100;not null"`
	Type        string  `gorm:"size:20;default:'personal';not null"` // personal/team
	Description string  `gorm:"type:text"`
	Price       float64 `gorm:"type:decimal(10,2);not null;default:0"`
	Status      string  `gorm:"size:20;default:'active';not null"`
	LayoutRef   string  `gorm:"size:100;default:''"` // 前端 Layout 组件引用
	MaxSeats    int     `gorm:"default:0;not null"`  // 最大席位数量，0=无限制
	TrialDays   int     `gorm:"default:0;not null"`  // 试用天数

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
		Code:        entity.Code,
		Name:        entity.Name,
		Type:        string(entity.Type),
		Description: entity.Description,
		Price:       entity.Price,
		Status:      string(entity.Status),
		LayoutRef:   entity.LayoutRef,
		MaxSeats:    entity.MaxSeats,
		TrialDays:   entity.TrialDays,
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
		Code:        m.Code,
		Name:        m.Name,
		Type:        product.Type(m.Type),
		Description: m.Description,
		Price:       m.Price,
		Status:      product.Status(m.Status),
		LayoutRef:   m.LayoutRef,
		MaxSeats:    m.MaxSeats,
		TrialDays:   m.TrialDays,
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
