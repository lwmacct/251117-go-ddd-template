package container

import (
	"go.uber.org/fx"

	"github.com/lwmacct/251117-go-ddd-template/internal/config"
	"github.com/lwmacct/251117-go-ddd-template/internal/domain/auth"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/twofa"

	infra_auth "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/auth"
	infra_captcha "github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/captcha"
)

// ServicesModule aggregates all domain and infrastructure services.
type ServicesModule struct {
	// Domain Services
	Auth auth.Service

	// Infrastructure Services
	JWT             *infra_auth.JWTManager
	TokenGenerator  auth.TokenGenerator
	LoginSession    *infra_auth.LoginSessionService
	PermissionCache *infra_auth.PermissionCacheService
	PAT             *infra_auth.PATService
	Captcha         *infra_captcha.Service
	TwoFA           *twofa.Service
}

// ServiceModule provides all domain and infrastructure services.
//
// Services handle business logic and technical concerns:
//   - Auth: password hashing, JWT token generation
//   - PermissionCache: user permission caching with cache-aside pattern
//   - TwoFA: TOTP-based two-factor authentication
//   - Captcha: image captcha generation
var ServiceModule = fx.Module("service",
	fx.Provide(
		// Infrastructure services
		newJWTManager,
		newTokenGenerator,
		newLoginSessionService,
		newAuthPermissionCacheService,
		newPATService,
		newCaptchaService,
		newTwoFAService,

		// Domain services
		newAuthService,

		// Aggregated module
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

// newServicesModule creates the aggregated services module.
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
