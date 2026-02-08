<template>
  <div class="w-full h-full flex flex-col overflow-hidden">
    <!-- 相册列表 -->
    <div class="flex-1 min-h-0 overflow-y-auto pr-2">
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        <AlbumItem
          class="w-full max-w-md"
          v-for="album in albums"
          :key="album.id"
          :album="album"
          @view="openAlbumDetails"
          @edit="openEditDialog"
          @delete="confirmDelete = $event"
        />
      </div>
    </div>

    <!-- 分页 -->
    <div class="pt-3 flex-shrink-0">
      <Pagination
        :currentPage="currentPage"
        :totalPages="totalPages"
        @prev="handlePrevPage"
        @next="handleNextPage"
      />
    </div>

    <!-- 相册编辑对话框 -->
    <AlbumEditDialog
      v-model:open="showEditDialog"
      :album="currentAlbum"
      :loading="uiStore.loading"
      @confirm="saveAlbum"
      @cancel="closeEditDialog"
    />

    <!--// 相册详情对话框 -->
    <AlbumDetailsDialog
      v-model:open="showDetailsDialog"
      :album="currentAlbumDetails"
      :loading="uiStore.loading"
      :photos="albumPhotos"
      :has-more-photos="hasMorePhotos"
      :loading-photos="loadingPhotos"
      @confirm-cover="handleSetCover"
      @upload="handleImageUpload"
      @delete-photo="handleDeletePhoto"
      @load-more="loadMorePhotos"
      @close="closeDetailsDialog"
    />

    <!-- 删除确认对话框 -->
    <GenericDialog
      :open="!!confirmDelete"
      title="删除确认"
      :loading="uiStore.loading"
      size-class="max-w-md"
      @update:open="confirmDelete = null"
      @cancel="confirmDelete = null"
    >
      <template #content>
        <p class="text-gray-700">您确定要删除这个相册吗？此操作不可恢复。</p>
      </template>
      <template #actions>
        <div class="w-full flex">
          <div class="flex-1 text-center cursor-pointer" @click="confirmDelete = null">取消</div>
          <div
            class="w-1/2 border-l border-gray-300 text-center cursor-pointer text-red-500"
            @click="performDelete"
          >
            确定删除
          </div>
        </div>
      </template>
    </GenericDialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'

import GenericDialog from '@/components/ui/GenericDialog.vue'
import Pagination from '@/components/ui/Pagination.vue'
import { type Album, albumApi, type Photo } from '@/services/albumApi'
import { uploadApi } from '@/services/upload'
import { useUIStore } from '@/stores/ui'
import { calculateFileHash } from '@/utils/fileUtils'
import { useToast } from '@/utils/toastUtils'

import AlbumDetailsDialog from './AlbumsManagement/AlbumDetailsDialog.vue'
import AlbumEditDialog from './AlbumsManagement/AlbumEditDialog.vue'
import AlbumItem from './AlbumsManagement/AlbumItem.vue'

// 错误响应接口
interface ErrorResponse {
  response?: {
    data?: {
      code?: number
      message?: string
      data?: unknown
    }
  }
}

const props = defineProps<{ triggerAdd: boolean }>()
watch(
  () => props.triggerAdd,
  val => {
    if (val) openAddDialog()
  }
)

const uiStore = useUIStore()
const showToast = useToast()

const albums = ref<Album[]>([])
const totalAlbums = ref(0)
const currentPage = ref(1)
const pageSize = ref(5)
const totalPages = computed(() => Math.ceil(totalAlbums.value / pageSize.value) || 1)

// 编辑相关
const showEditDialog = ref(false)
const currentAlbum = ref<Album | null>(null)

// 详情相关
const showDetailsDialog = ref(false)
const currentAlbumDetails = ref<Album | null>(null)

// 相册照片相关
const albumPhotos = ref<Photo[]>([])
const currentPhotoPage = ref(1)
const photoPageSize = ref(12)
const loadingPhotos = ref(false)
const hasMorePhotos = ref(true)

// 监听详情对话框关闭，重置照片相关状态
watch(
  () => showDetailsDialog.value,
  newVal => {
    if (!newVal) {
      albumPhotos.value = []
      currentPhotoPage.value = 1
      hasMorePhotos.value = true
    }
  }
)

