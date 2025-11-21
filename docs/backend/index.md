# 后端能力总览

面向后端开发者的功能手册，聚焦身份体系与访问控制实现细节。架构基础与数据平面已拆分到 `/architecture/`，这里专注于具体能力与使用方式。

## 内容导航

- **认证机制**：JWT 生命周期、刷新策略、中间件接入 → [查看](/backend/authentication)
- **RBAC 权限系统**：三段式权限格式、通配符匹配、审计日志 → [查看](/backend/rbac)
- **Personal Access Token**：长周期 Token、权限子集、最佳实践 → [查看](/backend/pat)

## 推荐阅读顺序

1. 先读 [认证机制](/backend/authentication) 理解令牌与中间件
2. 再读 [RBAC 权限系统](/backend/rbac) 明确授权模型
3. 最后看 [Personal Access Token](/backend/pat) 处理自动化和集成场景

## 相关入口

- **架构蓝图**：分层、依赖、演进 → [Architecture](/architecture/)
- **数据平面**：PostgreSQL / Redis 方案 → [Data Stack](/architecture/data-postgresql)
- **配置与部署**：运行参数、部署模式 → [Guide](/guide/)

## 适用读者

- 后端开发者：需要快速定位认证/授权代码与用法
- 平台工程师：希望落地接入、配置或扩展权限模型
- 审计/安全同事：评估 Token 与权限体系的落地情况
