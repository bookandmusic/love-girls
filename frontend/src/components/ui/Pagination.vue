<template>
  <div class="flex justify-center py-4 flex-shrink-0 w-full">
    <div class="glass-thick rounded-full p-1 flex items-center border border-white/40 shadow-sm">
      <button
        @click="onPrev"
        :disabled="currentPage === 1"
        class="p-2 rounded-full ios-transition tap-feedback disabled:opacity-30 disabled:pointer-events-none"
      >
        <BaseIcon name="left" size="w-5 h-5" color="var(--fe-text-primary)" />
      </button>

      <div class="px-4 py-1 flex items-center text-sm font-semibold text-[var(--fe-text-primary)]">
        <span class="mr-1">{{ currentPage }}</span>
        <span class="opacity-30 mx-1">/</span>
        <span class="opacity-60">{{ totalPages }}</span>
      </div>

      <button
        @click="onNext"
        :disabled="currentPage === totalPages"
        class="p-2 rounded-full ios-transition tap-feedback disabled:opacity-30 disabled:pointer-events-none"
      >
        <BaseIcon name="right" size="w-5 h-5" color="var(--fe-text-primary)" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import BaseIcon from '@/components/ui/BaseIcon.vue'
// 添加 name 选项以符合多词组件名称规范
defineOptions({
  name: 'CommonPagination',
})

interface Props {
  currentPage: number
  totalPages: number
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'prev'): void
  (e: 'next'): void
}>()

const onPrev = () => {
  if (props.currentPage > 1) {
    emit('prev')
  }
}

const onNext = () => {
  if (props.currentPage < props.totalPages) {
    emit('next')
  }
}
</script>
