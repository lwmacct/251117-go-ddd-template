package seeds

import (
	"context"
	"errors"

	"gorm.io/gorm"

	productDomain "github.com/lwmacct/251117-go-ddd-template/internal/domain/product"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
)

// ProductSeeder 产品种子数据
type ProductSeeder struct{}

// Seed 种子产品数据
func (s *ProductSeeder) Seed(ctx context.Context, db *gorm.DB) error {
	products := []persistence.ProductModel{
		{
			Code:        "task",
			Name:        "任务管理",
			Type:        string(productDomain.TypeTeam),
			Description: "团队协作任务管理，支持项目跟踪、任务分配和进度管理",
			Price:       99,
			Status:      string(productDomain.StatusActive),
			LayoutRef:   "TaskLayout",
			MaxSeats:    100,
			TrialDays:   14,
		},
		{
			Code:        "crm",
			Name:        "客户管理",
			Type:        string(productDomain.TypeTeam),
			Description: "CRM 客户关系管理，支持销售线索、商机跟进和客户分析",
			Price:       199,
			Status:      string(productDomain.StatusActive),
			LayoutRef:   "CRMLayout",
			MaxSeats:    500,
			TrialDays:   14,
		},
		{
			Code:        "lms",
			Name:        "设备租赁",
			Type:        string(productDomain.TypeTeam),
			Description: "电子产品租赁系统，支持设备管理、租单跟踪和费用结算",
			Price:       299,
			Status:      string(productDomain.StatusActive),
			LayoutRef:   "LMSLayout",
			MaxSeats:    200,
			TrialDays:   30,
		},
	}

	for i := range products {
		// 检查产品是否已存在
		var existingProduct persistence.ProductModel
		err := db.WithContext(ctx).Where("code = ?", products[i].Code).First(&existingProduct).Error
		if err == nil {
			// 产品已存在，跳过
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		// 创建新产品
		if err := db.WithContext(ctx).Create(&products[i]).Error; err != nil {
			return err
		}
	}

	return nil
}

// Name 返回种子名称
func (s *ProductSeeder) Name() string {
	return "ProductSeeder"
}
