# 开发文档

本节包含项目开发过程中的技术文档、升级记录和最佳实践。

## VitePress 文档系统

### 快速参考

- **[VitePress 快速参考](./quick-reference)** - VitePress 多环境部署的快速上手指南

### 详细指南

- **[VitePress 部署指南](./deployment)** - 完整的多环境部署文档
- **[文档集成说明](./docs-integration)** - Go API 服务器集成 VitePress 文档服务
- **[Mermaid 集成](./mermaid-integration)** - VitePress 中使用 Mermaid 图表的技术实现
- **[功能展示](./features)** - VitePress 2.0 原生功能完整示例

### 升级记录

- **[VitePress 2.0 升级](./upgrade)** - VitePress 2.0.0-alpha.13 升级总结

## 文档概览

### VitePress 快速参考

最快速了解 VitePress 多环境部署的方案：

- 🏠 本地 Go 服务器：`http://localhost:8080/docs/`
- 🌐 GitHub Pages：自动部署，零配置

**核心特性**：

- ✅ 零硬编码：仓库名自动从 GitHub 获取
- ✅ 统一命令：所有环境使用 `npm run docs:build`
- ✅ 可移植：Fork 项目后无需修改配置

### VitePress 部署指南

完整的部署文档，包括：

- 多环境配置方案
- GitHub Actions 自动部署
- 常见问题解答
- Windows 环境支持

### 文档集成说明

Go API 服务器如何集成 VitePress 文档：

- 路由配置 (`/docs/*`)
- 清洁 URL 支持
- SPA 路由回退
- 性能优化建议

### VitePress 2.0 升级

升级到 VitePress 2.0.0-alpha.13 的详细记录：

- 版本变更：1.6.4 → 2.0.0-alpha.13
- Vite 7 + Node.js 20.19+ 要求
- 新特性：CJK 支持、图片懒加载
- 破坏性变更和迁移指南

### Mermaid 集成

VitePress 中使用 Mermaid 图表的完整说明：

- 技术方案：markdown-it + Vue 组件
- 支持 10+ 种图表类型 (流程图、时序图、类图等)
- 自动主题切换 (亮色/暗色)
- 标准 Markdown 代码块语法

### 功能展示

VitePress 2.0 原生功能完整示例：

- Badge 徽章、代码高亮、代码差异
- 自定义容器、任务列表、Emoji
- 代码组、文件名显示、表格对齐
- 所有功能无需安装任何插件

### 高级功能

VitePress 高级功能和自定义组件：

- Medium Zoom 图片缩放
- 自定义 Vue 组件 (ApiEndpoint、FeatureCard、StepsGuide)
- 主题自定义 (品牌颜色、UI 增强)
- 实用组件库，可直接在文档中使用

## 相关链接

- [VitePress 官方文档](https://vitepress.dev/)
- [GitHub Actions 文档](https://docs.github.com/actions)
- [项目 GitHub 仓库](https://github.com/lwmacct/251117-go-ddd-template)
