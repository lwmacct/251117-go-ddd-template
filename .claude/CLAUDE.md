# 项目概览

基于 Go 的 DDD 模板应用，采用四层架构 + CQRS 模式。Monorepo 结构包含后端(Go)、前端(Vue 3)、文档(VitePress)。

## 架构概览

```
internal/
├── adapters/        # 适配器层 - HTTP Handler、中间件、路由
├── application/     # 应用层 - Use Cases (Command/Query Handler)
├── domain/          # 领域层 - 业务模型、Repository 接口
├── infrastructure/  # 基础设施层 - Repository 实现、数据库
├── container/       # 依赖注入容器
├── manualtest/      # 手动测试脚本
├── precommit/       # 预提交钩子脚本
└── commands/        # CLI 命令
```

> 详细规范见 `.claude/rules/backend/` 和 `.claude/rules/frontend/`，编辑对应目录时自动加载。
