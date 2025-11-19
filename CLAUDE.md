# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此仓库中工作时提供指导。

## 📋 项目概览

基于 Go 的 DDD (领域驱动设计) 模板应用，采用四层架构 + CQRS 模式，提供认证、RBAC 权限、审计日志等特性。Monorepo 结构包含后端(Go)、前端(Vue 3)、文档(VitePress)。

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

## 💻 添加新功能

### 标准开发流程（Use Case 模式）

#### 1. Domain 层定义

```go
// internal/domain/xxx/model.go
type Xxx struct {
    ID   uint
    Name string
}

// 业务行为方法（富领域模型）
func (x *Xxx) IsValid() bool { ... }
func (x *Xxx) Activate() { ... }

// internal/domain/xxx/command_repository.go
type CommandRepository interface {
    Create(ctx context.Context, entity *Xxx) error
    Update(ctx context.Context, entity *Xxx) error
    Delete(ctx context.Context, id uint) error
}

// internal/domain/xxx/query_repository.go
type QueryRepository interface {
    GetByID(ctx context.Context, id uint) (*Xxx, error)
    List(ctx context.Context, offset, limit int) ([]*Xxx, error)
    ExistsByName(ctx context.Context, name string) (bool, error)
}

// internal/domain/xxx/errors.go
var ErrXxxNotFound = errors.New("xxx not found")
```

#### 2. Infrastructure 层实现

```go
// internal/infrastructure/persistence/xxx_command_repository.go
type xxxCommandRepository struct { db *gorm.DB }
func NewXxxCommandRepository(db *gorm.DB) xxx.CommandRepository {
    return &xxxCommandRepository{db: db}
}
func (r *xxxCommandRepository) Create(ctx, entity) error { ... }

// internal/infrastructure/persistence/xxx_query_repository.go
type xxxQueryRepository struct { db *gorm.DB }
func NewXxxQueryRepository(db *gorm.DB) xxx.QueryRepository {
    return &xxxQueryRepository{db: db}
}
func (r *xxxQueryRepository) GetByID(ctx, id) (*xxx.Xxx, error) { ... }
```

#### 3. Application 层创建 Use Case

```go
// internal/application/xxx/command/create_xxx.go
type CreateXxxCommand struct {
    Name string
}

// internal/application/xxx/command/create_xxx_handler.go
type CreateXxxHandler struct {
    xxxCommandRepo xxx.CommandRepository
    xxxQueryRepo   xxx.QueryRepository
}

func (h *CreateXxxHandler) Handle(ctx context.Context, cmd CreateXxxCommand) (*CreateXxxResult, error) {
    // 1. 业务验证
    exists, _ := h.xxxQueryRepo.ExistsByName(ctx, cmd.Name)
    if exists {
        return nil, errors.New("name already exists")
    }

    // 2. 创建实体
    entity := &xxx.Xxx{Name: cmd.Name}
    h.xxxCommandRepo.Create(ctx, entity)

    return &CreateXxxResult{ID: entity.ID}, nil
}

// internal/application/xxx/query/get_xxx_handler.go
type GetXxxHandler struct {
    xxxQueryRepo xxx.QueryRepository
}
func (h *GetXxxHandler) Handle(ctx, query GetXxxQuery) (*XxxResponse, error) {
    return h.xxxQueryRepo.GetByID(ctx, query.ID)
}

// internal/application/xxx/dto.go
type XxxResponse struct {
    ID   uint   `json:"id"`
    Name string `json:"name"`
}
```

#### 4. Adapters 层创建 HTTP Handler

```go
// internal/adapters/http/handler/xxx_handler.go
type XxxHandler struct {
    createXxxHandler *command.CreateXxxHandler
    getXxxHandler    *query.GetXxxHandler
}

func (h *XxxHandler) Create(c *gin.Context) {
    var req CreateXxxRequest
    c.ShouldBindJSON(&req)

    // 调用 Use Case Handler
    result, err := h.createXxxHandler.Handle(c.Request.Context(), command.CreateXxxCommand{
        Name: req.Name,
    })

    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    c.JSON(201, gin.H{"message": "created", "data": result})
}
```

#### 5. Bootstrap 注册依赖

```go
// internal/bootstrap/container.go

// Repositories
xxxCommandRepo := persistence.NewXxxCommandRepository(db)
xxxQueryRepo := persistence.NewXxxQueryRepository(db)

// Use Case Handlers
createXxxHandler := command.NewCreateXxxHandler(xxxCommandRepo, xxxQueryRepo)
getXxxHandler := query.NewGetXxxHandler(xxxQueryRepo)

// HTTP Handler
xxxHandler := handler.NewXxxHandler(createXxxHandler, getXxxHandler)
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
9. **向前兼容** - 不需要考虑向后兼容，可以破坏现有功能

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

- `docs/architecture/ddd-cqrs.md` - DDD + CQRS 四层架构详解
- `docs/architecture/migration-guide.md` - 架构迁移指南和最佳实践
- `docs/architecture/overview.md` - 三层架构（遗留）

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
- ❌ Command 和 Query Repository 混用（写操作用 Command，读操作用 Query）
- ❌ 跳过 Use Case 直接从 Handler 调用 Repository

## 开发环境

- 当前系统环境为 ubuntu 22.04, 你可以使用 apt 安装任意软件包来完成工作
- 在完成每一个任务后进行 git commit -m "<COMMIT MESSAGE>" 来提交代码
- 环境中可能有多个 AI Agent 在工作，git commit 时不必在意其他被修改的文件
