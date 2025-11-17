import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  // GitHub Pages 部署的 base 路径
  // 如果部署到 https://username.github.io/，设置为 '/'
  // 如果部署到 https://username.github.io/repo/，设置为 '/repo/'
  base: '/251117-bd-vmalert/',

  title: "Go DDD Template",
  description: "基于 Go 的领域驱动设计（DDD）模板应用文档",

  // 主题配置
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: '首页', link: '/' },
      { text: '指南', link: '/guide/getting-started' },
      { text: 'API 文档', link: '/api/' }
    ],

    sidebar: {
      '/guide/': [
        {
          text: '指南',
          items: [
            { text: '快速开始', link: '/guide/getting-started' },
            { text: '项目架构', link: '/guide/architecture' },
            { text: '配置系统', link: '/guide/configuration' },
            { text: '部署文档', link: '/guide/deployment' }
          ]
        },
        {
          text: '核心功能',
          items: [
            { text: '认证授权', link: '/guide/authentication' },
            { text: 'PostgreSQL', link: '/guide/postgresql' },
            { text: 'Redis', link: '/guide/redis' }
          ]
        }
      ],
      '/api/': [
        {
          text: 'API 参考',
          items: [
            { text: '概览', link: '/api/' },
            { text: '认证接口', link: '/api/auth' },
            { text: '用户接口', link: '/api/users' }
          ]
        }
      ]
    },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/lwmacct/251117-bd-vmalert' }
    ],

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2025'
    },

    // 搜索配置
    search: {
      provider: 'local'
    },

    // 编辑链接
    editLink: {
      pattern: 'https://github.com/lwmacct/251117-bd-vmalert/edit/main/docs/:path',
      text: '在 GitHub 上编辑此页'
    },

    // 最后更新时间
    lastUpdated: {
      text: '最后更新于',
      formatOptions: {
        dateStyle: 'short',
        timeStyle: 'medium'
      }
    }
  },

  // Markdown 配置
  markdown: {
    lineNumbers: true
  },

  // 语言配置
  lang: 'zh-CN'
})
