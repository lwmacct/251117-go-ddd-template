// Package seeds 提供各种领域模型的种子数据
package seeds

import (
	"context"

	"github.com/lwmacct/251117-go-ddd-template/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserSeeder 用户种子数据
type UserSeeder struct{}

// Seed 执行用户种子数据填充
func (s *UserSeeder) Seed(ctx context.Context, db *gorm.DB) error {
	// 生成密码哈希 (默认密码：password123)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	users := []user.User{
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

	// 使用 FirstOrCreate 避免重复创建
	for _, u := range users {
		var existing user.User
		result := db.Where("username = ?", u.Username).First(&existing)

		if result.Error == gorm.ErrRecordNotFound {
			// 用户不存在，创建新用户
			if err := db.Create(&u).Error; err != nil {
				return err
			}
		}
		// 用户已存在，跳过
	}

	return nil
}
