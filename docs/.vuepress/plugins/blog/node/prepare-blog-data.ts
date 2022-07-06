import * as cheerio from 'cheerio'

import type { BlogData, BlogPage, BlogPluginOptions } from '../shared'
import type { GitPluginPageData } from '@vuepress/plugin-git'
import type { App, Page } from 'vuepress'

const HMR_CODE = `

if (import.meta.webpackHot) {
  import.meta.webpackHot.accept()
  if (__VUE_HMR_RUNTIME__.updateBlogData) {
    __VUE_HMR_RUNTIME__.updateBlogData(blogData)
  }
}
if (import.meta.hot) {
  import.meta.hot.accept(({ blogData }) => {
    __VUE_HMR_RUNTIME__.updateBlogData(blogData)
  })
}
`

const getPageDate = (page: Page): Date => {
  if (page.date) {
    return new Date(page.date)
  }

  if (page.data?.['git']) {
    return new Date((page as Page<GitPluginPageData>).data.git.createdTime)
  }

  return new Date(0)
}

const getPageDescription = (page: Page): string => {
  let description = page.frontmatter.description
  if (!description) {
    const $ = cheerio.load(page.contentRendered)
    description = $('p').text()
  }

  if (description.length > 250) {
    const match = description.slice(0, 250).match(/(.*)\s+/)
    description = `${match[1]}…`
  }

  return description
}

const replacer = (key: string, value: unknown) => {
  if (value instanceof Date) {
    return value.toISOString()
  }

  return value
}

export const prepareBlogData = async (
  app: App,
  { path }: BlogPluginOptions
) => {
  const pages: BlogPage[] = app.pages
    .filter((p) => p.path.startsWith(path) && p.path !== path)
    .sort((l, r) => getPageDate(r).getTime() - getPageDate(l).getTime())
    .map((p) => ({
      key: p.key,
      title: p.title,
      description: getPageDescription(p),
      date: getPageDate(p),
    }))

  console.log(pages.map((p) => p.title))

  const blogData: BlogData = { pages }

  let content = `export const blogData = ${JSON.stringify(blogData, replacer)}`

  // inject HMR code
  if (app.env.isDev) {
    content += HMR_CODE
  }

  await app.writeTemp('internal/blog-data.js', content)
}
