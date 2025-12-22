# 后端能力总览

面向后端开发者的功能手册，聚焦身份体系与访问控制实现细节。架构基础与数据平面已拆分到 `/architecture/`，这里专注于具体能力与使用方式。

<!--TOC-->

## Table of Contents

- [内容导航](#内容导航) `:16+6`
- [推荐阅读顺序](#推荐阅读顺序) `:22+6`
- [相关入口](#相关入口) `:28+6`
- [适用读者](#适用读者) `:34+5`

<!--TOC-->

## 内容导航

- **认证机制**：JWT 生命周期、刷新策略、中间件接入 → [查看](./identity-authentication.md)
- **RBAC 权限系统**：三段式权限格式、通配符匹配、审计日志 → [查看](./identity-rbac.md)
- **Personal Access Token**：长周期 Token、权限子集、最佳实践 → [查看](./identity-pat.md)

## 推荐阅读顺序

1. 先读 [认证机制](./identity-authentication.md) 理解令牌与中间件
2. 再读 [RBAC 权限系统](./identity-rbac.md) 明确授权模型
3. 最后看 [Personal Access Token](./identity-pat.md) 处理自动化和集成场景

## 相关入口

- **架构蓝图**：分层、依赖、演进 → [Architecture](/architecture/)
- **数据平面**：PostgreSQL / Redis 方案 → [Data Stack](/architecture/data-postgresql)
- **配置与部署**：运行参数、部署模式 → [Guide](/guide/)

## 适用读者

- 后端开发者：需要快速定位认证/授权代码与用法
- 平台工程师：希望落地接入、配置或扩展权限模型
- 审计/安全同事：评估 Token 与权限体系的落地情况
