# 项目架构

本项目采用领域驱动设计 (DDD) 和整洁架构原则，确保代码的可维护性、可测试性和可扩展性。

## 架构概览

项目分为以下几个主要层次：

```
internal/
├── commands/          # CLI 命令层 (入口点)
├── adapters/          # 适配器层 (HTTP/GraphQL/WebSocket 等)
├── application/       # 应用层 (Command/Query + Handler + DTO)
├── domain/            # 领域层 (核心业务逻辑与接口)
├── infrastructure/    # 基础设施层 (技术实现)
└── bootstrap/         # 引导层 (依赖注入)
```

## 分层说明

### 1. Commands 层 (CLI 命令)

位于 `internal/commands/`，负责定义应用的入口点。

```
commands/
├── api/              # API 服务器命令
│   └── api.go        # 启动 HTTP 服务器
├── migrate/          # 数据库迁移命令
│   └── migrate.go    # 迁移管理 (up/status/fresh)
├── seed/             # 数据库种子命令
│   └── seed.go       # 种子数据填充
└── worker/           # 后台任务处理器
    └── worker.go     # 队列任务处理
```

**职责：**

- 解析命令行参数
- 初始化应用容器
- 启动服务器或执行特定任务
- 处理信号和优雅关闭

**可用命令：**

- `api` - 启动 REST API 服务器
- `migrate` - 数据库迁移管理
- `seed` - 填充种子数据
- `worker` - 启动后台任务处理器

### 2. Adapters 层 (适配器)

位于 `internal/adapters/`，负责处理外部通信。

```
adapters/
└── http/
    ├── handler/          # HTTP 处理器
    │   ├── auth.go       # 认证接口
    │   ├── user.go       # 用户接口
    │   ├── health.go     # 健康检查
    │   └── cache.go      # 缓存操作
    ├── middleware/       # 中间件
    │   └── jwt.go        # JWT 认证中间件
    ├── router.go         # 路由配置
    └── server.go         # HTTP 服务器封装
```

**职责：**

- 接收和解析 HTTP 请求
- 调用 Application 层的 Use Case Handler
- 构造统一响应
- 中间件处理 (认证、日志等)

### 3. Application 层 (应用)

位于 `internal/application/`，围绕 Use Case 编排业务。

```
application/
└── user/
    ├── command/
    │   ├── create_user.go
    │   └── create_user_handler.go
    ├── query/
    │   ├── list_users.go
    │   └── list_users_handler.go
    ├── dto.go
    └── mapper.go
```

**职责：**

- 定义 Command/Query 数据对象
- 实现 Handler 编排（依赖 Domain 接口 + Domain Service）
- 定义 DTO 与 Mapper（Domain ↔ DTO）
- 进行跨领域的应用级校验

**特点：**

- 不直接访问数据库，仅依赖 Domain 定义的接口
- Command 使用 CommandRepository，Query 使用 QueryRepository
- DTO 和 HTTP 请求/响应映射只在这一层处理

### 4. Domain 层 (领域)

位于 `internal/domain/`，包含核心业务逻辑和规则。

```
domain/
└── user/                     # 用户领域
    ├── entity_user.go        # 用户模型
    ├── command_repository.go # 写接口（Create/Update/Delete...）
    └── query_repository.go   # 读接口（Get/List/Search...）
```

**职责：**

- 定义领域模型 (实体、值对象)，包含业务行为
- 定义 Repository 接口 (Command/Query) 与 Domain Service 接口
- 保持纯净：无 GORM Tag、无 JSON Tag、无基础设施依赖
- 暴露领域错误，供 Application 层处理

**特点：**

- 业务行为收敛在实体方法或 Domain Service
- 所有持久化细节 (索引、外键) 完全放在 Infrastructure
- DTO、HTTP 请求等概念不会出现在 Domain

### 5. Infrastructure 层 (基础设施)

位于 `internal/infrastructure/`，提供技术实现。

```
infrastructure/
├── auth/                 # 认证领域服务实现
├── config/               # 配置管理
├── database/             # 数据库引导、迁移
├── persistence/          # GORM 持久化模型 + 仓储
│   ├── user_model.go
│   ├── user_command_repository.go
│   └── user_query_repository.go
├── redis/                # Redis 集成
└── queue/                # 任务/消息
```

**职责：**

- 实现 Domain 定义的接口 (Repository + Domain Service)
- 定义 `*_model.go` GORM 模型并负责 Model ↔ Entity 映射
- 管理数据库和缓存连接、迁移、种子数据
- 集成外部系统（Redis、队列、第三方 API）
- 暴露技术层错误给应用层处理

### 6. Bootstrap 层 (引导)

位于 `internal/bootstrap/`，负责依赖注入和初始化。

```
bootstrap/
└── container.go      # 依赖注入容器
```

**职责：**

- 初始化所有依赖
- 配置依赖关系
- 提供统一的容器接口
- 管理资源生命周期
- 条件性执行数据库迁移

**Container 包含：**

- Config (配置)
- DB (数据库连接)
- RedisClient (Redis 客户端)
- UserRepository (用户仓储)
- JWTManager (JWT 管理器)
- AuthService (认证服务)
- Router (HTTP 路由器)

