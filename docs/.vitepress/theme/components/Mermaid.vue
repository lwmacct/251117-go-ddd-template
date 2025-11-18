<template>
  <div>
    <div ref="codeContainer" style="display: none">
      <slot />
    </div>
    <div ref="mermaidContainer" class="mermaid-container"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { useData } from "vitepress";
import mermaid from "mermaid";

const { isDark } = useData();
const codeContainer = ref<HTMLElement>();
const mermaidContainer = ref<HTMLElement>();

// 初始化 Mermaid 配置
const initMermaid = (theme: "dark" | "default") => {
  mermaid.initialize({
    startOnLoad: false,
    theme: theme,
    securityLevel: "loose",
    fontFamily: "var(--vp-font-family-base)",
  });
};

// 渲染图表
const renderDiagram = async () => {
  if (!mermaidContainer.value || !codeContainer.value) return;

  try {
    // 查找 pre 标签并获取其文本内容
    const preElement = codeContainer.value.querySelector("pre");
    if (!preElement) {
      mermaidContainer.value.innerHTML = '<pre class="error">未找到代码内容</pre>';
      return;
    }

    // 使用 textContent 获取纯文本（会自动解码 HTML 实体并保留换行符）
    const code = preElement.textContent?.trim() || "";

    if (!code) {
      mermaidContainer.value.innerHTML = '<pre class="error">未提供 Mermaid 代码</pre>';
      return;
    }

    // 生成唯一 ID
    const id = `mermaid-${Math.random().toString(36).substring(2, 11)}`;

    // 初始化主题
    initMermaid(isDark.value ? "dark" : "default");

    // 渲染图表
    const { svg } = await mermaid.render(id, code);
    mermaidContainer.value.innerHTML = svg;
  } catch (error) {
    console.error("Mermaid rendering error:", error);
    if (mermaidContainer.value) {
      mermaidContainer.value.innerHTML = `<pre class="error">渲染失败: ${error}</pre>`;
    }
  }
};

// 组件挂载时渲染
onMounted(() => {
  renderDiagram();
});

// 监听主题变化
watch(isDark, () => {
  renderDiagram();
});
</script>

<style scoped>
.mermaid-container {
  margin: 16px 0;
  overflow-x: auto;
}

.mermaid-container :deep(svg) {
  max-width: 100%;
  height: auto;
  display: block;
}

.error {
  color: var(--vp-c-danger-1);
  background: var(--vp-c-danger-soft);
  padding: 12px;
  border-radius: 8px;
}
</style>
