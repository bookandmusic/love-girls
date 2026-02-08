AnniversariesView.vue
<script setup lang="ts">
import { getSolarDateFromLunar } from 'chinese-days'
import { computed, onMounted, ref } from 'vue'

import BaseIcon from '@/components/ui/BaseIcon.vue'
import MainLayout from '@/layouts/MainLayout.vue'
import { type Anniversary, anniversaryApi } from '@/services/anniversaryApi'
import { useSystemStore } from '@/stores/system'
import { useUIStore } from '@/stores/ui'
import { useToast } from '@/utils/toastUtils'

interface NormalizedAnniversary {
  anniversary: Anniversary
  nextDate: Date
}

const uiStore = useUIStore()
const systemStore = useSystemStore()

const anniversaries = ref<Anniversary[]>([])
const systemInfo = computed(() => systemStore.getSystemInfo)
const showToast = useToast()
/* ======================
 * 数据获取
 * ====================== */
const fetchAnniversaries = async () => {
  try {
    uiStore.setLoading(true)
    const res = await anniversaryApi.getAnniversaries(1, 100) // 获取所有纪念日
    anniversaries.value = res.data.anniversaries
  } catch {
    showToast('获取纪念日失败', 'error')
  } finally {
    uiStore.setLoading(false)
  }
}

/* ======================
 * 日期工具
 * ====================== */
const parseMonthDay = (date: string) => {
  const [, month, day] = date.split('-').map(Number)
  if (!month || !day || isNaN(month) || isNaN(day)) return null
  return { month, day }
}

const toSolarDate = (anniversary: Anniversary, year: number): Date | null => {
  const parsed = parseMonthDay(anniversary.date)
  if (!parsed) return null

  const { month, day } = parsed

  try {
    if (anniversary.calendar === 'lunar') {
      // 拼接成官方要求的 yyyy-MM-dd 形式
      const lunarStr = `${year.toString().padStart(4, '0')}-${month
        .toString()
        .padStart(2, '0')}-${day.toString().padStart(2, '0')}`

      const result = getSolarDateFromLunar(lunarStr)

      // result.date 是 YYYY-MM-DD
      if (result && result.date) {
        const [y, m, d] = result.date.split('-').map(Number)
        if (y != null && m != null && d != null && !isNaN(y) && !isNaN(m) && !isNaN(d)) {
          return new Date(y, m - 1, d)
        }
      }
      return null
    }

    // 阳历直接用
    return new Date(year, month - 1, day)
  } catch (err) {
    console.error('转换农历失败:', err)
    return null
  }
}

const getNextDate = (anniversary: Anniversary): Date | null => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  const year = today.getFullYear()

  let date = toSolarDate(anniversary, year)
  if (!date) return null

  if (date < today) {
    date = toSolarDate(anniversary, year + 1)
  }

  return date
}

/* ======================
 * 标准化纪念日（唯一数据源）
 * ====================== */
const normalizedAnniversaries = computed<NormalizedAnniversary[]>(() => {
  const list: NormalizedAnniversary[] = []

  for (const anniversary of anniversaries.value) {
    const nextDate = getNextDate(anniversary)
    if (nextDate) {
      list.push({ anniversary, nextDate })
    }
  }

  return list.sort((a, b) => a.nextDate.getTime() - b.nextDate.getTime())
})

/* ======================
 * 页面使用数据
 * ====================== */
const sortedAnniversaries = computed(() =>
  normalizedAnniversaries.value.map(item => item.anniversary)
)

const upcomingAnniversaryInfo = computed(() => {
  if (normalizedAnniversaries.value.length === 0) return null

  const today = new Date()
  today.setHours(0, 0, 0, 0)

  const next = normalizedAnniversaries.value[0]
  if (!next) return null

  // 计算天数差
  const diffTime = next.nextDate.getTime() - today.getTime()
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))

  return {
    title: next.anniversary.title,
    date: next.nextDate,
    diffDays,
    description: next.anniversary.description,
    isToday: diffDays === 0,
    anniversary: next.anniversary,
    days: diffDays,
    message:
      diffDays === 0
        ? `今天是${next.anniversary.title}`
        : `距离${next.anniversary.title}还有${diffDays}天`,
  }
})

/* ======================
 * 展示辅助
 * ====================== */
const formatLunarDate = (month: number, day: number) => {
  const months = [
    '正月',
    '二月',
    '三月',
    '四月',
    '五月',
    '六月',
    '七月',
    '八月',
    '九月',
    '十月',
    '十一月',
    '十二月',
  ]
  const days = [
    '初一',
    '初二',
    '初三',
    '初四',
    '初五',
    '初六',
    '初七',
    '初八',
    '初九',
    '初十',
    '十一',
    '十二',
    '十三',
    '十四',
    '十五',
    '十六',
    '十七',
    '十八',
    '十九',
    '二十',
    '廿一',
    '廿二',
    '廿三',
    '廿四',
    '廿五',
    '廿六',
    '廿七',
    '廿八',
    '廿九',
    '三十',
    '卅一',
  ]
  return `${months[month - 1]}${days[day - 1]}`
}

