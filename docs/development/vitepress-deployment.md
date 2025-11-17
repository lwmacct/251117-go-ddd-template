# VitePress 多环境部署指南

## 问题说明

VitePress 的 `base` 配置在不同部署环境中需要不同的值：

- **本地 Go 服务器**: `base: "/docs/"` - 文档通过 `http://localhost:8080/docs/` 访问
- **GitHub Pages**: `base: "/251117-go-ddd-template/"` - 仓库名作为路径

## 解决方案（推荐）

✅ **自动化方式**：GitHub Actions 自动获取仓库名

- 本地开发和 Go 服务器：使用 `npm run docs:build`（默认 base="/docs/"）
- GitHub Pages 部署：GitHub Actions 自动设置环境变量 `VITEPRESS_BASE`（值为仓库名）

### 工作原理

1. **VitePress 配置** 支持环境变量：

   ```typescript
   // docs/.vitepress/config.ts:12
   base: process.env.VITEPRESS_BASE || "/docs/",
   ```

2. **统一的构建命令**：

   ```json
   {
     "scripts": {
       "docs:build": "vitepress build docs" // 唯一的构建命令
     }
   }
   ```

3. **GitHub Actions** 自动设置仓库名：
   ```yaml
   # .github/workflows/deploy-docs.yml:54-56
   - name: Build with VitePress for GitHub Pages
     env:
       # 使用 GitHub Actions 上下文自动获取仓库名
       VITEPRESS_BASE: /仓库名/
     run: npm run docs:build
   ```

### 优势

- ✅ **零硬编码**：仓库名自动从 GitHub 获取，无需手动配置
- ✅ **统一命令**：本地和 CI 使用相同的 `npm run docs:build`
- ✅ **可移植性**：Fork 项目后无需修改任何配置
- ✅ **开发者友好**：本地开发只需 `npm run docs:build`
- ✅ **自动化**：推送代码后自动部署到 GitHub Pages
- ✅ **不易出错**：环境自动匹配，不会混淆

## 使用方法

### 1. 本地开发（默认）

```bash
# 开发服务器
npm run docs:dev
# 访问 http://localhost:5173

# 构建（用于 Go 服务器）
npm run docs:build
# 生成到 docs/.vitepress/dist/，base="/docs/"
```

### 2. GitHub Pages 部署（自动）

**无需特殊命令！** 推送代码后 GitHub Actions 自动处理：

```bash
# 推送代码
git push origin main

# GitHub Actions 自动：
# 1. 检测到 docs/** 变更
# 2. 设置 VITEPRESS_BASE=/仓库名/
# 3. 运行 npm run docs:build
# 4. 部署到 GitHub Pages
```

### 3. 本地测试 GitHub Pages 构建

如果需要在本地测试 GitHub Pages 的构建结果：

```bash
# 使用仓库名
VITEPRESS_BASE=/your-repo-name/ npm run docs:build

# 检查输出
cat docs/.vitepress/dist/index.html | grep '<base'
# 应该看到: <base href="/your-repo-name/">
```

## 部署流程

### 部署到 Go API 服务器

```bash
# 1. 构建文档（默认 base="/docs/"）
npm run docs:build

# 2. 构建 Go 应用
task go:build

# 3. 启动服务器
.local/bin/go-ddd-template api

# 4. 访问文档
open http://localhost:8080/docs/
```

### 部署到 GitHub Pages（自动化）✨

**推荐方式：推送代码自动部署**

```bash
# 1. 修改文档
vim docs/guide/getting-started.md

# 2. 提交并推送到 main 分支
git add docs/
git commit -m "docs: update getting started guide"
git push origin main

# 3. GitHub Actions 自动触发
# - 检测到 docs/** 变更
# - 自动运行 npm run docs:build:github
# - 自动部署到 GitHub Pages

# 4. 访问部署的文档（几分钟后）
open https://你的用户名.github.io/251117-go-ddd-template/
```

**手动方式（不推荐）：**

```bash
# 1. 使用 GitHub base 构建
npm run docs:build:github

# 2. 手动部署（通常不需要，GitHub Actions 会自动处理）
# ...
```

## GitHub Actions 自动部署

✅ **已配置完成**：`.github/workflows/deploy-docs.yml`

### 工作流配置说明

当前配置会在以下情况触发：

- 推送到 `main` 分支
- 修改了 `docs/**` 目录下的文件
- 修改了 `package.json` 或 `package-lock.json`
- 修改了工作流文件本身

关键配置：

```yaml
# .github/workflows/deploy-docs.yml
jobs:
  build:
    steps:
      - name: Build with VitePress for GitHub Pages
        env:
          # GitHub Actions 上下文：自动获取仓库名
          VITEPRESS_BASE: /仓库名/
        run: npm run docs:build # ← 使用统一的构建命令

      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: docs/.vitepress/dist

  deploy:
    needs: build
    steps:
      - name: Deploy to GitHub Pages
        uses: actions/deploy-pages@v4
```

