<template>
  <GenericDialog
    :open="open"
    :title="props.anniversary?.id ? '编辑纪念日' : '添加纪念日'"
    @cancel="closeDialog"
    :loading="props.loading"
  >
    <template #content>
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1"
          >标题 <span class="text-red-500">*</span></label
        >
        <input
          v-model="form.title"
          type="text"
          class="w-full win11-input"
          placeholder="请输入纪念日标题"
          :disabled="props.loading"
        />
      </div>

      <div class="mt-4 grid grid-cols-2 gap-x-4">
        <div class="relative">
          <label class="block text-sm font-medium text-gray-700 mb-1"
            >日期 <span class="text-red-500">*</span></label
          >
          <!-- 关键：加上 ref 和 readonly -->
          <input
            ref="inputRef"
            v-model="form.date"
            type="text"
            class="w-full win11-input"
            placeholder="请选择日期"
            readonly
            @focus="openCalendar"
            :disabled="props.loading"
          />
          <!-- 日历面板：关键加 @mousedown.prevent 和 ref -->
          <div
            v-if="isCalendarVisible"
            ref="calendarRef"
            class="absolute mt-1 z-10 w-80"
            @mousedown.prevent
          >
            <LunarCalendar
              :default-date="calendarDate"
              :show-lunar="true"
              @onSelect="handleDateChange"
            />
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1"
            >日历类型 <span class="text-red-500">*</span></label
          >
          <select
            v-model="form.calendar"
            class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            :disabled="props.loading"
          >
            <option value="solar">公历</option>
            <option value="lunar">农历</option>
          </select>
        </div>
      </div>

      <div class="mt-4">
        <label class="block text-sm font-medium text-gray-700 mb-1">描述</label>
        <textarea
          v-model="form.description"
          class="w-full win11-input"
          rows="4"
          placeholder="请输入纪念日描述"
          :disabled="props.loading"
        ></textarea>
      </div>
    </template>
    <template #actions>
      <div class="w-full flex">
        <div class="flex-1 text-center cursor-pointer" @click="closeDialog">取消</div>
        <div
          class="w-1/2 border-l border-gray-300 text-center cursor-pointer text-blue-500"
          @click="handleSave"
        >
          确认
        </div>
      </div>
    </template>
  </GenericDialog>

  <!-- 确认对话框 -->
  <GenericDialog
    :open="showConfirmDialog"
    title="确认保存"
    :loading="props.loading"
    size-class="w-md h-md"
    @cancel="cancelConfirm"
  >
    <template #content>
      <p class="text-gray-700">
        {{
          props.anniversary?.id
            ? '您确定要保存对这个纪念日的更改吗？'
            : '您确定要添加这个新纪念日吗？'
        }}
      </p>
    </template>
    <template #actions>
      <div class="w-full flex">
        <div class="flex-1 text-center cursor-pointer" @click="cancelConfirm">取消</div>
        <div
          class="w-1/2 border-l border-gray-300 text-center cursor-pointer text-blue-500"
          @click="confirmSave"
        >
          确定
        </div>
      </div>
    </template>
  </GenericDialog>
</template>

<script setup lang="ts">
import { getLunarDate, getSolarDateFromLunar } from 'chinese-days'
import { computed, ref, watch } from 'vue'

import GenericDialog from '@/components/ui/GenericDialog.vue'
import LunarCalendar from '@/components/ui/LunarCalendar.vue'
import { type Anniversary } from '@/services/anniversaryApi'
import { useToast } from '@/utils/toastUtils'

const showToast = useToast()

// refs
const inputRef = ref<HTMLInputElement | null>(null)

// props & emits
interface Props {
  open: boolean
  loading: boolean
  anniversary?: Anniversary | null
}

const props = withDefaults(defineProps<Props>(), {
  anniversary: null,
  loading: false,
})

interface Emits {
  (e: 'update:open', open: boolean): void
  (e: 'confirm', anniversary: Anniversary): void
  (e: 'cancel'): void
}
const emit = defineEmits<Emits>()

const closeDialog = () => {
  emit('update:open', false)
  emit('cancel')
}

// form state
const DEFAULT_ANNIVERSARY: Anniversary = {
  id: 0,
  title: '',
  date: '',
  description: '',
  calendar: 'solar' as 'solar' | 'lunar',
}

const form = ref<Anniversary>({ ...DEFAULT_ANNIVERSARY })

const isCalendarVisible = ref(false)
const showConfirmDialog = ref(false)

const calendarDate = computed(() => {
  switch (form.value.calendar) {
    case 'lunar':
      const solarDate = getSolarDateFromLunar(form.value.date)
      return solarDate.date
    default:
      return form.value.date
  }
})

// 打开日历
const openCalendar = () => {
  if (!props.loading) {
    isCalendarVisible.value = true
  }
}

// 选择日期
const handleDateChange = (selectedDate: string) => {
  let tmpDate = selectedDate
  switch (form.value.calendar) {
    case 'lunar':
      const lunarDate = getLunarDate(selectedDate)
      tmpDate = `${lunarDate.lunarYear}-${lunarDate.lunarMon}-${lunarDate.lunarDay}`
    default:
      break
  }
  form.value.date = tmpDate
  isCalendarVisible.value = false
  inputRef.value?.blur() // 可选：主动收起焦点
}

// 保存逻辑
const handleSave = async () => {
  // 校验必填项
  if (!form.value.title.trim()) {
    showToast('请输入标题', 'error')
    return
  }
  if (!form.value.date) {
    showToast('请选择日期', 'error')
    return
  }
  if (!form.value.calendar) {
    showToast('请选择日历类型', 'error')
    return
  }
  showConfirmDialog.value = true
}

const cancelConfirm = () => {
  showConfirmDialog.value = false
}

const confirmSave = async () => {
  showConfirmDialog.value = false
  const anniversary: Anniversary = {
    id: props.anniversary?.id ?? 0,
    title: form.value.title,
    date: form.value.date,
    description: form.value.description,
    calendar: form.value.calendar,
  }
  emit('confirm', anniversary)
}

// 监听 anniversary 变化初始化表单
watch(
  () => props.anniversary,
  newVal => {
    if (newVal) {
      form.value = {
        id: newVal.id,
        title: newVal.title,
        date: newVal.date,
        description: newVal.description,
        calendar: newVal.calendar,
      }
    } else {
      const today = new Date()
      const formattedDate = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}`
      form.value = {
        ...DEFAULT_ANNIVERSARY,
        date: formattedDate,
      }
    }
  },
  { deep: true, immediate: true }
)
</script>
