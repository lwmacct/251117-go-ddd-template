# 开发文档概览

本章节聚焦于**Go DDD Template** 的开发者体验：如何运行 VitePress 文档系统、如何与 Go API 服务器联动，以及如何扩展自定义组件。所有内容均以当前仓库的目录结构与代码实现为准，避免出现过期的实践。

## 仓库定位

| 路径        | 角色                    | 关键说明                                                                             |
| ----------- | ----------------------- | ------------------------------------------------------------------------------------ |
| `internal/` | Go DDD 四层 + CQRS 实现 | adapters → application → domain ← infrastructure，所有后端接口与依赖注入位于此目录。 |
| `docs/`     | VitePress 2.0 文档系统  | 包含 `package.json`、`.vitepress/` 主题扩展、自定义组件与 Markdown 内容。            |
| `web/`      | 前端 SPA                | 与 backend 通过 `/api` 对接，与本章节无直接耦合。                                    |

> 🌱 **单一事实来源**：架构细节以 `docs/architecture/*`（分层、数据平面、身份与访问控制）及源代码为准；如发现旧的三层描述或 `/backend/` 残留路径，需立即清理。

## 文档系统技术栈

- **框架**：VitePress `^2.0.0-alpha.13`
- **运行时**：Node.js `>= 20.19.0`（参见 `docs/package.json#engines`）
- **依赖**：`mermaid`、`medium-zoom`、Vue 3.5
- **自定义组件**：位于 `docs/.vitepress/theme/components/`（`ApiEndpoint.vue`、`FeatureCard.vue`、`Mermaid.vue`、`StepsGuide.vue`）
- **构建命令**：所有命令都在 `docs/` 目录下运行，例如 `npm --prefix docs run build`

## 常用命令

```bash
# 安装依赖（首次）
npm --prefix docs install

# 本地开发（VitePress Dev Server）
npm --prefix docs run dev -- --host

# 构建静态文件（默认 base=/docs/）
npm --prefix docs run build

# 预览生产构建
npm --prefix docs run preview

# Go API + 文档一体化启动
task go:run -- api  # 需先执行构建，输出位于 docs/.vitepress/dist
```

## 文档目录索引

| 文档                                         | 说明                                                               |
| -------------------------------------------- | ------------------------------------------------------------------ |
| [文档部署指南](./deployment.md)              | 统一的构建命令、Go 服务器与 GitHub Pages 的部署策略。              |
| [文档与 Go API 集成](./docs-integration.md)  | 解释 `internal/adapters/http/router.go` 如何代理 `/docs` 路由。    |
| [Mermaid 集成说明](./mermaid-integration.md) | 自定义 Markdown 渲染与 `Mermaid.vue` 组件的实现细节。              |
| [文档功能示例](./features.md)                | 常用自定义组件（ApiEndpoint、FeatureCard、StepsGuide）的组合示例。 |
| [主题与高级能力](./advanced.md)              | 主题覆写、Medium Zoom、UI 增强等高阶实践。                         |
| [升级记录](./upgrade.md)                     | VitePress 与依赖升级的记录及回归步骤。                             |

## 协作准则

1. **遵循 CLAUDE.md**：仅记录 DDD + CQRS 架构，避免提及已移除的三层模式。
2. **引用真实路径**：示例中的文件路径需存在，例如 `internal/infrastructure/config/config.go`。
3. **可验证命令**：命令必须直接可执行（`npm --prefix docs`, `task go:run -- api` 等）。
4. **文档即代码**：在 PR 中同步更新本文档和相关 Markdown，保持开发体验的一致性。
