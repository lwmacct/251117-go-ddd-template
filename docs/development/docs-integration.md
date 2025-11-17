# VitePress 文档集成到 Go API 服务器

## 功能概述

Go API 服务器现在可以同时提供 VitePress 构建的文档服务，通过 `/docs` 路由访问。这样您可以在同一个服务器上同时提供：

- **REST API** - `/api/*` 路由
- **在线文档** - `/docs` 路由
- **静态文件** - `/` 其他路由（可选）

## 实现方式

### 1. 配置结构

在 `internal/infrastructure/config/config.go:38` 添加了新的配置字段：

```go
type ServerConfig struct {
    Addr      string `koanf:"addr"`       // 监听地址
    Env       string `koanf:"env"`        // 运行环境
    StaticDir string `koanf:"static_dir"` // 静态资源目录
    DocsDir   string `koanf:"docs_dir"`   // 文档目录 ← 新增
}
```

默认值设置为 `docs/.vitepress/dist`。

### 2. 路由配置

在 `internal/adapters/http/router.go:70-106` 添加了文档服务路由：

```go
// 提供 VitePress 文档服务（通过 /docs 路由访问）
if cfg.Server.DocsDir != "" {
    docsGroup := r.Group("/docs")
    docsGroup.Use(func(c *gin.Context) {
        // 移除 /docs 前缀，因为 VitePress base 已包含它
        path := strings.TrimPrefix(c.Request.URL.Path, "/docs")
        if path == "" {
            path = "/"
        }

        // 构建文件路径
        filePath := filepath.Join(cfg.Server.DocsDir, path)

        // 检查文件是否存在
        if _, err := os.Stat(filePath); err == nil {
            c.File(filePath)
            return
        }

        // 如果路径不存在，尝试添加 .html 扩展名（VitePress 清洁 URL）
        if !strings.HasSuffix(path, ".html") && !strings.Contains(path, ".") {
            htmlPath := filepath.Join(cfg.Server.DocsDir, path+".html")
            if _, err := os.Stat(htmlPath); err == nil {
                c.File(htmlPath)
                return
            }
        }

        // 文件不存在，返回 index.html（用于 SPA 路由）
        indexPath := filepath.Join(cfg.Server.DocsDir, "index.html")
        if _, err := os.Stat(indexPath); err == nil {
            c.File(indexPath)
        } else {
            c.Status(404)
        }
    })
}
```

### 3. VitePress 配置

在 `docs/.vitepress/config.ts:8` 设置了正确的 base 路径：

```typescript
export default defineConfig({
  base: "/docs/", // 通过 Go API 服务器访问时使用 '/docs/'
  // ...其他配置
});
```

## 使用方法

### 1. 构建 VitePress 文档

```bash
npm run docs:build
```

这会在 `docs/.vitepress/dist/` 目录生成静态文件。

### 2. 启动 Go API 服务器

```bash
# 方式 1: 使用 task
task go:run -- api

# 方式 2: 直接运行编译后的二进制
.local/bin/go-ddd-template api
```

### 3. 访问文档

打开浏览器访问：

```
http://localhost:8080/docs/                  # 文档首页
http://localhost:8080/docs/guide/getting-started  # 指南页面
http://localhost:8080/docs/api/              # API 文档
```

## 路由结构

```
http://localhost:8080/
├── /health                    # 健康检查
├── /api/                      # REST API
│   ├── /api/auth/register     # 用户注册
│   ├── /api/auth/login        # 用户登录
│   ├── /api/users             # 用户列表
│   └── ...
├── /docs/                     # VitePress 文档 ← 新增
│   ├── /docs/                 # 文档首页
│   ├── /docs/guide/           # 指南
│   ├── /docs/api/             # API 文档
│   └── /docs/assets/          # 静态资源
└── /*                         # 其他静态文件（如果配置了 StaticDir）
```

## 配置选项

### 环境变量

可以通过环境变量覆盖文档目录路径：

```bash
# 使用自定义文档目录
APP_SERVER_DOCS_DIR=/path/to/docs .local/bin/go-ddd-template api

# 禁用文档服务（留空）
APP_SERVER_DOCS_DIR="" .local/bin/go-ddd-template api
```

### 配置文件

在 `config.yaml` 或 `configs/config.yaml` 中配置：

```yaml
server:
  addr: "0.0.0.0:8080"
  env: "development"
  static_dir: "web/dist"
  docs_dir: "docs/.vitepress/dist" # 文档目录
```

## 部署建议

### 开发环境

开发时推荐使用 VitePress 自带的开发服务器（支持热更新）：

```bash
# 终端 1: 启动 VitePress 开发服务器
npm run docs:dev
# 访问 http://localhost:5173

# 终端 2: 启动 Go API 服务器
task go:run -- api
# 访问 http://localhost:8080/api
```

### 生产环境

生产环境可以统一使用 Go API 服务器：

```bash
# 1. 构建 VitePress
npm run docs:build

# 2. 构建 Go 应用
task go:build

# 3. 启动服务（包含 API + 文档）
.local/bin/go-ddd-template api
```

## 技术细节

### VitePress 清洁 URL 支持

路由器实现了对 VitePress 清洁 URL 的支持：

- `/docs/guide/getting-started` → `guide/getting-started.html`
- `/docs/api/` → `api/index.html`
- `/docs/` → `index.html`

### SPA 路由回退

如果请求的文件不存在，会返回 `index.html`，由 VitePress 的客户端路由处理。

### 性能优化

- 文件直接由 Gin 的 `c.File()` 方法提供，避免额外的内存拷贝
- 支持浏览器缓存（HTTP 标准头）
- 静态文件不经过任何中间件处理（除了 CORS）

## 故障排除

### 文档返回 404

1. 检查文档是否已构建：

   ```bash
   ls docs/.vitepress/dist/
   ```

2. 检查配置：

   ```bash
   # 查看当前配置
   cat configs/config.example.yaml
   ```

3. 检查路径权限：
   ```bash
   ls -ld docs/.vitepress/dist
   ```

### 资源文件无法加载

检查 VitePress base 配置是否正确：

```typescript
// docs/.vitepress/config.ts
export default defineConfig({
  base: "/docs/", // 必须以 / 开头和结尾
});
```

### CSS/JS 文件路径错误

确保 VitePress 构建时使用了正确的 base 路径：

```bash
# 重新构建
npm run docs:build

# 检查生成的 HTML 中的资源路径
grep -r 'assets' docs/.vitepress/dist/index.html
```

## 示例请求

```bash
# 访问文档首页
curl http://localhost:8080/docs/

# 访问指南页面
curl http://localhost:8080/docs/guide/getting-started

# 访问 API 文档
curl http://localhost:8080/docs/api/auth

# 同时测试 API 和文档
curl http://localhost:8080/health
curl http://localhost:8080/api/auth/login
curl http://localhost:8080/docs/
```

## 相关文件

- 配置定义: `internal/infrastructure/config/config.go:38`
- 路由实现: `internal/adapters/http/router.go:70-106`
- VitePress 配置: `docs/.vitepress/config.ts:8`
- 配置示例: `configs/config.example.yaml:12`

## 后续改进

- [ ] 添加文档访问权限控制（JWT 认证）
- [ ] 支持多版本文档
- [ ] 添加文档搜索 API
- [ ] 集成 API 文档自动生成（Swagger/OpenAPI）
- [ ] 添加文档缓存策略
