# VitePress Mermaid 集成说明

## 概述

本项目已成功集成 Mermaid.js，支持在 VitePress 2.0 文档中使用 Mermaid 图表。

## 技术方案

- **不使用第三方插件**：`vitepress-plugin-mermaid` 仅支持 VitePress 1.x
- **自定义实现**：使用 markdown-it 自定义渲染器 + Vue 3 组件
- **完全兼容**：与 VitePress 2.0.0-alpha.13 完美兼容

## 架构

### 1. Markdown-it 配置 (`.vitepress/config.ts`)

```typescript
markdown: {
  config: (md) => {
    const fence = md.renderer.rules.fence!;
    md.renderer.rules.fence = (...args) => {
      const [tokens, idx] = args;
      const token = tokens[idx];
      const lang = token.info.trim();

      if (lang === "mermaid") {
        const code = md.utils.escapeHtml(token.content.trim());
        return `<Mermaid>${code}</Mermaid>\n`;
      }

      return fence(...args);
    };
  },
}
```

该配置将 ` ```mermaid ` 代码块转换为 `<Mermaid>` Vue 组件。

### 2. Mermaid Vue 组件 (`.vitepress/theme/components/Mermaid.vue`)

- 读取 slot 中的 Mermaid 代码（已转义的 HTML 实体）
- 使用 `innerHTML` 自动解码 HTML 实体
- 调用 `mermaid.render()` 渲染图表
- 自动适配亮色/暗色主题

### 3. 全局组件注册 (`.vitepress/theme/index.ts`)

```typescript
export default {
  extends: DefaultTheme,
  enhanceApp({ app }) {
    app.component("Mermaid", Mermaid);
  },
} satisfies Theme;
```

## 使用方法

在任何 Markdown 文件中使用标准的 Mermaid 代码块语法：

````markdown
```mermaid
flowchart LR
    A[开始] --> B[处理]
    B --> C[结束]
```
````

## 支持的图表类型

- ✅ 流程图 (Flowchart)
- ✅ 时序图 (Sequence Diagram)
- ✅ 类图 (Class Diagram)
- ✅ 状态图 (State Diagram)
- ✅ ER 图 (Entity Relationship)
- ✅ Git 分支图 (Gitgraph)
- ✅ 甘特图 (Gantt)
- ✅ 饼图 (Pie Chart)
- ✅ 思维导图 (Mindmap)
- ✅ 时间线 (Timeline)

## 特性

- ✅ 自动主题切换（亮色/暗色）
- ✅ 响应式设计
- ✅ 标准 Markdown 语法
- ✅ 无需第三方插件
- ✅ 完全类型安全

## 示例

查看完整示例：`docs/guide/mermaid-examples.md`

## 故障排除

### 渲染失败

如果图表渲染失败，请检查：

1. **语法错误**：使用 [Mermaid Live Editor](https://mermaid.live/) 验证语法
2. **特殊字符**：确保没有使用未转义的 HTML 特殊字符
3. **版本兼容**：当前使用 Mermaid v11.12.1

### 主题不匹配

Mermaid 组件会自动监听 VitePress 主题变化。如果主题不匹配：

1. 检查浏览器控制台是否有错误
2. 尝试手动切换主题
3. 清除浏览器缓存

## 依赖

- `mermaid`: ^11.12.1
- `vitepress`: ^2.0.0-alpha.13
- `vue`: ^3.5.24

## 参考资料

- [Mermaid 官方文档](https://mermaid.js.org/)
- [VitePress 官方文档](https://vitepress.dev/)
- [Markdown-it 文档](https://markdown-it.github.io/)
