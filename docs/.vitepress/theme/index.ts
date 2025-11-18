// .vitepress/theme/index.ts
import { h } from "vue";
import type { Theme } from "vitepress";
import DefaultTheme from "vitepress/theme";
import { onMounted, watch, nextTick } from "vue";
import { useRoute } from "vitepress";
import mediumZoom from "medium-zoom";

import Mermaid from "./components/Mermaid.vue";
import ApiEndpoint from "./components/ApiEndpoint.vue";
import FeatureCard from "./components/FeatureCard.vue";
import StepsGuide from "./components/StepsGuide.vue";
import "./style.css";

export default {
  extends: DefaultTheme,
  Layout: () => {
    return h(DefaultTheme.Layout, null, {
      // https://vitepress.dev/guide/extending-default-theme#layout-slots
    });
  },
  enhanceApp({ app }) {
    // 注册全局组件
    app.component("Mermaid", Mermaid);
    app.component("ApiEndpoint", ApiEndpoint);
    app.component("FeatureCard", FeatureCard);
    app.component("StepsGuide", StepsGuide);
  },
  setup() {
    const route = useRoute();
    const initZoom = () => {
      // 为所有文档中的图片添加缩放功能
      mediumZoom(".main img", {
        background: "var(--vp-c-bg)",
      });
    };
    onMounted(() => {
      initZoom();
    });
    watch(
      () => route.path,
      () => nextTick(() => initZoom()),
    );
  },
} satisfies Theme;
