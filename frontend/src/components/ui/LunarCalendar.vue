<template>
  <div class="max-w-sm mx-auto p-6 bg-white/90 rounded-lg shadow-md">
    <!-- 头部：年月导航 -->
    <div class="flex justify-between items-center mb-4">
      <button @click="previousYear" class="p-2 text-xl rounded hover:bg-gray-200">
        <BaseIcon name="left-double" size="w-6 h-6" />
      </button>
      <button @click="previousMonth" class="p-2 text-xl rounded hover:bg-gray-200">
        <BaseIcon name="left" size="w-6 h-6" />
      </button>
      <span class="text-xl font-medium">{{ currentYear }}年{{ currentMonth + 1 }}月</span>
      <button @click="nextMonth" class="p-2 text-xl rounded hover:bg-gray-200">
        <BaseIcon name="right" size="w-6 h-6" />
      </button>
      <button @click="nextYear" class="p-2 text-xl rounded hover:bg-gray-200">
        <BaseIcon name="right-double" size="w-6 h-6" />
      </button>
    </div>

    <!-- 星期标题 -->
    <div class="grid grid-cols-7 gap-2 mb-2">
      <div
        v-for="(day, i) in ['日', '一', '二', '三', '四', '五', '六']"
        :key="i"
        class="text-center font-bold text-gray-700"
      >
        {{ day }}
      </div>
    </div>

    <!-- 日历格子 -->
    <div class="grid grid-cols-7 gap-2">
      <div
        v-for="(item, index) in calendarDays"
        :key="index"
        class="h-10 flex flex-col items-center justify-center cursor-pointer rounded transition-colors relative"
        :class="{
          'text-gray-300': item && item.currentMonth === false, // 非当前月淡灰色
          'bg-green-500 text-white': item && isSelected(item.fullDate),
          'hover:bg-green-200': item && !isSelected(item.fullDate) && item.currentMonth,
          'hover:bg-gray-100':
            item && !isSelected(item.fullDate) && !showLunar && item.currentMonth,
          'cursor-default': !item,
        }"
        @click="item && item.currentMonth && selectDate(item.date)"
      >
        <span v-if="item">{{ item.date }}</span>
        <span v-if="item && showLunar" class="text-xs text-gray-500 mt-0.5">
          {{
            getLunarDayCN(item.fullDate) === '初一'
              ? getLunarMonCN(item.fullDate)
              : getLunarDayCN(item.fullDate)
          }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { getLunarDate } from 'chinese-days'
import { computed, ref } from 'vue'

import BaseIcon from './BaseIcon.vue'

const emit = defineEmits<{
  (e: 'onSelect', date: string): void
}>()

const props = defineProps<{
  defaultDate?: string
  showLunar?: boolean
}>()

const formatDate = (y: number, m: number, d: number): string => {
  return `${y}-${String(m).padStart(2, '0')}-${String(d).padStart(2, '0')}`
}

const year = ref(new Date().getFullYear())
const month = ref(new Date().getMonth())
const today = new Date()
const todayStr = formatDate(today.getFullYear(), today.getMonth() + 1, today.getDate())

const selectedDate = ref(props.defaultDate || todayStr)

const currentYear = computed(() => year.value)
const currentMonth = computed(() => month.value)

const daysInMonth = computed(() => new Date(year.value, month.value + 1, 0).getDate())

// 获取上个月的天数
const daysInPrevMonth = computed(() => {
  const prevMonth = month.value === 0 ? 11 : month.value - 1
  const prevYear = month.value === 0 ? year.value - 1 : year.value
  return new Date(prevYear, prevMonth + 1, 0).getDate()
})

// 构建 42 格日历
const calendarDays = computed(() => {
  const result: Array<{ date: number; fullDate: string; currentMonth: boolean } | null> = []

  const firstDayOfWeek = new Date(year.value, month.value, 1).getDay()

  // 上月填充
  for (let i = firstDayOfWeek - 1; i >= 0; i--) {
    const prevDay = daysInPrevMonth.value - i
    const prevMonth = month.value === 0 ? 11 : month.value - 1
    const prevYear = month.value === 0 ? year.value - 1 : year.value
    result.push({
      date: prevDay,
      fullDate: formatDate(prevYear, prevMonth + 1, prevDay),
      currentMonth: false,
    })
  }

  // 当月日期
  for (let day = 1; day <= daysInMonth.value; day++) {
    result.push({
      date: day,
      fullDate: formatDate(year.value, month.value + 1, day),
      currentMonth: true,
    })
  }

  // 下月填充
  while (result.length < 42) {
    const nextDay = result.length - firstDayOfWeek - daysInMonth.value + 1
    const nextMonth = month.value === 11 ? 0 : month.value + 1
    const nextYear = month.value === 11 ? year.value + 1 : year.value
    result.push({
      date: nextDay,
      fullDate: formatDate(nextYear, nextMonth + 1, nextDay),
      currentMonth: false,
    })
  }

  return result
})

const isSelected = (fullDate: string) => selectedDate.value === fullDate

const selectDate = (day: number) => {
  const newDate = formatDate(year.value, month.value + 1, day)
  selectedDate.value = newDate
  emit('onSelect', newDate)
}

const getLunarDayCN = (fullDate: string) => getLunarDate(fullDate).lunarDayCN
const getLunarMonCN = (fullDate: string) => getLunarDate(fullDate).lunarMonCN

// 导航
const previousMonth = () => {
  if (month.value === 0) {
    year.value--
    month.value = 11
  } else {
    month.value--
  }
}

const nextMonth = () => {
  if (month.value === 11) {
    year.value++
    month.value = 0
  } else {
    month.value++
  }
}

const previousYear = () => year.value--
const nextYear = () => year.value++

// 初始化 defaultDate
if (props.defaultDate) {
  const [yStr, mStr] = props.defaultDate.split('-')
  const y = Number(yStr)
  const m = Number(mStr)
  if (!isNaN(y) && !isNaN(m) && m >= 1 && m <= 12) {
    year.value = y
    month.value = m - 1
  } else {
    console.error('Invalid defaultDate prop:', props.defaultDate)
  }
}
</script>
