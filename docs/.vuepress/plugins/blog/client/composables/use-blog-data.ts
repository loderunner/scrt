// eslint-disable-next-line import/no-unresolved
import { blogData as blogDataRaw } from '@internal/blog-data'
import { ref } from 'vue'

import type { BlogData } from '../../shared'
import type { Ref } from 'vue'

declare const __VUE_HMR_RUNTIME__: Record<string, unknown>

export type BlogDataRef = Ref<BlogData>

export const blogData: BlogDataRef = ref(blogDataRaw)

export const useBlogData = (): BlogDataRef => blogData

if (import.meta.webpackHot || import.meta.hot) {
  __VUE_HMR_RUNTIME__.updateBlogData = (data: BlogData) => {
    blogData.value = data
  }
}
