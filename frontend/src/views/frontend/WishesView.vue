<template>
  <MainLayout
    title="祝福语"
    subtitle="留下对我们的美好祝愿"
    :start-date="systemInfo?.site.startDate"
    :show-empty-state="false"
  >
    <template #main-content>
      <div class="flex flex-col h-full bg-[var(--fe-bg-gray)]/30">
        <!-- 标签页头部 - iOS 分段控制器风格 -->
        <div class="p-4 flex justify-center flex-shrink-0">
          <div
            class="glass-thick p-1 rounded-xl flex w-full max-w-sm border border-white/40 shadow-sm"
          >
            <button
              @click="activeTab = 'form'"
              :class="[
                activeTab === 'form'
                  ? 'bg-white shadow-sm text-[var(--fe-text-primary)]'
                  : 'text-[var(--fe-text-secondary)]',
                'flex-1 py-1.5 px-4 text-center rounded-lg font-bold text-sm ios-transition tap-feedback',
              ]"
            >
              留下祝福
            </button>
            <button
              @click="activeTab = 'list'"
              :class="[
                activeTab === 'list'
                  ? 'bg-white shadow-sm text-[var(--fe-text-primary)]'
                  : 'text-[var(--fe-text-secondary)]',
                'flex-1 py-1.5 px-4 text-center rounded-lg font-bold text-sm ios-transition tap-feedback relative',
              ]"
            >
              祝福列表
              <span
                v-if="totalCount"
                class="ml-1.5 px-1.5 py-0.5 text-[10px] font-bold rounded-full bg-[var(--fe-primary)] text-white"
              >
                {{ totalCount }}
              </span>
            </button>
          </div>
        </div>

        <!-- 标签页内容 -->
        <div class="flex-grow overflow-hidden relative">
          <!-- 祝福表单 -->
          <div
            v-if="activeTab === 'form'"
            class="h-full overflow-y-auto p-4 md:p-8 custom-scrollbar"
          >
            <div
              class="max-w-2xl mx-auto glass-thick rounded-[var(--fe-radius-card)] p-6 md:p-10 border border-white/40 shadow-lg"
            >
              <h2 class="text-2xl font-bold mb-6 text-[var(--fe-text-primary)] text-center">
                留下您的祝福
              </h2>
              <form @submit.prevent="confirmSubmit" class="space-y-6">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <label
                      for="authorName"
                      class="block text-xs font-bold text-[var(--fe-text-secondary)] uppercase tracking-widest mb-2 ml-1"
                      >您的姓名</label
                    >
                    <input
                      id="authorName"
                      v-model="newWish.authorName"
                      type="text"
                      required
                      class="w-full glass-ultra-thin border border-white/60 rounded-xl px-4 py-3 text-sm focus:border-[var(--fe-primary)] focus:ring-2 focus:ring-[var(--fe-primary)]/20 outline-none ios-transition"
                      placeholder="如何称呼您"
                    />
                  </div>
                  <div>
                    <label
                      for="email"
                      class="block text-xs font-bold text-[var(--fe-text-secondary)] uppercase tracking-widest mb-2 ml-1"
                      >邮箱地址</label
                    >
                    <input
                      id="email"
                      v-model="newWish.email"
                      type="email"
                      class="w-full glass-ultra-thin border border-white/60 rounded-xl px-4 py-3 text-sm focus:border-[var(--fe-primary)] focus:ring-2 focus:ring-[var(--fe-primary)]/20 outline-none ios-transition"
                      placeholder="仅用于接收回复"
                    />
                  </div>
                </div>
                <div>
                  <label
                    for="content"
                    class="block text-xs font-bold text-[var(--fe-text-secondary)] uppercase tracking-widest mb-2 ml-1"
                    >祝福内容</label
                  >
                  <div class="relative">
                    <textarea
                      id="content"
                      v-model="newWish.content"
                      required
                      rows="5"
                      class="w-full glass-ultra-thin border border-white/60 rounded-xl px-4 py-3 text-sm focus:border-[var(--fe-primary)] focus:ring-2 focus:ring-[var(--fe-primary)]/20 outline-none ios-transition resize-none"
                      placeholder="写下对我们的祝愿..."
                      maxlength="500"
                    ></textarea>
                    <span
                      class="absolute bottom-3 right-3 text-[10px] font-bold text-[var(--fe-text-secondary)] opacity-50"
                    >
                      {{ newWish.content.length }}/500
                    </span>
                  </div>
                </div>

                <div class="flex space-x-4 pt-4">
                  <button
                    type="button"
                    @click="resetForm"
                    class="flex-1 py-3 text-sm font-bold text-[var(--fe-text-secondary)] ios-transition tap-feedback"
                  >
                    重置
                  </button>
                  <button
                    type="submit"
                    class="flex-[2] bg-[var(--fe-primary)] text-white py-3 rounded-xl text-sm font-bold shadow-md shadow-[var(--fe-primary)]/20 ios-transition tap-feedback"
                  >
                    提交祝福
                  </button>
                </div>
              </form>
            </div>
          </div>

          <!-- 祝福列表 -->
          <div v-else class="h-full flex flex-col overflow-hidden">
            <div
              v-if="wishes.length === 0 && !loadingMore"
              class="flex-1 flex flex-col items-center justify-center py-20 text-[var(--fe-text-secondary)]"
            >
              <BaseIcon name="wish" size="w-24" />
              <p class="text-xl font-bold mt-4 text-[var(--fe-text-primary)]">暂无祝福</p>
              <p class="text-md mt-2">成为第一个留下祝福的人吧</p>
            </div>

            <div
              v-else
              class="flex-grow overflow-y-auto px-4 md:px-8 space-y-4 py-4 custom-scrollbar"
              @scroll="handleScroll"
            >
              <div
                v-for="wish in wishes"
                :key="wish.id"
                class="glass-thick p-5 rounded-2xl border border-white/40 shadow-sm ios-transition"
              >
                <div class="flex items-start">
                  <!-- 用户头像 -->
                  <div
                    class="w-12 h-12 rounded-xl bg-gradient-to-br from-[#f8c9c0] to-[var(--fe-primary)] flex items-center justify-center text-white font-bold mr-4 flex-shrink-0 shadow-sm"
                  >
                    <span>{{ wish.authorName.substring(0, 1) }}</span>
                  </div>

                  <!-- 祝福内容 -->
                  <div class="flex-grow min-w-0">
                    <div class="flex justify-between items-center mb-1">
                      <h3 class="font-bold text-[var(--fe-text-primary)] truncate">
                        {{ wish.authorName }}
                      </h3>
                      <span class="text-[10px] font-bold text-[var(--fe-text-secondary)] uppercase">
                        {{ wish.createdAt }}
                      </span>
                    </div>
                    <p class="text-sm text-[var(--fe-text-primary)] leading-relaxed mt-2">
                      {{ wish.content }}
                    </p>
                  </div>
                </div>
              </div>

              <!-- 加载状态指示器 -->
              <div v-if="loadingMore || hasMore" class="py-10 flex justify-center">
                <div
                  v-if="loadingMore"
                  class="flex items-center space-x-2 text-[var(--fe-text-secondary)]"
                >
                  <div
                    class="w-5 h-5 border-2 border-[var(--fe-primary)] border-t-transparent rounded-full animate-spin"
                  ></div>
                  <span class="text-xs font-bold uppercase tracking-widest">正在加载更多...</span>
                </div>
                <div
                  v-else-if="!hasMore && wishes.length > 0"
                  class="text-xs font-bold text-[var(--fe-text-secondary)] uppercase tracking-widest opacity-30"
                >
                  没有更多祝福了
                </div>
              </div>

              <!-- 占位 -->
              <div class="h-20 md:hidden"></div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </MainLayout>

  <!-- 提交确认对话框 -->
  <GenericDialog
    :open="showConfirmDialog"
    title="提交确认"
    :loading="uiStore.loading"
    size-class="max-w-md"
    @cancel="cancelSubmit"
  >
    <template #content>
      <p class="text-gray-700">您确定要提交这条祝福吗？</p>
    </template>
    <template #actions>
      <div class="w-full flex">
        <div class="flex-1 text-center cursor-pointer" @click="cancelSubmit">取消</div>
        <div
          class="w-1/2 border-l border-gray-300 text-center cursor-pointer text-pink-500"
          @click="submitWishConfirmed"
        >
          确定提交
        </div>
      </div>
    </template>
  </GenericDialog>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

