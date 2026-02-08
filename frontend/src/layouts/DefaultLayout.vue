<template>
  <div class="w-full h-screen flex flex-col">
    <LoadingSpinner v-show="uiStore.loading" />
    <div v-show="!uiStore.loading" class="flex-grow flex flex-col">
      <WaveBackground :isPlaying="uiStore.playing" :waveColor="'#f0ada0'" class="flex-grow">
        <ToastNotification v-model:show="toast.show" :message="toast.message" :type="toast.type" />
        <RouterView />
      </WaveBackground>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { provide, ref } from 'vue'

import LoadingSpinner from '@/components/ui/LoadingSpinner.vue'
import ToastNotification from '@/components/ui/ToastNotification.vue'
import WaveBackground from '@/components/ui/WaveBackground.vue'
import { useUIStore } from '@/stores/ui'

const uiStore = useUIStore()

console.log(uiStore.playing)

// Toast 通知状态
const toast = ref({
  show: false,
  message: '',
  type: undefined as 'success' | 'error' | 'info' | undefined,
})

// 提供显示 toast 的函数给子组件使用
const showToast = (message: string, type: 'success' | 'error' | 'info' | undefined = 'info') => {
  toast.value = {
    show: true,
    message,
    type,
  }
}

// 提供给子组件使用的 toast 相关函数
provide('showToast', showToast)
</script>
