<script setup lang="ts">
import { getLunarDate } from 'chinese-days'
import { computed, onMounted, ref } from 'vue'

import MainLayout from '@/layouts/MainLayout.vue'
import { useSystemStore } from '@/stores/system'
import { useUIStore } from '@/stores/ui'
import { useToast } from '@/utils/toastUtils'

import CoupleAvatar from './components/CoupleAvatar.vue'

const uiStore = useUIStore()
const systemStore = useSystemStore()
// 使用store中的系统信息
const systemInfo = computed(() => systemStore.getSystemInfo)
const date = ref<string>('2025-12-23')
const showToast = useToast()

// 在组件挂载时设置date值，而不是在计算属性中修改
onMounted(() => {
  uiStore.setLoading(true)
  // 从store获取系统信息
  systemStore
    .fetchSystemInfo()
    .then(() => {
      if (systemInfo.value) {
        date.value = systemInfo.value.site.startDate
      }
    })
    .catch(() => {
      showToast('获取系统信息失败', 'error')
    })
    .finally(() => {
      uiStore.setLoading(false)
    })
})

const lunarDate = computed(() => {
  return getLunarDate(date.value)
})
</script>

<template>
  <MainLayout
    :title="systemInfo?.site.name"
    :subtitle="systemInfo?.site.description"
    :start-date="systemInfo?.site.startDate"
  >
    <template #main-content>
      <div class="flex items-center justify-center h-full">
        <div v-if="systemInfo" class="w-full max-w-4xl">
          <!-- 情侣信息 -->
          <CoupleAvatar :boy="systemInfo.couple.boy" :girl="systemInfo.couple.girl"></CoupleAvatar>
          <div
            class="md:hidden p-6 m-4 h-32 flex flex-col items-center justify-center text-center text-[#ff7500] font-(family-name:--font-heading)"
          >
            <div class="text-lg font-semibold mb-1">我们的故事开始于</div>
            <div class="text-2xl font-bold text-secondary-color font-(family-name:--font-math)">
              {{ date }}
            </div>
            <div class="mt-1 font-(family-name:--font-decor)">
              {{ lunarDate.lunarYearCN + '年' + lunarDate.lunarMonCN + lunarDate.lunarDayCN }}
            </div>
          </div>
        </div>
      </div>
    </template>
  </MainLayout>
</template>
