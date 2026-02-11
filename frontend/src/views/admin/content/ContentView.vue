<template>
  <div class="h-full flex flex-col">
    <!-- 标签页导航 -->
    <div class="h-full flex-1 flex flex-col min-h-0">
      <!-- PC端标签页导航 -->
      <div class="border-b border-gray-200 flex justify-between items-center">
        <nav class="-mb-px flex space-x-2 md:space-x-8" aria-label="Tabs">
          <router-link
            v-for="tab in currentTabs"
            :key="tab.name"
            :to="tab.path"
            :class="[
              tab.current
                ? 'border-[#CCA4E3] text-[#FFB61E]'
                : 'border-transparent text-[#FFC773] hover:text-[#EEDEB0] hover:border-[#E4C6D0]',
              'whitespace-nowrap py-2 md:py-4 px-1 md:px-3 border-b-2 font-medium text-sm',
            ]"
            :aria-current="tab.current ? 'page' : undefined"
          >
            <BaseIcon :name="tab.icon" size="w-6 h-6" />
          </router-link>
        </nav>
        <button v-if="currentTab !== 'wishes'" @click="handleAddClick">
          <BaseIcon name="add" size="w-6 h-6" color="text-[#FFB61E]" />
        </button>
      </div>

      <!-- 路由视图内容区域 -->
      <div class="h-full flex flex-1 flex-col pt-2 min-h-0">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import BaseIcon from '@/components/ui/BaseIcon.vue'

const route = useRoute()
const router = useRouter()

// 定义标签页
const tabs = ref([
  { name: 'moments', label: '动态', path: '/admin/content/moments', icon: 'moment' },
  {
    name: 'anniversaries',
    label: '纪念日',
    path: '/admin/content/anniversaries',
    icon: 'anniversary',
  },
  { name: 'places', label: '足迹', path: '/admin/content/places', icon: 'place' },
  { name: 'albums', label: '相册', path: '/admin/content/albums', icon: 'camera' },
  { name: 'wishes', label: '留言', path: '/admin/content/wishes', icon: 'wish' },
])

// 计算当前标签页
const currentTab = computed(() => {
  const path = route.path
  const currentTabInfo = tabs.value.find(tab => tab.path === path)
  return currentTabInfo?.name || 'moments'
})

// 根据路由计算当前标签页状态
const currentTabs = computed(() => {
  const path = route.path
  return tabs.value.map(tab => ({
    ...tab,
    current: tab.path === path,
  }))
})

// 处理添加按钮点击
const handleAddClick = () => {
  if (currentTab.value !== 'wishes') {
    // 导航到当前页面并添加查询参数
    router.push({ path: route.path, query: { action: 'add' } })
  }
}
</script>
