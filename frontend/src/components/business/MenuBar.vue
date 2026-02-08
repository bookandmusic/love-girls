<script setup lang="ts">
// 导入菜单图标
import anniversariesIcon from '@/assets/menu/anniversaries.png'
import homeIcon from '@/assets/menu/home.png'
import momentsIcon from '@/assets/menu/moments.png'
import photosIcon from '@/assets/menu/photos.png'
import placesIcon from '@/assets/menu/places.png'
import wishesIcon from '@/assets/menu/wishes.png'

// 定义菜单项类型
interface MenuItem {
  icon: string
  label: string
  path: string
}

// 定义菜单项
const menuItems: MenuItem[] = [
  { icon: homeIcon, label: '首页', path: '/' },
  { icon: momentsIcon, label: '动态', path: '/moments' },
  { icon: photosIcon, label: '相册', path: '/albums' },
  { icon: placesIcon, label: '足迹', path: '/places' },
  { icon: anniversariesIcon, label: '纪念日', path: '/anniversaries' },
  { icon: wishesIcon, label: '祝福', path: '/wishes' },
]

defineProps<{
  isMobile?: boolean
}>()
</script>

<template>
  <div v-if="!isMobile" class="hidden md:flex flex-col items-center justify-center space-y-4 w-32">
    <div
      v-for="(item, index) in menuItems"
      :key="index"
      class="flex flex-col items-center cursor-pointer group"
    >
      <RouterLink :to="item.path">
        <img
          :src="item.icon"
          :alt="item.label"
          class="w-20 h-20 object-contain transition-all duration-300 group-hover:scale-125"
        />
      </RouterLink>
    </div>
  </div>

  <!-- 手机端底部图标栏 -->
  <div v-else class="md:hidden grid grid-cols-6 gap-3 py-3">
    <RouterLink v-for="(item, index) in menuItems" :key="index" :to="item.path">
      <div class="flex flex-col items-center cursor-pointer group">
        <img
          :src="item.icon"
          :alt="item.label"
          class="w-16 h-16 object-contain transition-all duration-300 group-hover:scale-110"
        />
      </div>
    </RouterLink>
  </div>
</template>
