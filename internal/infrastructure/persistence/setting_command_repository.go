package persistence

import (
	"context"
	"fmt"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/setting"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// settingCommandRepository 配置命令仓储的 GORM 实现
// 嵌入 GenericCommandRepository 以复用基础 CRUD 操作
type settingCommandRepository struct {
	*GenericCommandRepository[setting.Setting, *SettingModel]
}

// NewSettingCommandRepository 创建配置命令仓储实例
func NewSettingCommandRepository(db *gorm.DB) setting.CommandRepository {
	return &settingCommandRepository{
		GenericCommandRepository: NewGenericCommandRepository[setting.Setting, *SettingModel](
			db, newSettingModelFromEntity,
		),
	}
}

// Create、Update、Delete 方法由 GenericCommandRepository 提供

// BatchUpsert 批量插入或更新配置
func (r *settingCommandRepository) BatchUpsert(ctx context.Context, settings []*setting.Setting) error {
	// 使用 GORM 的 Clauses 实现 Upsert (ON CONFLICT DO UPDATE)
	models := make([]*SettingModel, 0, len(settings))
	for _, s := range settings {
		if model := newSettingModelFromEntity(s); model != nil {
			models = append(models, model)
		}
	}

	if err := r.DB().WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(models).Error; err != nil {
		return fmt.Errorf("failed to batch upsert settings: %w", err)
	}

	for i := range models {
		if entity := models[i].ToEntity(); entity != nil {
			*settings[i] = *entity
		}
	}

	return nil
}