// 监听当前相册变化，加载照片
watch(
  () => currentAlbumDetails.value?.id,
  async albumId => {
    if (albumId) {
      // 重置照片状态
      albumPhotos.value = []
      currentPhotoPage.value = 1
      hasMorePhotos.value = true
      await loadPhotos(albumId, 1)
    }
  },
  { immediate: true }
)

// 删除确认
const confirmDelete = ref<Album | null>(null)

// 加载相册列表
const loadAlbums = async () => {
  uiStore.setLoading(true)
  try {
    const response = await albumApi.getAlbums(currentPage.value, pageSize.value)

    albums.value = response.data.albums

    totalAlbums.value =
      response.data.total || response.data.totalCount || response.data.albums.length
  } catch (error) {
    console.error('加载相册失败:', error)
    showToast('加载相册失败', 'error')
  } finally {
    uiStore.setLoading(false)
  }
}

// 删除相册
const performDelete = async () => {
  if (!confirmDelete.value) return

  const albumId = confirmDelete.value.id

  uiStore.setLoading(true)
  try {
    await albumApi.deleteAlbum(albumId)

    const index = albums.value.findIndex(a => a.id === albumId)
    if (index !== -1) {
      albums.value.splice(index, 1)
      totalAlbums.value--
    }
    showToast('相册删除成功', 'success')
    confirmDelete.value = null
  } catch (error: unknown) {
    console.error('删除相册失败:', error)
    // 显示后端返回的错误信息
    const errorMessage = (error as ErrorResponse)?.response?.data?.message || '删除相册失败'
    showToast(errorMessage, 'error')
  } finally {
    uiStore.setLoading(false)
  }
}

// 打开添加对话框
const openAddDialog = () => {
  currentAlbum.value = null
  showEditDialog.value = true
}

// 打开编辑对话框
const openEditDialog = (album: Album) => {
  currentAlbum.value = { ...album }
  showEditDialog.value = true
}

// 打开相册详情对话框
const openAlbumDetails = (album: Album) => {
  currentAlbumDetails.value = { ...album }
  showDetailsDialog.value = true
}

// 保存相册
const saveAlbum = async (albumData: Album) => {
  uiStore.setLoading(true)

  try {
    let response

    if (albumData.id && albumData.id > 0) {
      // 更新现有相册
      response = await albumApi.updateAlbum(albumData.id, {
        name: albumData.name,
        description: albumData.description,
      })
      if (response.code === 0 && response.data) {
        // 更新本地列表中的对应项
        const index = albums.value.findIndex(a => a.id === albumData.id)
        if (index !== -1) {
          albums.value[index] = response.data
        }
        showToast('相册更新成功', 'success')
      }
    } else {
      // 创建新相册
      response = await albumApi.createAlbum({
        name: albumData.name,
        description: albumData.description,
        photoCount: 0,
        createdAt: new Date().toISOString(),
      })
      if (response.code === 0 && response.data) {
        // 添加到本地列表
        albums.value.unshift(response.data)
        totalAlbums.value++
        showToast('相册添加成功', 'success')
      }
    }

    if (response.code !== 0) {
      throw new Error(response.message || '操作失败')
    }

    closeEditDialog()
  } catch (error) {
    console.error('保存相册失败:', error)
    showToast('保存相册失败', 'error')
  } finally {
    uiStore.setLoading(false)
  }
}

// 加载照片
const loadPhotos = async (albumId: number, page: number) => {
  if (loadingPhotos.value || !hasMorePhotos.value) return

  loadingPhotos.value = true
  try {
    const response = await albumApi.getPhotos(albumId, page, photoPageSize.value)

    if (response.code === 0 && response.data.photos) {
      const newPhotos = response.data.photos

      // 添加新照片
      if (page === 1) {
        albumPhotos.value = newPhotos
      } else {
        albumPhotos.value = [...albumPhotos.value, ...newPhotos]
      }

      // 检查是否还有更多照片
      hasMorePhotos.value = newPhotos.length === photoPageSize.value

      // 更新当前页码
      currentPhotoPage.value = page
    }
  } catch (error) {
    console.error('加载照片失败:', error)
    showToast('加载照片失败', 'error')
  } finally {
    loadingPhotos.value = false
  }
}

// 加载更多照片
const loadMorePhotos = () => {
  if (currentAlbumDetails.value && hasMorePhotos.value && !loadingPhotos.value) {
    loadPhotos(currentAlbumDetails.value.id, currentPhotoPage.value + 1)
  }
}

