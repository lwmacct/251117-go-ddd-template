// Package bootstrap 提供应用程序的依赖注入容器。
//
// 本包是 DDD+CQRS 架构的核心组装点，采用模块化设计：
//   - [InfrastructureModule]: 基础设施（数据库、Redis、事件总线）
//   - [RepositoriesModule]: CQRS 仓储（Command/Query 分离）
//   - [ServicesModule]: 领域服务和基础设施服务
//   - [UseCasesModule]: Use Case Handlers（业务逻辑编排）
//   - [HandlersModule]: HTTP Handlers（适配器层）
//
// 初始化顺序（显式依赖）：
//
//	Infra → Repos → Services → UseCases → Events → Handlers → Router
//
// 使用方式：
//
//	container, err := bootstrap.NewContainer(ctx, cfg, nil)
//	if err != nil { ... }
//	defer container.Close()
//	container.Router.Run(":8080")
package bootstrap

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/lwmacct/251117-go-ddd-template/internal/config"
	"github.com/lwmacct/251117-go-ddd-template/internal/infrastructure/persistence"
)

// ContainerOptions 容器初始化选项
type ContainerOptions struct {
	AutoMigrate bool // 是否自动执行数据库迁移 (仅开发环境推荐)
}

// DefaultOptions 返回默认选项
func DefaultOptions() *ContainerOptions {
	return &ContainerOptions{
		AutoMigrate: false, // 生产环境默认不自动迁移
	}
}

// Container DDD+CQRS 架构的模块化依赖注入容器
//
// 通过模块化设计，将原本 39 个扁平字段聚合为 6 个功能模块，
// 使依赖关系清晰、初始化顺序显式、测试更易于 mock。
type Container struct {
	Config *config.Config

	// 模块化依赖
	Infra    *InfrastructureModule // 基础设施：DB, Redis, EventBus
	Repos    *RepositoriesModule   // 仓储：所有 CQRS Repositories
	Services *ServicesModule       // 服务：Domain Services + Infrastructure Services
	UseCases *UseCasesModule       // 用例：所有 Use Case Handlers
	Handlers *HandlersModule       // HTTP Handlers：所有 HTTP Handlers

	Router *gin.Engine
}

// NewContainer 创建并初始化模块化依赖注入容器
//
// 初始化顺序（每步依赖前步）：
//  1. Infrastructure（DB, Redis, EventBus）
//  2. Repositories（依赖 DB）
//  3. Services（依赖 Repos, Redis）
//  4. UseCases（依赖 Repos, Services, EventBus）
//  5. EventHandlers（依赖 EventBus, Repos, Services）
//  6. Handlers（依赖 UseCases, Services）
//  7. Router（依赖 Handlers, Services）
func NewContainer(ctx context.Context, cfg *config.Config, opts *ContainerOptions) (*Container, error) {
	if opts == nil {
		opts = DefaultOptions()
	}

	c := &Container{Config: cfg}

	// 1. 基础设施
	var err error
	c.Infra, err = newInfrastructureModule(ctx, cfg, opts)
	if err != nil {
		return nil, err
	}

	// 2. 仓储
	c.Repos = newRepositoriesModule(c.Infra.DB)

	// 3. 服务
	c.Services = newServicesModule(cfg, c.Infra, c.Repos)

	// 4. 用例
	c.UseCases = newUseCasesModule(cfg, c.Infra, c.Repos, c.Services, c.Infra.EventBus)

	// 5. 事件处理器
	initEventHandlers(c.Infra.EventBus, c.Repos, c.Services)

	// 6. HTTP Handlers
	c.Handlers = newHandlersModule(cfg, c.Infra, c.UseCases)

	// 7. 路由
	c.Router = newRouter(cfg, c.Infra, c.Services, c.UseCases, c.Handlers)

	return c, nil
}

// Close 关闭容器中的所有资源
func (c *Container) Close() error {
	return c.Infra.Close()
}

// GetAllModels 返回所有需要迁移的领域模型
// 当添加新的领域模型时，需要在这里注册
func GetAllModels() []any {
	return []any{
		&persistence.UserModel{},
		&persistence.RoleModel{},
		&persistence.PermissionModel{},
		&persistence.PersonalAccessTokenModel{},
		&persistence.AuditLogModel{},
		&persistence.TwoFAModel{},
		&persistence.MenuModel{},
		&persistence.SettingModel{},
	}
}