const formatAnniversaryDate = (anniversary: Anniversary) => {
  const parsed = parseMonthDay(anniversary.date)
  if (!parsed) return '无效日期'

  const { month, day } = parsed

  return anniversary.calendar === 'lunar'
    ? formatLunarDate(month, day)
    : `${month}月${day.toString().padStart(2, '0')}日`
}

/* ======================
 * 状态判断（置灰 / 高亮）
 * ====================== */
const isPastAnniversary = (anniversary: Anniversary) => {
  if (!normalizedAnniversaries.value.length) return false

  // 获取当前日期
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  // 获取该纪念日在今年的日期
  const currentYear = today.getFullYear()
  const anniversaryInCurrentYear = toSolarDate(anniversary, currentYear)
  if (!anniversaryInCurrentYear) return false

  // 判断今年的纪念日是否已过且不是即将到来的那个
  const isPastThisYear = anniversaryInCurrentYear < today
  return isPastThisYear
}

/* ======================
 * 生命周期
 * ====================== */
onMounted(async () => {
  uiStore.setLoading(true)
  await systemStore.fetchSystemInfo()
  await fetchAnniversaries()
  uiStore.setLoading(false)
})
</script>

<template>
  <MainLayout
    title="专属纪念"
    subtitle="珍藏我们的浪漫约定"
    :start-date="systemInfo?.site.startDate"
    :show-empty-state="anniversaries.length === 0"
  >
    <template #empty-state>
      <BaseIcon name="anniversary" size="w-24" />
      <p class="text-xl font-medium mt-4">暂无纪念日数据</p>
      <p class="text-md mt-2">还没有添加任何纪念日</p>
    </template>

    <template #main-content>
      <div class="flex flex-col h-full">
        <!-- 下一个纪念日信息 -->
        <div
          v-if="upcomingAnniversaryInfo"
          class="p-3 flex justify-center items-center text-center border-b border-white/50 flex-shrink-0"
        >
          <div class="flex items-center space-x-2">
            <div class="text-xl font-bold font-[Ma_Shan_Zheng]">
              {{ upcomingAnniversaryInfo.message }}
            </div>

            <div v-if="!upcomingAnniversaryInfo.isToday">
              <BaseIcon name="clock" size="w-6" color="text-[#ff7500]" />
            </div>

            <div v-else>
              <BaseIcon name="calendar" size="w-6" color="text-[#ff7500]" />
            </div>
          </div>
        </div>

        <!-- 时间轴 -->
        <div class="p-4 flex-grow overflow-y-auto flex flex-col">
          <!-- 当有纪念日数据时显示时间轴 -->
          <div class="relative flex-grow">
            <!-- 时间轴线 -->
            <div
              class="absolute left-2 top-0 bottom-0 w-0.5 bg-pink-300 transform translate-x-[-1px]"
            ></div>

            <!-- 纪念日列表 -->
            <div class="space-y-8 pb-4">
              <div
                v-for="(anniversary, index) in sortedAnniversaries"
                :key="anniversary.id"
                class="relative flex items-start"
                :class="{ 'pt-8': index === 0 && upcomingAnniversaryInfo }"
              >
                <!-- 时间轴节点 -->
                <div
                  class="absolute left-4 w-4 h-4 rounded-full border-4 border-white transform translate-x-[-100%] z-10"
                  :class="{
                    'bg-pink-500':
                      upcomingAnniversaryInfo &&
                      anniversary.id === upcomingAnniversaryInfo.anniversary.id,
                    'bg-pink-300':
                      !upcomingAnniversaryInfo ||
                      (upcomingAnniversaryInfo &&
                        anniversary.id !== upcomingAnniversaryInfo.anniversary.id),
                  }"
                ></div>

                <!-- 箭头指示器 -->
                <div
                  v-if="
                    upcomingAnniversaryInfo &&
                    anniversary.id === upcomingAnniversaryInfo.anniversary.id
                  "
                  class="absolute left-2 transform translate-x-[-50%] translate-y-[10%] z-10"
                  :style="{ top: '-1rem' }"
                >
                  <BaseIcon name="arrow-down" size="w-12" color="text-[#ff7500]" />
                </div>

                <!-- 内容 -->
                <div class="ml-8 flex-grow">
                  <div
                    class="generic-card p-4"
                    :class="{
                      'opacity-50': isPastAnniversary(anniversary),
                    }"
                  >
                    <h3 class="font-bold text-lg flex items-center gap-2">
                      {{ anniversary.title }}
                    </h3>
                    <div class="mt-2 flex justify-between items-center">
                      <span class="text-sm text-gray-500 whitespace-nowrap ml-2">{{
                        formatAnniversaryDate(anniversary)
                      }}</span>
                      <span
                        v-if="anniversary.calendar === 'lunar'"
                        class="text-xs px-2 py-0.5 rounded-full bg-amber-100 text-amber-600"
                      >
                        农历
                      </span>
                    </div>
                    <p class="mt-2 text-md text-gray-700">{{ anniversary.description }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </MainLayout>
</template>
