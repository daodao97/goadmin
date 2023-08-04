import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "GoAdmin",
  description: "一个配置化开发后台管理系统",
  lastUpdated: true,
  cleanUrls: true,
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
    ],

    sidebar: [
      {
        text: '介绍',
        items: [
          { text: '系统架构', link: '/guide/design-philosophy' },
          { text: '快速开始', link: '/guide/getting-started' },
          { text: '示例', link: '/guide/example' },
        ]
      },
      {
        text: '页面配置',
        items: [
          { text: 'PageSchema', link: '/page/page-schema' },
          { text: '列表', link: '/page/table' },
          { text: '筛选项', link: '/page/filter' },
          { text: '表单', link: '/page/form' },
          { text: '按钮', link: '/page/button' },
          { text: '自定义组件', link: '/' },
        ]
      },
      {
        text: '后端开发',
        items: [
          { text: '依赖注入', link: '/' },
          { text: '路由', link: '/' },
          { text: '权限', link: '/' },
          { text: 'db', link: '/' },
          { text: '缓存', link: '/' },
        ]
      },
      {
        text: '其他',
        items: [
          { text: '文件上传', link: '/' },
          { text: '统一登录', link: '/' },
          { text: '模块', link: '/' },
        ]
      },
    ],

    editLink: {
      pattern: 'https://github.com/daodao97/goadmin/edit/dcos/docs/:path',
      text: 'Edit this page on GitHub'
    },


    socialLinks: [
      { icon: 'github', link: 'https://github.com/daodao97/goadmin' }
    ]
  }
})
