import * as path from 'path'

import { App, createPage } from 'vuepress'

import type { BlogPluginOptions } from '../shared'

export const createBlogPage = async (app: App, { base }: BlogPluginOptions) => {
  app.layouts['BlogMain'] = path.resolve(__dirname, '../client/BlogMain.vue')
  app.pages.push(
    await createPage(app, {
      path: base,
      frontmatter: {
        layout: 'BlogMain',
      },
    })
  )
}
