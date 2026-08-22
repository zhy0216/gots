import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: 'goTS',
  description: 'A TypeScript-like language that compiles to Go',
  lang: 'en-US',
  cleanUrls: true,
  lastUpdated: true,
  srcExclude: ['**/CLAUDE.md'],
  head: [['link', { rel: 'icon', href: '/favicon.ico' }]],

  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    logo: '/logo.svg',
    nav: [
      { text: 'Guide', link: '/guide/getting-started' },
      { text: 'Language Spec', link: '/language-spec' },
      { text: 'Built-in Reference', link: '/built-in-reference' },
      { text: 'Development', link: '/development/architecture' },
    ],

    sidebar: {
      '/guide/': [
        {
          text: 'Guide',
          items: [
            { text: 'Getting Started', link: '/guide/getting-started' },
            { text: 'CLI Reference', link: '/guide/cli' },
          ],
        },
      ],
      '/development/': [
        {
          text: 'Development',
          items: [
            { text: 'Architecture', link: '/development/architecture' },
            { text: 'Testing', link: '/development/testing' },
            { text: 'FAQ', link: '/faq' },
          ],
        },
      ],
    },

    outline: {
      level: [2, 3],
      label: 'On this page',
    },

    socialLinks: [{ icon: 'github', link: 'https://github.com/zhy0216/quickts' }],

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2026 goTS contributors',
    },

    search: {
      provider: 'local',
    },

    editLink: {
      pattern: 'https://github.com/zhy0216/quickts/edit/main/gots/docs/:path',
      text: 'Edit this page on GitHub',
    },

    lastUpdated: {
      text: 'Updated at',
      formatOptions: {
        dateStyle: 'medium',
        timeStyle: 'short',
      },
    },
  },
})
