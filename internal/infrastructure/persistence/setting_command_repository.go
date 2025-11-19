// Package persistence 提供配置命令仓储的 GORM 实现
package persistence

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// settingCommandRepository 配置命令仓储的 GORM 实现
type settingCommandRepository struct {
	db *gorm.DB
}

// NewSettingCommandRepository 创建配置命令仓储实例
func NewSettingCommandRepository(db *gorm.DB) setting.CommandRepository {
	return &settingCommandRepository{db: db}
}

// Create 创建配置
func (r *settingCommandRepository) Create(ctx context.Context, s *setting.Setting) error {
	model := newSettingModelFromEntity(s)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create setting: %w", err)
	}
	if saved := model.toEntity(); saved != nil {
		*s = *saved
	}
	return nil
}

// Update 更新配置
func (r *settingCommandRepository) Update(ctx context.Context, s *setting.Setting) error {
	model := newSettingModelFromEntity(s)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return fmt.Errorf("failed to update setting: %w", err)
	}
	if saved := model.toEntity(); saved != nil {
		*s = *saved
	}
	return nil
}

// Delete 删除配置
func (r *settingCommandRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&SettingModel{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete setting: %w", err)
	}
	return nil
}

// BatchUpsert 批量插入或更新配置
func (r *settingCommandRepository) BatchUpsert(ctx context.Context, settings []*setting.Setting) error {
	// 使用 GORM 的 Clauses 实现 Upsert (ON CONFLICT DO UPDATE)
	models := make([]*SettingModel, 0, len(settings))
	for _, s := range settings {
		if model := newSettingModelFromEntity(s); model != nil {
			models = append(models, model)
		}
	}

	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(models).Error; err != nil {
		return fmt.Errorf("failed to batch upsert settings: %w", err)
	}

	for i := range models {
		if entity := models[i].toEntity(); entity != nil {
			*settings[i] = *entity
		}
	}

	return nil
}
