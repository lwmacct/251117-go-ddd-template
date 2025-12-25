本文件为 AI Agent 在此仓库中工作时提供指导。

## 📋 项目概览

基于 Go 的 DDD (领域驱动设计) 模板应用，采用四层架构 + CQRS 模式，提供认证、RBAC 权限、审计日志等特性。Monorepo 结构包含后端(Go)、前端(Vue 3)、文档(VitePress)。

## 🏗️ 核心架构

### 🔷 DDD 四层架构 + CQRS

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

### 📦 各层职责

**1. Domain 层**（不依赖任何外层）

- 定义领域模型（富模型，包含业务行为方法；**不得出现任何 GORM Tag 或 `gorm` 依赖**）
- 定义 Repository 接口（CommandRepository、QueryRepository）
- 定义 Domain Service 接口（领域能力，如密码验证、Token 生成）
- 定义领域错误

**2. Infrastructure 层**（实现 Domain 接口）

- 在 `internal/infrastructure/persistence` 中为每个模块定义 `*_model.go`（GORM Model + 映射函数）
- 仓储实现中使用持久化 Model 与数据库交互，并在进入/返回领域层时进行映射
- 实现 CommandRepository（GORM 写操作）
- 实现 QueryRepository（GORM 读操作，可优化为 Redis/ES）
- 如需在依赖注入处同时传递读写仓储，可额外提供 `{模块}_repositories.go` 将 Command/Query 聚合
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

| 层级               | 文件类型            | 命名规范                                                 | 示例                                                     |
| ------------------ | ------------------- | -------------------------------------------------------- | -------------------------------------------------------- |
| **Domain**         | 实体模型            | `entity_{模块}.go`（仅含业务字段/行为，不允许 GORM Tag） | `entity_user.go`, `entity_role.go`                       |
|                    | Repository 接口     | `command_repository.go` / `query_repository.go`          | 每个模块固定命名                                         |
|                    | 值对象              | `value_objects.go`                                       | 复杂领域需要时使用                                       |
|                    | 错误定义            | `errors.go`                                              | 每个模块的领域错误                                       |
| **Infrastructure** | 持久化 Model        | `{模块}_model.go`（含 GORM Tag、映射函数）               | `user_model.go`, `role_model.go`, `pat_model.go`         |
|                    | Repository 实现     | `{模块}_{操作类型}_repository.go`（入/出都映射 Domain）  | `user_command_repository.go`, `user_query_repository.go` |
|                    | 仓储聚合            | `{模块}_repositories.go`（组合读写仓储，便于一次性注入） | `user_repositories.go`, `auditlog_repositories.go`       |
|                    | Domain Service 实现 | `service.go`                                             | 在各自子目录（如 `auth/service.go`）                     |
| **Application**    | Command 定义        | `cmd_{操作}.go`（结构体必须以 `Command` 结尾）           | `cmd_login.go`, `cmd_create_user.go`                     |
|                    | Command Handler     | `cmd_{操作}_handler.go`                                  | `cmd_login_handler.go`, `cmd_create_user_handler.go`     |
|                    | Query 定义          | `qry_{操作}.go`（结构体必须以 `Query` 结尾）             | `qry_get_user.go`, `qry_list_users.go`                   |
|                    | Query Handler       | `qry_{操作}_handler.go`                                  | `qry_get_user_handler.go`, `qry_list_users_handler.go`   |
|                    | DTO 定义            | `dto.go`（结构体必须以 `DTO` 结尾，包括 `*ResultDTO`）   | `dto.go`                                                 |
|                    | Mapper              | `mapper.go`                                              | `mapper.go`                                              |
| **Adapters**       | HTTP Handler        | `{模块}.go`（单数）                                      | `user.go`, `role.go`, `menu.go`                          |

### 📝 Go Doc 规范

#### 语言选择

**统一使用中文**编写文档注释，与项目整体风格保持一致。

#### 包注释（doc.go）

每个 Domain 模块**必须**包含 `doc.go` 文件，格式如下：

```go
// Package user 定义用户领域模型和仓储接口。
//
// 本包是用户管理的领域层核心，定义了：
//   - [User]: 用户实体（富领域模型）
//   - [CommandRepository]: 写仓储接口
//   - [QueryRepository]: 读仓储接口
//   - 用户领域错误（见 errors.go）
//
// 依赖倒置：
// 本包仅定义接口，实现位于 infrastructure/persistence 包。
package user
```

**要点**：

- 首行以 `// Package xxx` 开头，简述包职责
- 使用 `[TypeName]` 语法链接到同包类型（Go 1.19+）
- 列出包内关键类型和职责

#### 类型注释

```go
// User 用户实体，包含用户基本信息和 RBAC 角色关联。
//
// 业务行为：
//   - [User.CanLogin]: 检查用户是否可登录
//   - [User.HasRole]: 检查用户是否拥有指定角色
type User struct { ... }
```

