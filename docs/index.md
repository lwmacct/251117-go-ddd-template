---
layout: home

hero:
  name: "Go DDD Template"
  text: "领域驱动设计应用模板"
  tagline: 基于 Go 的整洁架构 DDD 模板，快速构建可维护的企业级应用
  actions:
    - theme: brand
      text: 快速开始
      link: /guide/getting-started
    - theme: alt
      text: 查看 GitHub
      link: https://github.com/lwmacct/251117-bd-vmalert

features:
  - icon: 🏗️
    title: 整洁架构
    details: 遵循 DDD 原则，清晰的分层架构，职责分离明确，易于维护和扩展
  - icon: 🔐
    title: JWT 认证
    details: 完整的用户认证授权系统，支持 Token 刷新、密码加密、用户状态管理
  - icon: 🗄️
    title: PostgreSQL 集成
    details: GORM ORM 支持，自动迁移，连接池管理，软删除，分页查询
  - icon: ⚡
    title: Redis 缓存
    details: 高性能缓存系统，JSON 自动序列化，分布式锁，健康检查
  - icon: ⚙️
    title: 灵活配置
    details: Koanf 配置管理，多层优先级支持（默认值/文件/环境变量/CLI）
  - icon: 🚀
    title: 生产就绪
    details: Docker 支持，优雅关闭，健康检查，连接池管理，开发热重载
---

## 快速开始

```bash
# 克隆项目
git clone https://github.com/lwmacct/251117-bd-vmalert.git
cd 251117-bd-vmalert

# 启动数据库和 Redis
docker-compose up -d

# 运行应用
task go:run -- api

# 健康检查
curl http://localhost:8080/health
```

## 技术栈

- **框架**: Gin（HTTP 服务器）
- **数据库**: PostgreSQL + GORM
- **缓存**: Redis
- **认证**: JWT (golang-jwt/jwt/v5)
- **配置**: Koanf
- **CLI**: urfave/cli v3
- **容器**: Docker & Docker Compose

## 项目特性

- ✅ DDD 领域驱动设计
- ✅ 整洁架构分层
- ✅ 依赖注入容器
- ✅ 仓储模式
- ✅ 用户认证授权
- ✅ 数据库迁移
- ✅ Redis 缓存
- ✅ 健康检查
- ✅ 优雅关闭
- ✅ 开发热重载
