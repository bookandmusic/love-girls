<template>
  <div class="h-full flex flex-col">
    <div class="h-full flex-1 flex flex-col min-h-0">
      <div class="border-b border-white/60 flex justify-between items-center pb-2">
        <nav class="-mb-px flex space-x-1 md:space-x-2" aria-label="Tabs">
          <router-link
            v-for="tab in currentTabs"
            :key="tab.name"
            :to="tab.path"
            class="flex items-center justify-center p-2 md:p-3 rounded-xl transition-all duration-200"
            :class="[
              tab.current
                ? 'bg-gradient-to-r from-[#f0ada0] to-[#d89388] text-white shadow-md'
                : 'text-[#f0ada0] hover:bg-white/40',
            ]"
            :aria-current="tab.current ? 'page' : undefined"
          >
            <BaseIcon :name="tab.icon" size="w-6 h-6" />
          </router-link>
        </nav>
        <button
          v-if="currentTab !== 'wishes'"
          @click="handleAddClick"
          class="p-2 rounded-xl bg-gradient-to-r from-[#f0ada0] to-[#d89388] text-white hover:shadow-md transition-all"
        >
          <BaseIcon name="add" size="w-6 h-6" />
        </button>
      </div>

      <div class="h-full flex flex-1 flex-col pt-4 min-h-0">
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

const currentTab = computed(() => {
  const path = route.path
  const currentTabInfo = tabs.value.find(tab => tab.path === path)
  return currentTabInfo?.name || 'moments'
})

const currentTabs = computed(() => {
  const path = route.path
  return tabs.value.map(tab => ({
    ...tab,
    current: tab.path === path,
  }))
})

const handleAddClick = () => {
  if (currentTab.value !== 'wishes') {
    router.push({ path: route.path, query: { action: 'add' } })
  }
}
</script>