**ContainerOptions：**

```go
type ContainerOptions struct {
    AutoMigrate bool  // 是否自动执行数据库迁移
}
```

- 默认不自动迁移 (生产环境安全)
- 可通过配置 `data.auto_migrate` 开启 (开发环境便利)
- 通过 `GetAllModels()` 获取所有需要迁移的模型

## 数据流

### 请求处理流程

```
HTTP 请求
    ↓
[Router] 路由匹配
    ↓
[Middleware] JWT 验证 (如需要)
    ↓
[Handler] 解析请求
    ↓
[Domain Service] 业务逻辑 (可选)
    ↓
[Repository] 数据访问
    ↓
[Infrastructure] 数据库/Redis 操作
    ↓
[Handler] 构造响应
    ↓
HTTP 响应
```

### 依赖关系

```
Commands → Bootstrap → Infrastructure → Domain
         ↓                   ↑
    Adapters → Application → Domain
```

**依赖原则：**

- 外层依赖内层
- Domain 层不依赖任何外层
- Infrastructure 实现 Domain 定义的接口
- Application 层仅依赖 Domain 接口与领域服务
- Adapters 通过 Use Case Handler 调用 Application 层

## 设计模式

### 1. 仓储模式 (Repository Pattern)

**定义：**在 `internal/domain/user/command_repository.go` 与 `query_repository.go`

```go
type CommandRepository interface {
    Create(ctx context.Context, user *User) error
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id uint) error
}

type QueryRepository interface {
    GetByID(ctx context.Context, id uint) (*User, error)
    GetByUsernameWithRoles(ctx context.Context, username string) (*User, error)
    List(ctx context.Context, offset, limit int) ([]*User, error)
}
```

**实现：**在 `internal/infrastructure/persistence/user_command_repository.go` 和 `user_query_repository.go`

**优点：**

- 解耦业务逻辑和数据访问
- 易于测试 (可 Mock)
- 可以轻松切换数据源

### 2. 依赖注入 (Dependency Injection)

使用 `bootstrap.Container` 集中管理依赖：

```go
type Container struct {
    Config      *config.Config
    DB          *gorm.DB
    RedisClient *redis.Client

    // CQRS repositories
    UserCommandRepo user.CommandRepository
    UserQueryRepo   user.QueryRepository

    // Domain services
    AuthService auth.Service

    // Use Case handlers
    LoginHandler    *authcommand.LoginHandler
    RegisterHandler *authcommand.RegisterHandler

    // HTTP handlers
    AuthHandler *handler.AuthHandler
}
```

**优点：**

- 松耦合
- 易于测试
- 便于管理生命周期

### 3. 适配器模式 (Adapter Pattern)

HTTP 层作为适配器，将外部请求转换为领域操作：

```go
func (h *AuthHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.ValidationError(c, err.Error())
        return
    }

    result, err := h.registerHandler.Handle(c.Request.Context(), command.RegisterCommand{
        Username: req.Username,
        Email:    req.Email,
        Password: req.Password,
    })
    if err != nil {
        response.BadRequest(c, err.Error())
        return
    }

    response.Created(c, gin.H{
        "user_id":       result.UserID,
        "access_token":  result.AccessToken,
        "refresh_token": result.RefreshToken,
    })
}
```

## 扩展应用

### 添加新功能

1. **Domain 层**：在 `internal/domain/<module>/entity_<module>.go` 定义实体（含业务行为），更新 `command_repository.go` 与 `query_repository.go`。
2. **Infrastructure 层**：在 `internal/infrastructure/persistence/<module>_model.go` 定义 GORM 模型，并分别实现 Command/Query Repository。
3. **Application 层**：在 `internal/application/<module>/command|query` 创建 Command/Query + Handler，补充 `dto.go`、`mapper.go`。
4. **Bootstrap**：在 `internal/bootstrap/container.go` 注入新的 Repository/Handler。
5. **Adapters**：在 `internal/adapters/http/handler/<module>.go` 增加 HTTP Handler，并在 `router.go` 注册路由。

### 添加新的数据源

1. 在 `internal/domain/` 定义仓储接口
2. 在 `internal/infrastructure/` 实现接口
3. 在 `bootstrap.Container` 中注入

### 添加新的外部接口

1. 在 `internal/adapters/` 创建新的适配器 (如 gRPC、GraphQL)
2. 复用 Application + Domain + Infrastructure 层
3. 在 `internal/commands/` 添加新的命令

## 最佳实践

1. **保持 Domain 层纯净**：不要引入外部框架依赖
2. **使用接口**：通过接口定义契约，便于测试和替换
3. **错误处理**：使用自定义错误类型，提供清晰的错误信息
4. **依赖注入**：在 Container 中集中管理依赖
5. **单一职责**：每个层次只负责自己的职责
6. **测试优先**：为核心业务逻辑编写单元测试

## 下一步

- 了解[配置系统](/guide/configuration)
- 学习[认证授权](/backend/authentication)
- 探索 [PostgreSQL 集成](/backend/postgresql)
- 查看 [Redis 缓存](/backend/redis)
