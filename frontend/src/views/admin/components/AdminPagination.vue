<template>
  <div class="flex justify-center py-6 flex-shrink-0 w-full sticky bottom-0 z-10">
    <div class="admin-pagination-container">
      <button @click="onFirst" :disabled="currentPage === 1" class="admin-pagination-btn">
        <BaseIcon name="left-double" size="w-4 h-4" />
      </button>

      <button @click="onPrev" :disabled="currentPage === 1" class="admin-pagination-btn">
        <BaseIcon name="left" size="w-4 h-4" />
      </button>

      <div class="flex items-center space-x-1 px-2 border-x border-black/5">
        <template v-for="page in visiblePages" :key="page">
          <button
            v-if="typeof page === 'number'"
            @click="onPageChange(page)"
            class="admin-pagination-btn"
            :class="{ active: page === currentPage }"
          >
            {{ page }}
          </button>
          <span v-else class="text-gray-400 px-1 text-xs">•••</span>
        </template>
      </div>

      <button @click="onNext" :disabled="currentPage === totalPages" class="admin-pagination-btn">
        <BaseIcon name="right" size="w-4 h-4" />
      </button>

      <button @click="onLast" :disabled="currentPage === totalPages" class="admin-pagination-btn">
        <BaseIcon name="right-double" size="w-4 h-4" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import BaseIcon from '@/components/ui/BaseIcon.vue'

defineOptions({
  name: 'AdminPagination',
})

interface Props {
  currentPage: number
  totalPages: number
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'prev'): void
  (e: 'next'): void
  (e: 'page', page: number): void
}>()

const visiblePages = computed(() => {
  const pages: (number | string)[] = []
  const { currentPage, totalPages } = props

  if (totalPages <= 7) {
    for (let i = 1; i <= totalPages; i++) {
      pages.push(i)
    }
  } else {
    if (currentPage <= 4) {
      for (let i = 1; i <= 5; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(totalPages)
    } else if (currentPage >= totalPages - 3) {
      pages.push(1)
      pages.push('...')
      for (let i = totalPages - 4; i <= totalPages; i++) {
        pages.push(i)
      }
    } else {
      pages.push(1)
      pages.push('...')
      for (let i = currentPage - 1; i <= currentPage + 1; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(totalPages)
    }
  }

  return pages
})

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

const onFirst = () => {
  if (props.currentPage > 1) {
    emit('page', 1)
  }
}

const onLast = () => {
  if (props.currentPage < props.totalPages) {
    emit('page', props.totalPages)
  }
}

const onPageChange = (page: number) => {
  if (page !== props.currentPage) {
    emit('page', page)
  }
}
</script>
