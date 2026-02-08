<template>
  <Teleport to="body">
    <transition name="fade">
      <div
        v-if="show"
        :class="['fixed top-5 right-5 px-4 py-2 rounded shadow z-[9999]', typeClass]"
      >
        {{ message }}
      </div>
    </transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue'

const props = withDefaults(
  defineProps<{
    show: boolean
    message: string
    type?: 'success' | 'error' | 'info'
    duration?: number
  }>(),
  {
    type: 'info',
    duration: 3000, // 默认3秒后自动隐藏
  }
)

const emit = defineEmits<{
  'update:show': [value: boolean]
}>()

const typeClass = computed(() => {
  switch (props.type) {
    case 'success':
      return 'bg-green-500 text-white'
    case 'error':
      return 'bg-red-500 text-white'
    default:
      return 'bg-gray-800 text-white'
  }
})

// 监听show属性变化，当变为true时启动定时器自动隐藏
watch(
  () => props.show,
  newVal => {
    if (newVal) {
      setTimeout(() => {
        // 通过emit事件通知父组件更新show状态
        emit('update:show', false)
      }, props.duration)
    }
  }
)
</script>

<style>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
