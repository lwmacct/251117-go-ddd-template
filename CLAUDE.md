# CLAUDE.md

本文件为 AI Agent 在此仓库中工作时提供指导。

## 📋 项目概览

基于 Go 的 DDD (领域驱动设计) 模板应用，采用四层架构 + CQRS 模式，提供认证、RBAC 权限、审计日志等特性。Monorepo 结构包含后端(Go)、前端(Vue 3)、文档(VitePress)。

> 🚫 仓库**仅**保留最新的 DDD + CQRS 架构。发现任何与此规范不符的遗留实现（如合并后的 Repository）时，应立即规划迁移，不允许在新功能中继续使用。

## 🏗️ 核心架构

### DDD 四层架构 + CQRS

```
internal/
├── adapters/        # 适配器层 - HTTP Handler、中间件、路由（仅做请求/响应转换）
├── application/     # 应用层 - Use Cases 业务编排（Command/Query Handler）
├── domain/          # 领域层 - 业务模型、Domain Service 接口、Repository 接口
├── infrastructure/  # 基础设施层 - Repository 实现、Domain Service 实现、数据库/Redis
├── bootstrap/       # 依赖注入容器
└── commands/        # CLI 命令
```

**依赖方向**: Adapters → Application → Domain ← Infrastructure (严格单向)

**CQRS 模式**:

- CommandRepository：写操作（Create, Update, Delete）
- QueryRepository：读操作（Get, List, Search, Count）

### 各层职责

**1. Domain 层**（不依赖任何外层）

- 定义领域模型（富模型，包含业务行为方法）
- 定义 Repository 接口（CommandRepository、QueryRepository）
- 定义 Domain Service 接口（领域能力，如密码验证、Token 生成）
- 定义领域错误

**2. Infrastructure 层**（实现 Domain 接口）

- 实现 CommandRepository（GORM 写操作）
- 实现 QueryRepository（GORM 读操作，可优化为 Redis/ES）
- 实现 Domain Service（技术实现，如 BCrypt、JWT）
- 数据库、Redis、外部 API

**3. Application 层**（业务编排）

- 定义 Command/Query（纯数据对象）
- 定义 Handler（协调 Domain Service 和 Repository 完成业务用例）
- 定义应用层 DTO

**4. Adapters 层**（接口适配）

- HTTP Handler：仅做请求绑定和响应转换
- 依赖 Application Use Case Handlers
- 不包含业务逻辑

### 📁 文件命名规范

| 层级               | 文件类型            | 命名规范                                        | 示例                                                     |
| ------------------ | ------------------- | ----------------------------------------------- | -------------------------------------------------------- |
| **Domain**         | 实体模型            | `entity_{模块}.go`                              | `entity_user.go`, `entity_role.go`                       |
|                    | Repository 接口     | `command_repository.go` / `query_repository.go` | 每个模块固定命名                                         |
|                    | 值对象              | `value_objects.go`                              | 复杂领域需要时使用                                       |
|                    | 错误定义            | `errors.go`                                     | 每个模块的领域错误                                       |
| **Infrastructure** | Repository 实现     | `{模块}_{操作类型}_repository.go`               | `user_command_repository.go`, `user_query_repository.go` |
|                    | Domain Service 实现 | `service.go`                                    | 在各自子目录（如 `auth/service.go`）                     |
| **Application**    | Command 定义        | `{操作}_xxx.go`                                 | `create_user.go`, `update_user.go`                       |
|                    | Command Handler     | `{操作}_xxx_handler.go`                         | `create_user_handler.go`                                 |
|                    | Query 定义          | `{操作}_xxx.go`                                 | `get_user.go`, `list_users.go`                           |
|                    | Query Handler       | `{操作}_xxx_handler.go`                         | `get_user_handler.go`                                    |
|                    | DTO 定义            | `dto.go`                                        | 模块根目录                                               |
|                    | Mapper 函数         | `mapper.go`                                     | 模块根目录                                               |
| **Adapters**       | HTTP Handler        | `{模块}.go`（单数）                             | `user.go`, `role.go`, `menu.go`                          |