#### 方法注释

```go
// HasRole 检查用户是否拥有指定角色。
func (u *User) HasRole(roleName string) bool { ... }

// CanLogin 报告用户是否可以登录。
// 当用户状态为 "active" 时返回 true。
func (u *User) CanLogin() bool { ... }
```

**要点**：

- 首句以方法名开头，使用动词描述功能
- 布尔方法使用 "报告..." 或 "检查..." 开头
- 可附加参数说明、返回值含义、错误条件

#### Go 1.19+ 文档特性

| 特性         | 语法          | 示例                     |
| ------------ | ------------- | ------------------------ |
| **类型链接** | `[TypeName]`  | `参见 [User] 实体定义`   |
| **跨包链接** | `[pkg.Type]`  | `使用 [context.Context]` |
| **标题**     | `// # 标题`   | 需前后空行               |
| **列表**     | `//   - item` | 缩进 2-3 空格            |
| **代码块**   | 缩进 4 空格   | 不会被重新换行           |

**目录结构示例（以 user 模块为例）**：

```
internal/domain/user/
├── entity_user.go                 # User 实体/领域行为
├── command_repository.go          # User 写仓储接口
├── query_repository.go            # User 读仓储接口
└── errors.go                      # User 领域错误

internal/infrastructure/persistence/
├── user_model.go                  # GORM Model + 映射函数
├── user_command_repository.go     # 写仓储实现（入参/返回都映射 Domain）
├── user_query_repository.go       # 读仓储实现
└── user_repositories.go           # Command/Query 聚合（可选）

internal/application/user/
├── cmd_create_user.go             # CreateUserCommand
├── cmd_create_user_handler.go     # CreateUserHandler
├── cmd_update_user.go             # UpdateUserCommand
├── cmd_update_user_handler.go     # UpdateUserHandler
├── qry_get_user.go                # GetUserQuery
├── qry_get_user_handler.go        # GetUserHandler
├── qry_list_users.go              # ListUsersQuery
├── qry_list_users_handler.go      # ListUsersHandler
├── dto.go                         # CreateUserResultDTO, UserWithRolesDTO 等
├── mapper.go                      # Entity => DTO
└── doc.go                         # 包文档

internal/adapters/http/handler/
└── user.go                        # User Handler（仅绑定/响应）
```

## 💻 添加新功能

### 🔄 标准开发流程（Use Case 模式）

#### 1️⃣ Domain 层定义

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

#### 2️⃣ Infrastructure 层实现

**所有 Repository 实现统一在 `internal/infrastructure/persistence/` 目录，并通过 Model 进行映射**

```go
// internal/infrastructure/persistence/xxx_model.go
type XxxModel struct {
    ID   uint   `gorm:"primaryKey"`
    Name string `gorm:"size:100;not null"`
    // ...
}

func newXxxModelFromEntity(entity *xxx.Xxx) *XxxModel { ... }
func (m *XxxModel) toEntity() *xxx.Xxx { ... }

// internal/infrastructure/persistence/xxx_command_repository.go
type xxxCommandRepository struct { db *gorm.DB }

func (r *xxxCommandRepository) Create(ctx context.Context, entity *xxx.Xxx) error {
    model := newXxxModelFromEntity(entity)
    if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
        return err
    }
    if saved := model.toEntity(); saved != nil {
        *entity = *saved
    }
    return nil
}

// internal/infrastructure/persistence/xxx_query_repository.go
type xxxQueryRepository struct { db *gorm.DB }

func (r *xxxQueryRepository) GetByID(ctx context.Context, id uint) (*xxx.Xxx, error) {
    var model XxxModel
    if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
        return nil, err
    }
    return model.toEntity(), nil
}

// internal/infrastructure/persistence/xxx_repositories.go（可选）
// 将 Command/Query 聚合，方便容器一次性注入
type XxxRepositories struct {
    Command xxx.CommandRepository
    Query   xxx.QueryRepository
}

func NewXxxRepositories(db *gorm.DB) XxxRepositories {
    return XxxRepositories{
        Command: NewXxxCommandRepository(db),
        Query:   NewXxxQueryRepository(db),
    }
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

#### 3️⃣ Application 层创建 Use Case

**目录结构**（扁平化，使用前缀区分）：

```
internal/application/xxx/
├── cmd_create_xxx.go           # Command 定义（仅含 CreateXxxCommand）
├── cmd_create_xxx_handler.go   # Command Handler
├── cmd_update_xxx.go
├── cmd_update_xxx_handler.go
├── cmd_delete_xxx.go
├── cmd_delete_xxx_handler.go
├── qry_get_xxx.go              # Query 定义（仅含 GetXxxQuery）
├── qry_get_xxx_handler.go      # Query Handler
├── qry_list_xxx.go
├── qry_list_xxx_handler.go
├── dto.go                      # DTO 定义（包括 *ResultDTO）
├── mapper.go                   # Entity → DTO 映射函数
└── doc.go                      # 包文档
```

**命名强制规范**（pre-commit 检查）：

| 文件模式   | 结构体后缀要求 |
| ---------- | -------------- |
| `cmd_*.go` | 仅 `*Command`  |
| `qry_*.go` | 仅 `*Query`    |
| `dto.go`   | 仅 `*DTO`      |

**Command 定义和 Handler**：

```go
// internal/application/xxx/cmd_create_xxx.go
package xxx

