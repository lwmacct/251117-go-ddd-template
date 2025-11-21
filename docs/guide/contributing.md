# 贡献指南

感谢你对本项目感兴趣！本文档说明如何为项目贡献代码和文档。

## 文档贡献

本项目使用 [VitePress](https://vitepress.dev/) 构建文档站点。

### 快速开始

#### 1. 克隆项目

```bash
git clone https://github.com/lwmacct/251117-go-ddd-template.git
cd 251117-go-ddd-template
```

#### 2. 安装依赖

```bash
npm install
```

#### 3. 本地运行文档

```bash
# 启动开发服务器 (支持热重载)
cd docs
npm run dev
```

文档将在本地端口 5173 启动，通常是 `http://localhost:5173`。

#### 4. 构建测试

```bash
# 构建生产版本
cd docs
npm run build

# 预览构建结果
npm run preview
```

### 文档结构

```
docs/
├── index.md                # 首页
├── guide/                  # 用户指南
│   ├── getting-started.md          # 快速开始
│   ├── configuration.md            # 配置系统
│   ├── cli-commands.md             # CLI 命令
│   ├── application-deployment.md   # 应用部署
│   ├── docs-deployment.md          # 文档部署
│   ├── testing.md                  # 测试指南
│   └── contributing.md             # 贡献指南
│
├── architecture/           # 系统架构
│   ├── index.md                    # 架构概览
│   ├── overview.md                 # 架构设计
│   ├── authentication.md           # 认证机制
│   ├── rbac.md                     # RBAC 权限系统
│   ├── pat.md                      # Personal Access Token
│   ├── postgresql.md               # PostgreSQL 架构
│   └── redis.md                    # Redis 架构
│
├── frontend/               # 前端文档
│   ├── index.md
│   ├── getting-started.md
│   ├── project-structure.md
│   └── api-integration.md
│
├── api/                    # API 参考文档
│   ├── index.md
│   ├── auth.md
│   ├── users.md
│   └── cache.md
│
└── development/            # 开发者文档
    ├── index.md
    ├── deployment.md
    ├── docs-integration.md
    ├── mermaid-integration.md
    ├── features.md
    ├── advanced.md
    └── upgrade.md
```

### 文档分类说明

#### 1. 用户指南 (`/guide/`)

**目标读者**: 应用的使用者、部署人员、新手开发者

**内容特点**:

- 面向实践，注重操作步骤
- 适合快速上手
- 包含大量示例代码

#### 2. 架构蓝图 (`/architecture/`)

**目标读者**: 架构师、高级开发者、技术决策者

**内容特点**:

-- 分层与依赖关系
-- 数据与基础设施方案
-- 包含架构图和流程图
-- 适合理解系统工作原理与演进

#### 3. 参考资料 (`/guide/reference`)

**目标读者**: API 使用者、前端开发者、集成开发者

**内容特点**:

- 对接说明与补充资料
- API 端点说明（可跳转 Swagger）
- 常见问题与最佳实践

#### 4. 开发者文档 (`/development/`)

**目标读者**: 项目维护者、文档编写者

**内容特点**:

- 面向文档系统本身
- 适合维护和扩展文档

### 编写文档

#### Markdown 格式

所有文档使用 Markdown 格式编写，支持：

- **标题**：`# H1`, `## H2`, `### H3`
- **代码块**：使用 ` ```语言 ` 包裹
- **链接**：`[文本](/guide/getting-started)`
- **表格**、列表、图片等

#### Frontmatter

每个页面可以添加 Frontmatter 配置：

```markdown
---
title: 页面标题
editLink: true
---

# 页面内容
```

#### 代码示例

````markdown
```bash
# 示例命令
npm install
```

```go
// Go 代码示例
func main() {
    fmt.Println("Hello, World!")
}
```
````

#### 内部链接

使用相对路径引用其他页面：

```markdown
- 查看 [快速开始](/guide/getting-started)
- 了解 [项目架构](/architecture/)
- 探索 [参考资料](/guide/reference)
```

### 修改配置

编辑 `docs/.vitepress/config.ts` 修改导航栏、侧边栏等配置：

```typescript
export default defineConfig({
  // 站点配置
  title: "Go DDD Template",
  description: "...",

  // 导航栏
  themeConfig: {
    nav: [
      {
        text: "首页",
        link: "/",
      },
      {
        text: "指南",
        link: "/guide/getting-started",
      },
    ],

    // 侧边栏
    sidebar: {
      "/guide/": [
        {
          text: "指南",
          items: [
            {
              text: "快速开始",
              link: "/guide/getting-started",
            },
          ],
        },
      ],
    },
  },
});
```

### 提交流程

#### 1. 创建分支

```bash
git checkout -b docs/update-guide
```

#### 2. 修改文档

编辑 `docs/` 目录下的 Markdown 文件。

#### 3. 本地测试

```bash
# 确保构建通过
cd docs
npm run build
```

#### 4. 提交更改

```bash
git add docs/
git commit -m "docs: update getting started guide"
git push origin docs/update-guide
```

#### 5. 创建 Pull Request

访问 GitHub 仓库创建 PR。

### 文档风格指南

#### 标题层级

- 一级标题 (`#`) ：页面标题，每页只有一个
- 二级标题 (`##`) ：主要章节
- 三级标题 (`###`) ：子章节
- 四级标题 (`####`) ：细节说明

#### 代码示例

- 总是指定语言类型
- 添加必要的注释
- 保持代码简洁易懂

```bash
# ✅ 好的示例 (有注释，清晰)
docker-compose up -d  # 启动服务
```

```bash
# ❌ 不好的示例 (无注释，冗长)
docker-compose -f docker-compose.yml up --detach
```

#### 术语一致性

| 术语   | 使用 | 不使用             |
| ------ | ---- | ------------------ |
| 认证   | ✅   | 验证、鉴权         |
| 仓储   | ✅   | 存储库、Repository |
| 中间件 | ✅   | 拦截器             |

### 文档分类原则

#### 用户指南 vs 架构文档

| 特性         | 用户指南     | 架构文档           |
| ------------ | ------------ | ------------------ |
| **目标**     | 如何使用     | 为什么这样设计     |
| **深度**     | 操作层面     | 设计层面           |
| **受众**     | 初学者、用户 | 架构师、高级开发者 |
| **内容**     | 步骤、示例   | 原理、流程、设计   |
| **更新频率** | 功能变化时   | 架构变化时         |

#### 示例：权限系统文档的分类

- **用户指南** (`/guide/`): 如何使用权限功能、如何分配角色
- **架构文档** (`/architecture/identity-rbac`): RBAC 系统的工作原理、三段式权限格式、中间件实现
- **API 文档** (`/api/`): 权限相关的 API 端点、请求格式

### 文档编写最佳实践

#### 架构文档编写要点

1. **解释"为什么"，而不仅仅是"是什么"**
   - ✅ "使用三段式权限格式是为了实现更细粒度的权限控制"
   - ✗ "权限格式是 domain:resource:action"

2. **包含设计决策和权衡**
   - ✅ "选择 JWT 而非 Session 是为了支持水平扩展"
   - ✗ "系统使用 JWT"

3. **提供架构图和流程图**
   - 使用 Mermaid 绘制序列图、流程图
   - 图文结合，易于理解

4. **深入技术细节**
   - 代码实现
   - 数据库查询优化
   - 性能考虑

#### 用户指南编写要点

1. **聚焦操作步骤**
   - 清晰的步骤说明
   - 命令行示例
   - 预期结果

2. **面向实践**
   - 真实场景
   - 常见问题
   - 故障排查

3. **易于上手**
   - 循序渐进
   - 避免过多技术细节
   - 提供快速路径

## 代码贡献

### 开发环境

#### 前置要求

- Go 1.25.4+
- Docker & Docker Compose
- Task (可选)

#### 设置开发环境

```bash
# 1. 克隆项目
git clone https://github.com/lwmacct/251117-go-ddd-template.git
cd 251117-go-ddd-template

# 2. 启动依赖服务
docker-compose up -d

# 3. 运行应用
task go:run -- api

# 或使用热重载
air
```

### 代码风格

- 遵循 Go 官方代码风格
- 使用 `gofmt` 格式化代码
- 运行 `golangci-lint` 检查

```bash
# 格式化代码
gofmt -w .

# 运行 linter
golangci-lint run
```

### 提交规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 格式：

```
<type>(<scope>): <subject>

<body>

<footer>
```

**类型 (type) **：

- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式 (不影响功能)
- `refactor`: 重构
- `test`: 测试
- `chore`: 构建/工具配置

**示例**：

```bash
feat(auth): add password reset functionality

- Add forgot password endpoint
- Implement email sending
- Add reset token validation

Closes #123
```

### Pull Request 流程

1. Fork 项目到你的 GitHub 账号
2. 克隆你的 fork：`git clone https://github.com/YOUR_USERNAME/251117-go-ddd-template.git`
3. 创建功能分支：`git checkout -b feature/amazing-feature`
4. 提交更改：`git commit -m 'feat: add amazing feature'`
5. 推送到分支：`git push origin feature/amazing-feature`
6. 创建 Pull Request

### 测试

```bash
# 运行所有测试
go test ./...

# 带覆盖率
go test -cover ./...

# 详细输出
go test -v ./...
```

## 部署文档

文档通过 GitHub Actions 自动部署到 GitHub Pages。

### 自动部署

推送到 `main` 分支时，如果修改了以下文件，会自动触发部署：

- `docs/**` - 任何文档文件变更
- `.github/workflows/deploy-docs.yml` - 部署流程配置变更

### 手动触发

访问 [Actions](https://github.com/lwmacct/251117-go-ddd-template/actions)，选择 "Deploy VitePress Docs to Pages"，点击 "Run workflow"。

### 查看部署结果

部署成功后访问：https://lwmacct.github.io/251117-go-ddd-template/

## 问题反馈

### 报告 Bug

在 [Issues](https://github.com/lwmacct/251117-go-ddd-template/issues) 页面创建新 Issue，包含：

- 问题描述
- 复现步骤
- 预期行为
- 实际行为
- 环境信息 (Go 版本、操作系统等)

### 功能建议

创建 Issue 时选择 "Feature Request" 模板，描述：

- 功能描述
- 使用场景
- 预期效果

## 许可证

本项目采用 MIT 许可证。贡献代码即表示你同意将你的贡献以相同许可证发布。

## 获取帮助

- 📚 查看[项目文档](https://lwmacct.github.io/251117-go-ddd-template/)
- 💬 在 [Issues](https://github.com/lwmacct/251117-go-ddd-template/issues) 提问
- 📧 联系维护者

---

感谢你的贡献！🎉
