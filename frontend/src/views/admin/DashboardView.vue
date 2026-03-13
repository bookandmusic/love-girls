<template>
  <div class="flex flex-col pb-10">
    <h2 class="text-2xl font-bold admin-text-primary mb-6 font-(family-name:--font-signature)">
      数据概览
    </h2>

    <!-- 统计卡片网格 -->
    <div class="grid grid-cols-2 lg:grid-cols-3 gap-4 md:gap-6 mb-8">
      <!-- 相册 -->
      <div class="admin-card p-6 flex flex-col justify-between group">
        <div class="flex items-center justify-between mb-4">
          <div class="p-3 rounded-2xl bg-blue-500/10 text-blue-500 transition-transform">
            <BaseIcon name="camera" size="w-6" color="currentColor" />
          </div>
          <span class="text-xs font-bold text-blue-500 bg-blue-500/10 px-2 py-1 rounded-full"
            >相册</span
          >
        </div>
        <div>
          <p class="text-3xl font-bold admin-text-primary tracking-tight mb-1">
            {{ albumStats.total }}
          </p>
          <p class="text-xs admin-text-secondary font-medium">
            包含 {{ albumStats.totalPhotos }} 张照片
          </p>
        </div>
      </div>

      <!-- 足迹 -->
      <div class="admin-card p-6 flex flex-col justify-between group">
        <div class="flex items-center justify-between mb-4">
          <div class="p-3 rounded-2xl bg-green-500/10 text-green-500 transition-transform">
            <BaseIcon name="place" size="w-6" color="currentColor" />
          </div>
          <span class="text-xs font-bold text-green-500 bg-green-500/10 px-2 py-1 rounded-full"
            >足迹</span
          >
        </div>
        <div>
          <p class="text-3xl font-bold admin-text-primary tracking-tight mb-1">
            {{ placeStats.total }}
          </p>
          <p class="text-xs admin-text-secondary font-medium">走过的每一座城市</p>
        </div>
      </div>

      <!-- 动态 -->
      <div class="admin-card p-6 flex flex-col justify-between group">
        <div class="flex items-center justify-between mb-4">
          <div class="p-3 rounded-2xl bg-orange-500/10 text-orange-500 transition-transform">
            <BaseIcon name="moment" size="w-6" color="currentColor" />
          </div>
          <span class="text-xs font-bold text-orange-500 bg-orange-500/10 px-2 py-1 rounded-full"
            >动态</span
          >
        </div>
        <div>
          <p class="text-3xl font-bold admin-text-primary tracking-tight mb-1">
            {{ momentStats.total }}
          </p>
          <p class="text-xs admin-text-secondary font-medium">
            {{ momentStats.private }} 条私密内容
          </p>
        </div>
      </div>

      <!-- 纪念日 -->
      <div class="admin-card p-6 flex flex-col justify-between group">
        <div class="flex items-center justify-between mb-4">
          <div class="p-3 rounded-2xl bg-pink-500/10 text-pink-500 transition-transform">
            <BaseIcon name="anniversary" size="w-6" color="currentColor" />
          </div>
          <span class="text-xs font-bold text-pink-500 bg-pink-500/10 px-2 py-1 rounded-full"
            >纪念日</span
          >
        </div>
        <div>
          <p class="text-3xl font-bold admin-text-primary tracking-tight mb-1">
            {{ anniversaryStats.total }}
          </p>
          <p class="text-xs admin-text-secondary font-medium">值得珍藏的里程碑</p>
        </div>
      </div>

      <!-- 留言 -->
      <div class="admin-card p-6 flex flex-col justify-between group">
        <div class="flex items-center justify-between mb-4">
          <div class="p-3 rounded-2xl bg-purple-500/10 text-purple-500 transition-transform">
            <BaseIcon name="wish" size="w-6" color="currentColor" />
          </div>
          <span class="text-xs font-bold text-purple-500 bg-purple-500/10 px-2 py-1 rounded-full"
            >留言</span
          >
        </div>
        <div>
          <p class="text-3xl font-bold admin-text-primary tracking-tight mb-1">
            {{ wishStats.total }}
          </p>
          <p class="text-xs admin-text-secondary font-medium">{{ wishStats.pending }} 条待审核</p>
        </div>
      </div>

      <!-- 用户 -->
      <div class="admin-card p-6 flex flex-col justify-between group">
        <div class="flex items-center justify-between mb-4">
          <div class="p-3 rounded-2xl bg-indigo-500/10 text-indigo-500 transition-transform">
            <BaseIcon name="user" size="w-6" color="currentColor" />
          </div>
          <span class="text-xs font-bold text-indigo-500 bg-indigo-500/10 px-2 py-1 rounded-full"
            >用户</span
          >
        </div>
        <div>
          <p class="text-3xl font-bold admin-text-primary tracking-tight mb-1">
            {{ userStats.total }}
          </p>
          <p class="text-xs admin-text-secondary font-medium">站点的共同维护者</p>
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

const albumStats = ref({
  total: 0,
  totalPhotos: 0,
})

const placeStats = ref({
  total: 0,
})

const momentStats = ref({
  total: 0,
  private: 0,
})

const anniversaryStats = ref({
  total: 0,
})

const wishStats = ref({
  total: 0,
  pending: 0,
})

const userStats = ref({
  total: 0,
})

const uiStore = useUIStore()
const showToast = useToast()

const loadDashboardData = async () => {
  uiStore.setLoading(true)
  try {
    const response = await dashboardApi.getDashboardStats()
    if (response.code === 0) {
      const data: DashboardData = response.data

      albumStats.value = data.albumStats
      placeStats.value = data.placeStats
      momentStats.value = data.momentStats
      anniversaryStats.value = data.anniversaryStats
      wishStats.value = data.wishStats
      userStats.value = data.userStats
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