// CreateXxxCommand 创建 Xxx 命令
type CreateXxxCommand struct {
    Name string
}

// internal/application/xxx/cmd_create_xxx_handler.go
package xxx

import (
    "context"
    "errors"
    domainXxx "your-project/internal/domain/xxx"
)

// CreateXxxHandler 创建 Xxx 命令处理器
type CreateXxxHandler struct {
    xxxCommandRepo domainXxx.CommandRepository
    xxxQueryRepo   domainXxx.QueryRepository
}

// NewCreateXxxHandler 创建处理器实例
func NewCreateXxxHandler(cmdRepo domainXxx.CommandRepository, queryRepo domainXxx.QueryRepository) *CreateXxxHandler {
    return &CreateXxxHandler{
        xxxCommandRepo: cmdRepo,
        xxxQueryRepo:   queryRepo,
    }
}

// Handle 处理创建命令
func (h *CreateXxxHandler) Handle(ctx context.Context, cmd CreateXxxCommand) (*CreateXxxResultDTO, error) {
    // 1. 业务验证
    exists, _ := h.xxxQueryRepo.ExistsByName(ctx, cmd.Name)
    if exists {
        return nil, errors.New("name already exists")
    }

    // 2. 创建领域实体
    entity := &domainXxx.Xxx{Name: cmd.Name}

    // 3. 调用 Command Repository
    if err := h.xxxCommandRepo.Create(ctx, entity); err != nil {
        return nil, err
    }

    return &CreateXxxResultDTO{ID: entity.ID}, nil
}
```

**Query 定义和 Handler**：

```go
// internal/application/xxx/qry_get_xxx.go
package xxx

// GetXxxQuery 获取 Xxx 查询
type GetXxxQuery struct {
    ID uint
}

// internal/application/xxx/qry_get_xxx_handler.go
package xxx

import (
    "context"
    domainXxx "your-project/internal/domain/xxx"
)

// GetXxxHandler 获取 Xxx 查询处理器
type GetXxxHandler struct {
    xxxQueryRepo domainXxx.QueryRepository
}

// NewGetXxxHandler 创建处理器实例
func NewGetXxxHandler(queryRepo domainXxx.QueryRepository) *GetXxxHandler {
    return &GetXxxHandler{xxxQueryRepo: queryRepo}
}

// Handle 处理查询
func (h *GetXxxHandler) Handle(ctx context.Context, query GetXxxQuery) (*domainXxx.Xxx, error) {
    return h.xxxQueryRepo.GetByID(ctx, query.ID)
}
```

**DTO 和 Mapper**：

```go
// internal/application/xxx/dto.go
package xxx

// CreateXxxDTO HTTP 创建请求
type CreateXxxDTO struct {
    Name string `json:"name" binding:"required"`
}

// UpdateXxxDTO HTTP 更新请求
type UpdateXxxDTO struct {
    Name string `json:"name"`
}

// CreateXxxResultDTO 创建结果（Handler 返回）
type CreateXxxResultDTO struct {
    ID uint `json:"id"`
}

// XxxResponseDTO HTTP 响应
type XxxResponseDTO struct {
    ID   uint   `json:"id"`
    Name string `json:"name"`
}

// internal/application/xxx/mapper.go
package xxx

import domainXxx "your-project/internal/domain/xxx"

// ToXxxResponseDTO 将领域实体转换为响应 DTO
func ToXxxResponseDTO(entity *domainXxx.Xxx) *XxxResponseDTO {
    return &XxxResponseDTO{
        ID:   entity.ID,
        Name: entity.Name,
    }
}
```

#### 4️⃣ Adapters 层创建 HTTP Handler

**文件位置**：`internal/adapters/http/handler/xxx.go`（使用单数命名）

```go