**关键改进**：

- ✅ 使用 GitHub Actions 上下文自动获取仓库名
- ✅ 统一使用 `npm run docs:build` 命令
- ✅ Fork 项目后无需修改任何配置

### 启用 GitHub Pages

1. 进入仓库设置：`Settings` → `Pages`
2. **Source** 选择：`GitHub Actions`
3. 保存后，推送代码即可自动部署

### 查看部署状态

- GitHub 仓库 → `Actions` 标签
- 查看 "Deploy VitePress Docs to Pages" 工作流
- 绿色勾选表示部署成功

### 手动触发部署

```bash
# 在 GitHub 网页上：
# Actions → Deploy VitePress Docs to Pages → Run workflow → Run workflow
```

## package.json 脚本说明

```json
{
  "scripts": {
    "docs:dev": "vitepress dev docs",
    "docs:build": "vitepress build docs",
    "docs:preview": "vitepress preview docs"
  }
}
```

| 脚本           | base 路径                  | 用途                                    |
| -------------- | -------------------------- | --------------------------------------- |
| `docs:dev`     | `/docs/`                   | 本地开发服务器                          |
| `docs:build`   | `/docs/`（默认）或环境变量 | 本地构建 & GitHub Pages（通过 Actions） |
| `docs:preview` | 根据上次构建               | 预览构建结果                            |

**简化优势**：

- ✅ 只需一个构建命令 `docs:build`
- ✅ base 路径通过环境变量自动配置
- ✅ 本地和 CI 使用相同的构建逻辑

## 验证部署

### 验证 Go 服务器部署

```bash
# 构建
npm run docs:build
task go:build

# 启动
.local/bin/go-ddd-template api

# 测试
curl http://localhost:8080/docs/ | grep '<base'
# 应该包含: <base href="/docs/">

# 浏览器访问
open http://localhost:8080/docs/
```

### 验证 GitHub Pages 部署

```bash
# 本地模拟 GitHub Pages 构建
VITEPRESS_BASE=/your-repo-name/ npm run docs:build

# 检查输出文件
cat docs/.vitepress/dist/index.html | grep '<base'
# 应该包含: <base href="/your-repo-name/">

# 部署后访问（GitHub Actions 会自动使用正确的仓库名）
open https://你的用户名.github.io/你的仓库名/
```

## 常见问题

### Q: 为什么资源加载失败？

A: 检查 base 路径是否正确：

```bash
# 查看构建后的 HTML
cat docs/.vitepress/dist/index.html | grep -E '(href|src)='
```

所有资源路径应该以 base 路径开头（如 `/docs/assets/` 或 `/251117-go-ddd-template/assets/`）。

### Q: 可以同时支持两个环境吗？

A: 不能在同一次构建中同时支持。需要针对不同环境分别构建：

- Go 服务器使用: `npm run docs:build`
- GitHub Pages 使用: `npm run docs:build:github`

### Q: 开发时如何测试 GitHub Pages 的路径？

A: 使用环境变量临时设置：

```bash
VITEPRESS_BASE=/your-repo-name/ npm run docs:dev
```

### Q: Windows 环境如何设置环境变量？

A: 使用 cross-env 或直接在 PowerShell 中设置：

```bash
# 方式 1: 使用 cross-env
npm install -D cross-env

# 临时在 package.json 添加测试脚本:
# "docs:build:test": "cross-env VITEPRESS_BASE=/test-repo/ vitepress build docs"

# 方式 2: PowerShell
$env:VITEPRESS_BASE="/your-repo-name/"; npm run docs:build
```

## 总结

| 环境         | base 路径          | 构建命令                                | 访问 URL                           |
| ------------ | ------------------ | --------------------------------------- | ---------------------------------- |
| 本地开发     | `/docs/`           | `npm run docs:dev`                      | `http://localhost:5173`            |
| Go 服务器    | `/docs/`           | `npm run docs:build`                    | `http://localhost:8080/docs/`      |
| GitHub Pages | `/仓库名/`（自动） | `npm run docs:build`（由 Actions 调用） | `https://用户名.github.io/仓库名/` |

**核心优势**：

- ✅ **零硬编码**: 仓库名自动从 GitHub 获取
- ✅ **统一命令**: 所有环境使用相同的 `npm run docs:build`
- ✅ **可移植**: Fork 项目后无需修改配置
- ✅ **自动化**: 推送代码自动部署到 GitHub Pages

## 相关文件

- VitePress 配置: `docs/.vitepress/config.ts:12` - 支持 `VITEPRESS_BASE` 环境变量
- package.json 脚本: `package.json:11` - 统一的 `docs:build` 命令
- GitHub Actions: `.github/workflows/deploy-docs.yml:54-56` - 自动获取仓库名
- Go 路由配置: `internal/adapters/http/router.go:71-107` - `/docs` 路由实现
