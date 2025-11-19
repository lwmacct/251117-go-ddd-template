# 文档与 Go API 集成

VitePress 构建出的静态文件由 Go API 服务器直接托管，路径为 `/docs`. 本文记录配置结构、路由实现以及常见扩展点，所有示例均来自当前代码库。

## 配置来源

`internal/infrastructure/config/config.go` 中的 `ServerConfig` 定义了文档目录：

```go
// internal/infrastructure/config/config.go
type ServerConfig struct {
    Addr      string `koanf:"addr"`
    Env       string `koanf:"env"`
    StaticDir string `koanf:"static_dir"`
    DocsDir   string `koanf:"docs_dir"` // VitePress 构建输出
}
```

默认值在同文件的 `DefaultConfig()` 中设置为 `docs/.vitepress/dist`。运行时可通过：

```bash
APP_SERVER_DOCS_DIR=/custom/path task go:run -- api
```

来覆盖此参数；若置为空字符串则自动关闭文档路由。

## 适配器层实现

`internal/adapters/http/router.go`（第 161 行起）注册了文档路由：

```go
if cfg.Server.DocsDir != "" {
    docs := r.Group("/docs")
    docs.Use(func(c *gin.Context) {
        reqPath := strings.TrimPrefix(c.Request.URL.Path, "/docs")
        if reqPath == "" {
            reqPath = "/"
        }
        fullPath := filepath.Join(cfg.Server.DocsDir, reqPath)

        if _, err := os.Stat(fullPath); err == nil {
            c.File(fullPath)
            return
        }

        if !strings.HasSuffix(reqPath, ".html") && !strings.Contains(reqPath, ".") {
            htmlPath := filepath.Join(cfg.Server.DocsDir, reqPath+".html")
            if _, err := os.Stat(htmlPath); err == nil {
                c.File(htmlPath)
                return
            }
        }

        indexPath := filepath.Join(cfg.Server.DocsDir, "index.html")
        if _, err := os.Stat(indexPath); err == nil {
            c.File(indexPath)
        } else {
            c.Status(http.StatusNotFound)
        }
    })
}
```

### 行为说明

1. **清洁 URL**：`/docs/backend/ddd-cqrs` 会被映射到 `docs/.vitepress/dist/backend/ddd-cqrs.html`。
2. **静态文件优先**：若请求恰好匹配物理文件直接返回。
3. **SPA 回退**：不存在的路径会回退到 `index.html`，由 VitePress 前端路由处理。
4. **完全隔离**：逻辑位于 adapters 层，未侵入 application/domain 层，符合 DDD 依赖方向。

## 目录结构

```
docs/.vitepress/dist/
├── index.html
├── guide/
├── backend/
├── api/
└── assets/
```

- `assets/` 中的静态资源通过同一路由返回。
- 需要确保 Go 进程对该目录具有读取权限。

## 禁用或替换

| 需求           | 做法                                                                                                  |
| -------------- | ----------------------------------------------------------------------------------------------------- |
| 禁用 `/docs`   | 将 `APP_SERVER_DOCS_DIR` 或 `server.docs_dir` 设为空字符串。                                          |
| 切换到 CDN     | 将 `DocsDir` 指向一个同步目录，并在 CDN 发布静态文件；同时可保留 `/docs` 作为回退。                   |
| 支持多版本文档 | 修改 `DocsDir` 指向版本化目录，例如 `docs/.vitepress/dist/v2`，并在 VitePress 内使用多语言/多基路径。 |

## 常见故障排查

| 现象                  | 解决方案                                                                                                                        |
| --------------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| 访问 `/docs` 提示 404 | 检查 `docs/.vitepress/dist/index.html` 是否存在，以及 `server.docs_dir` 是否配置正确。                                          |
| 静态资源丢失          | 构建命令必须保持 `base=/docs/`。若 `VITEPRESS_BASE` 设置错误，重新以 `VITEPRESS_BASE=/docs/ npm --prefix docs run build` 构建。 |
| 生产环境需要缓存控制  | 在 `docsGroup.Use` 中添加自定义中间件，或在上游 Nginx/CDN 层处理 Cache-Control 头。                                             |

## 相关命令速查

```bash
# 构建文档
npm --prefix docs run build

# 运行 Go API 并提供 /docs
task go:run -- api

# 使用自定义 DocsDir 进行验证
APP_SERVER_DOCS_DIR=/tmp/docs-dist task go:run -- api
```
