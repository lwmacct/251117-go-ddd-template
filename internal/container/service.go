package container

import (
	"go.uber.org/fx"

	"github.com/lwmacct/251117-go-ddd-template/internal/config"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/twofa"

	infra_auth "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
	infra_captcha "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/captcha"
)

// ServicesModule 聚合所有领域服务和基础设施服务。
type ServicesModule struct {
	// 领域服务
	Auth auth.Service

	// 基础设施服务
	JWT             *infra_auth.JWTManager
	TokenGenerator  auth.TokenGenerator
	LoginSession    *infra_auth.LoginSessionService
	PermissionCache *infra_auth.PermissionCacheService
	PAT             *infra_auth.PATService
	Captcha         *infra_captcha.Service
	TwoFA           *twofa.Service
}

// ServiceModule 提供所有领域服务和基础设施服务。
//
// 服务处理业务逻辑和技术关注点：
//   - Auth: 密码哈希、JWT Token 生成
//   - PermissionCache: 用户权限缓存（Cache-Aside 模式）
//   - TwoFA: 基于 TOTP 的双因素认证
//   - Captcha: 图形验证码生成
var ServiceModule = fx.Module("service",
	fx.Provide(
		// 基础设施服务
		newJWTManager,
		newTokenGenerator,
		newLoginSessionService,
		newAuthPermissionCacheService,
		newPATService,
		newCaptchaService,
		newTwoFAService,

		// 领域服务
		newAuthService,

		// 聚合模块
		newServicesModule,
	),
)

func newJWTManager(cfg *config.Config) *infra_auth.JWTManager {
	return infra_auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.AccessTokenExpiry, cfg.JWT.RefreshTokenExpiry)
}

func newTokenGenerator() *infra_auth.TokenGenerator {
	return infra_auth.NewTokenGenerator()
}

func newLoginSessionService() *infra_auth.LoginSessionService {
	return infra_auth.NewLoginSessionService()
}

func newAuthPermissionCacheService(
	cache *CacheServicesModule,
	repos *RepositoriesModule,
) *infra_auth.PermissionCacheService {
	return infra_auth.NewPermissionCacheService(
		cache.Permission,
		cache.UserWithRoles,
		repos.User.Query,
	)
}

func newAuthService(jwt *infra_auth.JWTManager, tokenGen *infra_auth.TokenGenerator) auth.Service {
	passwordPolicy := auth.DefaultPasswordPolicy()
	return infra_auth.NewAuthService(jwt, tokenGen, passwordPolicy)
}

func newCaptchaService() *infra_captcha.Service {
	return infra_captcha.NewService()
}

func newPATService(repos *RepositoriesModule, tokenGen *infra_auth.TokenGenerator) *infra_auth.PATService {
	return infra_auth.NewPATService(repos.PAT.Command, repos.PAT.Query, tokenGen)
}

func newTwoFAService(cfg *config.Config, repos *RepositoriesModule) *twofa.Service {
	return twofa.NewService(repos.TwoFA.Command, repos.TwoFA.Query, repos.User.Query, cfg.Auth.TwoFAIssuer)
}

// newServicesModule 创建聚合的服务模块。
func newServicesModule(
	authSvc auth.Service,
	jwt *infra_auth.JWTManager,
	tokenGen *infra_auth.TokenGenerator,
	loginSession *infra_auth.LoginSessionService,
	permCache *infra_auth.PermissionCacheService,
	pat *infra_auth.PATService,
	captchaSvc *infra_captcha.Service,
	twofaSvc *twofa.Service,
) *ServicesModule {
	return &ServicesModule{
		Auth:            authSvc,
		JWT:             jwt,
		TokenGenerator:  tokenGen,
		LoginSession:    loginSession,
		PermissionCache: permCache,
		PAT:             pat,
		Captcha:         captchaSvc,
		TwoFA:           twofaSvc,
	}
}
