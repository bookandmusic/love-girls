<template>
  <MainLayout
    title="祝福语"
    subtitle="留下对我们的美好祝愿"
    :start-date="systemInfo?.site.startDate"
    :show-empty-state="false"
  >
    <template #main-content>
      <!-- 标签页头部 -->
      <div class="border-b border-gray-200">
        <nav class="-mb-px flex" aria-label="Tabs">
          <button
            @click="activeTab = 'form'"
            :class="[
              activeTab === 'form'
                ? 'border-pink-500 text-pink-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300',
              'w-1/2 py-4 px-4 text-center border-b-2 font-medium text-sm',
            ]"
          >
            留下祝福
          </button>
          <button
            @click="activeTab = 'list'"
            :class="[
              activeTab === 'list'
                ? 'border-pink-500 text-pink-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300',
              'w-1/2 py-4 px-4 text-center border-b-2 font-medium text-sm',
            ]"
          >
            祝福列表
            <span
              v-if="totalCount"
              class="ml-2 inline-flex items-center justify-center px-2 py-1 text-xs font-bold rounded-full bg-pink-100 text-pink-800"
            >
              {{ totalCount }}
            </span>
          </button>
        </nav>
      </div>

      <!-- 标签页内容 -->
      <div class="p-4 h-full overflow-y-auto">
        <!-- 祝福表单 -->
        <div v-if="activeTab === 'form'" class="p-4 md:p-8 generic-card h-full flex flex-col">
          <h2 class="text-xl font-bold mb-4 font-[Ma_Shan_Zheng]">留下您的祝福</h2>
          <form @submit.prevent="confirmSubmit" class="flex-1 flex flex-col">
            <div class="flex-grow">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label for="authorName" class="block text-sm font-medium text-gray-700 mb-1"
                    >您的姓名</label
                  >
                  <input
                    id="authorName"
                    v-model="newWish.authorName"
                    type="text"
                    required
                    class="w-full win11-input"
                    placeholder="请输入您的姓名"
                  />
                </div>
                <div>
                  <label for="email" class="block text-sm font-medium text-gray-700 mb-1"
                    >邮箱地址（可选）</label
                  >
                  <input
                    id="email"
                    v-model="newWish.email"
                    type="email"
                    class="w-full win11-input"
                    placeholder="请输入您的邮箱"
                  />
                </div>
              </div>
              <div class="mt-4">
                <label for="content" class="block text-sm font-medium text-gray-700 mb-1"
                  >祝福内容</label
                >
                <div class="relative">
                  <textarea
                    id="content"
                    v-model="newWish.content"
                    required
                    rows="4"
                    class="w-full win11-input"
                    placeholder="请输入您的祝福语..."
                    maxlength="500"
                  ></textarea>
                  <span class="absolute bottom-2 right-2 text-xs text-gray-500">
                    {{ newWish.content.length }}/500
                  </span>
                </div>
              </div>
            </div>
            <div class="mt-4 pt-4 border-t border-gray-300">
              <div class="w-full flex">
                <button
                  type="button"
                  @click="resetForm"
                  class="flex-1 text-sm font-medium text-gray-700 focus:outline-none"
                >
                  重置
                </button>
                <button
                  type="submit"
                  class="flex-1 border-l border-gray-300 text-sm font-medium text-pink-500 focus:outline-none"
                >
                  提交祝福
                </button>
              </div>
            </div>
          </form>
        </div>

        <!-- 祝福列表 -->
        <div v-else class="flex flex-col h-full overflow-y-auto">
          <div
            v-if="wishes.length === 0"
            class="h-full w-full p-4 flex flex-col items-center justify-center py-8 text-[#FF7500]"
          >
            <BaseIcon name="wish" size="w-24" />
            <p class="text-xl font-medium mt-4">暂无祝福</p>
            <p class="text-md mt-2">还没有收到任何祝福</p>
          </div>

          <div v-else class="flex-grow overflow-y-auto space-y-4 px-3">
            <div
              v-for="wish in wishes"
              :key="wish.id"
              class="generic-card p-4 hover:bg-white/80 transition-all duration-300"
            >
              <div class="flex items-start">
                <!-- 用户头像 -->
                <div
                  class="w-12 h-12 rounded-full bg-gradient-to-br from-red-400 to-pink-500 flex items-center justify-center text-white font-bold mr-3 flex-shrink-0"
                >
                  <span>{{ wish.authorName.substring(0, 1) }}</span>
                </div>

                <!-- 祝福内容 -->
                <div class="flex-grow">
                  <div class="flex justify-between items-start">
                    <div>
                      <h3 class="font-bold">
                        {{ wish.authorName }}
                        <span v-if="wish.email" class="text-xs text-gray-500 ml-2">{{
                          wish.email
                        }}</span>
                      </h3>
                      <p class="text-xs text-gray-500">{{ wish.createdAt }}</p>
                    </div>
                  </div>

                  <p class="my-2 text-gray-700">{{ wish.content }}</p>
                </div>
              </div>
            </div>
          </div>

          <!-- 分页组件 -->
          <div v-if="wishes.length > 0" class="mt-4">
            <Pagination
              :current-page="currentPage"
              :total-pages="totalPages"
              @prev="goToPage(currentPage - 1)"
              @next="goToPage(currentPage + 1)"
            />
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
import Pagination from '@/components/ui/Pagination.vue'
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
const fetchWishes = async (page: number) => {
  uiStore.setLoading(true)
  try {
    const response = await wishApi.getWishes(page, pageSize.value, true)

    wishes.value = response.data.wishes
    totalPages.value = response.data.totalPages
    totalCount.value = response.data.total
    currentPage.value = page
  } catch {
    showToast('获取祝福列表失败', 'error')
  } finally {
    uiStore.setLoading(false)
  }
}

// 页面跳转
const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    fetchWishes(page)
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
