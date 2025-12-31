package seeds

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
)

// ProductSeeder 产品种子数据
type ProductSeeder struct{}

// Seed 种子产品数据
func (s *ProductSeeder) Seed(ctx context.Context, db *gorm.DB) error {
	products := []persistence.ProductModel{
		{
			Name:        "free",
			Description: "免费版，适合个人用户体验",
			Price:       0,
			Status:      "active",
		},
		{
			Name:        "starter",
			Description: "入门版，适合小型团队",
			Price:       99,
			Status:      "active",
		},
		{
			Name:        "professional",
			Description: "专业版，适合成长型企业",
			Price:       299,
			Status:      "active",
		},
		{
			Name:        "enterprise",
			Description: "企业版，适合大型组织，提供专属支持",
			Price:       999,
			Status:      "active",
		},
	}

	for i := range products {
		if err := db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).Create(&products[i]).Error; err != nil {
			return err
		}
	}

	return nil
}

// Name 返回种子名称
func (s *ProductSeeder) Name() string {
	return "ProductSeeder"
}
