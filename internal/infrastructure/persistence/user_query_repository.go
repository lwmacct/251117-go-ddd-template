package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	"gorm.io/gorm"
)

// userQueryRepository 用户查询仓储的 GORM 实现
type userQueryRepository struct {
	db *gorm.DB
}

// NewUserQueryRepository 创建用户查询仓储实例
func NewUserQueryRepository(db *gorm.DB) user.QueryRepository {
	return &userQueryRepository{db: db}
}

// GetByID 根据 ID 获取用户
func (r *userQueryRepository) GetByID(ctx context.Context, id uint) (*user.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return model.ToEntity(), nil
}

// GetByUsername 根据用户名获取用户
func (r *userQueryRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return model.ToEntity(), nil
}

// GetByEmail 根据邮箱获取用户
func (r *userQueryRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return model.ToEntity(), nil
}

// GetByIDWithRoles 根据 ID 获取用户（包含角色和权限信息）
func (r *userQueryRepository) GetByIDWithRoles(ctx context.Context, id uint) (*user.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).
		Preload("Roles.Permissions").
		First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id with roles: %w", err)
	}
	return model.ToEntity(), nil
}

// GetByUsernameWithRoles 根据用户名获取用户（包含角色和权限信息）
func (r *userQueryRepository) GetByUsernameWithRoles(ctx context.Context, username string) (*user.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).
		Preload("Roles.Permissions").
		Where("username = ?", username).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username with roles: %w", err)
	}
	return model.ToEntity(), nil
}

// GetByEmailWithRoles 根据邮箱获取用户（包含角色和权限信息）
func (r *userQueryRepository) GetByEmailWithRoles(ctx context.Context, email string) (*user.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).
		Preload("Roles.Permissions").
		Where("email = ?", email).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email with roles: %w", err)
	}
	return model.ToEntity(), nil
}

// List 获取用户列表 (分页)
func (r *userQueryRepository) List(ctx context.Context, offset, limit int) ([]*user.User, error) {
	var models []UserModel
	query := r.db.WithContext(ctx).Offset(offset).Limit(limit)

	if err := query.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return mapUserModelsToEntities(models), nil
}

// Count 统计用户数量
func (r *userQueryRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&UserModel{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

// GetRoles 获取用户的所有角色 ID
func (r *userQueryRepository) GetRoles(ctx context.Context, userID uint) ([]uint, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).Preload("Roles").First(&model, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	roleIDs := make([]uint, 0, len(model.Roles))
	for _, r := range model.Roles {
		roleIDs = append(roleIDs, r.ID)
	}

	return roleIDs, nil
}

// Search 搜索用户（支持用户名、邮箱、全名模糊匹配）
func (r *userQueryRepository) Search(ctx context.Context, keyword string, offset, limit int) ([]*user.User, error) {
	var models []UserModel
	query := r.db.WithContext(ctx).
		Where("username LIKE ? OR email LIKE ? OR full_name LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").
		Offset(offset).
		Limit(limit)

	if err := query.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	return mapUserModelsToEntities(models), nil
}

// CountBySearch 统计搜索结果数量
func (r *userQueryRepository) CountBySearch(ctx context.Context, keyword string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Where("username LIKE ? OR email LIKE ? OR full_name LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count search results: %w", err)
	}
	return count, nil
}

// Exists 检查用户是否存在
func (r *userQueryRepository) Exists(ctx context.Context, id uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&UserModel{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return count > 0, nil
}

// ExistsByUsername 检查用户名是否存在
func (r *userQueryRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&UserModel{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}
	return count > 0, nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *userQueryRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&UserModel{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	return count > 0, nil
}

// GetUserIDsByRole 获取拥有指定角色的所有用户 ID
func (r *userQueryRepository) GetUserIDsByRole(ctx context.Context, roleID uint) ([]uint, error) {
	var userIDs []uint

	// 通过 user_roles 关联表查询
	if err := r.db.WithContext(ctx).
		Table("user_roles").
		Where("role_id = ?", roleID).
		Pluck("user_id", &userIDs).Error; err != nil {
		return nil, fmt.Errorf("failed to get user IDs by role: %w", err)
	}

	return userIDs, nil
}
