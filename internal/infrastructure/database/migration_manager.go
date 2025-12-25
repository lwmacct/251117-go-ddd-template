package database

import (
	"fmt"
	"strings"
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

// Fresh 删除所有表并重新迁移 (危险！仅开发环境使用)
func (m *MigrationManager) Fresh() error {
	if err := m.dropAllTablesWithSQL(); err != nil {
		return fmt.Errorf("failed to drop tables via SQL: %w", err)
	}

	// 重新执行迁移
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

// dropAllTablesWithSQL 使用原生 SQL 删除所有用户表
func (m *MigrationManager) dropAllTablesWithSQL() error {
	const listTablesSQL = `
SELECT table_schema, table_name
FROM information_schema.tables
WHERE table_type = 'BASE TABLE'
  AND table_schema NOT IN ('pg_catalog', 'information_schema')
`

	type tableInfo struct {
		Schema string `gorm:"column:table_schema"`
		Name   string `gorm:"column:table_name"`
	}

	var tables []tableInfo
	if err := m.db.Raw(listTablesSQL).Scan(&tables).Error; err != nil {
		return fmt.Errorf("failed to list tables: %w", err)
	}

	if len(tables) == 0 {
		return nil
	}

	return m.db.Transaction(func(tx *gorm.DB) error {
		for _, tbl := range tables {
			stmt := fmt.Sprintf("DROP TABLE IF EXISTS %s.%s CASCADE",
				quoteIdentifier(tbl.Schema), quoteIdentifier(tbl.Name))
			if err := tx.Exec(stmt).Error; err != nil {
				return fmt.Errorf("failed to drop table %s.%s: %w", tbl.Schema, tbl.Name, err)
			}
		}
		return nil
	})
}

// quoteIdentifier 安全地引用标识符
func quoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}
