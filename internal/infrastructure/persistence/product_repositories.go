package persistence

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/product"
	"gorm.io/gorm"
)

// ProductRepositories 聚合产品读写仓储
type ProductRepositories struct {
	Command product.CommandRepository
	Query   product.QueryRepository
}

// NewProductRepositories 创建产品仓储聚合实例
func NewProductRepositories(db *gorm.DB) ProductRepositories {
	return ProductRepositories{
		Command: NewProductCommandRepository(db),
		Query:   NewProductQueryRepository(db),
	}
}
