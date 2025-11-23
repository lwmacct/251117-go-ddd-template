import { defineConfig } from "vitepress";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  // 文档基础路径
  // - 本地开发和 Go API 服务器: '/docs/' (默认)
  // - GitHub Pages 部署: '/251117-go-ddd-template/'
  //
  // 使用方式：
  //   本地/Go服务器: npm run build
  //   GitHub Pages:  由 GitHub Actions 自动设置 BASE
  base: process.env.BASE || "/docs/",

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
        text: "架构",
        link: "/architecture/",
      },
      {
        text: "工程实践",
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
              text: "前端概览",
              link: "/guide/frontend-overview",
            },
            {
              text: "配置系统",
              link: "/guide/configuration",
            },
          ],
        },
        {
          text: "部署与运维",
          items: [
            {
              text: "生产环境快速部署",
              link: "/guide/production-quickstart",
            },
            {
              text: "应用部署（完整指南）",
              link: "/guide/application-deployment",
            },
            {
              text: "文档部署",
              link: "/guide/docs-deployment",
            },
          ],
        },
        {
          text: "工具与质量",
          items: [
            {
              text: "CLI 命令",
              link: "/guide/cli-commands",
            },
            {
              text: "测试指南",
              link: "/guide/testing",
            },
          ],
        },
        {
          text: "前端指南",
          items: [
            {
              text: "前端概览",
              link: "/guide/frontend-overview",
            },
            {
              text: "前端快速开始",
              link: "/guide/frontend-getting-started",
            },
            {
              text: "前端项目结构",
              link: "/guide/frontend-project-structure",
            },
            {
              text: "前端 API 集成",
              link: "/guide/frontend-api-integration",
            },
          ],
        },
        {
          text: "协作与贡献",
          items: [
            {
              text: "贡献指南",
              link: "/guide/contributing",
            },
            {
              text: "参考资料",
              link: "/guide/reference",
            },
          ],
        },
      ],
      "/architecture/": [
        {
          text: "架构蓝图",
          items: [
            {
              text: "概览",
              link: "/architecture/",
            },
            {
              text: "DDD + CQRS",
              link: "/architecture/ddd-cqrs",
            },
            {
              text: "分层与目录",
              link: "/architecture/architecture-layers",
            },
            {
              text: "迁移指南",
              link: "/architecture/migration-guide",
            },
          ],
        },
        {
          text: "数据与基础设施",
          items: [
            {
              text: "PostgreSQL 架构",
              link: "/architecture/data-postgresql",
            },
            {
              text: "Redis 架构",
              link: "/architecture/data-redis",
            },
          ],
        },
        {
          text: "身份与访问控制",
          items: [
            {
              text: "身份能力总览",
              link: "/architecture/identity-overview",
            },
            {
              text: "认证机制",
              link: "/architecture/identity-authentication",
            },
            {
              text: "RBAC 权限系统",
              link: "/architecture/identity-rbac",
            },
            {
              text: "Personal Access Token",
              link: "/architecture/identity-pat",
            },
          ],
        },
      ],
      "/development/": [
        {
          text: "开发指南",
          items: [
            {
              text: "概览",
              link: "/development/",
            },
            {
              text: "AI Agent",
              link: "/development/ai-agent",
            },
          ],
        },
        {
          text: "交付与升级",
          collapsed: false,
          items: [
            {
              text: "部署指南",
              link: "/development/deployment",
            },
            {
              text: "升级记录",
              link: "/development/upgrade",
            },
          ],
        },
        {
          text: "质量与效率",
          collapsed: false,
          items: [
            {
              text: "Pre-commit 代码检查",
              link: "/development/pre-commit",
            },
            {
              text: "功能示例",
              link: "/development/features",
            },
            {
              text: "主题与高级能力",
              link: "/development/advanced",
            },
          ],
        },
        {
          text: "文档与可视化",
          collapsed: false,
          items: [
            {
              text: "文档集成",
              link: "/development/docs-integration",
            },
            {
              text: "Mermaid 图表",
              link: "/development/mermaid-integration",
            },
          ],
        },
      ],
      "/reference/": [
        {
          text: "参考资料",
          items: [
            {
              text: "概览",
              link: "/reference/",
            },
            {
              text: "Admin Users API",
              link: "/reference/admin-users-api",
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

    // 右侧大纲展示 H2-H3
    outline: [2, 3],
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