**目录结构示例**：

```
internal/domain/user/
├── entity_user.go              # User 实体
├── command_repository.go       # 写操作接口
├── query_repository.go         # 读操作接口
└── errors.go                   # 领域错误

internal/infrastructure/persistence/
├── user_command_repository.go  # User 写操作实现
├── user_query_repository.go    # User 读操作实现
├── role_command_repository.go
├── role_query_repository.go
└── ...

internal/application/user/
├── command/
│   ├── create_user.go
│   ├── create_user_handler.go
│   ├── update_user.go
│   └── update_user_handler.go
├── query/
│   ├── get_user.go
│   ├── get_user_handler.go
│   ├── list_users.go
│   └── list_users_handler.go
├── dto.go                      # 所有 DTO
└── mapper.go                   # Entity → DTO 映射

internal/adapters/http/handler/
├── user.go                     # UserHandler
├── role.go                     # RoleHandler
└── menu.go                     # MenuHandler
```

> ⚠️ 任何组合仓储（如 `repository.go`）都被视为技术债务，新增或修改功能时应立即拆分出 `CommandRepository` 与 `QueryRepository`。

## 💻 添加新功能

### 标准开发流程（Use Case 模式）

#### 1. Domain 层定义

```go
// internal/domain/xxx/entity_xxx.go
// 实体文件使用 entity_ 前缀命名
type Xxx struct {
    ID   uint
    Name string
}

// 业务行为方法（富领域模型）
func (x *Xxx) IsValid() bool { ... }
func (x *Xxx) Activate() { ... }

// internal/domain/xxx/command_repository.go
// 写操作 Repository 接口
type CommandRepository interface {
    Create(ctx context.Context, entity *Xxx) error
    Update(ctx context.Context, entity *Xxx) error
    Delete(ctx context.Context, id uint) error
}

// internal/domain/xxx/query_repository.go
// 读操作 Repository 接口
type QueryRepository interface {
    GetByID(ctx context.Context, id uint) (*Xxx, error)
    List(ctx context.Context, offset, limit int) ([]*Xxx, error)
    ExistsByName(ctx context.Context, name string) (bool, error)
}

// internal/domain/xxx/errors.go
// 领域错误定义
var ErrXxxNotFound = errors.New("xxx not found")

// internal/domain/xxx/value_objects.go (可选)
// 复杂领域的值对象定义（如 pat、twofa 模块）
type XxxValueObject struct { ... }
```

#### 2. Infrastructure 层实现

**所有 Repository 实现统一在 `internal/infrastructure/persistence/` 目录**

```go
// internal/infrastructure/persistence/xxx_command_repository.go
// 命名规范：{模块}_{操作类型}_repository.go
type xxxCommandRepository struct { db *gorm.DB }

func NewXxxCommandRepository(db *gorm.DB) xxx.CommandRepository {
    return &xxxCommandRepository{db: db}
}

func (r *xxxCommandRepository) Create(ctx context.Context, entity *xxx.Xxx) error {
    return r.db.WithContext(ctx).Create(entity).Error
}

func (r *xxxCommandRepository) Update(ctx context.Context, entity *xxx.Xxx) error {
    return r.db.WithContext(ctx).Save(entity).Error
}

func (r *xxxCommandRepository) Delete(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Delete(&xxx.Xxx{}, id).Error
}

// internal/infrastructure/persistence/xxx_query_repository.go
// 读操作 Repository 实现
type xxxQueryRepository struct { db *gorm.DB }

func NewXxxQueryRepository(db *gorm.DB) xxx.QueryRepository {
    return &xxxQueryRepository{db: db}
}

func (r *xxxQueryRepository) GetByID(ctx context.Context, id uint) (*xxx.Xxx, error) {
    var entity xxx.Xxx
    err := r.db.WithContext(ctx).First(&entity, id).Error
    if err != nil {
        return nil, err
    }
    return &entity, nil
}

func (r *xxxQueryRepository) List(ctx context.Context, offset, limit int) ([]*xxx.Xxx, error) {
    var entities []*xxx.Xxx
    err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&entities).Error
    return entities, err
}

func (r *xxxQueryRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
    var count int64
    err := r.db.WithContext(ctx).Model(&xxx.Xxx{}).Where("name = ?", name).Count(&count).Error
    return count > 0, err
}
```

