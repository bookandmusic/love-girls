<template>
  <CardItem>
    <template #background>
      <div
        v-if="place.image"
        class="absolute inset-0 bg-cover bg-center"
        :style="{ backgroundImage: `url(${place.image?.file?.thumbnail})` }"
      />
    </template>

    <!-- 遮罩 -->
    <template #overlay>
      <div class="absolute inset-0 bg-black/45"></div>
    </template>
    <template #header>
      <div class="flex items-center justify-between">
        <p class="text-xl font-medium text-white truncate">{{ place.name }}</p>
      </div>
    </template>

    <template #content>
      <div class="mt-2 flex items-center space-x-4">
        <div class="flex text-sm text-white/90">
          <BaseIcon name="map" size="w-4" class="mr-1" color="text-[#FFC773]" />
          {{ place.latitude }}, {{ place.longitude }}
        </div>
        <div class="flex text-sm text-white/90">
          <BaseIcon name="calendar" size="w-4" class="mr-1" color="text-[#FFC773]" />
          {{ place.date }}
        </div>
      </div>
    </template>

    <template #footer>
      <div class="mt-4 flex justify-end space-x-3">
        <button @click="onEdit">
          <BaseIcon name="edit" size="w-6 h-6" color="text-[#FFB61E]" />
        </button>
        <button @click="onDelete">
          <BaseIcon name="delete" size="w-6 h-6" color="text-[#FFB61E]" />
        </button>
      </div>
    </template>
  </CardItem>
</template>

<script setup lang="ts">
import BaseIcon from '@/components/ui/BaseIcon.vue'
import CardItem from '@/components/ui/CardItem.vue'
import type { Place } from '@/services/placeApi'

interface Props {
  place: Place
}

interface Emits {
  (e: 'edit', place: Place): void
  (e: 'delete', place: Place): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const onEdit = () => {
  emit('edit', props.place)
}

const onDelete = () => {
  emit('delete', props.place)
}
</script>
