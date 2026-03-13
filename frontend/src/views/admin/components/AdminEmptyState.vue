<template>
  <div
    class="admin-card flex flex-col items-center justify-center py-16 px-6 text-center max-w-md mx-auto my-8"
  >
    <div
      class="w-20 h-20 rounded-3xl bg-black/5 flex items-center justify-center text-4xl mb-6 shadow-inner"
    >
      {{ iconEmoji }}
    </div>
    <h3 class="text-xl font-bold admin-text-primary mb-2">{{ title }}</h3>
    <p class="admin-text-secondary text-sm max-w-[240px] mb-8 leading-relaxed">{{ description }}</p>
    <button
      v-if="actionText"
      @click="$emit('action')"
      class="admin-btn active:scale-95 transition-transform"
    >
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