**Domain Service 实现示例**（如认证服务）：

```go
// internal/infrastructure/auth/service.go
// 实现 domain/auth.Service 接口
type authService struct {
    jwtManager *JWTManager
}

func NewAuthService(jwtManager *JWTManager) auth.Service {
    return &authService{jwtManager: jwtManager}
}

func (s *authService) HashPassword(password string) (string, error) { ... }
func (s *authService) VerifyPassword(hashedPassword, password string) error { ... }
func (s *authService) GenerateToken(userID uint) (string, error) { ... }
```

#### 3. Application 层创建 Use Case

**目录结构**：

```
internal/application/xxx/
├── command/              # 写操作 Use Cases
│   ├── create_xxx.go           # Command 定义
│   ├── create_xxx_handler.go   # Command Handler
│   ├── update_xxx.go
│   ├── update_xxx_handler.go
│   ├── delete_xxx.go
│   └── delete_xxx_handler.go
├── query/                # 读操作 Use Cases
│   ├── get_xxx.go              # Query 定义
│   ├── get_xxx_handler.go      # Query Handler
│   ├── list_xxx.go
│   └── list_xxx_handler.go
├── dto.go                # DTO 定义（请求/响应）
└── mapper.go             # Entity → DTO 映射函数
```

**Command 定义和 Handler**：

```go
// internal/application/xxx/command/create_xxx.go
package command

type CreateXxxCommand struct {
    Name string
}

type CreateXxxResult struct {
    ID uint
}

// internal/application/xxx/command/create_xxx_handler.go
package command

import (
    "context"
    "errors"
    "your-project/internal/domain/xxx"
)

type CreateXxxHandler struct {
    xxxCommandRepo xxx.CommandRepository
    xxxQueryRepo   xxx.QueryRepository
}

func NewCreateXxxHandler(cmdRepo xxx.CommandRepository, queryRepo xxx.QueryRepository) *CreateXxxHandler {
    return &CreateXxxHandler{
        xxxCommandRepo: cmdRepo,
        xxxQueryRepo:   queryRepo,
    }
}

func (h *CreateXxxHandler) Handle(ctx context.Context, cmd CreateXxxCommand) (*CreateXxxResult, error) {
    // 1. 业务验证
    exists, _ := h.xxxQueryRepo.ExistsByName(ctx, cmd.Name)
    if exists {
        return nil, errors.New("name already exists")
    }

    // 2. 创建领域实体
    entity := &xxx.Xxx{Name: cmd.Name}

    // 3. 调用 Command Repository
    if err := h.xxxCommandRepo.Create(ctx, entity); err != nil {
        return nil, err
    }

    return &CreateXxxResult{ID: entity.ID}, nil
}
```

**Query 定义和 Handler**：

```go
// internal/application/xxx/query/get_xxx.go
package query

type GetXxxQuery struct {
    ID uint
}

// internal/application/xxx/query/get_xxx_handler.go
package query

import (
    "context"
    "your-project/internal/domain/xxx"
)

type GetXxxHandler struct {
    xxxQueryRepo xxx.QueryRepository
}

func NewGetXxxHandler(queryRepo xxx.QueryRepository) *GetXxxHandler {
    return &GetXxxHandler{xxxQueryRepo: queryRepo}
}

func (h *GetXxxHandler) Handle(ctx context.Context, query GetXxxQuery) (*xxx.Xxx, error) {
    return h.xxxQueryRepo.GetByID(ctx, query.ID)
}
```

**DTO 和 Mapper**：

```go
// internal/application/xxx/dto.go
package xxx

type CreateXxxDTO struct {
    Name string `json:"name" binding:"required"`
}

type UpdateXxxDTO struct {
    Name string `json:"name"`
}

type XxxResponse struct {
    ID   uint   `json:"id"`
    Name string `json:"name"`
}

// internal/application/xxx/mapper.go
package xxx

import "your-project/internal/domain/xxx"

func ToXxxResponse(entity *xxx.Xxx) *XxxResponse {
    return &XxxResponse{
        ID:   entity.ID,
        Name: entity.Name,
    }
}
```

