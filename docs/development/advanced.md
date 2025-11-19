# 主题与高级能力

本页记录 VitePress 主题在本仓库中的扩展方式，方便在不破坏默认主题的前提下添加交互与品牌样式。

## 入口：`docs/.vitepress/theme/index.ts`

关键点：

1. **继承默认主题**：`extends: DefaultTheme`，保留搜索、侧边栏等默认行为。
2. **注册组件**：`enhanceApp` 中把 `Mermaid`、`ApiEndpoint`、`FeatureCard`、`StepsGuide` 设为全局组件，避免在 Markdown 中逐页引入。
3. **Medium Zoom**：通过 `setup()` + `medium-zoom` 实现图片点击放大，并在路由切换后重新初始化。

```ts
import mediumZoom from "medium-zoom";
import { onMounted, watch, nextTick } from "vue";
import { useRoute } from "vitepress";

const initZoom = () => {
  mediumZoom(".main img", { background: "var(--vp-c-bg)" });
};

onMounted(initZoom);
watch(
  () => route.path,
  () => nextTick(initZoom),
);
```

> 📌 `medium-zoom` 已写入 `docs/package.json`，无需额外引用。

## 全局样式：`docs/.vitepress/theme/style.css`

- **品牌色**：声明 `--vp-c-brand-*`，与 README、CLAUDE.md 中的设计保持一致。
- **Mermaid 容器**：`.mermaid-container` 加 margin，避免与段落贴合。
- **外链标识**：为所有外部链接自动添加 ↗，方便辨识跳转到外部站点。
- **代码块/表格圆角**：统一 8px 圆角，匹配前端与后端 UI。
- **滚动条样式**：自定义浅色/深色模式下的滚动条颜色。

所有样式都放在单文件，避免零散 CSS 难以追踪。

## 新增组件流程

1. 在 `docs/.vitepress/theme/components/` 编写 Vue SFC。
2. 在 `theme/index.ts` 中 `app.component("Foo", Foo)` 注册。
3. 在 Markdown 中直接 `<Foo />` 使用。

示例：`ApiEndpoint.vue` 支持 `method`、`path`、`version`、`deprecated` 等 props，方便 API 文档快速展现请求信息。

## 使用 `<script setup>`

VitePress 2.0 支持在 Markdown 中内联 `<script setup>`，本仓库的 `StepsGuide` 示例正是如此。推荐写成：

```vue
<script setup>
import type { Step } from '../.vitepress/theme/components/StepsGuide.vue'
const steps: Step[] = [...]
</script>

<StepsGuide :steps="steps" />
```

这样可以享受 TypeScript 类型提示，保持与应用代码一致的开发体验。

## 调整默认布局

如需向顶部/底部插槽注入内容，可修改 `theme/index.ts` 中的 `Layout`：

```ts
Layout: () =>
  h(DefaultTheme.Layout, null, {
    "layout-bottom": () => h(MyFooter),
  });
```

目前未做额外定制，保留扩展空间。

## 注意事项

- 主题层不得引入与后端耦合的逻辑，保持纯前端职责。
- 修改主题后需运行 `npm --prefix docs run dev` 验证 HMR 与生产构建。
- 若新增外部依赖（如图表组件），请同步更新 `docs/package.json` 与锁文件，并在《升级记录》中说明原因。
