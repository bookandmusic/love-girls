<template>
  <div class="h-full flex flex-col">
    <div class="h-full flex-1 flex flex-col min-h-0 relative">
      <!-- iOS Segmented Control 风格导航 -->
      <div class="flex justify-between items-center mb-6 px-1">
        <nav
          class="p-1 bg-black/5 backdrop-blur-md rounded-2xl flex flex-1 sm:flex-initial space-x-1 overflow-x-auto no-scrollbar"
          aria-label="Tabs"
        >
          <router-link
            v-for="tab in currentTabs"
            :key="tab.name"
            :to="tab.path"
            class="flex items-center justify-center px-3 sm:px-4 py-2 rounded-xl transition-all duration-300 text-sm font-bold whitespace-nowrap"
            :class="[
              tab.current
                ? 'bg-white admin-text-primary shadow-sm scale-100'
                : 'text-gray-500 hover:text-gray-700 hover:bg-white/40 scale-95',
            ]"
            :aria-current="tab.current ? 'page' : undefined"
          >
            <BaseIcon :name="tab.icon" size="w-5 h-5" class="mr-2" />
            <span class="hidden sm:inline">{{ tab.label }}</span>
          </router-link>
        </nav>

        <!-- PC端：保留原始位置的玻璃质感按钮 -->
        <button
          v-if="currentTab !== 'wishes'"
          @click="handleAddClick"
          class="hidden sm:flex w-12 h-12 flex-shrink-0 items-center justify-center rounded-2xl bg-white/40 backdrop-blur-md border border-white/40 text-[var(--admin-accent-color)] shadow-sm hover:bg-white/60 active:scale-95 transition-all ml-4"
        >
          <BaseIcon name="add" size="w-6 h-6" />
        </button>
      </div>

      <!-- 移动端：仅在手机端显示的悬浮动作按钮 (FAB) -->
      <Transition name="fab">
        <button
          v-if="currentTab !== 'wishes'"
          @click="handleAddClick"
          class="sm:hidden fixed bottom-24 right-6 w-14 h-14 z-[60] flex items-center justify-center rounded-2xl bg-[var(--admin-accent-color)]/80 backdrop-blur-xl border border-white/20 text-white shadow-[0_8px_24px_rgba(240,173,160,0.35)] active:scale-90 active:bg-[var(--admin-accent-color)]/90 transition-all"
        >
          <BaseIcon name="add" size="w-7 h-7" />
        </button>
      </Transition>

      <div class="flex-1 min-h-0">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
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
  const currentTabInfo = tabs.value.find(tab => path.startsWith(tab.path))
  return currentTabInfo?.name || 'moments'
})

const currentTabs = computed(() => {
  const path = route.path
  return tabs.value.map(tab => ({
    ...tab,
    current: path.startsWith(tab.path),
  }))
})

const handleAddClick = () => {
  if (currentTab.value !== 'wishes') {
    router.push({ path: route.path, query: { action: 'add' } })
  }
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* FAB 动画 */
.fab-enter-active,
.fab-leave-active {
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.fab-enter-from,
.fab-leave-to {
  opacity: 0;
  transform: scale(0.5) translateY(20px);
}

.no-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