#### 4. Adapters 层创建 HTTP Handler

**文件位置**：`internal/adapters/http/handler/xxx.go`（使用单数命名）

```go
// internal/adapters/http/handler/xxx.go
package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "your-project/internal/adapters/http/response"
    "your-project/internal/application/xxx"
    "your-project/internal/application/xxx/command"
    "your-project/internal/application/xxx/query"
)

type XxxHandler struct {
    // Use Case Handlers（依赖注入）
    createXxxHandler *command.CreateXxxHandler
    updateXxxHandler *command.UpdateXxxHandler
    deleteXxxHandler *command.DeleteXxxHandler
    getXxxHandler    *query.GetXxxHandler
    listXxxHandler   *query.ListXxxHandler
}

func NewXxxHandler(
    createHandler *command.CreateXxxHandler,
    updateHandler *command.UpdateXxxHandler,
    deleteHandler *command.DeleteXxxHandler,
    getHandler *query.GetXxxHandler,
    listHandler *query.ListXxxHandler,
) *XxxHandler {
    return &XxxHandler{
        createXxxHandler: createHandler,
        updateXxxHandler: updateHandler,
        deleteXxxHandler: deleteHandler,
        getXxxHandler:    getHandler,
        listXxxHandler:   listHandler,
    }
}

// Create 处理创建请求
func (h *XxxHandler) Create(c *gin.Context) {
    var req xxx.CreateXxxDTO
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid request", err)
        return
    }

    // 调用 Use Case Handler
    result, err := h.createXxxHandler.Handle(c.Request.Context(), command.CreateXxxCommand{
        Name: req.Name,
    })
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to create", err)
        return
    }

    response.Success(c, http.StatusCreated, "Created successfully", result)
}

// GetByID 处理获取单个资源请求
func (h *XxxHandler) GetByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid ID", err)
        return
    }

    // 调用 Query Handler
    entity, err := h.getXxxHandler.Handle(c.Request.Context(), query.GetXxxQuery{
        ID: uint(id),
    })
    if err != nil {
        response.Error(c, http.StatusNotFound, "Not found", err)
        return
    }

    // 使用 Mapper 转换为 DTO
    resp := xxx.ToXxxResponse(entity)
    response.Success(c, http.StatusOK, "Success", resp)
}

// Update 处理更新请求
func (h *XxxHandler) Update(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    var req xxx.UpdateXxxDTO
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid request", err)
        return
    }

    _, err := h.updateXxxHandler.Handle(c.Request.Context(), command.UpdateXxxCommand{
        ID:   uint(id),
        Name: req.Name,
    })
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to update", err)
        return
    }

    response.Success(c, http.StatusOK, "Updated successfully", nil)
}

// Delete 处理删除请求
func (h *XxxHandler) Delete(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

    err := h.deleteXxxHandler.Handle(c.Request.Context(), command.DeleteXxxCommand{
        ID: uint(id),
    })
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to delete", err)
        return
    }

    response.Success(c, http.StatusOK, "Deleted successfully", nil)
}
```

#### 5. Bootstrap 注册依赖

**在 `internal/bootstrap/container.go` 中按顺序注册**：

