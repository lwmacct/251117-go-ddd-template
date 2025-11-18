import { createApp } from "vue";
import { createPinia } from "pinia";

import App from "./App.vue";
import router from "./router";

// https://vuetifyjs.com/en/introduction/why-vuetify/
import "vuetify/styles";
import "@mdi/font/css/materialdesignicons.css";
import { createVuetify } from "vuetify";
import * as components from "vuetify/components";
import * as directives from "vuetify/directives";
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

const vuetify = createVuetify({
  components,
  directives,
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
app.mount("#app");
