<template>
  <Teleport to="body">
    <div
      v-if="open"
      class="fixed inset-0 z-50 flex items-center justify-center transition-opacity p-4"
      @click="handleCancel"
    >
      <div
        :class="[
          'backdrop-blur-2xl w-full p-6 rounded-lg shadow-xl transform transition-all flex flex-col',
          sizeClass,
        ]"
        @click.stop
      >
        <!-- 对话框标题和关闭按钮 -->
        <div class="flex justify-between items-start mb-4">
          <slot name="header">
            <h3 class="text-lg font-medium text-gray-900">{{ title }}</h3>
          </slot>
          <button @click="handleCancel" :disabled="loading" class="text-[var(--primary-color)]">
            <BaseIcon name="close" size="w-6 h-6" color="text-[#FFB61E]" />
          </button>
        </div>

        <div class="space-y-4 flex-1 overflow-y-hidden">
          <!-- 内容通过插槽传递 -->
          <slot name="content"></slot>
        </div>

        <!-- 底部动作区域 -->
        <div class="mt-6 flex flex-wrap gap-3 justify-end border-t border-gray-300 pt-4">
          <slot name="actions"></slot>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import BaseIcon from './BaseIcon.vue'

// 组件的 Props
interface Props {
  open: boolean
  title?: string
  loading?: boolean
  sizeClass?: string
}

withDefaults(defineProps<Props>(), {
  title: '对话框',
  loading: false,
  sizeClass: 'max-w-2xl h-full',
})

const emit = defineEmits<{
  'update:open': [value: boolean]
  cancel: []
}>()

const handleCancel = () => {
  // 触发取消事件，通知父组件关闭对话框
  emit('cancel')
  emit('update:open', false)
}
</script>

<style scoped>
/* 可以根据实际需求进行样式调整 */
</style>
