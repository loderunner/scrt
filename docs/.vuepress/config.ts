import { defineUserConfig } from 'vuepress'
import { defaultTheme } from '@vuepress/theme-default'
import { viteBundler } from 'vuepress-vite'
import path from 'path'

import type { HeadConfig, NavbarConfig, SidebarConfig } from 'vuepress-vite'

const title = 'scrt'
const description = 'The secret manager for the command line'
const baseURL = 'https://scrt.run'

const head: HeadConfig[] = [
  [
    'link',
    {
      rel: 'apple-touch-icon',
      sizes: '57x57',
      href: '/images/favicon-57.png',
    },
  ],
  [
    'link',
    {
      rel: 'apple-touch-icon',
      sizes: '60x60',
      href: '/images/favicon-60.png',
    },
  ],
  [
    'link',
    {
      rel: 'apple-touch-icon',
      sizes: '72x72',
      href: '/images/favicon-72.png',
    },
  ],
  [
    'link',
    {
      rel: 'apple-touch-icon',
      sizes: '76x76',
      href: '/images/favicon-76.png',
    },
  ],
  [
    'link',
    {
      rel: 'apple-touch-icon',
      sizes: '114x114',
      href: '/images/favicon-114.png',
    },
  ],
  [
    'link',
    {
      rel: 'apple-touch-icon',
      sizes: '120x120',
      href: '/images/favicon-120.png',
    },
  ],
  [
    'link',
    {
      rel: 'apple-touch-icon',
      sizes: '144x144',
      href: '/images/favicon-144.png',
    },
  ],
  [
    'link',
    {
      rel: 'apple-touch-icon',
      sizes: '152x152',
      href: '/images/favicon-152.png',
    },
  ],
  [
    'link',
    {
      rel: 'apple-touch-icon',
      sizes: '180x180',
      href: '/images/favicon-180.png',
    },
  ],
  [
    'link',
    {
      rel: 'icon',
      type: 'image/png',
      sizes: '192x192',
      href: '/images/favicon-192.png',
    },
  ],
  [
    'link',
    {
      rel: 'icon',
      type: 'image/png',
      sizes: '32x32',
      href: '/images/favicon-32.png',
    },
  ],
  [
    'link',
    {
      rel: 'icon',
      type: 'image/png',
      sizes: '96x96',
      href: '/images/favicon-96.png',
    },
  ],
  [
    'link',
    {
      rel: 'icon',
      type: 'image/png',
      sizes: '16x16',
      href: '/images/favicon-16.png',
    },
  ],
  ['meta', { name: 'msapplication-TileColor', content: '#ffffff' }],
  [
    'meta',
    { name: 'msapplication-TileImage', content: '/images/favicon-144.png' },
  ],
  ['meta', { name: 'theme-color', content: '#ffffff' }],
  ['meta', { property: 'og:type', content: 'website' }],
  ['meta', { property: 'og:url', content: baseURL }],
  ['meta', { property: 'og:title', content: title }],
  ['meta', { property: 'og:description', content: description }],
  ['meta', { property: 'og:image', content: `${baseURL}/images/social.png` }],
  ['meta', { property: 'twitter:card', content: 'summary_large_image' }],
  ['meta', { property: 'twitter:url', content: baseURL }],
  ['meta', { property: 'twitter:title', content: title }],
  ['meta', { property: 'twitter:description', content: description }],
  [
    'meta',
    {
      property: 'twitter:image',
      content: `${baseURL}/images/social.png`,
    },
  ],
]

const navbar: NavbarConfig = [
  {
    text: 'Guide',
    link: '/guide',
  },
  {
    text: 'Reference',
    children: [
      {
        text: 'Commands',
        children: [
          '/reference/commands/global.md',
          '/reference/commands/init.md',
          '/reference/commands/set.md',
          '/reference/commands/get.md',
          '/reference/commands/unset.md',
        ],
      },
      {
        text: 'Storage',
        children: [
          '/reference/storage/local.md',
          '/reference/storage/s3.md',
          '/reference/storage/git.md',
        ],
      },
      {
        text: 'Configuration',
        children: [
          { text: 'Reference', link: '/reference/configuration/README.md' },
        ],
      },
    ],
  },
]

const sidebar: SidebarConfig = {
  '/guide': [
    {
      text: 'Guide',
      children: [
        '/guide/README.md',
        '/guide/installation.md',
        '/guide/getting-started.md',
        '/guide/configuration.md',
        '/guide/faq.md',
      ],
    },
  ],
  '/reference': [
    {
      text: 'Commands',
      children: [
        '/reference/commands/global.md',
        '/reference/commands/init.md',
        '/reference/commands/set.md',
        '/reference/commands/get.md',
        '/reference/commands/unset.md',
      ],
    },
    {
      text: 'Storage',
      children: [
        '/reference/storage/local.md',
        '/reference/storage/s3.md',
        '/reference/storage/git.md',
      ],
    },
    '/reference/configuration/README.md',
  ],
}

export default defineUserConfig({
  lang: 'en-US',
  title,
  description,
  head,

  base: '/',

  bundler: viteBundler({}),
  theme: defaultTheme({
    logo: '/images/logo.svg',
    navbar,
    sidebar,
  }),
  markdown: {
    code: {
      lineNumbers: false,
    },
  },

  alias: {
    '@theme/HomeFeatures.vue': path.resolve(
      __dirname,
      './components/HomeFeatures.vue'
    ),
  },
})
