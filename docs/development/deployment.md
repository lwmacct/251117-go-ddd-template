# 文档部署指南

本指南说明 VitePress 文档在不同环境中的发布方式，覆盖本地预览、与 Go API 服务器的联动部署以及 GitHub Pages 自动化。所有步骤均以当前仓库结构为准。

## 运行要求

| 组件          | 版本/路径              | 说明                                                |
| ------------- | ---------------------- | --------------------------------------------------- |
| Node.js       | `>= 20.19.0`           | 受 `docs/package.json#engines` 限制。               |
| npm           | v10+                   | 直接驱动 VitePress 脚本。                           |
| Go            | 1.22+                  | 运行 `task go:run -- api` 以托管 `/docs` 静态文件。 |
| Docs 输出目录 | `docs/.vitepress/dist` | Go 服务器通过 `cfg.Server.DocsDir` 读取该目录。     |

## 本地开发流程

1. **安装依赖**：`npm --prefix docs install`
2. **开发模式**：`npm --prefix docs run dev -- --host`
3. **访问地址**：默认 `http://localhost:5173`
4. **跨域/接口调试**：文档中的示例 API 与 `task go:run -- api` 启动的后端一致，均在 `http://localhost:8080`

> 🧪 推荐在两个终端分别运行 `npm --prefix docs run dev` 与 `task go:run -- api`，即可同时调试文档与 API。

## 构建产物

```bash
# 生成静态文件（base 默认为 /docs/）
npm --prefix docs run build

# 输出目录：docs/.vitepress/dist
ls docs/.vitepress/dist
```

构建过程中，`docs/.vitepress/config.ts` 会根据 `process.env.VITEPRESS_BASE` 设置 `base`，默认 `/docs/`，与 Go 服务器的 `/docs` 前缀保持一致。

## 与 Go API 服务器的联动部署

1. **构建文档**：`npm --prefix docs run build`
2. **构建或运行 API**：`task go:run -- api`（或 `task go:build` + 运行二进制）
3. **配置文件**：`configs/config.yaml` 中的 `server.docs_dir` 默认为 `docs/.vitepress/dist`
4. **访问地址**：`http://localhost:8080/docs/`

在 `internal/adapters/http/router.go`（约第 161 行）中，如果 `cfg.Server.DocsDir` 不为空，会注册 `/docs` 路由组并提供清洁 URL、SPA 回退等能力。详情见《文档与 Go API 集成》章节。

## GitHub Pages 自动部署

- **工作流**：`.github/workflows/deploy-docs.yml`
- **触发条件**：推送到 `main` 且修改了 `docs/**`、`docs/package*.json` 或工作流本身
- **构建命令**：同样是 `npm --prefix docs run build`
- **关键环境变量**：

  ```yaml
  env:
    VITEPRESS_BASE: /${{ github.event.repository.name }}/
  ```

  GitHub Actions 会根据仓库名自动设置 `base`，因此无需硬编码。

- **产物上传**：`actions/upload-pages-artifact@v3` → `docs/.vitepress/dist`
- **发布**：`actions/deploy-pages@v4`

启用步骤：在仓库 `Settings → Pages` 中选择 `Source: GitHub Actions` 即可。

## 手动验证 GitHub Pages 构建

```bash
VITEPRESS_BASE=/your-repo/ npm --prefix docs run build
sed -n '1,5p' docs/.vitepress/dist/index.html | grep '<base'
```

若输出 `<base href="/your-repo/">` 即表示配置正确。

## 常见问题

| 现象                | 排查                                                                                                   |
| ------------------- | ------------------------------------------------------------------------------------------------------ |
| 文档 404            | 确认 `npm --prefix docs run build` 是否成功，`DocsDir` 是否指向 dist 目录。                            |
| 静态资源路径错误    | 检查 `VITEPRESS_BASE` 是否与部署路径一致，例如 Go 服务器必须保持 `/docs/`。                            |
| GitHub Actions 失败 | 查看工作流日志中的 `npm install`、`npm run build`、`upload` 步骤；大多为 Node 版本或锁文件不一致导致。 |
| 浏览器缓存旧版本    | 尝试访问 `http://localhost:8080/docs/index.html?t=$(date +%s)` 或清理缓存。                            |

## 发版 checklist

1. `npm --prefix docs ci`（或 `npm install`）确保 lock 文件与依赖同步
2. `npm --prefix docs run lint`（如需，可依据自定义脚本）
3. `npm --prefix docs run build`
4. `task go:run -- api` 并访问 `/docs`
5. 推送到 `main` 观察 GitHub Actions `Deploy VitePress Docs to Pages` 工作流结果
