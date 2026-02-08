MomentsView.vue
<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import VueEasyLightbox from 'vue-easy-lightbox'

import BaseIcon from '@/components/ui/BaseIcon.vue'
import Pagination from '@/components/ui/Pagination.vue'
import MainLayout from '@/layouts/MainLayout.vue'
import { type Moment, momentApi } from '@/services/momentApi'
import { useSystemStore } from '@/stores/system'
import { useUIStore } from '@/stores/ui'
import { useToast } from '@/utils/toastUtils'

const uiStore = useUIStore()
const systemStore = useSystemStore()

// 获取系统信息
const systemInfo = computed(() => systemStore.getSystemInfo)

const showToast = useToast()

// 动态相关状态
const moments = ref<Moment[]>([])
const currentPage = ref(1)
const totalPages = ref(0)
const pageSize = ref(5) // 每页显示5条动态

// 计算当前显示的标题和副标题
const pageTitle = computed(() => {
  return '时光动态'
})

const pageSubtitle = computed(() => {
  return '记录我们的点点滴滴'
})

// 获取动态列表
const fetchMoments = async (page: number) => {
  uiStore.setLoading(true)
  try {
    const response = await momentApi.getMoments(page, pageSize.value)

    moments.value = response.data.moments
    totalPages.value = response.data.totalPages
    currentPage.value = page
  } catch {
    showToast('获取动态列表失败', 'error')
  } finally {
    uiStore.setLoading(false)
  }
}

// 页面跳转
const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    fetchMoments(page)
  }
}

// 点赞功能
const likeMoment = async (momentId: number) => {
  try {
    const response = await momentApi.likeMoment(momentId)
    if (response.code === 0) {
      const moment = moments.value.find(m => m.id === momentId)
      if (moment) {
        moment.likes = response.data.likes
      }
      showToast('点赞成功', 'success')
    }
  } catch {
    showToast('点赞失败', 'error')
  }
}

// 查看图片
const visibleRef = ref(false)
const indexRef = ref(0) // default 0
const imgsRef = ref('')
const onShow = () => {
  visibleRef.value = true
}
const onHide = () => (visibleRef.value = false)
const viewImage = (imageUrl: string) => {
  imgsRef.value = imageUrl
  onShow()
}

onMounted(async () => {
  await systemStore.fetchSystemInfo()
  fetchMoments(1)
})
</script>

<template>
  <div class="h-full w-full">
    <MainLayout
      :title="pageTitle"
      :subtitle="pageSubtitle"
      :start-date="systemInfo?.site.startDate"
      :show-empty-state="moments.length === 0"
    >
      <template #empty-state>
        <BaseIcon name="moment" size="w-24" />
        <p class="text-xl font-medium mt-4">暂无动态</p>
        <p class="text-md mt-2">还没有发布任何动态</p>
      </template>

      <template #main-content>
        <vue-easy-lightbox
          :visible="visibleRef"
          :imgs="imgsRef"
          :index="indexRef"
          @hide="onHide"
        ></vue-easy-lightbox>
        <div class="flex flex-col h-full">
          <!-- 动态列表 -->
          <div class="flex-grow overflow-y-auto p-4 space-y-4">
            <div
              v-for="moment in moments"
              :key="moment.id"
              class="generic-card p-4 hover:bg-white/80 transition-all duration-300"
            >
              <div class="flex items-start">
                <!-- 用户头像 -->
                <div
                  class="w-12 h-12 rounded-full bg-gradient-to-br from-red-400 to-pink-500 flex items-center justify-center text-white font-bold mr-3 flex-shrink-0"
                >
                  <img
                    v-if="moment.author.avatar?.thumbnail"
                    :src="moment.author.avatar.thumbnail"
                    :alt="moment.author.name"
                    class="w-full h-full object-cover rounded-full"
                    @error="
                      $event.target && (($event.target as HTMLImageElement).style.display = 'none')
                    "
                  />
                  <img
                    v-else-if="moment.author.avatar?.url"
                    :src="moment.author.avatar.url"
                    :alt="moment.author.name"
                    class="w-full h-full object-cover rounded-full"
                    @error="
                      $event.target && (($event.target as HTMLImageElement).style.display = 'none')
                    "
                  />
                  <span
                    v-if="
                      !moment.author.avatar ||
                      (!moment.author.avatar.thumbnail && !moment.author.avatar.url)
                    "
                  >
                    {{ moment.author.name.substring(0, 1) }}
                  </span>
                </div>

                <!-- 动态内容 -->
                <div class="flex-grow">
                  <div class="flex justify-between items-start">
                    <div>
                      <h3 class="font-bold">{{ moment.author.name }}</h3>
                      <p class="text-sm text-gray-500">{{ moment.createdAt }}</p>
                    </div>
                  </div>

                  <p class="my-2 text-sm md:text-md text-gray-700">{{ moment.content }}</p>

                  <!-- 动态图片列表 -->
                  <div v-if="moment.images && moment.images.length > 0" class="mb-3">
                    <div class="grid grid-cols-3 md:grid-cols-6 gap-2">
                      <div
                        v-for="(image, index) in moment.images"
                        :key="index"
                        class="aspect-square overflow-hidden rounded-lg border border-gray-200 cursor-pointer"
                        @click="viewImage(image.file?.url || '')"
                      >
                        <img
                          :src="image.file?.thumbnail"
                          alt="动态图片"
                          class="w-full h-full object-cover"
                        />
                      </div>
                    </div>
                  </div>

                  <!-- 点赞按钮 -->
                  <div class="w-full flex justify-end space-x-2">
                    <button
                      @click="likeMoment(moment.id)"
                      class="flex items-center text-gray-500 hover:text-red-500 transition-colors"
                    >
                      <BaseIcon name="like" size="w-4" color="text-[#ff7500]" />
                      <span class="ml-1">{{ moment.likes }}</span>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 分页组件 -->
          <Pagination
            :current-page="currentPage"
            :total-pages="totalPages"
            @prev="goToPage(currentPage - 1)"
            @next="goToPage(currentPage + 1)"
          />
        </div>
      </template>
    </MainLayout>
  </div>
</template>
