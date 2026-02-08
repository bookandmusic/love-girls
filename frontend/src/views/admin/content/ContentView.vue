<template>
  <div class="h-full flex flex-col">
    <h3 class="text-2xl font-bold text-gray-800 mb-6">{{ pageTitle }}</h3>

    <!-- 标签页 -->
    <div class="h-full flex-1 flex flex-col min-h-0">
      <!-- PC端标签页导航 -->
      <div class="border-b border-gray-200 flex justify-between items-center">
        <nav class="-mb-px flex space-x-2 md:space-x-8" aria-label="Tabs">
          <button
            v-for="tab in tabs"
            :key="tab.name"
            :class="[
              tab.current
                ? 'border-[#CCA4E3] text-[#FFB61E]'
                : 'border-transparent text-[#FFC773] hover:text-[#EEDEB0] hover:border-[#E4C6D0]',
              'whitespace-nowrap py-2 md:py-4 px-1 md:px-3 border-b-2 font-medium text-sm',
            ]"
            :aria-current="tab.current ? 'page' : undefined"
            @click="switchTab(tab.name)"
          >
            <BaseIcon :name="tab.icon" size="w-6 h-6" />
          </button>
        </nav>
        <button v-if="currentTab !== 'wishes'" @click="handleAddClick">
          <BaseIcon name="add" size="w-6 h-6" color="text-[#FFB61E]" />
        </button>
      </div>

      <!-- 标签页内容区域 -->
      <div class="h-full flex flex-1 flex-col pt-2 min-h-0">
        <!-- 动态管理 -->
        <MomentsManagement v-if="currentTab === 'moments'" :trigger-add="addTrigger" />

        <!-- 纪念日管理 -->
        <AnniversariesManagement v-if="currentTab === 'anniversaries'" :trigger-add="addTrigger" />

        <!-- 旅游地点管理 -->
        <PlacesManagement v-if="currentTab === 'places'" :trigger-add="addTrigger" />

        <!-- 相册管理 -->
        <AlbumsManagement v-if="currentTab === 'albums'" :trigger-add="addTrigger" />

        <!-- 访客祝福管理 -->
        <WishesManagement v-if="currentTab === 'wishes'" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

import BaseIcon from '@/components/ui/BaseIcon.vue'
import { useAddTrigger } from '@/utils/useAddTrigger'
import AlbumsManagement from '@/views/admin/content/components/AlbumsManagement.vue'
import AnniversariesManagement from '@/views/admin/content/components/AnniversariesManagement.vue'
import MomentsManagement from '@/views/admin/content/components/MomentsManagement.vue'
import PlacesManagement from '@/views/admin/content/components/PlacesManagement.vue'
import WishesManagement from '@/views/admin/content/components/WishesManagement.vue'

// 定义标签页
const tabs = ref([
  { name: 'moments', label: '动态', current: true, icon: 'moment' },
  { name: 'anniversaries', label: '纪念日', current: false, icon: 'anniversary' },
  { name: 'places', label: '足迹', current: false, icon: 'place' },
  { name: 'albums', label: '相册', current: false, icon: 'camera' },
  { name: 'wishes', label: '留言', current: false, icon: 'wish' },
])

const currentTab = ref('moments')

const { trigger: addTrigger, fire: handleAddClick } = useAddTrigger()

// 添加计算页面标题的方法
const pageTitle = computed(() => {
  const currentTabInfo = tabs.value.find(tab => tab.name === currentTab.value)
  return `${currentTabInfo?.label}管理`
})

// 切换标签页
const switchTab = (tabName: string) => {
  currentTab.value = tabName
  tabs.value.forEach(tab => {
    tab.current = tab.name === tabName
  })
}
</script>
