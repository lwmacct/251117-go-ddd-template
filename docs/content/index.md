---
layout: home

hero:
  name: "Go DDD Template"
  text: "领域驱动设计应用模板"
  tagline: 基于 Go 的整洁架构 DDD 模板，快速构建可维护的企业级应用
  actions:
    - theme: brand
      text: 快速开始
      link: /guide/quick-start
    - theme: alt
      text: 查看 GitHub
      link: https://github.com/lwmacct/251117-go-ddd-template

features:
  - icon: 🏗️
    title: DDD 四层架构 + CQRS
    details: 完整实现领域驱动设计，读写分离的 CQRS 模式，清晰的分层架构和职责分离
  - icon: 🔐
    title: JWT 认证 + PAT
    details: 完整的用户认证授权系统，支持 JWT Token 刷新、PAT 永久令牌、密码加密、用户状态管理
  - icon: 🗄️
    title: PostgreSQL 集成
    details: GORM ORM 支持，自动迁移，连接池管理，软删除，分页查询
  - icon: ⚡
    title: Redis 缓存
    details: 高性能缓存系统，JSON 自动序列化，分布式锁，健康检查
  - icon: ⚙️
    title: 灵活配置
    details: Koanf 配置管理，多层优先级支持 (默认值/文件/环境变量/CLI)
  - icon: 🚀
    title: 生产就绪
    details: Docker 支持，优雅关闭，健康检查，连接池管理，开发热重载
---

## 快速开始

```bash
# 克隆项目
git clone https://github.com/lwmacct/251117-go-ddd-template.git
cd 251117-go-ddd-template

# 启动数据库和 Redis
docker-compose up -d

# 运行应用
task go:run -- api

# 健康检查
curl http://localhost:8080/health
```

## 技术栈

- **框架**: Gin (HTTP 服务器)
- **数据库**: PostgreSQL + GORM
- **缓存**: Redis
- **认证**: JWT (golang-jwt/jwt/v5)
- **配置**: Koanf
- **CLI**: urfave/cli v3
- **容器**: Docker & Docker Compose

## 项目特性

- ✅ DDD 四层架构 (Adapters → Application → Domain ← Infrastructure)
- ✅ CQRS 模式 (CommandRepository / QueryRepository)
- ✅ Use Case Pattern (业务编排集中在 Application 层)
- ✅ 富领域模型 (业务逻辑封装在 Domain 实体中)
- ✅ 依赖注入容器
- ✅ 用户认证授权 (JWT + PAT 双重认证)
- ✅ RBAC 权限系统 (三段式细粒度权限)
- ✅ 数据库迁移 (PostgreSQL + GORM)
- ✅ Redis 缓存 (查询优化 + 分布式锁)
- ✅ 审计日志系统
- ✅ 健康检查
- ✅ 优雅关闭
- ✅ 开发热重载
