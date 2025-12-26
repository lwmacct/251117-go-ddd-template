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
	var model SettingModel
	err := r.db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //nolint:nilnil // returns nil for not found, valid pattern
		}
		return nil, fmt.Errorf("failed to find setting by ID: %w", err)
	}
	return model.ToEntity(), nil
}

// FindByKey 根据 Key 查找配置
func (r *settingQueryRepository) FindByKey(ctx context.Context, key string) (*setting.Setting, error) {
	var model SettingModel
	err := r.db.WithContext(ctx).Where("key = ?", key).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //nolint:nilnil // returns nil for not found, valid pattern
		}
		return nil, fmt.Errorf("failed to find setting by key: %w", err)
	}
	return model.ToEntity(), nil
}

// FindByKeys 根据多个 Key 批量查找配置
func (r *settingQueryRepository) FindByKeys(ctx context.Context, keys []string) ([]*setting.Setting, error) {
	if len(keys) == 0 {
		return []*setting.Setting{}, nil
	}
	var models []SettingModel
	err := r.db.WithContext(ctx).Where("key IN ?", keys).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find settings by keys: %w", err)
	}
	return mapSettingModelsToEntities(models), nil
}

// FindByCategory 根据分类查找配置列表
func (r *settingQueryRepository) FindByCategory(ctx context.Context, category string) ([]*setting.Setting, error) {
	var models []SettingModel
	err := r.db.WithContext(ctx).Where("category = ?", category).Order("key ASC").Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find settings by category: %w", err)
	}
	return mapSettingModelsToEntities(models), nil
}

// FindAll 查找所有配置
func (r *settingQueryRepository) FindAll(ctx context.Context) ([]*setting.Setting, error) {
	var models []SettingModel
	err := r.db.WithContext(ctx).Order("category ASC, key ASC").Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find all settings: %w", err)
	}
	return mapSettingModelsToEntities(models), nil
}
