import { defineConfig } from "vitepress";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  // 文档基础路径
  // - 本地开发和 Go API 服务器: '/docs/' (默认)
  // - GitHub Pages 部署: '/251117-go-ddd-template/'
  //
  // 使用方式：
  //   本地/Go服务器: npm run docs:build
  //   GitHub Pages:  由 GitHub Actions 自动设置 VITEPRESS_BASE
  base: process.env.VITEPRESS_BASE || "/docs/",

  title: "Go DDD Template",
  description: "基于 Go 的领域驱动设计 (DDD) 模板应用文档",

  // 主题配置
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      {
        text: "首页",
        link: "/",
      },
      {
        text: "指南",
        link: "/guide/",
      },
      {
        text: "后端",
        link: "/backend/",
      },
      {
        text: "前端",
        link: "/frontend/",
      },
      {
        text: "API 文档",
        link: "/api/",
      },
      {
        text: "开发文档",
        link: "/development/",
      },
    ],

    sidebar: {
      "/guide/": [
        {
          text: "快速入门",
          items: [
            {
              text: "指南概览",
              link: "/guide/",
            },
            {
              text: "快速开始",
              link: "/guide/getting-started",
            },
            {
              text: "配置系统",
              link: "/guide/configuration",
            },
            {
              text: "CLI 命令",
              link: "/guide/cli-commands",
            },
          ],
        },
        {
          text: "部署运维",
          items: [
            {
              text: "应用部署",
              link: "/guide/application-deployment",
            },
            {
              text: "文档部署",
              link: "/guide/docs-deployment",
            },
          ],
        },
        {
          text: "开发指南",
          items: [
            {
              text: "测试指南",
              link: "/guide/testing",
            },
            {
              text: "贡献指南",
              link: "/guide/contributing",
            },
          ],
        },
        {
          text: "示例",
          items: [
            {
              text: "Mermaid 图表",
              link: "/guide/mermaid-examples",
            },
            {
              text: "任务列表示例",
              link: "/guide/task-list-examples",
            },
          ],
        },
      ],
      "/backend/": [
        {
          text: "系统架构",
          items: [
            {
              text: "架构概览",
              link: "/backend/",
            },
            {
              text: "架构设计概览",
              link: "/backend/overview",
            },
            {
              text: "DDD + CQRS 架构",
              link: "/backend/ddd-cqrs",
            },
            {
              text: "架构迁移指南",
              link: "/backend/migration-guide",
            },
          ],
        },
        {
          text: "认证与授权",
          items: [
            {
              text: "认证机制",
              link: "/backend/authentication",
            },
            {
              text: "RBAC 权限系统",
              link: "/backend/rbac",
            },
            {
              text: "Personal Access Token",
              link: "/backend/pat",
            },
          ],
        },
        {
          text: "数据层",
          items: [
            {
              text: "PostgreSQL 架构",
              link: "/backend/postgresql",
            },
            {
              text: "Redis 架构",
              link: "/backend/redis",
            },
          ],
        },
      ],
      "/frontend/": [
        {
          text: "前端文档",
          items: [
            {
              text: "概览",
              link: "/frontend/",
            },
            {
              text: "快速开始",
              link: "/frontend/getting-started",
            },
            {
              text: "项目结构",
              link: "/frontend/project-structure",
            },
            {
              text: "API 集成",
              link: "/frontend/api-integration",
            },
          ],
        },
      ],
      "/api/": [
        {
          text: "API 参考",
          items: [
            {
              text: "概览",
              link: "/api/",
            },
            {
              text: "认证接口",
              link: "/api/auth",
            },
            {
              text: "用户接口",
              link: "/api/users",
            },
            {
              text: "缓存接口",
              link: "/api/cache",
            },
          ],
        },
      ],
      "/development/": [
        {
          text: "开发文档",
          items: [
            {
              text: "概览",
              link: "/development/",
            },
          ],
        },
        {
          text: "VitePress 文档系统",
          collapsed: false,
          items: [
            {
              text: "快速参考",
              link: "/development/quick-reference",
            },
            {
              text: "部署指南",
              link: "/development/deployment",
            },
            {
              text: "文档集成",
              link: "/development/docs-integration",
            },
            {
              text: "升级记录",
              link: "/development/upgrade",
            },
          ],
        },
        {
          text: "VitePress 功能扩展",
          collapsed: false,
          items: [
            {
              text: "Mermaid 图表",
              link: "/development/mermaid-integration",
            },
            {
              text: "功能展示",
              link: "/development/features",
            },
            {
              text: "高级功能",
              link: "/development/advanced",
            },
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
      pattern: "https://github.com/lwmacct/251117-go-ddd-template/edit/main/docs/:path",
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
    // VitePress 2.0 新增：CJK 友好的强调语法 (默认启用)
    cjkFriendlyEmphasis: true,
    // 图片懒加载
    image: {
      lazyLoading: true,
    },
    // 自定义 markdown-it 配置
    config: (md) => {
      // 自定义渲染 mermaid 代码块
      const fence = md.renderer.rules.fence!;
      md.renderer.rules.fence = (...args) => {
        const [tokens, idx] = args;
        const token = tokens[idx];
        const lang = token.info.trim();

        if (lang === "mermaid") {
          // 转换为 Mermaid 组件，使用 pre 标签保留换行符
          const code = md.utils.escapeHtml(token.content.trim());
          return `<Mermaid><pre style="display:none;">${code}</pre></Mermaid>\n`;
        }

        return fence(...args);
      };
    },
  },

  // 语言配置
  lang: "zh-CN",

  // Vite 配置 (如需自定义)
  vite: {
    // Vite 7 配置选项
  },
});
