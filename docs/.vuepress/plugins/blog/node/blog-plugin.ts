import * as path from 'path'

import gitPlugin from '@vuepress/plugin-git'
import * as chokidar from 'chokidar'

import { createBlogPage } from './create-blog-page'
import { prepareBlogData } from './prepare-blog-data'

import type { BlogPluginOptions } from '../shared'
import type { App, Plugin, PluginObject } from 'vuepress'

const defaultOptions: BlogPluginOptions = { path: '/blog', base: '/blog' }

export const blogPlugin =
  (options: BlogPluginOptions): Plugin =>
  (app: App): PluginObject => {
    options = { ...defaultOptions, ...options }

    app.use(gitPlugin({}))

    return {
      name: 'blog-plugin',

      onInitialized: (app) => createBlogPage(app, options),
      onPrepared: (app) => prepareBlogData(app, options),
      onWatched: (app, watchers) => {
        console.log(path.join(options.path, '*'))
        const blogPageWatcher = chokidar.watch('blog/*', {
          cwd: app.dir.source(),
          ignoreInitial: true,
        })
        console.log(blogPageWatcher)
        blogPageWatcher.on('add', () => prepareBlogData(app, options))
        blogPageWatcher.on('change', () => prepareBlogData(app, options))
        blogPageWatcher.on('unlink', () => prepareBlogData(app, options))
        watchers.push(blogPageWatcher)
      },
    }
  }
