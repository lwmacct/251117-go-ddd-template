package seeds

import (
	"context"
	"log/slog"

	_persistence "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SettingCategorySeeder 配置分类种子数据
type SettingCategorySeeder struct{}

// Seed 执行配置分类种子数据填充
func (s *SettingCategorySeeder) Seed(ctx context.Context, db *gorm.DB) error {
	db = db.WithContext(ctx)

	categories := []_persistence.SettingCategoryModel{
		{
			Key:   "general",
			Label: "常规设置",
			Icon:  "mdi-cog",
			Order: 1,
		},
		{
			Key:   "security",
			Label: "安全设置",
			Icon:  "mdi-shield-lock",
			Order: 2,
		},
		{
			Key:   "notification",
			Label: "通知设置",
			Icon:  "mdi-bell",
			Order: 3,
		},
		{
			Key:   "backup",
			Label: "备份设置",
			Icon:  "mdi-backup-restore",
			Order: 4,
		},
	}

	result := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"label", "icon", "sort_order",
		}), // 更新 UI 元数据，保持数据一致
	}).Create(&categories)
	if result.Error != nil {
		return result.Error
	}

	slog.Info("Seeded setting categories", "attempted", len(categories), "inserted", result.RowsAffected)
	return nil
}
