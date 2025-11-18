// Package database 提供数据库迁移功能
package database

import (
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

// Migrator 数据库迁移器
type Migrator struct {
	db *gorm.DB
}

// NewMigrator 创建迁移器
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

// AutoMigrate 自动迁移模型
// models 参数为需要迁移的模型指针切片
func (m *Migrator) AutoMigrate(models ...interface{}) error {
	if len(models) == 0 {
		slog.Info("No models to migrate")
		return nil
	}

	slog.Info("Starting database migration", "model_count", len(models))

	if err := m.db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("auto migration failed: %w", err)
	}

	slog.Info("Database migration completed successfully")
	return nil
}

// DropTables 删除表 (谨慎使用！)
func (m *Migrator) DropTables(models ...interface{}) error {
	if len(models) == 0 {
		return nil
	}

	slog.Warn("Dropping tables", "model_count", len(models))

	migrator := m.db.Migrator()
	for _, model := range models {
		if migrator.HasTable(model) {
			if err := migrator.DropTable(model); err != nil {
				return fmt.Errorf("failed to drop table: %w", err)
			}
		}
	}

	slog.Info("Tables dropped successfully")
	return nil
}

// HasTable 检查表是否存在
func (m *Migrator) HasTable(model interface{}) bool {
	return m.db.Migrator().HasTable(model)
}

// CreateIndexes 创建索引
func (m *Migrator) CreateIndexes(model interface{}, indexes []string) error {
	migrator := m.db.Migrator()

	for _, idx := range indexes {
		if !migrator.HasIndex(model, idx) {
			if err := migrator.CreateIndex(model, idx); err != nil {
				return fmt.Errorf("failed to create index %s: %w", idx, err)
			}
			slog.Info("Index created", "index", idx)
		}
	}

	return nil
}
