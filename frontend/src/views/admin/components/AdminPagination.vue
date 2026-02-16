<template>
  <div class="flex justify-center py-4 flex-shrink-0 w-full">
    <div class="flex items-center space-x-2">
      <button @click="onFirst" :disabled="currentPage === 1" class="admin-pagination-btn">
        首页
      </button>

      <button @click="onPrev" :disabled="currentPage === 1" class="admin-pagination-btn">
        上一页
      </button>

      <div class="flex items-center space-x-1">
        <template v-for="page in visiblePages" :key="page">
          <button
            v-if="typeof page === 'number'"
            @click="onPageChange(page)"
            class="admin-pagination-btn"
            :class="{ active: page === currentPage }"
          >
            {{ page }}
          </button>
          <span v-else class="admin-text-muted px-2">...</span>
        </template>
      </div>

      <button @click="onNext" :disabled="currentPage === totalPages" class="admin-pagination-btn">
        下一页
      </button>

      <button @click="onLast" :disabled="currentPage === totalPages" class="admin-pagination-btn">
        末页
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

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
