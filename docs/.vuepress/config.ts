import { defineUserConfig } from 'vuepress'
import { defaultTheme } from '@vuepress/theme-default'
import { viteBundler } from 'vuepress-vite'
import path from 'path'

export default defineUserConfig({
  lang: 'en-US',
  title: 'scrt',
  description: 'scrt – a secret manager for the command line',

  base: '/',

  bundler: viteBundler({}),
  theme: defaultTheme({
    logo: '/images/logo.svg',
  }),
  alias: {
    '@theme/HomeFeatures.vue': path.resolve(
      __dirname,
      './components/HomeFeatures.vue'
    ),
  },
})
