import { fileURLToPath, URL } from "node:url";

import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import vueDevTools from "vite-plugin-vue-devtools";

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), vueDevTools()],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  server: {
    port: 40013, // 前端端口
    host: true, // 或者使用 '0.0.0.0' 来监听所有网络接口
    proxy: {
      // 代理验证码请求，添加固定参数
      "/api/platform/auth/captcha": {
        target: "http://localhost:40008",
        changeOrigin: true,
        secure: false,
        rewrite: (path) => "/api/platform/auth/captcha?code=9999&secret=lwmacct",
      },
      // 代理所有 /api/platform 请求到后端服务器
      "/api/platform": {
        target: "http://localhost:40008",
        changeOrigin: true,
        secure: false,
      },
      "/api/itam": {
        target: "http://localhost:40007",
        changeOrigin: true,
        secure: false,
      },
      // 代理 Swagger 文档请求
      "/swagger": {
        target: "http://localhost:40008",
        changeOrigin: true,
        secure: false,
      },
    },
  },
});
