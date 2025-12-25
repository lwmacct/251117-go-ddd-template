import { defineConfig, type DefaultTheme } from "vitepress";
import nav from "./config/nav.json";
import sidebarGuide from "./config/sidebar.guide.json";
import sidebarArchitecture from "./config/sidebar.architecture.json";
import sidebarDevelopment from "./config/sidebar.development.json";
import sidebarApi from "./config/sidebar.api.json";
import sidebarRoadmap from "./config/sidebar.roadmap.json";
import cfgSearch from "./config/search.json";
import viteConfig from "./config/vite";
import markdownConfig from "./config/markdown";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Go DDD Template",
  description: "基于 Go 的领域驱动设计应用模板",
  base: process.env.BASE || "/",
  srcDir: "content",

  // Vite 构建优化配置 (从 ./config/vite.ts 导入)
  vite: viteConfig,

  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav,
    sidebar: [...sidebarGuide, ...sidebarArchitecture, ...sidebarDevelopment, ...sidebarApi, ...sidebarRoadmap],

    // 本地搜索 - 使用 MiniSearch 实现浏览器内索引
    search: cfgSearch as DefaultTheme.Config["search"],

    socialLinks: [{ icon: "github", link: "https://github.com/lwmacct/251117-go-ddd-template" }],

    footer: {
      message: "基于 DDD + CQRS 架构的企业级应用模板",
      copyright: "Copyright © 2024 Go DDD Template",
    },

    editLink: {
      pattern: "https://github.com/lwmacct/251117-go-ddd-template/edit/main/docs/:path",
      text: "在 GitHub 上编辑此页",
    },

    lastUpdated: {
      text: "最后更新于",
      formatOptions: {
        dateStyle: "short",
        timeStyle: "medium",
      },
    },

    outline: {
      label: "页面导航",
      level: [2, 3],
    },

    docFooter: {
      prev: "上一页",
      next: "下一页",
    },

    returnToTopLabel: "回到顶部",
    sidebarMenuLabel: "菜单",
    darkModeSwitchLabel: "主题",
    lightModeSwitchTitle: "切换到浅色模式",
    darkModeSwitchTitle: "切换到深色模式",
  },

  // Markdown 渲染配置 (从 ./config/markdown.ts 导入)
  markdown: markdownConfig,

  // 忽略死链接检查（开发阶段）
  ignoreDeadLinks: [
    // 忽略本地主机链接
    /^https?:\/\/localhost/,
    // 忽略内部锚点
    /^#/,
  ],
});
