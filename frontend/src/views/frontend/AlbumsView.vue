<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

import BaseIcon from '@/components/ui/BaseIcon.vue'
import Pagination from '@/components/ui/Pagination.vue'
import MainLayout from '@/layouts/MainLayout.vue'
import { type Album, albumApi, type Photo } from '@/services/albumApi'
import { useSystemStore } from '@/stores/system'
import { useUIStore } from '@/stores/ui'
import { useToast } from '@/utils/toastUtils'

import AlbumList from './components/AlbumList.vue'
import PhotoList from './components/PhotoList.vue'

const uiStore = useUIStore()
const systemStore = useSystemStore()

// 获取系统信息
const systemInfo = computed(() => systemStore.getSystemInfo)

const showToast = useToast()
// 相册相关状态
const albums = ref<Album[]>([])
const currentAlbum = ref<Album | null>(null)
const currentPage = ref(1)
const totalPages = ref(0)
const pageSize = ref(6) // 每页显示6个相册

// 照片相关状态
const photos = ref<Photo[]>([])
const currentAlbumId = ref<number | null>(null)
const currentPhotoPage = ref(1)
const totalPhotoPages = ref(0)
const photoPageSize = ref(9) // 每页显示9张照片

// 计算当前显示的标题和副标题
const pageTitle = computed(() => {
  if (currentAlbum.value) {
    return currentAlbum.value.name
  }
  return '记忆相册'
})

const pageSubtitle = computed(() => {
  if (currentAlbum.value) {
    return currentAlbum.value.description
  }
  return '珍藏我们的美好瞬间'
})

// 获取相册列表
const fetchAlbums = async (page: number) => {
  uiStore.setLoading(true)
  try {
    const response = await albumApi.getAlbums(page, pageSize.value)

    albums.value = response.data.albums
    totalPages.value = response.data.totalPages
    currentPage.value = page
  } catch {
    showToast('获取相册列表失败', 'error')
  } finally {
    uiStore.setLoading(false)
  }
}

// 获取相册中的照片
const fetchPhotos = async (albumId: number, page: number) => {
  uiStore.setLoading(true)
  try {
    // 先获取相册信息
    const albumResponse = await albumApi.getAlbums(1, 100) // 获取所有相册以便找到当前相册

    const album = albumResponse.data.albums.find((a: Album) => a.id === albumId)
    if (album) {
      currentAlbum.value = album
    }

    const response = await albumApi.getPhotos(albumId, page, photoPageSize.value)

    photos.value = response.data.photos
    totalPhotoPages.value = response.data.totalPages
    currentPhotoPage.value = page
    currentAlbumId.value = albumId
  } catch {
    showToast('获取照片列表失败', 'error')
  } finally {
    uiStore.setLoading(false)
  }
}

// 返回相册列表
const backToAlbums = () => {
  currentAlbumId.value = null
  currentAlbum.value = null
  photos.value = []
  currentPhotoPage.value = 1
  totalPhotoPages.value = 0
}

// 页面跳转
const goToAlbumPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    fetchAlbums(page)
  }
}

const goToPhotoPage = (page: number) => {
  if (currentAlbumId.value && page >= 1 && page <= totalPhotoPages.value) {
    fetchPhotos(currentAlbumId.value, page)
  }
}

// 处理相册选择
const handleSelectAlbum = (album: Album) => {
  fetchPhotos(album.id, 1)
}

// 处理返回相册列表
const handleBack = () => {
  backToAlbums()
}

onMounted(async () => {
  await systemStore.fetchSystemInfo()
  fetchAlbums(1)
})
</script>

<template>
  <MainLayout
    :title="pageTitle"
    :subtitle="pageSubtitle"
    :start-date="systemInfo?.site.startDate"
    :show-empty-state="!currentAlbumId && albums.length === 0"
  >
    <template #empty-state>
      <BaseIcon name="camera" size="w-24" />
      <p class="text-xl font-medium mt-4">暂无相册</p>
      <p class="text-md mt-2">还没有创建任何相册</p>
    </template>

    <template #main-content>
      <!-- 相册列表视图 -->
      <div v-if="!currentAlbumId" class="flex flex-col h-full">
        <AlbumList :albums="albums" @select-album="handleSelectAlbum" />
        <Pagination
          :current-page="currentPage"
          :total-pages="totalPages"
          @prev="goToAlbumPage(currentPage - 1)"
          @next="goToAlbumPage(currentPage + 1)"
        />
      </div>

      <!-- 照片列表视图 -->
      <div v-else class="flex flex-col h-full">
        <PhotoList :photos="photos" @back="handleBack" />
        <Pagination
          :current-page="currentPhotoPage"
          :total-pages="totalPhotoPages"
          @prev="goToPhotoPage(currentPhotoPage - 1)"
          @next="goToPhotoPage(currentPhotoPage + 1)"
        />
      </div>
    </template>
  </MainLayout>
</template>
