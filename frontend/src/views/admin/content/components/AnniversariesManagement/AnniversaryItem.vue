<template>
  <CardItem>
    <template #header>
      <div class="flex items-center justify-between">
        <p class="text-sm font-medium text-indigo-600 truncate">{{ anniversary.title }}</p>
        <div class="ml-2 flex flex-shrink-0">
          <span
            :class="[
              'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
              anniversary.calendar === 'lunar'
                ? 'bg-purple-100 text-purple-800'
                : 'bg-blue-100 text-blue-800',
            ]"
          >
            {{ anniversary.calendar === 'lunar' ? '农历' : '公历' }}
          </span>
        </div>
      </div>
    </template>

    <template #content>
      <div class="mt-2 sm:flex sm:justify-between">
        <div class="sm:flex">
          <div class="mr-6 flex items-center text-sm text-gray-500">
            <BaseIcon name="calendar" size="w-4" class="mr-1" color="text-[#FFC773]" />
            {{ formatDate(anniversary.date, anniversary.calendar) }}
          </div>
        </div>
      </div>
      <div class="mt-4 flex justify-between">
        <p class="text-sm text-gray-500">{{ anniversary.description }}</p>
      </div>
    </template>

    <template #footer>
      <div class="flex justify-end space-x-3">
        <button @click="$emit('edit', anniversary)">
          <BaseIcon name="edit" size="w-6 h-6" color="text-[#FFB61E]" />
        </button>
        <button @click="$emit('delete', anniversary)">
          <BaseIcon name="delete" size="w-6 h-6" color="text-[#FFB61E]" />
        </button>
      </div>
    </template>
  </CardItem>
</template>

<script setup lang="ts">
import { getLunarDate, getSolarDateFromLunar } from 'chinese-days'

import BaseIcon from '@/components/ui/BaseIcon.vue'
import CardItem from '@/components/ui/CardItem.vue'
import { type Anniversary } from '@/services/anniversaryApi'

interface Props {
  anniversary: Anniversary
}

defineProps<Props>()

defineEmits(['edit', 'delete'])

// 格式化日期显示
const formatDate = (dateStr: string, calendarType: 'solar' | 'lunar'): string => {
  switch (calendarType) {
    case 'lunar':
      const solarDate = getSolarDateFromLunar(dateStr)
      const lunarDate = getLunarDate(solarDate.date)
      // 农历日期格式化
      return `${lunarDate.lunarMonCN}${lunarDate.lunarDayCN}`
    case 'solar':
      const [, month, day] = dateStr.split('-')
      return `${month}月${day}日`
  }
}
</script>