// Update 处理更新请求
func (h *XxxHandler) Update(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    var req xxx.UpdateXxxDTO
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid request", err)
        return
    }

    _, err := h.updateXxxHandler.Handle(c.Request.Context(), xxx.UpdateXxxCommand{
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

    err := h.deleteXxxHandler.Handle(c.Request.Context(), xxx.DeleteXxxCommand{
        ID: uint(id),
    })
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to delete", err)
        return
    }

    response.Success(c, http.StatusOK, "Deleted successfully", nil)
}
```

#### 5️⃣ Bootstrap 注册依赖

**在 `internal/bootstrap/container.go` 中按顺序注册**：

```go
// internal/bootstrap/container.go
package bootstrap

import (
    "your-project/internal/adapters/http/handler"
    "your-project/internal/application/xxx"
    domainXxx "your-project/internal/domain/xxx"
    "your-project/internal/infrastructure/persistence"
)

type Container struct {
    // ... 其他字段

    // Repositories
    XxxCommandRepo domainXxx.CommandRepository
    XxxQueryRepo   domainXxx.QueryRepository

    // Use Case Handlers（扁平化后直接引用）
    CreateXxxHandler *xxx.CreateXxxHandler
    UpdateXxxHandler *xxx.UpdateXxxHandler
    DeleteXxxHandler *xxx.DeleteXxxHandler
    GetXxxHandler    *xxx.GetXxxHandler
    ListXxxHandler   *xxx.ListXxxHandler

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
    c.CreateXxxHandler = xxx.NewCreateXxxHandler(c.XxxCommandRepo, c.XxxQueryRepo)
    c.UpdateXxxHandler = xxx.NewUpdateXxxHandler(c.XxxCommandRepo, c.XxxQueryRepo)
    c.DeleteXxxHandler = xxx.NewDeleteXxxHandler(c.XxxCommandRepo)
    c.GetXxxHandler = xxx.NewGetXxxHandler(c.XxxQueryRepo)
    c.ListXxxHandler = xxx.NewListXxxHandler(c.XxxQueryRepo)

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

## ⚠️ 核心原则

1. **依赖倒置** - Domain 层定义接口，Infrastructure 层实现，Application 层依赖接口
2. **领域纯度** - Domain 模型仅承载业务语义，不得引用 GORM 或其它 ORM Tag；Infra 通过 `*_model.go` 负责映射
3. **CQRS 分离** - 写操作用 CommandRepository，读操作用 QueryRepository
4. **Use Case 模式** - 业务逻辑在 Application 层的 Handler 中处理，HTTP Handler 只做入参/出参绑定
5. **富领域模型** - 业务行为通过方法体现（如 `entity.Activate()`），禁止直接修改结构体字段
6. **单一职责** - Handler 仅做 HTTP 转换，Use Case Handler 编排业务，Repository 访问数据
7. **依赖注入** - 所有依赖在 `container.go` 中注册
8. **统一响应** - HTTP 响应使用 `adapters/http/response` 包
9. **接口优先** - 先定义 Domain 接口，再实现 Infrastructure
10. **统一架构** - 所有模块必须遵循最新 DDD+CQRS 约定，发现旧式实现立即拆分重构，禁止新增兼容层

## 🔑 关键文件位置

- **依赖注入**: `internal/bootstrap/container.go`
- **路由定义**: `internal/adapters/http/router.go`
- **配置管理**: `internal/infrastructure/config/config.go`
- **数据库迁移**: `internal/infrastructure/database/migrations.go`

## 📚 项目文档

- `internal/adapters/http/docs/`：Swagger API 文档, 不需要修改，自动生成

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

### ✍️ 添加新的 Command（写操作）

1. Domain: 定义 `CommandRepository` 接口方法
2. Infrastructure: 实现该方法（GORM）
3. Application: 创建 `XxxCommand` + `XxxHandler`
4. Adapters: HTTP Handler 调用 Use Case Handler
5. Bootstrap: 注册 Handler

### 🔍 添加新的 Query（读操作）

1. Domain: 定义 `QueryRepository` 接口方法
2. Infrastructure: 实现该方法（GORM，可优化为 Redis）
3. Application: 创建 `XxxQuery` + `XxxHandler`
4. Adapters: HTTP Handler 调用 Query Handler
5. Bootstrap: 注册 Handler

### 🔧 添加 Domain Service（领域能力）

1. Domain: 定义 `Service` 接口（如 `auth.Service`）
2. Infrastructure: 实现接口（技术实现，如 BCrypt、JWT）
3. Application: Use Case Handler 依赖该接口
4. Bootstrap: 注册 Domain Service 实现

## 🚫 禁止操作

- ❌ 在 HTTP Handler 中编排业务逻辑或直接调用 Repository
- ❌ 在 Application 层直接依赖 Infrastructure 实现（只能依赖 Domain 接口）
- ❌ Domain 层 import 外层代码（禁止 `gorm`/Infra 依赖）
- ❌ Command 和 Query Repository 混用，或复用旧的 `repository.go`
- ❌ 跳过 Use Case，直接从 Handler 或 Infra 操作数据库
