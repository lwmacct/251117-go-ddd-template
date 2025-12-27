package bootstrap

import (
	"gorm.io/gorm"

	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/captcha"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/stats"
)

// newRepositoriesModule 初始化仓储模块
// 依赖：InfrastructureModule.DB
func newRepositoriesModule(db *gorm.DB) *RepositoriesModule {
	// Captcha Repository（内存实现，组合接口）
	captchaRepo := captcha.NewRepository()

	return &RepositoriesModule{
		// CQRS 仓储（数据库实现）
		User:        persistence.NewUserRepositories(db),
		AuditLog:    persistence.NewAuditLogRepositories(db),
		Role:        persistence.NewRoleRepositories(db),
		Permission:  persistence.NewPermissionRepositories(db),
		PAT:         persistence.NewPATRepositories(db),
		Menu:        persistence.NewMenuRepositories(db),
		Setting:     persistence.NewSettingRepositories(db),
		UserSetting: persistence.NewUserSettingRepositories(db),
		TwoFA:       persistence.NewTwoFARepositories(db),

		// 特殊仓储（内存实现）
		CaptchaCommand: captchaRepo,
		CaptchaQuery:   captchaRepo,

		// 只读仓储
		StatsQuery: stats.NewQueryRepository(db),
	}
}
