package bootstrap

import (
	"github.com/lwmacct/251117-go-ddd-template/internal/config"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"

	authInfra "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/captcha"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/twofa"
)

// newServicesModule 初始化服务模块
// 依赖：InfrastructureModule, RepositoriesModule, config.Config
func newServicesModule(cfg *config.Config, infra *InfrastructureModule, repos *RepositoriesModule) *ServicesModule {
	m := &ServicesModule{}

	// Infrastructure 组件
	m.JWT = authInfra.NewJWTManager(cfg.JWT.Secret, cfg.JWT.AccessTokenExpiry, cfg.JWT.RefreshTokenExpiry)
	tokenGenerator := authInfra.NewTokenGenerator()
	m.TokenGenerator = tokenGenerator
	m.LoginSession = authInfra.NewLoginSessionService()
	m.PermissionCache = authInfra.NewPermissionCacheService(infra.RedisClient, repos.User.Query, cfg.Data.RedisKeyPrefix)

	// Domain Services
	passwordPolicy := auth.DefaultPasswordPolicy()
	m.Auth = authInfra.NewAuthService(m.JWT, tokenGenerator, passwordPolicy)

	// Captcha Service
	m.Captcha = captcha.NewService()

	// PAT Service（需要仓储）
	m.PAT = authInfra.NewPATService(repos.PAT.Command, repos.PAT.Query, repos.User.Query, tokenGenerator)

	// TwoFA Service（需要仓储）
	m.TwoFA = twofa.NewService(repos.TwoFA.Command, repos.TwoFA.Query, repos.User.Query, cfg.Auth.TwoFAIssuer)

	return m
}
