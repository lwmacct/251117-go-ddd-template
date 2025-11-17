# VitePress 文档部署 - 快速开始

## ✅ 已完成的配置

- [x] VitePress 项目结构（`docs/` 目录）
- [x] VitePress 配置文件（`docs/.vitepress/config.ts`）
- [x] GitHub Actions workflow（`.github/workflows/deploy-docs.yml`）
- [x] Base 路径配置（`/251117-bd-vmalert/`）
- [x] 文档内容（首页、指南、API 参考）
- [x] npm 依赖安装
- [x] 构建测试通过 ✓
- [x] `.gitignore` 更新

## 📦 项目文件

```
/apps/data/workspace/251117-go-ddd-template/
├── .github/workflows/
│   └── deploy-docs.yml          # GitHub Actions 配置
├── docs/                        # 文档根目录
│   ├── .vitepress/
│   │   ├── config.ts           # VitePress 配置
│   │   └── dist/               # 构建产物（已忽略）
│   ├── index.md                # 首页
│   ├── guide/                  # 指南
│   │   ├── getting-started.md
│   │   ├── architecture.md
│   │   ├── configuration.md
│   │   └── deployment.md       # 部署指南
│   └── api/                    # API 文档
│       ├── index.md
│       ├── auth.md
│       └── users.md
├── package.json                # npm 配置
├── package-lock.json           # npm 锁文件（需要提交）
├── DEPLOYMENT.md               # 部署快速参考
└── .gitignore                  # Git 忽略规则
```

## 🚀 部署步骤

### 1. 在 GitHub 启用 Pages

访问仓库设置：https://github.com/lwmacct/251117-bd-vmalert/settings/pages

1. 进入 **Settings** → **Pages**
2. **Source** 选择：**GitHub Actions**
3. 点击 **Save**

### 2. 提交并推送代码

```bash
# 添加所有文件
git add .

# 提交（包含 package-lock.json）
git commit -m "Add VitePress documentation with GitHub Pages deployment

- Setup VitePress with Chinese locale
- Add comprehensive documentation (Guide + API)
- Configure GitHub Actions workflow for deployment
- Set base path for GitHub Pages

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"

# 推送到 main 分支
git push origin main
```

### 3. 监控部署

访问 Actions 页面：https://github.com/lwmacct/251117-bd-vmalert/actions

等待 "Deploy VitePress Docs to Pages" workflow 完成（约 2-3 分钟）

### 4. 访问文档

部署成功后访问：**https://lwmacct.github.io/251117-bd-vmalert/**

## 🧪 本地测试

```bash
# 安装依赖
npm install

# 开发服务器（http://localhost:5173）
npm run docs:dev

# 构建测试
npm run docs:build

# 预览生产构建
npm run docs:preview
```

## 📝 文档内容

### 首页 (/)
- Hero 布局
- 功能特性展示
- 快速开始指南

### 指南 (/guide/)
- **快速开始** - 安装和运行指南
- **项目架构** - DDD 分层架构说明
- **配置系统** - Koanf 配置管理
- **部署文档** - GitHub Pages 部署详细说明

### API 文档 (/api/)
- **API 概览** - 接口总览和规范
- **认证接口** - 注册、登录、刷新、当前用户
- **用户接口** - CRUD 操作和管理

## 🔧 关键配置

### Base 路径
`docs/.vitepress/config.ts`:
```typescript
base: '/251117-bd-vmalert/'
```

### 死链接忽略
```typescript
ignoreDeadLinks: [
  '/guide/authentication',
  '/guide/postgresql',
  '/guide/redis'
]
```
这些页面的内容暂时从 `docs_old/` 目录复制过来即可。

### Workflow 触发条件
- 推送到 `main` 分支
- 修改 `docs/**` 目录
- 修改 `package.json` 或 `package-lock.json`
- 手动触发

## 📊 构建状态

本地构建测试：✅ 通过

```
vitepress v1.6.4
build complete in 2.25s.
✓ building client + server bundles...
✓ rendering pages...
```

## 💡 后续优化建议

1. **补充缺失文档**：
   - `/guide/authentication.md` - 从 docs_old 复制
   - `/guide/postgresql.md` - 从 docs_old 复制
   - `/guide/redis.md` - 从 docs_old 复制

2. **增强功能**：
   - 添加代码示例
   - 添加图表和流程图
   - 添加 API 测试用例
   - 配置搜索功能（已启用本地搜索）

3. **SEO 优化**：
   - 添加 meta 标签
   - 配置 sitemap
   - 添加 robots.txt

4. **自定义主题**：
   - 自定义主题颜色
   - 添加自定义组件
   - 配置深色模式

## 🐛 故障排查

### 构建失败
- 检查 `package-lock.json` 是否提交
- 查看 Actions 日志
- 确认 Markdown 语法正确

### 页面 404
- 检查 `base: '/251117-bd-vmalert/'` 配置
- 确认 GitHub Pages 已启用
- 等待部署完成（约 2-3 分钟）

### 样式加载失败
- 清除浏览器缓存
- 检查 base 路径配置
- 确认静态资源路径正确

## 📚 相关文档

- 详细部署指南：`docs/guide/deployment.md`
- 快速参考：`DEPLOYMENT.md`
- VitePress 官方文档：https://vitepress.dev/

---

**准备好了吗？** 执行上面的 git 命令开始部署吧！🚀
