<script setup lang="ts">
import ParentLayout from '@vuepress/theme-default/lib/client/layouts/Layout.vue'

import { useBlogData } from './composables/use-blog-data'

const blogData = useBlogData()
</script>

<template>
  <ParentLayout>
    <template #page-content-top>
      <div>
        <div v-for="page of blogData.pages" :key="page.key">
          <router-link :to="{ name: page.key }">{{ page.title }}</router-link>
        </div>
      </div>
    </template>
  </ParentLayout>
</template>
