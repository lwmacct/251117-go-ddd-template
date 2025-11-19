# 项目架构

本项目采用领域驱动设计 (DDD) 和整洁架构原则，确保代码的可维护性、可测试性和可扩展性。

## 架构概览

项目分为以下几个主要层次：

```
internal/
├── commands/          # CLI 命令层 (入口点)
├── adapters/         # 适配器层 (外部接口)
├── domain/           # 领域层 (核心业务逻辑)
├── infrastructure/   # 基础设施层 (技术实现)
├── bootstrap/        # 引导层 (依赖注入)
└── shared/          # 共享工具层
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
- 调用领域服务或仓储
- 构造 HTTP 响应
- 中间件处理 (认证、日志等)

### 3. Domain 层 (领域)

位于 `internal/domain/`，包含核心业务逻辑和规则。

```
domain/
└── user/             # 用户领域
    ├── model.go      # 用户模型和 DTO
    └── repository.go # 用户仓储接口
```

**职责：**

- 定义领域模型 (实体、值对象)
- 定义仓储接口
- 实现业务规则和验证
- 不依赖任何外部框架

**特点：**

- 用户模型包含业务逻辑 (如密码加密)
- DTO 用于数据传输
- 仓储接口定义数据访问契约

### 4. Infrastructure 层 (基础设施)

位于 `internal/infrastructure/`，提供技术实现。

```
infrastructure/
├── auth/             # 认证基础设施
│   ├── jwt.go        # JWT 管理器
│   └── service.go    # 认证服务
├── config/           # 配置管理
│   └── config.go     # Koanf 配置
├── database/         # 数据库
│   ├── connection.go # PostgreSQL 连接
│   ├── migrator.go   # 基础迁移器
│   ├── migration_manager.go  # 迁移管理器
│   ├── seeder.go     # 种子管理器
│   └── seeds/        # 种子数据
│       └── user_seeder.go  # 用户种子
├── persistence/      # 持久化
│   └── user_repository.go  # 用户仓储实现
├── queue/            # 队列系统
│   ├── redis_queue.go   # Redis 队列
│   └── processor.go     # 任务处理器
└── redis/            # Redis
    ├── client.go     # Redis 客户端
    └── cache_repository.go  # 缓存仓储
```

**职责：**

- 实现领域层定义的仓储接口
- 提供数据库连接和管理
- 提供外部服务集成
- 实现认证授权机制
- 提供队列和后台任务处理

### 5. Bootstrap 层 (引导)

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

### 6. Shared 层 (共享)

位于 `internal/shared/`，提供通用工具。

```
shared/
└── errors/           # 自定义错误类型
    └── errors.go
```

**职责：**

- 提供通用错误类型
- 提供工具函数
- 提供常量定义

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
         ↓
    Adapters → Domain
```

**依赖原则：**

- 外层依赖内层
- Domain 层不依赖任何外层
- Infrastructure 实现 Domain 定义的接口
- Adapters 通过接口调用 Domain

## 设计模式

### 1. 仓储模式 (Repository Pattern)

**定义：**在 `internal/domain/user/repository.go`

```go
type Repository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id uint) (*User, error)
    // ...
}
```

**实现：**在 `internal/infrastructure/persistence/user_repository.go`

**优点：**

- 解耦业务逻辑和数据访问
- 易于测试 (可 Mock)
- 可以轻松切换数据源

### 2. 依赖注入 (Dependency Injection)

使用 `bootstrap.Container` 集中管理依赖：

```go
type Container struct {
    Config         *config.Config
    DB             *database.Connection
    RedisClient    *redis.Client
    UserRepository domain.Repository
    JWTManager     *auth.JWTManager
    AuthService    *auth.Service
    Router         *http.Router
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
    // 解析请求
    var req RegisterRequest
    // 调用领域服务
    user, tokens, err := h.authService.Register(...)
    // 返回响应
    c.JSON(http.StatusOK, RegisterResponse{...})
}
```

## 扩展应用

### 添加新功能

1. **定义领域模型**：在 `internal/domain/<name>/model.go`
2. **定义仓储接口**：在 `internal/domain/<name>/repository.go`
3. **实现仓储**：在 `internal/infrastructure/persistence/`
4. **创建 Handler**：在 `internal/adapters/http/handler/`
5. **注册路由**：在 `internal/adapters/http/router.go`
6. **注入依赖**：在 `internal/bootstrap/container.go`

### 添加新的数据源

1. 在 `internal/domain/` 定义仓储接口
2. 在 `internal/infrastructure/` 实现接口
3. 在 `bootstrap.Container` 中注入

### 添加新的外部接口

1. 在 `internal/adapters/` 创建新的适配器 (如 gRPC、GraphQL)
2. 复用 Domain 层和 Infrastructure 层
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
- 学习[认证授权](/architecture/authentication)
- 探索 [PostgreSQL 集成](/architecture/postgresql)
- 查看 [Redis 缓存](/architecture/redis)
