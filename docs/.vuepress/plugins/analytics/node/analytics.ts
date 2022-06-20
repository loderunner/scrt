import path from 'path'

import type { App, Page, Plugin, PluginObject } from 'vuepress'

interface Options {
  dataDomain: string
  src: string
}

export interface AnalyticsPageData {
  analytics?: Options
}

export const analyticsPlugin =
  (options: Options): Plugin =>
  (app: App): PluginObject => {
    const plugin: PluginObject = {
      name: 'analytics-plugin',
    }

    if (app.env.isDev) {
      return plugin
    }

    return {
      ...plugin,
      clientConfigFile: path.resolve(__dirname, '../client/enhance.ts'),
      extendsPage: (page: Page<AnalyticsPageData>) => {
        page.data.analytics = { ...options }
      },
    }
  }