import BaseIcon from '@/components/ui/BaseIcon.vue'
import GenericDialog from '@/components/ui/GenericDialog.vue'
import MainLayout from '@/layouts/MainLayout.vue'
import { type Wish, wishApi } from '@/services/wishApi'
import { useSystemStore } from '@/stores/system'
import { useUIStore } from '@/stores/ui'
import { useToast } from '@/utils/toastUtils'

const uiStore = useUIStore()
const systemStore = useSystemStore()

// 获取系统信息
const systemInfo = computed(() => systemStore.getSystemInfo)

// Tab标签页状态
const activeTab = ref<'form' | 'list'>('form') // 默认显示表单

// Toast 通知状态
const showToast = useToast()

// 祝福相关状态
const wishes = ref<Wish[]>([])
const currentPage = ref(1)
const totalPages = ref(0)
const totalCount = ref(0)
const pageSize = ref(5) // 每页显示5条祝福

// 新祝福状态
const newWish = ref({
  authorName: '',
  email: '',
  content: '',
})

// 提交确认对话框
const showConfirmDialog = ref(false)

// 获取祝福列表
const fetchWishes = async (page: number, append = false) => {
  if (loadingMore.value) return
  loadingMore.value = true

  try {
    const response = await wishApi.getWishes(page, pageSize.value, true)

    if (append) {
      wishes.value = [...wishes.value, ...response.data.wishes]
    } else {
      wishes.value = response.data.wishes
    }

    totalPages.value = response.data.totalPages
    totalCount.value = response.data.total
    currentPage.value = page
  } catch {
    showToast('获取祝福列表失败', 'error')
  } finally {
    loadingMore.value = false
    uiStore.setLoading(false)
  }
}

