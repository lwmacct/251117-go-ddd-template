// Package database 提供数据库相关的基础设施
package database

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Migration 迁移记录模型
type Migration struct {
	ID        uint      `gorm:"primarykey"`
	Version   string    `gorm:"uniqueIndex;size:50;comment:迁移版本号"`
	Name      string    `gorm:"size:255;comment:迁移名称"`
	AppliedAt time.Time `gorm:"comment:执行时间"`
}

// TableName 指定表名
func (Migration) TableName() string {
	return "migrations"
}

// MigrationManager 迁移管理器
type MigrationManager struct {
	db     *gorm.DB
	models []any
}

// NewMigrationManager 创建迁移管理器
func NewMigrationManager(db *gorm.DB, models []any) *MigrationManager {
	return &MigrationManager{
		db:     db,
		models: models,
	}
}

// Up 执行向上迁移
func (m *MigrationManager) Up() error {
	// 1. 确保迁移记录表存在
	if err := m.db.AutoMigrate(&Migration{}); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// 2. 执行所有模型的迁移
	if err := m.db.AutoMigrate(m.models...); err != nil {
		return fmt.Errorf("failed to migrate models: %w", err)
	}

	// 3. 记录迁移
	migration := &Migration{
		Version:   time.Now().Format("20060102150405"),
		Name:      "auto_migrate",
		AppliedAt: time.Now(),
	}
	if err := m.db.Create(migration).Error; err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	return nil
}

// Status 查看迁移状态
func (m *MigrationManager) Status() ([]Migration, error) {
	// 确保迁移记录表存在
	if !m.db.Migrator().HasTable(&Migration{}) {
		return []Migration{}, nil
	}

	var migrations []Migration
	if err := m.db.Order("applied_at DESC").Find(&migrations).Error; err != nil {
		return nil, fmt.Errorf("failed to query migrations: %w", err)
	}

	return migrations, nil
}

// Fresh 删除所有表并重新迁移（危险！仅开发环境使用）
func (m *MigrationManager) Fresh() error {
	migrator := m.db.Migrator()

	// 1. 删除所有模型表
	for _, model := range m.models {
		if migrator.HasTable(model) {
			if err := migrator.DropTable(model); err != nil {
				return fmt.Errorf("failed to drop table for model %T: %w", model, err)
			}
		}
	}

	// 2. 删除迁移记录表
	if migrator.HasTable(&Migration{}) {
		if err := migrator.DropTable(&Migration{}); err != nil {
			return fmt.Errorf("failed to drop migrations table: %w", err)
		}
	}

	// 3. 重新执行迁移
	return m.Up()
}

// HasMigrations 检查是否有迁移记录
func (m *MigrationManager) HasMigrations() bool {
	var count int64
	if !m.db.Migrator().HasTable(&Migration{}) {
		return false
	}
	m.db.Model(&Migration{}).Count(&count)
	return count > 0
}
