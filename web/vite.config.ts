import { fileURLToPath, URL } from "node:url";

import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import vueDevTools from "vite-plugin-vue-devtools";
import vuetify from "vite-plugin-vuetify";
import AutoImport from "unplugin-auto-import/vite";
import Components from "unplugin-vue-components/vite";

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    // Vuetify 按需自动导入（处理 Vuetify 组件和指令）
    vuetify({ autoImport: true }),
    // Vue API 自动导入
    AutoImport({
      imports: [
        "vue",
        "vue-router",
        "pinia",
        "@vueuse/core",
        {
          // 自定义导入
          "@/stores/auth": ["useAuthStore"],
        },
      ],
      dts: "src/auto-imports.d.ts",
      vueTemplate: true,
    }),
    // 自定义组件自动注册（Vuetify 组件由 vite-plugin-vuetify 处理）
    Components({
      dts: "src/components.d.ts",
      dirs: ["src/components"],
      // 排除 Vuetify 组件，避免与 vite-plugin-vuetify 冲突
      resolvers: [],
    }),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  // 生产构建时移除 console 和 debugger
  esbuild: {
    drop: process.env.NODE_ENV === "production" ? ["console", "debugger"] : [],
  },
  server: {
    port: 40013, // 前端端口
    host: true, // 或者使用 '0.0.0.0' 来监听所有网络接口
    proxy: {
      // 代理验证码请求，添加开发模式固定参数
      "/api/auth/captcha": {
        target: "http://localhost:40012",
        changeOrigin: true,
        secure: false,
        rewrite: (_path) => "/api/auth/captcha?code=9999&secret=dev-secret-change-me",
      },
      // 代理所有 /api 请求到后端服务器
      "/api": {
        target: "http://localhost:40012",
        changeOrigin: true,
        secure: false,
      },
      // 代理 Swagger 文档请求
      "/swagger": {
        target: "http://localhost:40012",
        changeOrigin: true,
        secure: false,
      },
    },
  },
});
