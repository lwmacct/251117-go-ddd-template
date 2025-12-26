---
layout: home

hero:
  name: "Go DDD Template"
  text: "领域驱动设计应用模板"
  tagline: 基于 Go 的整洁架构 DDD 模板，快速构建可维护的企业级应用
  actions:
    - theme: brand
      text: 快速开始
      link: /getting-started
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
