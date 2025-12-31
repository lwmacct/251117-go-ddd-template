package persistence

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/organization"
	"gorm.io/gorm"
)

// organizationCommandRepository 组织命令仓储的 GORM 实现
type organizationCommandRepository struct {
	*GenericCommandRepository[organization.Organization, *OrganizationModel]
}

// NewOrganizationCommandRepository 创建组织命令仓储实例
func NewOrganizationCommandRepository(db *gorm.DB) organization.CommandRepository {
	return &organizationCommandRepository{
		GenericCommandRepository: NewGenericCommandRepository(
			db, newOrganizationModelFromEntity,
		),
	}
}

// Create、Update、Delete 方法由 GenericCommandRepository 提供

// UpdateStatus 更新组织状态
func (r *organizationCommandRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	if err := r.DB().WithContext(ctx).Model(&OrganizationModel{}).
		Where("id = ?", id).
		Update("status", status).Error; err != nil {
		return fmt.Errorf("failed to update organization status: %w", err)
	}
	return nil
}
