package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/organization"
	"gorm.io/gorm"
)

// orgMemberQueryRepository 组织成员查询仓储的 GORM 实现
type orgMemberQueryRepository struct {
	db *gorm.DB
}

// NewOrgMemberQueryRepository 创建组织成员查询仓储实例
func NewOrgMemberQueryRepository(db *gorm.DB) organization.MemberQueryRepository {
	return &orgMemberQueryRepository{db: db}
}

// GetByOrgAndUser 获取指定组织的指定用户成员信息
func (r *orgMemberQueryRepository) GetByOrgAndUser(ctx context.Context, orgID, userID uint) (*organization.Member, error) {
	var model OrganizationMemberModel
	if err := r.db.WithContext(ctx).
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, organization.ErrMemberNotFound
		}
		return nil, fmt.Errorf("failed to get org member: %w", err)
	}
	return model.ToEntity(), nil
}

// ListByOrg 获取组织的所有成员
func (r *orgMemberQueryRepository) ListByOrg(ctx context.Context, orgID uint, offset, limit int) ([]*organization.Member, error) {
	var models []*OrganizationMemberModel
	if err := r.db.WithContext(ctx).
		Where("organization_id = ?", orgID).
		Offset(offset).Limit(limit).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list org members: %w", err)
	}
	return mapOrgMemberModelPtrsToEntities(models), nil
}

// CountByOrg 统计组织成员数量
func (r *orgMemberQueryRepository) CountByOrg(ctx context.Context, orgID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&OrganizationMemberModel{}).
		Where("organization_id = ?", orgID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count org members: %w", err)
	}
	return count, nil
}

// ListByUser 获取用户加入的所有组织的成员记录
func (r *orgMemberQueryRepository) ListByUser(ctx context.Context, userID uint) ([]*organization.Member, error) {
	var models []*OrganizationMemberModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list user org memberships: %w", err)
	}
	return mapOrgMemberModelPtrsToEntities(models), nil
}

// IsMember 检查用户是否为组织成员
func (r *orgMemberQueryRepository) IsMember(ctx context.Context, orgID, userID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&OrganizationMemberModel{}).
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check membership: %w", err)
	}
	return count > 0, nil
}

// GetOwner 获取组织所有者
func (r *orgMemberQueryRepository) GetOwner(ctx context.Context, orgID uint) (*organization.Member, error) {
	var model OrganizationMemberModel
	if err := r.db.WithContext(ctx).
		Where("organization_id = ? AND role = ?", orgID, organization.MemberRoleOwner).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, organization.ErrMemberNotFound
		}
		return nil, fmt.Errorf("failed to get org owner: %w", err)
	}
	return model.ToEntity(), nil
}

// CountOwners 统计组织所有者数量
func (r *orgMemberQueryRepository) CountOwners(ctx context.Context, orgID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&OrganizationMemberModel{}).
		Where("organization_id = ? AND role = ?", orgID, organization.MemberRoleOwner).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count org owners: %w", err)
	}
	return count, nil
}
