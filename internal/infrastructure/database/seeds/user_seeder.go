// Package seeds 提供各种领域模型的种子数据
package seeds

import (
	"context"
	"log/slog"

	_persistence "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UserSeeder 用户种子数据
type UserSeeder struct{}

// Seed 执行用户种子数据填充
func (s *UserSeeder) Seed(ctx context.Context, db *gorm.DB) error {
	db = db.WithContext(ctx)

	// 生成密码哈希 (默认密码：password123)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	users := []_persistence.UserModel{
		{
			Username: "admin",
			Email:    "admin@example.com",
			Password: string(hashedPassword),
			FullName: "System Administrator",
			Avatar:   "https://api.dicebear.com/9.x/micah/svg?seed=admin",
			Status:   "active",
		},
		{
			Username: "testuser",
			Email:    "test@example.com",
			Password: string(hashedPassword),
			FullName: "Test User",
			Avatar:   "https://api.dicebear.com/9.x/micah/svg?seed=testuser",
			Status:   "active",
		},
		{
			Username: "demo",
			Email:    "demo@example.com",
			Password: string(hashedPassword),
			FullName: "Demo User",
			Avatar:   "https://api.dicebear.com/9.x/micah/svg?seed=demo",
			Status:   "active",
		},
	}

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "username"}},
		DoUpdates: clause.AssignmentColumns([]string{"email", "password", "full_name", "avatar", "bio", "status"}),
	}).Create(&users)
	if result.Error != nil {
		return result.Error
	}

	slog.Info("Seeded demo users", "attempted", len(users), "inserted", result.RowsAffected)

	return nil
}

// Name 返回种子器名称。
func (s *UserSeeder) Name() string {
	return "UserSeeder"
}
