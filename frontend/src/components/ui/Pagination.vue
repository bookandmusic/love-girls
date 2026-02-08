<template>
  <div class="flex justify-center py-2 flex-shrink-0 w-full">
    <div class="flex space-x-2">
      <button
        @click="onPrev"
        :disabled="currentPage === 1"
        class="px-3 py-1 rounded-md bg-white/50 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-white/80"
      >
        上一页
      </button>

      <span class="px-3 py-1 flex items-center"> {{ currentPage }} / {{ totalPages }} </span>

      <button
        @click="onNext"
        :disabled="currentPage === totalPages"
        class="px-3 py-1 rounded-md bg-white/50 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-white/80"
      >
        下一页
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
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
