<template>
  <div class="flex flex-col items-center justify-center py-12 px-4">
    <div class="text-6xl mb-4">{{ iconEmoji }}</div>
    <h3 class="text-xl font-semibold admin-text-primary mb-2">{{ title }}</h3>
    <p class="admin-text-secondary text-center mb-6">{{ description }}</p>
    <button v-if="actionText" @click="$emit('action')" class="admin-btn">
      {{ actionText }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

defineOptions({
  name: 'AdminEmptyState',
})

interface Props {
  icon?: string
  title: string
  description: string
  actionText?: string
}

const props = withDefaults(defineProps<Props>(), {
  icon: 'moment',
})

defineEmits<{
  (e: 'action'): void
}>()

const iconEmoji = computed(() => {
  const iconMap: Record<string, string> = {
    moment: '💕',
    album: '📷',
    anniversary: '💕',
    place: '📍',
    wish: '💌',
    user: '👥',
    box: '📦',
  }
  return iconMap[props.icon] || '📦'
})
</script>
