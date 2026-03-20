<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterView, useRouter } from 'vue-router'

import SplashScreen from '@/components/SplashScreen.vue'
import { useUIStore } from '@/stores/ui'
import { isDesktopMode } from '@/utils/platform'

const router = useRouter()
const uiStore = useUIStore()
const showSplash = ref(isDesktopMode())

onMounted(() => {
  if (isDesktopMode()) {
    const unsubscribe = router.afterEach(() => {
      uiStore.setAppReady(true)
      setTimeout(() => {
        showSplash.value = false
      }, 300)
      unsubscribe()
    })
  }
})
</script>

<template>
  <div class="min-h-screen bg-background font-(family-name:--font-body) overflow-x-hidden">
    <SplashScreen v-if="showSplash" />
    <RouterView />
  </div>
</template>

<style scoped>
.bg-background {
  background-color: var(--background);
}

.overflow-x-hidden {
  overflow-x: hidden;
  overscroll-behavior: none;
}
</style>
