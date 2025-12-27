---
paths:
  - "internal/bootstrap/**/*.go"
---

# Bootstrap 依赖注入规范

<!--TOC-->

## Table of Contents

- [Container 结构体字段顺序](#container-结构体字段顺序) `:19+29`
- [初始化顺序（严格遵守）](#初始化顺序严格遵守) `:48+36`
- [新增模块检查清单](#新增模块检查清单) `:84+10`
- [关键文件位置](#关键文件位置) `:94+4`

<!--TOC-->

## Container 结构体字段顺序

```go
type Container struct {
    // 1. 基础设施
    DB    *gorm.DB
    Redis *redis.Client

    // 2. Domain Services
    AuthService auth.Service

    // 3. Repositories（按模块分组）
    UserCommandRepo user.CommandRepository
    UserQueryRepo   user.QueryRepository
    RoleCommandRepo role.CommandRepository
    RoleQueryRepo   role.QueryRepository

    // 4. Use Case Handlers（按模块分组）
    CreateUserHandler *userapp.CreateUserHandler
    UpdateUserHandler *userapp.UpdateUserHandler
    GetUserHandler    *userapp.GetUserHandler
    ListUsersHandler  *userapp.ListUsersHandler

    // 5. HTTP Handlers
    UserHandler *handler.UserHandler
    RoleHandler *handler.RoleHandler
}
```

## 初始化顺序（严格遵守）

```go
func NewContainer(cfg *config.Config) (*Container, error) {
    c := &Container{}

    // 1️⃣ 初始化基础设施（数据库、Redis、外部服务）
    c.DB = initDatabase(cfg)
    c.Redis = initRedis(cfg)

    // 2️⃣ 创建 Domain Services
    c.AuthService = auth.NewAuthService(jwtManager)

    // 3️⃣ 创建 Repositories
    c.UserCommandRepo = persistence.NewUserCommandRepository(c.DB)
    c.UserQueryRepo = persistence.NewUserQueryRepository(c.DB)

    // 4️⃣ 创建 Use Case Handlers（依赖 Repositories + Domain Services）
    c.CreateUserHandler = userapp.NewCreateUserHandler(
        c.UserCommandRepo,
        c.UserQueryRepo,
        c.AuthService,
    )

    // 5️⃣ 创建 HTTP Handlers（依赖 Use Case Handlers）
    c.UserHandler = handler.NewUserHandler(
        c.CreateUserHandler,
        c.UpdateUserHandler,
        c.GetUserHandler,
        c.ListUsersHandler,
    )

    return c, nil
}
```

## 新增模块检查清单

添加新模块时，按以下顺序在 `container.go` 中注册：

1. [ ] 在 Container 结构体中添加字段（按类型分组）
2. [ ] 创建 Command/Query Repository
3. [ ] 创建 Use Case Handlers
4. [ ] 创建 HTTP Handler
5. [ ] 在 `router.go` 中注册路由

## 关键文件位置

- **依赖注入**: `internal/bootstrap/container.go`
- **路由定义**: `internal/adapters/http/router.go`
