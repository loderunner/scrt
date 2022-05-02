const { viteBundler } = require('@vuepress/bundler-vite');
const { defaultTheme } = require('@vuepress/theme-default');

module.exports = {
  lang: 'en-US',
  title: 'scrt - Documentation',
  description: 'Documentation site for scrt, the command-line secret manager',
  bundler: viteBundler({}),
  theme: defaultTheme({}),
};
