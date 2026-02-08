<template>
  <div class="h-full w-full overflow-y-hidden flex flex-col">
    <h2 class="text-2xl font-bold text-gray-800 mb-6">数据概览</h2>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-2 xl:grid-cols-4 gap-3 md:gap-6 mb-4">
      <div class="bg-white/30 rounded-xl shadow-sm px-3 md:px-6 py-6 border border-gray-100">
        <div class="flex items-center">
          <div class="p-3 rounded-lg bg-blue-50">
            <BaseIcon name="camera" size="w-6" color="text-blue-500" />
          </div>
          <div class="ml-4">
            <h3 class="text-md md:text-lg font-semibold text-gray-600">相册</h3>
            <p class="text-xl md:text-2xl font-bold text-gray-900 font-(family-name:--font-math)">
              {{ albumStats.total }}
            </p>
          </div>
        </div>
        <div class="mt-4">
          <p class="text-sm text-gray-500">包含 {{ albumStats.totalPhotos }} 张照片</p>
        </div>
      </div>

      <div class="bg-white/30 rounded-xl shadow-sm px-3 md:px-6 py-6 border border-gray-100">
        <div class="flex items-center">
          <div class="p-3 rounded-lg bg-green-50">
            <BaseIcon name="place" size="w-6" color="text-green-500" />
          </div>
          <div class="ml-4">
            <h3 class="text-md md:text-lg font-semibold text-gray-600">足迹</h3>
            <p class="text-xl md:text-2xl font-bold text-gray-900 font-(family-name:--font-math)">
              {{ placeStats.total }}
            </p>
          </div>
        </div>
      </div>

      <div class="bg-white/30 rounded-xl shadow-sm px-3 md:px-6 py-6 border border-gray-100">
        <div class="flex items-center">
          <div class="p-3 rounded-lg bg-pink-50">
            <BaseIcon name="anniversary" size="w-6" color="text-pink-500" />
          </div>
          <div class="ml-4">
            <h3 class="text-md md:text-lg font-semibold text-gray-600">纪念日</h3>
            <p class="text-xl md:text-2xl font-bold text-gray-900 font-(family-name:--font-math)">
              {{ momentStats.total }}
            </p>
          </div>
        </div>
      </div>

      <div class="bg-white/30 rounded-xl shadow-sm px-3 md:px-6 py-6 border border-gray-100">
        <div class="flex items-center">
          <div class="p-3 rounded-lg bg-purple-50">
            <BaseIcon name="wish" size="w-6" color="text-purple-500" />
          </div>
          <div class="ml-4">
            <h3 class="text-md md:text-lg font-semibold text-gray-600">祝福</h3>
            <p class="text-xl md:text-2xl font-bold text-gray-900 font-(family-name:--font-math)">
              {{ wishStats.total }}
            </p>
          </div>
        </div>
        <div class="mt-4">
          <p class="text-sm text-gray-500">{{ wishStats.pending }} 条待审核</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'

import BaseIcon from '@/components/ui/BaseIcon.vue'
import { dashboardApi, type DashboardData } from '@/services/dashboardApi'
import { useUIStore } from '@/stores/ui'
import { useToast } from '@/utils/toastUtils'

// 统计数据
const albumStats = ref({
  total: 0,
  totalPhotos: 0,
})

const placeStats = ref({
  total: 0,
})

const momentStats = ref({
  total: 0,
})

const wishStats = ref({
  total: 0,
  pending: 0,
})

const uiStore = useUIStore()
const showToast = useToast()

// 加载仪表盘数据
const loadDashboardData = async () => {
  uiStore.setLoading(true)
  try {
    const response = await dashboardApi.getDashboardStats()
    if (response.code === 0) {
      const data: DashboardData = response.data

      albumStats.value = data.albumStats
      placeStats.value = data.placeStats
      momentStats.value = data.momentStats
      wishStats.value = data.wishStats
    } else {
      showToast('获取仪表盘数据失败', 'error')
    }
  } catch {
    showToast('获取仪表盘数据失败', 'error')
  } finally {
    uiStore.setLoading(false)
  }
}

onMounted(() => {
  loadDashboardData()
})
</script>
