package seeds

import (
	"context"
	"log/slog"
	"time"

	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// OrganizationSeeder 创建示例组织和团队。
// 依赖 RBACSeeder 和 UserSeeder 已创建 admin 用户。
type OrganizationSeeder struct{}

// Seed implements database.Seeder interface.
func (s *OrganizationSeeder) Seed(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 创建组织 acme
		org := &persistence.OrganizationModel{
			Name:        "acme",
			DisplayName: "Acme Corporation",
			Description: "示例组织",
			Status:      "active",
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).Create(org).Error; err != nil {
			return err
		}

		// 如果是冲突跳过，需要重新查询获取 ID
		if org.ID == 0 {
			if err := tx.Where("name = ?", "acme").First(org).Error; err != nil {
				return err
			}
		}
		slog.Info("seeded organization", "name", org.Name, "id", org.ID)

		// 2. 查找 admin 用户
		var admin persistence.UserModel
		if err := tx.Where("username = ?", "admin").First(&admin).Error; err != nil {
			slog.Warn("admin user not found, skipping member seeding", "err", err)
			return nil
		}

		// 3. 将 admin 设为 owner
		member := &persistence.OrganizationMemberModel{
			OrganizationID: org.ID,
			UserID:         admin.ID,
			Role:           "owner",
			JoinedAt:       time.Now(),
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "organization_id"}, {Name: "user_id"}},
			DoNothing: true,
		}).Create(member).Error; err != nil {
			return err
		}
		slog.Info("seeded organization member", "org", org.Name, "user", admin.Username, "role", member.Role)

		// 4. 创建 engineering 团队
		team := &persistence.TeamModel{
			OrganizationID: org.ID,
			Name:           "engineering",
			DisplayName:    "Engineering Team",
			Description:    "工程团队",
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "organization_id"}, {Name: "name"}},
			DoNothing: true,
		}).Create(team).Error; err != nil {
			return err
		}

		// 如果是冲突跳过，需要重新查询获取 ID
		if team.ID == 0 {
			if err := tx.Where("organization_id = ? AND name = ?", org.ID, "engineering").First(team).Error; err != nil {
				return err
			}
		}
		slog.Info("seeded team", "name", team.Name, "org", org.Name, "id", team.ID)

		// 5. 将 admin 加入团队
		teamMember := &persistence.TeamMemberModel{
			TeamID:   team.ID,
			UserID:   admin.ID,
			Role:     "lead",
			JoinedAt: time.Now(),
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "team_id"}, {Name: "user_id"}},
			DoNothing: true,
		}).Create(teamMember).Error; err != nil {
			return err
		}
		slog.Info("seeded team member", "team", team.Name, "user", admin.Username, "role", teamMember.Role)

		return nil
	})
}

// Name implements database.Seeder interface.
func (s *OrganizationSeeder) Name() string {
	return "OrganizationSeeder"
}
