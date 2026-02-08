<template>
  <CardItem>
    <template #header>
      <div class="flex items-center justify-between">
        <p class="text-base font-medium text-indigo-600 truncate">{{ wish.authorName }}</p>
        <div class="ml-2 flex flex-shrink-0">
          <span class="text-xs text-gray-500">{{ wish.email }}</span>
        </div>
      </div>
    </template>
    <template #content>
      <div class="mt-2 flex flex-wrap justify-between">
        <div class="flex items-center text-sm text-gray-500">
          <BaseIcon name="calendar" class="mr-1" size="w-4" color="text-[#FFC773]" />
          {{ wish.createdAt }}
        </div>
        <div class="flex items-center text-sm text-gray-500">
          <span
            :class="[
              wish.approved ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800',
              'inline-flex items-center px-3 py-0.5 rounded-full text-sm font-medium',
            ]"
          >
            {{ wish.approved ? '已发布' : '待审核' }}
          </span>
        </div>
      </div>
      <div class="mt-4">
        <p class="text-sm text-gray-700">{{ wish.content }}</p>
      </div>
    </template>

    <template #footer>
      <div class="mt-4 flex justify-end space-x-3">
        <button v-if="!wish.approved" @click="approveWish(wish)">
          <BaseIcon name="pass" size="w-6 h-6" color="text-[#FFB61E]" />
        </button>
        <button @click="deleteWish(wish)">
          <BaseIcon name="delete" size="w-6 h-6" color="text-[#FFB61E]" />
        </button>
      </div>
    </template>
  </CardItem>
</template>

<script setup lang="ts">
import BaseIcon from '@/components/ui/BaseIcon.vue'
import CardItem from '@/components/ui/CardItem.vue'
import { type Wish } from '@/services/wishApi'

interface Props {
  wish: Wish
}

interface Emits {
  (e: 'approve', wish: Wish): void
  (e: 'delete', wish: Wish): void
}

defineProps<Props>()
const emit = defineEmits<Emits>()

const approveWish = (wish: Wish) => {
  emit('approve', wish)
}

const deleteWish = (wish: Wish) => {
  emit('delete', wish)
}
</script>
