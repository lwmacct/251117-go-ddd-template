---
name: vitepress-docs
description: 当用户要求"创建文档"、"更新文档"、"写文档"、"添加指南"、"写 API 文档"、"新建文档页面" 或提到 Docs/VitePress 时使用。
---

# VitePress 文档管理器

创建或更新符合 VitePress 规范的 Markdown 文档，并自动维护导航配置。

## 何时使用此技能

- ✅ 用户要求"创建/更新文档"
- ✅ 用户要求"添加指南/API 文档"
- ✅ 用户提到 VitePress 相关操作

## 核心工作流程

### 1. 探索项目结构

首先使用 `Glob` 或 `Bash` 探索项目的文档目录结构：

```bash
# 查看文档根目录结构
tree docs -L 2 -I node_modules

# 或查看 VitePress 配置
cat docs/.vitepress/config.ts
```

了解现有的文档分类和目录组织方式。

### 2. 创建或更新 Markdown 文件

#### 文件命名规范

- 使用小写字母和短横线分隔
- 扩展名 `.md`
- ✅ `quick-start.md` / `api-reference.md`
- ❌ `Getting_Started.md` / `快速开始.md`

#### Frontmatter

```yaml
---
title: 文档标题
description: 文档描述
outline: [2, 3]
---
```

### 3. 更新 VitePress 配置

配置文件位于 `docs/.vitepress/config/` 目录。

`nav.json`:

```json
[
  { "text": "首页", "link": "/" },
  { "text": "快速开始", "link": "/guide/quick-start" },
  { "text": "新分类", "link": "/new-category/" }
]
```

### ❌ 使用绝对 URL

```markdown
<!-- ❌ 错误 -->

[文档](http://localhost:8080/docs/guide)

<!-- ✅ 正确 -->

[文档](./guide)
```

### ❌ 使用模板语法

```markdown
<!-- ❌ 错误 -->

变量：${{ variable }}

<!-- ✅ 正确 -->

变量：使用环境变量
```

### ❌ 忘记更新侧边栏

创建文档后需要在配置文件中添加侧边栏链接。

### ❌ 中文文件名

```markdown
<!-- ❌ 错误 -->

docs/快速开始.md

<!-- ✅ 正确 -->

docs/quick-start.md
```