const loadingMore = ref(false)
const hasMore = computed(() => currentPage.value < totalPages.value)

const handleScroll = (e: Event) => {
  if (activeTab.value !== 'list') return
  const target = e.target as HTMLElement
  if (loadingMore.value || !hasMore.value) return

  const bottomDistance = target.scrollHeight - target.scrollTop - target.clientHeight
  if (bottomDistance < 100) {
    fetchWishes(currentPage.value + 1, true)
  }
}

// 重置表单
const resetForm = () => {
  newWish.value = {
    authorName: '',
    email: '',
    content: '',
  }
}

// 显示提交确认对话框
const confirmSubmit = (e: Event) => {
  e.preventDefault()

  // 表单验证
  if (!newWish.value.authorName.trim()) {
    showToast('请输入您的姓名', 'error')
    return
  }

  if (newWish.value.email && !validateEmail(newWish.value.email)) {
    showToast('请输入有效的邮箱地址', 'error')
    return
  }

  if (!newWish.value.content.trim()) {
    showToast('请输入祝福内容', 'error')
    return
  }

  showConfirmDialog.value = true
}

// 邮箱验证函数
const validateEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

// 取消提交
const cancelSubmit = () => {
  showConfirmDialog.value = false
}

// 确认提交新祝福
const submitWishConfirmed = async () => {
  // 先关闭确认对话框
  showConfirmDialog.value = false

  uiStore.setLoading(true)
  try {
    await wishApi.createWish({
      authorName: newWish.value.authorName,
      email: newWish.value.email,
      content: newWish.value.content,
    })

    // 重置表单
    resetForm()

    // 重新获取列表，显示最新数据
    fetchWishes(1)
    showToast('祝福提交成功！', 'success')

    // 切换到列表标签页显示最新祝福
    activeTab.value = 'list'
  } catch {
    showToast('祝福提交失败', 'error')
  } finally {
    uiStore.setLoading(false)
  }
}

onMounted(async () => {
  await systemStore.fetchSystemInfo()
  fetchWishes(1)
})
</script>