```go
// internal/bootstrap/container.go
package bootstrap

import (
    "your-project/internal/adapters/http/handler"
    "your-project/internal/application/xxx/command"
    "your-project/internal/application/xxx/query"
    "your-project/internal/infrastructure/persistence"
)

type Container struct {
    // ... 其他字段

    // Repositories
    XxxCommandRepo xxx.CommandRepository
    XxxQueryRepo   xxx.QueryRepository

    // Use Case Handlers
    CreateXxxHandler *command.CreateXxxHandler
    UpdateXxxHandler *command.UpdateXxxHandler
    DeleteXxxHandler *command.DeleteXxxHandler
    GetXxxHandler    *query.GetXxxHandler
    ListXxxHandler   *query.ListXxxHandler

    // HTTP Handler
    XxxHandler *handler.XxxHandler
}

func NewContainer(cfg *config.Config) (*Container, error) {
    c := &Container{}

    // 1. 初始化数据库等基础设施
    db := initDatabase(cfg)

    // 2. 创建 Repositories
    c.XxxCommandRepo = persistence.NewXxxCommandRepository(db)
    c.XxxQueryRepo = persistence.NewXxxQueryRepository(db)

    // 3. 创建 Use Case Handlers
    c.CreateXxxHandler = command.NewCreateXxxHandler(c.XxxCommandRepo, c.XxxQueryRepo)
    c.UpdateXxxHandler = command.NewUpdateXxxHandler(c.XxxCommandRepo, c.XxxQueryRepo)
    c.DeleteXxxHandler = command.NewDeleteXxxHandler(c.XxxCommandRepo)
    c.GetXxxHandler = query.NewGetXxxHandler(c.XxxQueryRepo)
    c.ListXxxHandler = query.NewListXxxHandler(c.XxxQueryRepo)

    // 4. 创建 HTTP Handler
    c.XxxHandler = handler.NewXxxHandler(
        c.CreateXxxHandler,
        c.UpdateXxxHandler,
        c.DeleteXxxHandler,
        c.GetXxxHandler,
        c.ListXxxHandler,
    )

    return c, nil
}
```

> 🧠 实际 wiring 位于 `internal/bootstrap/container.go`。新增模块时务必遵循其中的顺序：先构建 Repository，再创建 Use Case Handler，最后初始化 HTTP Handler 并将其实例通过 `http.SetupRouter` 注册到路由层。

**在 `internal/adapters/http/router.go` 中注册路由**（容器会将依赖逐一传入，而不是直接传递 Container）：

```go
// internal/adapters/http/router.go
func SetupRouter(
    cfg *config.Config,
    db *gorm.DB,
    redisClient *redis.Client,
    userCommandRepo user.CommandRepository,
    userQueryRepo user.QueryRepository,
    auditLogCommandRepo auditlog.CommandRepository,
    captchaCommandRepo captcha.CommandRepository,
    jwtManager *infraauth.JWTManager,
    patService *infraauth.PATService,
    authService *infraauth.Service,
    captchaService *infracaptcha.Service,
    twofaService *infratwofa.Service,
    authHandler *handler.AuthHandler,
    roleHandler *handler.RoleHandler,
    menuHandler *handler.MenuHandler,
    settingHandler *handler.SettingHandler,
    patHandler *handler.PATHandler,
    auditLogHandler *handler.AuditLogHandler,
) *gin.Engine {
    r := gin.New()
    r.Use(gin.Recovery(), middleware.CORS())

    healthHandler := handler.NewHealthHandler(db, redisClient)
    r.GET("/health", healthHandler.Check)

    api := r.Group("/api")
    {
        captchaHandler := handler.NewCaptchaHandler(captchaCommandRepo, captchaService, cfg.Auth.DevSecret)
        auth := api.Group("/auth")
        {
            auth.POST("/register", authHandler.Register)
            auth.POST("/login", authHandler.Login)
            auth.POST("/refresh", authHandler.RefreshToken)
            auth.GET("/captcha", captchaHandler.GetCaptcha)
        }

        twofaHandler := handler.NewTwoFAHandler(twofaService)
        twofa := api.Group("/auth/2fa")
        twofa.Use(middleware.Auth(jwtManager, patService, userQueryRepo))
        {
            twofa.POST("/setup", twofaHandler.Setup)
            // ...
        }

        admin := api.Group("/admin")
        admin.Use(middleware.Auth(jwtManager, patService, userQueryRepo))
        admin.Use(middleware.AuditMiddleware(auditLogCommandRepo))
        admin.Use(middleware.RequireRole("admin"))
        {
            adminUserHandler := handler.NewAdminUserHandler(userCommandRepo, userQueryRepo)
            admin.POST("/users", middleware.RequirePermission("admin:users:create"), adminUserHandler.CreateUser)
            admin.GET("/roles", middleware.RequirePermission("admin:roles:read"), roleHandler.ListRoles)
            // ...
        }
    }

    return r
}
```

