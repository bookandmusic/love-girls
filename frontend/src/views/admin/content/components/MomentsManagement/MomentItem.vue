<template>
  <CardItem>
    <template #header>
      <div class="mt-2 flex flex-wrap sm:flex-nowrap sm:justify-between gap-2">
        <div class="flex flex-wrap gap-2 min-w-0 flex-1">
          <div class="flex items-center text-sm text-gray-500 whitespace-nowrap">
            <BaseIcon name="user" size="w-4" class="mr-1" color="text-[#FFC773]" />
            {{ moment.author.name }}
          </div>
          <div class="flex items-center text-sm text-gray-500 whitespace-nowrap">
            <BaseIcon name="calendar" size="w-4" class="mr-1" color="text-[#FFC773]" />
            {{ moment.createdAt }}
          </div>
          <div class="flex flex-shrink-0 whitespace-nowrap">
            <span class="flex items-center">
              <BaseIcon name="like" size="w-4" class="mr-1" color="text-[#FFC773]" />
              {{ moment.likes }}
            </span>
          </div>
        </div>
        <div class="flex items-center text-sm text-gray-500 flex-shrink-0">
          <span
            :class="[
              'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
              moment.isPublic ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800',
            ]"
          >
            {{ moment.isPublic ? '公开' : '私有' }}
          </span>
        </div>
      </div>
    </template>

    <template #content>
      <div class="flex items-center justify-between">
        <p class="text-md font-medium">{{ moment.content }}</p>
      </div>
      <div class="mt-2 flex flex-wrap">
        <div v-for="image in moment.images" :key="image.id" class="w-1/2 md:w-1/6 p-1">
          <img
            :src="image.file?.thumbnail || image.file?.url || ''"
            class="w-full h-full object-cover"
          />
        </div>
      </div>
    </template>

    <template #footer>
      <div class="pt-2 flex justify-end space-x-3">
        <button @click="$emit('togglePublic', moment)">
          <BaseIcon
            :name="moment.isPublic ? 'lock' : 'unlock'"
            size="w-6 h-6"
            color="text-[#FFB61E]"
          />
        </button>
        <button @click="$emit('edit', moment)">
          <BaseIcon name="edit" size="w-6 h-6" color="text-[#FFB61E]" />
        </button>
        <button @click="$emit('delete', moment)">
          <BaseIcon name="delete" size="w-6 h-6" color="text-[#FFB61E]" />
        </button>
      </div>
    </template>
  </CardItem>
</template>

<script setup lang="ts">
import BaseIcon from '@/components/ui/BaseIcon.vue'
import CardItem from '@/components/ui/CardItem.vue'
import { type Moment } from '@/services/momentApi'

defineProps({
  moment: {
    type: Object as () => Moment,
    required: true,
  },
})

defineEmits(['togglePublic', 'edit', 'delete'])
</script>
