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

- 使用 **小写字母** 和 **短横线** (kebab-case)
- 扩展名必须是 `.md`
- ✅ `getting-started.md`
- ❌ `Getting_Started.md` / `快速开始.md`

#### Frontmatter (可选)

```yaml
---
title: 文档标题
description: 文档描述
outline: [2, 3]
---
```

#### 内容规范

- 使用清晰的标题层级 (# H1、## H2、### H3)
- 代码块指定语言 (\`\`\`typescript、\`\`\`bash)
- 使用相对链接引用其他文档
- 避免使用绝对 URL
- 避免使用 `${{ }}` 语法 (会导致 Vue 编译错误)

### 3. 更新 VitePress 配置

查找并编辑 VitePress 配置文件（通常是 `docs/.vitepress/config.ts` 或 `config.mts`）

#### 添加顶部导航

```typescript
nav: [
  { text: "首页", link: "/" },
  { text: "新分类", link: "/new-category/" }, // ← 添加
];
```

#### 添加侧边栏链接（必须）

```typescript
sidebar: {
  "/category/": [
    {
      text: "分类名",
      items: [
        { text: "文档标题", link: "/category/doc-name" }, // ← 添加
      ],
    },
  ],
};
```

### 4. 验证

```bash
# 开发模式预览
npm run dev

# 构建验证
npm run build
```

## VitePress 特性支持

### 代码块增强

**语法高亮**：

````markdown
```typescript
const config = defineConfig({
  title: "My Site",
});
```
````

**行高亮**：

````markdown
```typescript{2}
export default {
  title: "highlighted line", // [!code highlight]
};
```
````

### 容器块

```markdown
::: tip 提示
这是一个提示
:::

::: warning 警告
这是一个警告
:::

::: danger 危险
这是危险提示
:::
```

### 链接规范

**内部链接**（推荐）：

```markdown
[同级文档](./other-doc)
[其他分类](/category/doc)
```

**外部链接**：

```markdown
[VitePress 官网](https://vitepress.dev/)
```

## 必须避免的错误

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

创建文档后必须在配置文件中添加侧边栏链接，否则无法导航。

### ❌ 中文文件名

```markdown
<!-- ❌ 错误 -->
docs/快速开始.md

<!-- ✅ 正确 -->
docs/getting-started.md
```

## 文档模板

### 用户指南模板

````markdown
# 文档标题

简短介绍 (1-2 句话说明这个功能是什么)

## 前提条件

- 需要的环境或依赖
- 需要的权限或配置

## 步骤

### 1. 第一步

详细说明...

```bash
# 示例命令
command here
```

### 2. 第二步

详细说明...

## 常见问题

### Q: 问题描述？

A: 解答...

## 相关链接

- [相关文档](./related-doc)
- [外部资源](https://example.com)
````

### API 文档模板

````markdown
# API 名称

简短描述

## 端点

### POST /api/endpoint

描述

**请求参数**：

| 参数 | 类型   | 必填 | 说明 |
| ---- | ------ | ---- | ---- |
| name | string | 是   | 名称 |

**请求示例**：

```json
{
  "name": "example"
}
```

**响应示例**：

```json
{
  "success": true,
  "data": {}
}
```

**错误码**：

| 错误码 | 说明     |
| ------ | -------- |
| 400    | 参数错误 |
| 401    | 未授权   |
````

## 完成清单

完成后提供以下信息：

- [ ] 创建/更新的文件路径
- [ ] 修改的配置文件（如果有）
- [ ] 预览命令（如 `npm run dev`）
