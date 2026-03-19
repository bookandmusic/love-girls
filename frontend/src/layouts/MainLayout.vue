<template>
  <div
    class="w-full h-screen flex flex-col frontend-root relative overflow-hidden"
    :style="{
      backgroundImage: `url(${backgroundImage})`,
      backgroundSize: 'cover',
      backgroundPosition: 'center',
    }"
  >
    <!-- 背景遮罩 - 可选，用于提升文字清晰度 -->
    <div class="absolute inset-0 bg-white/10 pointer-events-none"></div>

    <!-- 标题区域 -->
    <div class="w-full max-w-6xl mx-auto p-4 px-6 flex-shrink-0 relative z-10 pt-8">
      <div>
        <h1
          v-if="title"
          class="text-4xl md:text-5xl text-left font-bold font-(family-name:--font-signature) text-[var(--fe-text-primary)]"
        >
          {{ title }}
        </h1>
        <p
          v-if="subtitle"
          class="text-md text-left mt-2 font-(family-name:--font-decor) text-[var(--fe-text-secondary)]"
        >
          {{ subtitle }}
        </p>
      </div>
    </div>

    <!-- 主内容区域 -->
    <div class="flex-1 flex flex-grow min-h-0 overflow-hidden relative z-10">
      <!-- PC端左侧图标栏 -->
      <MenuBar class="flex-shrink-0" />

      <!-- 页面内容插槽 -->
      <div class="flex-1 flex-grow min-h-0 overflow-y-auto md:p-8 w-full flex justify-center">
        <div
          class="w-full max-w-6xl h-full flex flex-col overflow-hidden bg-transparent md:glass-regular md:rounded-[var(--fe-radius-card)] md:border md:border-white/40 md:shadow-xl"
        >
          <div v-if="showEmptyState" class="h-full w-full p-4 flex items-center justify-center">
            <div class="text-[#FF7500] flex flex-col justify-center items-center text-center">
              <slot name="empty-state">
                <!-- 默认的空状态内容 -->
                <BaseIcon name="box" size="w-16" style="color: #ff7500" />
                <h3 class="font-bold text-lg mt-4">暂无数据</h3>
                <p class="mt-2 text-center">还没有任何内容</p>
              </slot>
            </div>
          </div>
          <div v-else class="h-full flex flex-col">
            <slot name="main-content"></slot>
          </div>
        </div>
      </div>
    </div>

    <!-- 手机端底部图标栏 -->
    <MenuBar :is-mobile="true" class="flex-shrink-0 relative z-20" />

    <!-- 页脚 (PC可见) -->
    <div
      class="glass-ultra-thin w-full flex-shrink-0 hidden md:block border-t border-white/20 relative z-10"
    >
      <div class="max-w-4xl mx-auto py-3 text-center">
        <p class="text-sm text-[var(--fe-text-secondary)]">
          我们的故事开始于:
          <span class="font-bold text-[var(--fe-text-primary)]">{{ computedStartDate }}</span>
        </p>
      </div>
    </div>
  </div>

  <!-- 服务切换按钮 - Teleport 到 body 避免 overflow-hidden 影响 -->
  <Teleport to="body">
    <div v-if="isDesktop" class="fixed top-8 right-6 z-[9999]">
      <DesktopMenu />
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import '@/assets/frontend-theme.css'

import { computed, ref } from 'vue'

import bgSrc from '@/assets/images/bg.png'
import MenuBar from '@/components/business/MenuBar.vue'
import BaseIcon from '@/components/ui/BaseIcon.vue'
import DesktopMenu from '@/components/ui/DesktopMenu.vue'
import { useSystemStore } from '@/stores/system'
import { isDesktopMode } from '@/utils/platform'

const isDesktop = isDesktopMode()

// 定义组件属性
interface Props {
  title?: string
  subtitle?: string
  startDate?: string
  showEmptyState?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  title: '',
  subtitle: '',
  startDate: '',
  showEmptyState: false,
})

const systemStore = useSystemStore()

// 设置背景图片
const backgroundImage = ref(bgSrc)

// 计算属性：优先使用传入的startDate，否则使用store中的数据
const computedStartDate = computed(() => {
  if (props.startDate) {
    return props.startDate
  }
  return systemStore.getSystemInfo?.site.startDate || ''
})
</script>

<style scoped>
/* 局部样式，全局滚动条已在 frontend-theme.css 处理 */
</style>