## ⚠️ 核心原则

1. **依赖倒置** - Domain 层定义接口，Infrastructure 层实现，Application 层依赖接口
2. **CQRS 分离** - 写操作用 CommandRepository，读操作用 QueryRepository
3. **Use Case 模式** - 业务逻辑在 Application 层的 Handler 中，不在 HTTP Handler
4. **富领域模型** - Domain 模型包含业务行为（`entity.Activate()` 而非 `entity.Status = "active"`）
5. **单一职责** - Handler 仅做 HTTP 转换，Use Case Handler 编排业务，Repository 访问数据
6. **依赖注入** - 所有依赖在 `container.go` 中注册
7. **统一响应** - HTTP 响应使用 `adapters/http/response` 包
8. **接口优先** - 先定义 Domain 接口，再实现 Infrastructure
9. **统一架构** - 所有模块必须遵循最新 DDD+CQRS 约定，发现旧式实现立即拆分重构，禁止新增兼容层

## 🔑 关键文件位置

- **依赖注入**: `internal/bootstrap/container.go`
- **路由定义**: `internal/adapters/http/router.go`
- **配置管理**: `internal/infrastructure/config/config.go`
- **数据库迁移**: `internal/infrastructure/database/migrations.go`

## 📚 项目文档

**VitePress 文档系统**（位于 `docs/` 目录）：

- 文档索引：`docs/.vitepress/config.ts`（定义所有可用文档页面）
- 架构文档：`docs/architecture/`
- API 文档：`docs/api/`
- 开发指南：`docs/development/`

**架构文档参考**：

- `docs/architecture/ddd-cqrs.md` - DDD + CQRS 四层架构详解（主架构标准）

**查看文档时**：

1. 先查 `docs/.vitepress/config.ts` 了解有哪些文档
2. 读取 `docs/architecture/` 下对应的 Markdown 文件
3. 架构变更时同步更新 VitePress 文档

## 🎯 常见任务

### 添加新的 Command（写操作）

1. Domain: 定义 `CommandRepository` 接口方法
2. Infrastructure: 实现该方法（GORM）
3. Application: 创建 `XxxCommand` + `XxxHandler`
4. Adapters: HTTP Handler 调用 Use Case Handler
5. Bootstrap: 注册 Handler

### 添加新的 Query（读操作）

1. Domain: 定义 `QueryRepository` 接口方法
2. Infrastructure: 实现该方法（GORM，可优化为 Redis）
3. Application: 创建 `XxxQuery` + `XxxHandler`
4. Adapters: HTTP Handler 调用 Query Handler
5. Bootstrap: 注册 Handler

### 添加 Domain Service（领域能力）

1. Domain: 定义 `Service` 接口（如 `auth.Service`）
2. Infrastructure: 实现接口（技术实现，如 BCrypt、JWT）
3. Application: Use Case Handler 依赖该接口
4. Bootstrap: 注册 Domain Service 实现

## 🚫 禁止操作

- ❌ 在 HTTP Handler 中写业务逻辑
- ❌ 在 Application 层直接依赖 Infrastructure 实现（只依赖 Domain 接口）
- ❌ 在 Domain 层依赖外层（Domain 不能 import Infrastructure/Application）
- ❌ Command 和 Query Repository 混用（写操作必须调用 CommandRepository，读操作必须调用 QueryRepository，遇到旧的组合接口要第一时间拆分）
- ❌ 跳过 Use Case 直接从 Handler 调用 Repository

## 开发环境

- 当前系统环境为 ubuntu 22.04, 你可以使用 apt 安装任意软件包来完成工作
- 你可以使用常用工具如 `ripgrep fd-find tree` 等来辅助你完成任务
- 在完成每一个任务后进行 git commit 来提交工作报告
- 环境中可能有多个 AI Agent 在工作，git commit 时不必在意其他被修改的文件
