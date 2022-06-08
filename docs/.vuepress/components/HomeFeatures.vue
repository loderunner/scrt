<script setup lang="ts">
import { usePageFrontmatter } from '@vuepress/client'
import { isArray } from '@vuepress/shared'
import { computed } from 'vue'

import type { DefaultThemeHomePageFrontmatter } from '@vuepress/theme-default/lib/shared/index'

type Feature = DefaultThemeHomePageFrontmatter['features'][number] & {
  image?: string
}
type Features = Feature[]

const frontmatter = usePageFrontmatter<DefaultThemeHomePageFrontmatter>()
const features = computed<Features>(() => {
  if (isArray(frontmatter.value.features)) {
    return frontmatter.value.features
  }
  return []
})
</script>

<template>
  <div v-if="features.length" class="features">
    <div v-for="feature in features" :key="feature.title" class="feature">
      <img v-if="feature.image" :src="feature.image" />
      <h2>{{ feature.title }}</h2>
      <p>{{ feature.details }}</p>
    </div>
  </div>
</template>
