// Package persistence 提供配置查询仓储的 GORM 实现
package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
	"gorm.io/gorm"
)

// settingQueryRepository 配置查询仓储的 GORM 实现
type settingQueryRepository struct {
	db *gorm.DB
}

// NewSettingQueryRepository 创建配置查询仓储实例
func NewSettingQueryRepository(db *gorm.DB) setting.QueryRepository {
	return &settingQueryRepository{db: db}
}

// FindByID 根据 ID 查找配置
func (r *settingQueryRepository) FindByID(ctx context.Context, id uint) (*setting.Setting, error) {
	var s setting.Setting
	err := r.db.WithContext(ctx).First(&s, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find setting by ID: %w", err)
	}
	return &s, nil
}

// FindByKey 根据 Key 查找配置
func (r *settingQueryRepository) FindByKey(ctx context.Context, key string) (*setting.Setting, error) {
	var s setting.Setting
	err := r.db.WithContext(ctx).Where("key = ?", key).First(&s).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find setting by key: %w", err)
	}
	return &s, nil
}

// FindByCategory 根据分类查找配置列表
func (r *settingQueryRepository) FindByCategory(ctx context.Context, category string) ([]*setting.Setting, error) {
	var settings []*setting.Setting
	err := r.db.WithContext(ctx).Where("category = ?", category).Order("key ASC").Find(&settings).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find settings by category: %w", err)
	}
	return settings, nil
}

// FindAll 查找所有配置
func (r *settingQueryRepository) FindAll(ctx context.Context) ([]*setting.Setting, error) {
	var settings []*setting.Setting
	err := r.db.WithContext(ctx).Order("category ASC, key ASC").Find(&settings).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find all settings: %w", err)
	}
	return settings, nil
}
