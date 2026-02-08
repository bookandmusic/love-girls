<template>
  <div
    class="w-full h-screen flex flex-col"
    :style="{
      backgroundImage: `url(${backgroundImage})`,
      backgroundSize: 'cover',
      backgroundPosition: 'center',
    }"
  >
    <!-- 标题区域 -->
    <div class="w-full max-w-6xl mx-auto p-4 px-6 flex-shrink-0">
      <h1
        v-if="title"
        class="text-4xl md:text-5xl text-left font-bold font-(family-name:--font-signature)"
      >
        {{ title }}
      </h1>
      <p v-if="subtitle" class="text-md text-left mt-1 font-(family-name:--font-decor)">
        {{ subtitle }}
      </p>
    </div>

    <!-- 主内容区域 -->
    <div class="flex-1 flex flex-grow min-h-0 overflow-hidden">
      <!-- PC端左侧图标栏 -->
      <MenuBar class="flex-shrink-0" />

      <!-- 页面内容插槽 -->
      <div class="flex-1 flex-grow min-h-0 overflow-y-auto md:p-8 w-full flex justify-center">
        <div class="w-full max-w-6xl h-full flex flex-col overflow-hidden generic-card">
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
    <MenuBar :is-mobile="true" class="flex-shrink-0" />

    <!-- 页脚 -->
    <div class="glass-footer w-full flex-shrink-0 hidden md:block">
      <div class="max-w-4xl mx-auto py-3 xl:2 text-center">
        <p class="text-md text-gray-500">
          我们的故事开始于:
          <span class="font-bold">{{ computedStartDate }}</span>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

import bgSrc from '@/assets/images/bg.png'
import MenuBar from '@/components/business/MenuBar.vue'
import BaseIcon from '@/components/ui/BaseIcon.vue'
import { useSystemStore } from '@/stores/system'

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
.glass-footer {
  background: rgba(255, 255, 255, 0.4);
  backdrop-filter: blur(15px);
  border-top: 1px solid rgba(229, 231, 235, 0.3);
  box-shadow: 0 -2px 15px rgba(0, 0, 0, 0.05);
}
</style>
