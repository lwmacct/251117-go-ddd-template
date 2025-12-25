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
			FullName: "Admin User",
			Status:   "active",
		},
		{
			Username: "testuser",
			Email:    "test@example.com",
			Password: string(hashedPassword),
			FullName: "Test User",
			Status:   "active",
		},
		{
			Username: "demo",
			Email:    "demo@example.com",
			Password: string(hashedPassword),
			FullName: "Demo User",
			Status:   "active",
		},
	}

	result := db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&users)
	if result.Error != nil {
		return result.Error
	}

	slog.Info("Seeded demo users", "attempted", len(users), "inserted", result.RowsAffected)

	return nil
}