// 设置相册封面
const handleSetCover = async (albumId: number, photoId: number) => {
  uiStore.setLoading(true)
  try {
    const response = await albumApi.setCover(albumId, photoId)

    if (response.code === 0 && response.data) {
      showToast('封面更新成功', 'success')

      // 重新加载相册列表以更新封面信息
      await loadAlbums()

      // 更新当前相册详情的封面
      const index = albums.value.findIndex(a => a.id === albumId)
      if (index !== -1 && albums.value[index]) {
        currentAlbumDetails.value = albums.value[index]
      }
    } else {
      console.error('设置封面失败:', response.message || '未知错误')
      showToast('设置封面失败', 'error')
    }
  } catch (error) {
    console.error('设置封面失败:', error)
    showToast('设置封面失败', 'error')
  } finally {
    uiStore.setLoading(false)
  }
}

// 处理图片上传
const handleImageUpload = async (event: Event): Promise<void> => {
  const target = event.target as HTMLInputElement
  if (!target.files || target.files.length === 0) return

  const files = Array.from(target.files)
  const uploadedImageIds: number[] = []
  let failedCount = 0

  // 生成路径: albums/{albumName}
  const albumName = currentAlbumDetails.value?.name || 'unknown'
  const path = `albums/${albumName}`

  uiStore.setLoading(true)
  try {
    // 1. 顺序上传图片到服务器（确保按顺序完成）
    for (const file of files) {
      const formData = new FormData()
      // 计算文件哈希值
      const hash = await calculateFileHash(file)
      formData.append('file', file)
      formData.append('hash', hash)
      formData.append('path', path)

      try {
        const response = await uploadApi.uploadImage(formData)
        if (response.data.code === 0) {
          // 收集上传成功的图片ID
          uploadedImageIds.push(response.data.data.file.id)
        } else {
          showToast(`上传图片 ${file.name} 失败`, 'error')
          failedCount++
        }
      } catch {
        showToast(`上传图片 ${file.name} 失败`, 'error')
        failedCount++
      }
    }

    // 2. 将上传的图片追加到当前相册
    if (uploadedImageIds.length > 0 && currentAlbumDetails.value) {
      const addResponse = await albumApi.addPhotos(currentAlbumDetails.value.id, uploadedImageIds)

      if (addResponse.code === 0) {
        const successMsg =
          failedCount > 0
            ? `${uploadedImageIds.length} 张图片上传成功，${failedCount} 张失败`
            : `${uploadedImageIds.length} 张图片上传成功`
        showToast(successMsg, 'success')

        // 重置加载状态
        loadingPhotos.value = false
        hasMorePhotos.value = true

        // 重新加载照片列表以确保数据最新
        await loadPhotos(currentAlbumDetails.value.id, 1)

        // 重新加载相册列表以更新照片数量
        await loadAlbums()
      } else {
        showToast('添加照片到相册失败', 'error')
      }
    } else if (files.length > 0) {
      showToast('所有图片上传失败', 'error')
    }
  } catch (error) {
    console.error('图片上传流程失败:', error)
    showToast('图片上传失败', 'error')
  } finally {
    uiStore.setLoading(false)
    // 清空 input，允许重复选择相同文件
    target.value = ''
  }
}

// 删除相册中的照片
const handleDeletePhoto = async (photoId: number) => {
  if (!currentAlbumDetails.value) return

  uiStore.setLoading(true)
  try {
    const albumId = currentAlbumDetails.value.id
    const response = await albumApi.deletePhoto(albumId, photoId)

    if (response.code === 0) {
      showToast('照片删除成功', 'success')

      // 重置加载状态
      loadingPhotos.value = false
      hasMorePhotos.value = true

      // 重新加载照片列表以确保数据最新
      await loadPhotos(currentAlbumDetails.value.id, 1)

      // 重新加载相册列表以更新照片数量
      await loadAlbums()
    } else {
      showToast('删除照片失败', 'error')
    }
  } catch (error) {
    console.error('删除照片失败:', error)
    showToast('删除照片失败', 'error')
  } finally {
    uiStore.setLoading(false)
  }
}

// 关闭编辑对话框
const closeEditDialog = () => {
  showEditDialog.value = false
  currentAlbum.value = null
}

// 关闭详情对话框
const closeDetailsDialog = () => {
  showDetailsDialog.value = false
  currentAlbumDetails.value = null
}

const handlePrevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    loadAlbums()
  }
}

const handleNextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    loadAlbums()
  }
}

onMounted(() => {
  loadAlbums()
})
</script>
