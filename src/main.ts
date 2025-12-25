import { createApp } from "vue";
import { createPinia } from "pinia";

import App from "./App.vue";
import router from "./router";
import { useAuthStore } from "@/stores/auth";
import { loadingBarPlugin } from "@/plugins/loadingBar";

// Vuetify 样式和图标
import "vuetify/styles";
import "@mdi/font/css/materialdesignicons.css";
import { createVuetify } from "vuetify";
import { aliases, mdi } from "vuetify/iconsets/mdi";

const app = createApp(App);

// 获取默认主题，处理 'auto' 情况
function getDefaultTheme(): string {
  const storedTheme = localStorage.getItem("theme") || "light";
  if (storedTheme === "auto") {
    // 如果是 'auto'，根据系统主题偏好选择
    const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
    return prefersDark ? "dark" : "light";
  }
  // 确保只有 'light' 或 'dark'
  return storedTheme === "dark" ? "dark" : "light";
}

// Vuetify 配置 - 组件和指令由 vite-plugin-vuetify 自动按需导入
const vuetify = createVuetify({
  icons: {
    defaultSet: "mdi",
    aliases,
    sets: {
      mdi,
    },
  },
  theme: {
    defaultTheme: getDefaultTheme(),
    themes: {
      dark: {
        dark: true,
        colors: {
          primary: "#1867C0",
          secondary: "#5CBBF6",
          error: "#CF6679",
          success: "#4CAF50",
          warning: "#FF9800",
        },
      },
      light: {
        dark: false,
        colors: {
          primary: "#42A5F5",
          secondary: "#757575",
          error: "#B00020",
          success: "#4CAF50",
          warning: "#FB8C00",
        },
      },
    },
  },
});

app.use(vuetify);

const pinia = createPinia();
app.use(pinia);

app.use(router);

// 路由加载进度条插件
app.use(loadingBarPlugin, { router });

/**
 * 初始化应用
 * 在挂载前恢复认证状态，确保路由守卫能正确判断用户是否已登录
 */
async function initializeApp() {
  const authStore = useAuthStore();

  // 从 localStorage 恢复认证状态
  await authStore.initAuth();

  // 挂载应用
  app.mount("#app");
}

// 启动应用
initializeApp().catch((error) => {
  console.error("[App Init] Failed to initialize app:", error);
  // 即使初始化失败也要挂载应用，让用户能访问登录页
  app.mount("#app");
});
