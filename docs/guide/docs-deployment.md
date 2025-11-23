# GitHub Pages 部署指南

本文档说明如何将 VitePress 文档部署到 GitHub Pages (**快速操作指南**) 。

> 💡 **技术细节**：如需了解多环境部署的技术原理和自动化方案，请查看 [VitePress 多环境部署指南](/development/deployment)

## 前提条件

- GitHub 仓库：`lwmacct/251117-go-ddd-template`
- 已配置 GitHub Actions workflow (`.github/workflows/deploy-docs.yml`)
- VitePress base 路径自动配置 (无需手动设置)

## 部署步骤

### 1. 在 GitHub 仓库中启用 GitHub Pages

1. 访问你的 GitHub 仓库：https://github.com/lwmacct/251117-go-ddd-template
2. 点击 **Settings** (设置)
3. 在左侧菜单中找到 **Pages**
4. 在 **Source** 下拉菜单中选择：
   - **Source**: GitHub Actions

   ![GitHub Pages 设置](https://docs.github.com/assets/cb-47267/images/help/pages/publishing-source-drop-down.png)

5. 点击 **Save** (保存)

### 2. 推送代码触发部署

GitHub Actions workflow 会在以下情况自动触发：

- 推送到 `main` 分支时
- 修改了 `docs/**` 目录下的文件
- 修改了 workflow 文件本身 (`.github/workflows/deploy-docs.yml`)

#### 首次部署

```bash
# 1. 添加所有文件
git add .

# 2. 提交
git commit -m "Add VitePress documentation with GitHub Pages deployment

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"

# 3. 推送到 main 分支
git push origin main
```

### 3. 查看部署状态

1. 访问仓库的 **Actions** 标签页
2. 查看 "Deploy VitePress Docs to Pages" workflow 的运行状态
3. 等待构建和部署完成 (通常需要 1-3 分钟)

### 4. 访问文档站点

部署成功后，文档将发布到：

**https://lwmacct.github.io/251117-go-ddd-template/**

## 手动触发部署

如果需要手动触发部署：

1. 访问仓库的 **Actions** 标签页
2. 选择 "Deploy VitePress Docs to Pages" workflow
3. 点击 **Run workflow** 按钮
4. 选择分支 (通常是 `main`)
5. 点击绿色的 **Run workflow** 按钮

## Workflow 说明

### 触发条件

```yaml
on:
  push:
    branches: [main]
    paths:
      - "docs/**" # 文档文件变更
      - ".github/workflows/deploy-docs.yml" # workflow 自身变更
  workflow_dispatch: # 手动触发
```

### 构建流程

1. **Checkout** - 检出代码 (包含完整历史记录)
2. **Setup Node** - 安装 Node.js (v20)
3. **Install dependencies** - 安装项目依赖 (使用 npm)
4. **Build** - 构建 VitePress 站点
5. **Upload artifact** - 上传构建产物

### 部署流程

1. **Deploy to GitHub Pages** - 将构建产物部署到 GitHub Pages

## 配置说明

### VitePress Base 路径

本项目使用**环境变量**自动管理不同部署环境的 base 路径：

```typescript
// docs/.vitepress/config.ts
export default defineConfig({
  // 自动适配：本地 /docs/，GitHub Pages /仓库名/
  base: process.env.BASE || "/docs/",
  // ...
});
```

**工作原理：**

- **本地开发/Go 服务器**：使用默认值 `/docs/`
- **GitHub Pages**：GitHub Actions 自动设置 `BASE` 为仓库名
- **零配置**：无需手动修改配置文件

> 📖 **详细技术说明**：查看 [VitePress 多环境部署指南](/development/deployment) 了解自动化方案的完整原理

### 权限配置

Workflow 需要以下权限：

```yaml
permissions:
  contents: read # 读取仓库内容
  pages: write # 写入 Pages
  id-token: write # 写入 ID Token (用于部署验证)
```

## 本地预览

在推送到 GitHub 之前，可以本地预览：

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 或构建并预览生产版本
npm run build
npm run preview
```

**注意：** 本地开发时不需要设置 base 路径，VitePress 会自动处理。

## 更新文档

更新文档非常简单：

1. 编辑 `docs/` 目录下的 Markdown 文件
2. 提交并推送到 `main` 分支
3. GitHub Actions 会自动构建和部署

```bash
# 编辑文档
vim docs/guide/getting-started.md

# 提交并推送
git add docs/guide/getting-started.md
git commit -m "Update getting started guide"
git push origin main
```

## 故障排查

### 构建失败

如果构建失败，检查：

1. **查看 Actions 日志**：在 GitHub Actions 标签页查看详细错误信息
2. **依赖问题**：确保 `package-lock.json` 已提交
3. **Markdown 语法**：检查是否有 Markdown 语法错误
4. **链接问题**：检查内部链接是否正确

### 页面 404

如果访问页面出现 404：

1. **检查 base 路径**：确保 `config.ts` 中的 `base` 设置正确
2. **检查 Pages 设置**：确保 GitHub Pages 已启用且 Source 设置为 "GitHub Actions"
3. **等待部署完成**：首次部署可能需要几分钟

### 样式或资源加载失败

如果样式或图片无法加载：

1. **检查 base 路径**：确保 `base` 配置正确
2. **使用相对路径**：在 Markdown 中使用相对路径引用资源
3. **静态资源**：将静态资源放在 `docs/public/` 目录

### 本地构建正常，但部署后有问题

1. **清除缓存**：在浏览器中清除缓存后重试
2. **检查 base 路径**：确保生产环境的 base 路径正确
3. **检查链接**：确保所有链接都是相对路径或包含 base 路径

## 自定义域名 (可选)

如果你想使用自定义域名：

1. 在 `docs/public/` 目录下创建 `CNAME` 文件
2. 在文件中写入你的域名 (如 `docs.example.com`)
3. 在 DNS 提供商处配置 CNAME 记录指向 `lwmacct.github.io`
4. 在 GitHub Pages 设置中验证域名

示例 `CNAME` 文件：

```
docs.example.com
```

## 进阶配置

### 添加自定义 404 页面

在 `docs/` 目录下创建 `404.md`：

```markdown
---
layout: page
---

# 页面未找到

抱歉，您访问的页面不存在。

[返回首页](/)
```

### 配置缓存

修改 workflow 以启用依赖缓存：

```yaml
- name: Setup Node
  uses: actions/setup-node@v4
  with:
    node-version: 20
    cache: npm # 已启用
```

### 部署预览环境

可以为 PR 创建预览环境：

```yaml
on:
  pull_request:
    branches: [main]
```

## 监控和维护

### 查看部署历史

1. 访问 **Actions** 标签页
2. 查看历史 workflow 运行记录
3. 每次运行都会显示构建时间、状态和日志

### 更新依赖

定期更新 VitePress 和相关依赖：

```bash
# 更新所有依赖到最新版本
npm update

# 或只更新 VitePress
npm update vitepress

# 提交更新
git add package.json package-lock.json
git commit -m "Update dependencies"
git push
```

## 相关资源

- [VitePress 官方文档](https://vitepress.dev/)
- [GitHub Pages 文档](https://docs.github.com/en/pages)
- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [本项目仓库](https://github.com/lwmacct/251117-go-ddd-template)
- [文档站点](https://lwmacct.github.io/251117-go-ddd-template/)

## 下一步

- ✅ 配置完成
- ⏳ 推送代码到 GitHub
- ⏳ 在 GitHub 仓库中启用 Pages
- ⏳ 等待首次部署完成
- ⏳ 访问并验证文档站点
