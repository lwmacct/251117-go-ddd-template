# VitePress 多环境部署 - 快速参考

## 🎯 核心概念

同一份 VitePress 代码，支持两种部署环境：

| 环境 | base 路径 | 访问 URL | 构建命令 |
|------|-----------|----------|----------|
| 🏠 **Go API 服务器** | `/docs/` | `http://localhost:8080/docs/` | `npm run docs:build` |
| 🌐 **GitHub Pages** | `/251117-go-ddd-template/` | `https://用户名.github.io/251117-go-ddd-template/` | `npm run docs:build:github` |

## 🚀 快速开始

### 本地开发 + Go 服务器

```bash
# 1. 构建文档
npm run docs:build

# 2. 启动 Go 服务器
task go:run -- api

# 3. 访问
open http://localhost:8080/docs/
```

### GitHub Pages（自动部署）

```bash
# 1. 修改文档
vim docs/guide/getting-started.md

# 2. 提交推送
git add docs/
git commit -m "docs: update guide"
git push

# 3. 自动部署 ✨
# GitHub Actions 自动：
# - 检测 docs/** 变更
# - 运行 npm run docs:build:github
# - 部署到 GitHub Pages
```

## 📐 架构说明

### 1. VitePress 配置 (docs/.vitepress/config.ts:12)

```typescript
base: process.env.VITE_BASE_PATH || "/docs/"
```

- 默认: `/docs/`（Go 服务器）
- 环境变量可覆盖（GitHub Pages）

### 2. npm 脚本 (package.json)

```json
{
  "scripts": {
    "docs:build": "vitepress build docs",
    "docs:build:github": "VITE_BASE_PATH=/251117-go-ddd-template/ vitepress build docs"
  }
}
```

### 3. GitHub Actions (.github/workflows/deploy-docs.yml:54)

```yaml
- name: Build with VitePress for GitHub Pages
  run: npm run docs:build:github  # ← 使用 GitHub base
```

### 4. Go 路由 (internal/adapters/http/router.go:71-107)

```go
docs := r.Group("/docs")
docs.GET("/*filepath", handler)
```

## ✅ 验证

### Go 服务器

```bash
npm run docs:build
grep '/docs/assets' docs/.vitepress/dist/index.html
# ✅ href="/docs/assets/style.css"
```

### GitHub Pages

```bash
npm run docs:build:github
grep '/251117-go-ddd-template/assets' docs/.vitepress/dist/index.html
# ✅ href="/251117-go-ddd-template/assets/style.css"
```

## 🔧 启用 GitHub Pages

1. GitHub 仓库 → **Settings** → **Pages**
2. **Source**: `GitHub Actions`
3. 保存 → 完成！

## 📝 相关文件

| 文件 | 作用 |
|------|------|
| `docs/.vitepress/config.ts:12` | 支持环境变量配置 base |
| `package.json:12` | 两个构建脚本 |
| `.github/workflows/deploy-docs.yml:54` | 自动使用 GitHub base |
| `internal/adapters/http/router.go:71-107` | Go 服务器 /docs 路由 |

## 📚 详细文档

- **完整部署指南**: `VITEPRESS_DEPLOYMENT.md`
- **文档集成说明**: `DOCS_INTEGRATION.md`
- **VitePress 2.0 升级**: `VITEPRESS_2.0_UPGRADE.md`
