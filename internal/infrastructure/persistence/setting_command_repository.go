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
	if err := r.db.WithContext(ctx).Create(s).Error; err != nil {
		return fmt.Errorf("failed to create setting: %w", err)
	}
	return nil
}

// Update 更新配置
func (r *settingCommandRepository) Update(ctx context.Context, s *setting.Setting) error {
	if err := r.db.WithContext(ctx).Save(s).Error; err != nil {
		return fmt.Errorf("failed to update setting: %w", err)
	}
	return nil
}

// Delete 删除配置
func (r *settingCommandRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&setting.Setting{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete setting: %w", err)
	}
	return nil
}

// BatchUpsert 批量插入或更新配置
func (r *settingCommandRepository) BatchUpsert(ctx context.Context, settings []*setting.Setting) error {
	// 使用 GORM 的 Clauses 实现 Upsert (ON CONFLICT DO UPDATE)
	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(settings).Error; err != nil {
		return fmt.Errorf("failed to batch upsert settings: %w", err)
	}
	return nil
}
