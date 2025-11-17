import { defineConfig } from "vitepress";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  // 文档基础路径
  // - 通过 Go API 服务器访问时使用 '/docs/'
  // - GitHub Pages 部署时使用 '/251117-go-ddd-template/'
  base: "/docs/",

  title: "Go DDD Template",
  description: "基于 Go 的领域驱动设计（DDD）模板应用文档",

  // 主题配置
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: "首页", link: "/" },
      { text: "指南", link: "/guide/getting-started" },
      { text: "API 文档", link: "/api/" },
    ],

    sidebar: {
      "/guide/": [
        {
          text: "指南",
          items: [
            { text: "快速开始", link: "/guide/getting-started" },
            { text: "项目架构", link: "/guide/architecture" },
            { text: "配置系统", link: "/guide/configuration" },
            { text: "部署文档", link: "/guide/deployment" },
            { text: "贡献指南", link: "/guide/contributing" },
          ],
        },
        {
          text: "核心功能",
          items: [
            { text: "认证授权", link: "/guide/authentication" },
            { text: "PostgreSQL", link: "/guide/postgresql" },
            { text: "Redis", link: "/guide/redis" },
          ],
        },
      ],
      "/api/": [
        {
          text: "API 参考",
          items: [
            { text: "概览", link: "/api/" },
            { text: "认证接口", link: "/api/auth" },
            { text: "用户接口", link: "/api/users" },
          ],
        },
      ],
    },

    socialLinks: [
      {
        icon: "github",
        link: "https://github.com/lwmacct/251117-go-ddd-template",
      },
    ],

    footer: {
      message: "Released under the MIT License.",
      copyright: "Copyright © 2025",
    },

    // 搜索配置
    search: {
      provider: "local",
    },

    // 编辑链接
    editLink: {
      pattern:
        "https://github.com/lwmacct/251117-go-ddd-template/edit/main/docs/:path",
      text: "在 GitHub 上编辑此页",
    },

    // 最后更新时间
    lastUpdated: {
      text: "最后更新于",
      formatOptions: {
        dateStyle: "short",
        timeStyle: "medium",
      },
    },
  },

  // Markdown 配置
  markdown: {
    lineNumbers: true,
    // VitePress 2.0 新增：CJK 友好的强调语法（默认启用）
    cjkFriendlyEmphasis: true,
    // 图片懒加载
    image: {
      lazyLoading: true,
    },
  },

  // 语言配置
  lang: "zh-CN",

  // Vite 配置（如需自定义）
  vite: {
    // Vite 7 配置选项
  },
});
