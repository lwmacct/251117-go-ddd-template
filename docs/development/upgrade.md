# 文档依赖升级记录

记录 VitePress 及其伴生依赖的升级背景、步骤与回归测试，确保文档系统与 Go DDD Template 的代码基保持同步。

## 当前版本

| 包 | 版本 | 位置 |
| -- | ---- | ---- |
| `vitepress` | `^2.0.0-alpha.13` | `docs/package.json` |
| `vue` | `^3.5.24` | `docs/package.json` |
| `mermaid` | `^11.12.1` | `docs/package.json` |
| `medium-zoom` | `^1.1.0` | 主题增强 |
| `@types/node` | `^24.10.1` | 类型支持 |
| Node.js | `>= 20.19.0` | 由 `engines.node` 强制 |

> 最后一次重大升级：2025-11-18，从 VitePress 1.6.4 迁移至 2.0.0-alpha.13。

## 升级流程

1. **评估变更**
   - 阅读 VitePress Release Notes 与 Breaking Changes。
   - 对比 `docs/.vitepress/config.ts` 与 `theme/` 中可能受影响的配置（如 `markdown`、`search`）。
2. **建立分支**
   ```bash
   git checkout -b chore/vitepress-2.0.x
   npm --prefix docs install vitepress@next
   ```
3. **更新依赖**
   - 保持 `npm --prefix docs install` 使用 `package-lock.json` 记录。
   - 如需新增依赖（例如新的代码高亮器），请在本页新增条目。
4. **同步配置**
   - 依据官方迁移指南更新 `config.ts`（如 `cjkFriendlyEmphasis`）。
   - 若需要额外的 Vite 配置，可利用 `defineConfig({ vite: { ... } })` 区块。
5. **验证**
   ```bash
   npm --prefix docs run build
   npm --prefix docs run preview
   task go:run -- api  # 验证 /docs 路由
   ```
6. **提交说明**
   - 在 `docs/development/upgrade.md` 新增“升级条目”。
   - 提交信息示例：`docs: upgrade vitepress to 2.0.0-beta.1`。

## 回归清单

- [ ] `npm --prefix docs run dev` HMR 正常，Markdown/组件无报错。
- [ ] `npm --prefix docs run build` 通过，无 `mermaid` 渲染错误。
- [ ] `npm --prefix docs run preview` 可访问。
- [ ] `task go:run -- api` 下访问 `http://localhost:8080/docs/backend/` 正常。
- [ ] GitHub Actions `Deploy VitePress Docs to Pages` 成功。

## 常见问题

| 问题 | 说明 |
| ---- | ---- |
| Node 版本过低 | 升级到 LTS ≥ 20.19.0，或者在 CI 中显式指定 `actions/setup-node@v4` 的版本。 |
| Shiki/高亮冲突 | VitePress 2 默认内置 Shiki，若出现冲突请移除旧的 `markdown-it-attrs` 相关配置。 |
| Mermaid 报错 | 通常是因为语法错误或在 SSR 构建阶段无法渲染。使用 `<Mermaid>` 组件的实现可避免直接在 Node 环境运行 mermaid。 |
| 包含实验性版本 | `2.0.0-alpha.*` 与 `beta` 均为预发布版，升级后请关注 issue 列表，必要时固定具体版本号而不是 `^`。 |

## 历史记录

| 日期 | 变动 | 影响 |
| ---- | ---- | ---- |
| 2025-11-18 | VitePress `1.6.4 → 2.0.0-alpha.13` | 获得 CJK 友好强调、图片懒加载、单次 git 时间戳查询，Node 最低版本提升至 20.19。 |

未来如有新的升级，请在上表继续追加并描述对文档与 Go 服务的影响。
