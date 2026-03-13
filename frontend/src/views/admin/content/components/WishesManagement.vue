<template>
  <div class="w-full h-full flex flex-col overflow-hidden">
    <!-- 祝福列表 -->
    <ul class="flex-1 min-h-0 overflow-y-auto pr-2">
      <WishItem
        v-for="wish in wishes"
        :key="wish.id"
        :wish="wish"
        @approve="showApproveConfirmation"
        @delete="showDeleteConfirmation"
      />

      <!-- 观测点 -->
      <div
        ref="loadMoreTrigger"
        class="py-4 flex flex-col items-center justify-center text-gray-500 text-sm"
      >
        <template v-if="uiStore.loading && currentPage > 1">
          <div
            class="animate-spin rounded-full h-6 w-6 border-b-2 border-primary-500 mb-2"
            style="border-color: var(--primary-color) transparent transparent transparent"
          ></div>
          <span>正在加载更多...</span>
        </template>
        <template v-else-if="!hasMore && wishes.length > 0">
          <span>已经到底啦</span>
        </template>
      </div>
    </ul>

    <!-- 确认删除对话框 -->
    <GenericDialog
      variant="admin"
      :open="showConfirmDialog"
      title="确认删除"
      size-class="max-w-md"
      @update:open="showConfirmDialog = $event"
      @cancel="showConfirmDialog = false"
    >
      <template #content>
        <p class="text-gray-700">确定要删除这条祝福吗？此操作不可撤销。</p>
      </template>
      <template #actions>
        <div class="w-full flex">
          <div class="flex-1 text-center cursor-pointer" @click="showConfirmDialog = false">
            取消
          </div>
          <div
            class="w-1/2 border-l border-gray-300 text-center cursor-pointer text-red-500"
            @click="confirmDeleteWish"
          >
            删除
          </div>
        </div>
      </template>
    </GenericDialog>

    <!-- 确认批准对话框 -->
    <GenericDialog
      variant="admin"
      :open="showApproveDialog"
      title="确认批准"
      size-class="max-w-md"
      @update:open="showApproveDialog = $event"
      @cancel="showApproveDialog = false"
    >
      <template #content>
        <p class="text-gray-700">确定要批准这条祝福吗？</p>
      </template>
      <template #actions>
        <div class="w-full flex">
          <div class="flex-1 text-center cursor-pointer" @click="showApproveDialog = false">
            取消
          </div>
          <div
            class="w-1/2 border-l border-gray-300 text-center cursor-pointer text-blue-500"
            @click="confirmApproveWish"
          >
            批准
          </div>
        </div>
      </template>
    </GenericDialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'

import GenericDialog from '@/components/ui/GenericDialog.vue'
import { type Wish, wishApi } from '@/services/wishApi'
import { useUIStore } from '@/stores/ui'
import { useToast } from '@/utils/toastUtils'

import WishItem from './WishesManagement/WishItem.vue'

const uiStore = useUIStore()
const showToast = useToast()

const wishes = ref<Wish[]>([])
const totalWishes = ref(0)
const currentPage = ref(1)
const pageSize = ref(5)
const totalPages = computed(() => Math.ceil(totalWishes.value / pageSize.value))
const hasMore = computed(() => currentPage.value < totalPages.value)

// 观测点相关
const loadMoreTrigger = ref<HTMLElement | null>(null)
let observer: IntersectionObserver | null = null

const showConfirmDialog = ref(false)
const showApproveDialog = ref(false)
let wishToDelete: Wish | null = null
let wishToApprove: Wish | null = null

// 加载祝福列表
const loadWishes = async (append = false) => {
  uiStore.setLoading(true)
  try {
    const response = await wishApi.getWishes(currentPage.value, pageSize.value)

    if (append) {
      wishes.value = [...wishes.value, ...response.data.wishes]
    } else {
      wishes.value = response.data.wishes
    }

    totalWishes.value =
      response.data.total || response.data.totalCount || response.data.wishes.length
  } catch (error) {
    console.error('加载祝福失败:', error)
    showToast('加载祝福失败', 'error')
  } finally {
    uiStore.setLoading(false)
  }
}

// 处理交叉观测
const handleIntersect = (entries: IntersectionObserverEntry[]) => {
  const entry = entries[0]
  if (entry && entry.isIntersecting && hasMore.value && !uiStore.loading) {
    currentPage.value++
    loadWishes(true)
  }
}

// 显示删除确认对话框
const showDeleteConfirmation = (wish: Wish) => {
  wishToDelete = wish
  showConfirmDialog.value = true
}

// 显示批准确认对话框
const showApproveConfirmation = (wish: Wish) => {
  wishToApprove = wish
  showApproveDialog.value = true
}

// 确认删除祝福
const confirmDeleteWish = async () => {
  if (!wishToDelete) return

  uiStore.setLoading(true)
  try {
    // 发送API请求删除祝福
    await wishApi.deleteWish(wishToDelete.id)

    // 删除后重置到第一页
    currentPage.value = 1
    await loadWishes(false)
    showToast('删除成功', 'success')
  } catch (error) {
    console.error('删除祝福失败:', error)
    showToast('删除祝福失败', 'error')
  } finally {
    uiStore.setLoading(false)
    showConfirmDialog.value = false
    wishToDelete = null
  }
}

// 确认批准祝福
const confirmApproveWish = async () => {
  if (!wishToApprove) return

  uiStore.setLoading(true)
  try {
    // 发送API请求批准祝福
    await wishApi.approveWish(wishToApprove.id)

    // 更新本地状态
    const index = wishes.value.findIndex(w => w.id === wishToApprove!.id)
    if (index !== -1 && wishes.value[index]) {
      wishes.value[index].approved = true
    }

    showToast('批准成功', 'success')
  } catch (error) {
    console.error('批准祝福失败:', error)
    showToast('批准祝福失败', 'error')
  } finally {
    uiStore.setLoading(false)
    wishToApprove = null
    showApproveDialog.value = false
  }
}

onMounted(() => {
  loadWishes()
  observer = new IntersectionObserver(handleIntersect, { threshold: 0.1 })
  if (loadMoreTrigger.value) {
    observer.observe(loadMoreTrigger.value)
  }
})

onUnmounted(() => {
  if (observer) {
    observer.disconnect()
  }
})
</script>
