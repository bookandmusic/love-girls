<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import VueEasyLightbox from 'vue-easy-lightbox'

import BaseIcon from '@/components/ui/BaseIcon.vue'
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
const pageSize = ref(8)
const loadingMore = ref(false)
const hasMore = computed(() => currentPage.value < totalPages.value)

// 获取动态列表
const fetchMoments = async (page: number, append = false) => {
  if (loadingMore.value) return
  loadingMore.value = true

  try {
    const response = await momentApi.getMoments(page, pageSize.value)
    if (append) {
      moments.value = [...moments.value, ...response.data.moments]
    } else {
      moments.value = response.data.moments
    }
    totalPages.value = response.data.totalPages
    currentPage.value = page
  } catch {
    showToast('获取动态列表失败', 'error')
  } finally {
    loadingMore.value = false
    uiStore.setLoading(false)
  }
}

// 滚动加载
const handleScroll = (e: Event) => {
  const target = e.target as HTMLElement
  if (!target || loadingMore.value || !hasMore.value) return

  const bottomDistance = target.scrollHeight - target.scrollTop - target.clientHeight
  if (bottomDistance < 100) {
    fetchMoments(currentPage.value + 1, true)
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
const indexRef = ref(0)
const imgsRef = ref('')
const onShow = () => (visibleRef.value = true)
const onHide = () => (visibleRef.value = false)
const viewImage = (imageUrl: string) => {
  imgsRef.value = imageUrl
  onShow()
}

onMounted(async () => {
  uiStore.setLoading(true)
  await systemStore.fetchSystemInfo()
  await fetchMoments(1)
})
</script>

<template>
  <div class="h-full w-full">
    <MainLayout
      title="时光动态"
      subtitle="记录我们的点点滴滴"
      :start-date="systemInfo?.site.startDate"
      :show-empty-state="moments.length === 0 && !loadingMore"
    >
      <template #empty-state>
        <div
          class="flex flex-col items-center justify-center py-20 text-[var(--fe-text-secondary)]"
        >
          <BaseIcon name="moment" size="w-24" />
          <p class="text-xl font-bold mt-4 text-[var(--fe-text-primary)]">暂无动态</p>
        </div>
      </template>

      <template #main-content>
        <vue-easy-lightbox
          :visible="visibleRef"
          :imgs="imgsRef"
          :index="indexRef"
          @hide="onHide"
        ></vue-easy-lightbox>

        <div class="flex flex-col h-full glass-regular">
          <!-- 动态列表 - 监听滚动事件 -->
          <div
            class="flex-grow overflow-y-auto p-4 md:p-8 space-y-0 custom-scrollbar"
            @scroll="handleScroll"
          >
            <div
              v-for="moment in moments"
              :key="moment.id"
              class="py-6 border-b border-black/5 last:border-0 ios-transition"
            >
              <div class="flex items-start">
                <!-- 用户头像 -->
                <div
                  class="w-12 h-12 rounded-lg overflow-hidden bg-white/50 border border-white/60 flex items-center justify-center text-[var(--fe-primary)] font-bold mr-4 flex-shrink-0"
                >
                  <img
                    v-if="moment.author.avatar?.thumbnail || moment.author.avatar?.url"
                    :src="moment.author.avatar?.thumbnail || moment.author.avatar?.url"
                    :alt="moment.author.name"
                    class="w-full h-full object-cover"
                  />
                  <span v-else>{{ moment.author.name.substring(0, 1) }}</span>
                </div>

                <!-- 动态内容 -->
                <div class="flex-grow min-w-0">
                  <div class="mb-1">
                    <h3 class="font-bold text-[#576b95] text-base truncate">
                      {{ moment.author.name }}
                    </h3>
                  </div>

                  <p
                    class="text-[var(--fe-text-primary)] leading-relaxed mb-3 text-sm md:text-base"
                  >
                    {{ moment.content }}
                  </p>

                  <!-- 动态图片列表 -->
                  <div v-if="moment.images && moment.images.length > 0" class="mb-3">
                    <div
                      class="grid gap-1.5"
                      :class="{
                        'grid-cols-1 w-max': moment.images.length === 1,
                        'grid-cols-2 w-full md:max-w-[280px]':
                          moment.images.length === 2 || moment.images.length === 4,
                        'grid-cols-3 w-full max-w-[320px] md:max-w-[420px]':
                          moment.images.length === 3 || moment.images.length >= 5,
                      }"
                    >
                      <div
                        v-for="(image, index) in moment.images"
                        :key="index"
                        class="overflow-hidden cursor-pointer tap-feedback ios-transition"
                        :class="[
                          moment.images.length === 1
                            ? 'rounded-lg max-w-[240px] max-h-[320px]'
                            : 'w-full aspect-square rounded-md',
                        ]"
                        @click="viewImage(image.file?.url || '')"
                      >
                        <img
                          :src="image.file?.thumbnail || image.file?.url"
                          alt="动态图片"
                          class="w-full h-full object-cover"
                          loading="lazy"
                        />
                      </div>
                    </div>
                  </div>

                  <!-- 底部信息与交互栏 -->
                  <div class="flex justify-between items-center mt-3">
                    <span class="text-xs font-medium text-[var(--fe-text-secondary)] opacity-60">
                      {{ moment.createdAt }}
                    </span>
                    <button
                      @click="likeMoment(moment.id)"
                      class="flex items-center space-x-1.5 px-2.5 py-1 rounded-md bg-black/5 tap-feedback ios-transition"
                    >
                      <BaseIcon name="like" size="w-3.5 h-3.5" color="var(--fe-primary)" />
                      <span class="text-xs font-bold text-[var(--fe-text-primary)]">
                        {{ moment.likes }}
                      </span>
                    </button>
                  </div>
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
                v-else-if="!hasMore && moments.length > 0"
                class="text-xs font-bold text-[var(--fe-text-secondary)] uppercase tracking-widest opacity-30"
              >
                没有更多动态了
              </div>
            </div>

            <!-- 占位防止底部 TabBar 遮挡 -->
            <div class="h-20 md:hidden"></div>
          </div>
        </div>
      </template>
    </MainLayout>
  </div>
</template>
