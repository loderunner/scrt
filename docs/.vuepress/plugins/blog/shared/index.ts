export interface BlogPluginOptions {
  base?: string
  path?: string
}

export interface BlogPage {
  key: string
  title: string
  description: string
  date: Date
}

export interface BlogData {
  pages: BlogPage[]
}
