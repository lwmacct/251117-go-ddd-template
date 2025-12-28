---
paths:
  - "internal/container/**/*.go"
---

# Uber Fx 依赖注入规范

## 核心原则

1. **直接依赖** - 构造函数直接声明所需类型，不通过聚合模块间接访问
2. **精确注入** - 只注入需要的依赖，不多不少
3. **类型安全** - 使用 `fx.In`/`fx.Out` 保持类型安全

## 文件结构

```
internal/container/
├── infra.go       # 基础设施（DB、Redis、Telemetry）
├── cache.go       # 缓存服务
├── repository.go  # 仓储
├── service.go     # 领域服务 + 基础设施服务
├── usecase.go     # 用例处理器
├── http.go        # HTTP Handler + 路由
├── hooks.go       # 生命周期钩子
└── types.go       # 共享类型定义
```

## fx.Module 定义规范

```go
var CacheModule = fx.Module("cache",
    fx.Provide(
        // 直接使用包构造函数（简单情况）
        persistence.NewAuditLogRepositories,

        // 需要参数转换时使用包装函数
        newSettingCacheService,
    ),
)
```

## fx.Out 批量返回

**适用场景**：一个构造函数返回多个相关依赖

```go
type CacheServicesResult struct {
    fx.Out
    Setting         cache.SettingCacheService
    SettingCategory cache.SettingCategoryCacheService
    UserSetting     cache.UserSettingCacheService
    // ...
}

func NewAllCacheServices(client *redis.Client, cfg *config.Config) CacheServicesResult {
    prefix := cfg.Data.RedisKeyPrefix
    return CacheServicesResult{
        Setting:         infracache.NewSettingCacheService(client, prefix),
        SettingCategory: infracache.NewSettingCategoryCacheService(client, prefix),
        // ...
    }
}
```

## fx.In 聚合参数

**适用场景**：构造函数参数超过 5 个

```go
type authUseCasesParams struct {
    fx.In
    UserRepos    persistence.UserRepositories
    AuthSvc      auth.Service
    LoginSession *infra_auth.LoginSessionService
    TwoFASvc     *twofa.Service
    AuditLog     *AuditLogUseCases
}

func newAuthUseCases(p authUseCasesParams) *AuthUseCases {
    return &AuthUseCases{
        Login: auth.NewLoginHandler(p.UserRepos.Query, ...),
    }
}
```

## 禁止事项

### 禁止聚合模块传递

```go
// ❌ 禁止：通过聚合模块间接访问
func newAuthUseCases(repos *RepositoriesModule) *AuthUseCases {
    repos.User.Query  // 间接访问
}

// ✅ 正确：直接依赖所需类型
func newAuthUseCases(userRepos persistence.UserRepositories) *AuthUseCases {
    userRepos.Query  // 直接访问
}
```

### 禁止一行包装函数

```go
// ❌ 禁止：无意义的包装
func newAuditLogRepositories(db *gorm.DB) persistence.AuditLogRepositories {
    return persistence.NewAuditLogRepositories(db)
}

// ✅ 正确：直接使用
fx.Provide(persistence.NewAuditLogRepositories)
```

### 禁止聚合结构体

```go
// ❌ 禁止：聚合结构体
type RepositoriesModule struct {
    User    persistence.UserRepositories
    Role    persistence.RoleRepositories
    // ...
}

// ✅ 正确：让 Fx 直接注入各个仓储
```

## 允许的包装函数

当需要参数转换或额外逻辑时，允许使用包装函数：

```go
// ✅ 允许：需要从 config 提取 prefix
func newSettingCacheService(client *redis.Client, cfg *config.Config) cache.SettingCacheService {
    return infracache.NewSettingCacheService(client, cfg.Data.RedisKeyPrefix)
}

// ✅ 允许：需要缓存装饰器
func newUserRepositoriesWithCache(
    db *gorm.DB,
    userCache cache.UserWithRolesCacheService,
) persistence.UserRepositories {
    rawRepos := persistence.NewUserRepositories(db)
    cachedQuery := persistence.NewCachedUserQueryRepository(rawRepos.Query, userCache)
    return persistence.UserRepositories{
        Command: rawRepos.Command,
        Query:   cachedQuery,
    }
}
```

## 生命周期钩子

```go
fx.Invoke(func(lc fx.Lifecycle, db *gorm.DB) {
    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            return db.Ping()
        },
        OnStop: func(ctx context.Context) error {
            sqlDB, _ := db.DB()
            return sqlDB.Close()
        },
    })
})
```

## UseCase 子结构体

保留按领域分组的子结构体（有实际意义）：

```go
// ✅ 保留：按领域分组的处理器
type AuthUseCases struct {
    Login        *auth.LoginHandler
    Login2FA     *auth.Login2FAHandler
    Register     *auth.RegisterHandler
    RefreshToken *auth.RefreshTokenHandler
}

type UserUseCases struct {
    Create      *user.CreateHandler
    Update      *user.UpdateHandler
    Delete      *user.DeleteHandler
    // ...
}
```

## 模块依赖顺序

```
InfraModule (DB, Redis, Telemetry)
    ↓
CacheModule
    ↓
RepositoryModule
    ↓
ServiceModule
    ↓
UseCaseModule
    ↓
HTTPModule
    ↓
Hooks (EventHandlers, Warmup, HTTPServer)
```
