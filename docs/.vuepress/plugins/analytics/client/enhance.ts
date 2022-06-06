import { defineClientAppSetup, usePageData } from '@vuepress/client'

import type { AnalyticsPageData } from '../node'

export default defineClientAppSetup(() => {
  const pageData = usePageData<AnalyticsPageData>()
  const { analytics } = pageData.value

  if (!analytics) {
    return
  }

  const script = document.createElement('script')
  script.defer = true
  script.src = analytics.src
  script.setAttribute('data-domain', analytics.dataDomain)
  document.head.appendChild(script)
})
