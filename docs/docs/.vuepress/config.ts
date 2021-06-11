import { defineUserConfig } from 'vuepress';
import type { DefaultThemeOptions } from 'vuepress';

const productionMode = process.env.NODE_ENV === 'production';

export default defineUserConfig<DefaultThemeOptions>({
  base: '/',

  lang: 'en-US',
  title: 'scrt',
  description: 'Secret manager for the command line',

  bundler: productionMode ? '@vuepress/webpack' : '@vuepress/vite',
});
