# 项目介绍

Go DDD Template 是一个基于领域驱动设计（DDD）和 CQRS 模式的企业级应用模板。

<!--TOC-->

## Table of Contents

- [项目背景](#项目背景) `:22+4`
- [核心价值](#核心价值) `:26+7`
- [适用场景](#适用场景) `:33+7`
- [技术特点](#技术特点) `:40+30`
  - [架构设计](#架构设计) `:42+7`
  - [身份认证](#身份认证) `:49+7`
  - [数据管理](#数据管理) `:56+7`
  - [开发体验](#开发体验) `:63+7`
- [项目结构](#项目结构) `:70+17`
- [社区贡献](#社区贡献) `:87+3`

<!--TOC-->

## 项目背景

本项目旨在提供一个生产就绪的 Go 应用模板，通过整合最佳实践和现代化的架构设计，帮助开发团队快速构建可维护、可扩展的企业级应用。

## 核心价值

- **架构清晰**: 四层架构分离关注点，职责明确
- **开发高效**: 丰富的工具链和自动化脚本
- **质量保证**: 完整的测试策略和代码规范
- **生产就绪**: Docker 容器化、健康检查、优雅关闭

## 适用场景

- 企业级 Web 应用开发
- RESTful API 服务
- 微服务架构基础
- DDD 实践学习项目

## 技术特点

### 架构设计

- DDD 四层架构（Adapters → Application → Domain ← Infrastructure）
- CQRS 读写分离模式
- Use Case Pattern 业务编排
- 富领域模型设计

### 身份认证

- JWT 短期令牌认证
- PAT 永久访问令牌
- RBAC 细粒度权限控制
- 审计日志追踪

### 数据管理

- PostgreSQL 主数据存储
- Redis 缓存优化
- GORM 对象关系映射
- 数据迁移管理

### 开发体验

- 热重载开发模式
- Swagger API 文档
- Task 任务自动化
- Docker Compose 本地环境

## 项目结构

```
251117-go-ddd-template/
├── internal/           # 核心业务代码
│   ├── adapters/      # 适配器层（HTTP Handler）
│   ├── application/   # 应用层（Use Case Handler）
│   ├── domain/        # 领域层（业务模型）
│   ├── infrastructure/# 基础设施层（技术实现）
│   └── bootstrap/     # 依赖注入容器
├── web/               # 前端应用
├── docs/              # VitePress 文档
├── configs/           # 配置文件
├── testing/           # 测试脚本
└── docker-compose.yml # 本地环境
```

## 社区贡献

欢迎提交 Issue 和 Pull Request，共同完善这个项目模板。
