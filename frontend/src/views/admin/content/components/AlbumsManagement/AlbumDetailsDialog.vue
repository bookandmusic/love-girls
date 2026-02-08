<template>
  <GenericDialog :open="open" :title="album?.name" @cancel="closeDialog" :loading="loading">
    <template #content>
      <div class="space-y-4 h-full">
        <div class="flex flex-col h-full">
          <p class="text-sm text-gray-500 mb-4">照片列表</p>
          <!-- 照片列表 -->
          <div class="mt-1 mb-4 flex-1 overflow-y-auto" @scroll="handleScroll">
            <div class="px-2 grid grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-3 py-1">
              <div
                v-for="photo in photos"
                :key="photo.id"
                class="relative aspect-square group"
                @click="selectPhoto(photo.id)"
              >
                <img
                  :src="photo.file?.thumbnail || photo.file?.url || ''"
                  :alt="photo.alt || '照片'"
                  class="w-full h-full object-cover rounded border border-gray-300"
                  :class="{ 'ring-2 ring-indigo-500': selectedCoverId === photo.id }"
                />
                <!-- 删除按钮 -->
                <button
                  @click.stop="confirmDeletePhoto(photo)"
                  class="absolute top-1 right-1 p-1 bg-white rounded-full opacity-0 group-hover:opacity-100 transition-opacity shadow-sm hover:bg-red-50"
                  title="删除照片"
                >
                  <BaseIcon name="delete" size="w-4" color="text-red-500" />
                </button>
              </div>
              <!-- 添加图片按钮 -->
              <div
                @click="triggerImageUpload"
                class="relative aspect-square border-2 border-dashed border-gray-300 rounded flex items-center justify-center cursor-pointer hover:border-[var(--primary-color)] hover:bg-gray-50"
                :disabled="loading"
              >
                <span class="text-2xl text-gray-500">+</span>
              </div>
            </div>

            <!-- 加载更多指示器 -->
            <div v-if="loadingPhotos" class="mt-4 text-center py-2">
              <span
                class="inline-block animate-spin rounded-full h-4 w-4 border-b-2 border-primary"
              ></span>
            </div>

            <!-- 隐藏的文件输入框 -->
            <input
              ref="imageInputRef"
              type="file"
              accept="image/*"
              @change="handleSelectedImageUpload"
              class="hidden"
              multiple
            />
          </div>
        </div>
      </div>
    </template>
    <template #actions>
      <div class="w-full flex">
        <div class="flex-1 text-center cursor-pointer" @click="closeDialog">取消</div>
        <div
          class="w-1/2 border-l border-gray-300 text-center cursor-pointer text-blue-500"
          @click="handleSetAsCover"
          :disabled="!selectedCoverId || loading"
        >
          设置封面
        </div>
      </div>
    </template>
  </GenericDialog>

  <!-- 确认对话框 - 设置封面 -->
  <GenericDialog
    :open="showConfirmDialog"
    title="确认设置封面"
    :loading="loading"
    size-class="max-w-md"
  >
    <template #content>
      <p class="text-gray-700">您确定要将选中的照片设置为相册封面吗？</p>
    </template>
    <template #actions>
      <div class="w-full flex">
        <div class="flex-1 text-center cursor-pointer" @click="cancelConfirm">取消</div>
        <div
          class="w-1/2 border-l border-gray-300 text-center cursor-pointer text-blue-500"
          @click="confirmSetCover"
        >
          确定
        </div>
      </div>
    </template>
  </GenericDialog>

  <!-- 确认对话框 - 删除照片 -->
  <GenericDialog
    :open="!!photoToDelete"
    title="确认删除照片"
    :loading="loading"
    size-class="max-w-md"
    @cancel="photoToDelete = null"
  >
    <template #content>
      <p class="text-gray-700">您确定要删除这张照片吗？此操作不可恢复。</p>
    </template>
    <template #actions>
      <div class="w-full flex">
        <div class="flex-1 text-center cursor-pointer" @click="photoToDelete = null">取消</div>
        <div
          class="w-1/2 border-l border-gray-300 text-center cursor-pointer text-red-500"
          @click="handleDeletePhoto"
        >
          确定删除
        </div>
      </div>
    </template>
  </GenericDialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

import BaseIcon from '@/components/ui/BaseIcon.vue'
import GenericDialog from '@/components/ui/GenericDialog.vue'
import { type Album, type Photo } from '@/services/albumApi'

interface Props {
  open: boolean
  loading: boolean
  album: Album | null
  photos: Photo[]
  hasMorePhotos: boolean
  loadingPhotos: boolean
}

interface Emits {
  (e: 'update:open', open: boolean): void
  (e: 'confirm-cover', albumId: number, photoId: number): void
  (e: 'upload', event: Event): void
  (e: 'delete-photo', photoId: number): void
  (e: 'load-more'): void
  (e: 'close'): void
}

const props = withDefaults(defineProps<Props>(), {
  album: null,
  loading: false,
  photos: () => [],
  hasMorePhotos: false,
  loadingPhotos: false,
})

const emit = defineEmits<Emits>()

const closeDialog = () => {
  emit('update:open', false)
  emit('close')
}

const imageInputRef = ref<HTMLInputElement>()
const showConfirmDialog = ref(false)
const photoToDelete = ref<Photo | null>(null)

const selectedCoverId = ref<number | null>(null)

// 监听相册变化，设置封面图片为选中状态
watch(
  () => props.album,
  newAlbum => {
    if (newAlbum && newAlbum.coverImage) {
      const coverPhoto = props.photos.find(p => p.id === newAlbum.coverImage?.id)
      if (coverPhoto) {
        selectedCoverId.value = coverPhoto.id
      } else {
        selectedCoverId.value = null
      }
    } else {
      selectedCoverId.value = null
    }
  },
  { immediate: true }
)

// 监听照片变化，确保封面图片选中状态正确
watch(
  () => props.photos,
  () => {
    if (props.album && props.album.coverImage) {
      const coverPhoto = props.photos.find(p => p.id === props.album?.coverImage?.id)
      if (coverPhoto) {
        selectedCoverId.value = coverPhoto.id
      } else {
        selectedCoverId.value = null
      }
    }
  },
  { immediate: true }
)

// 处理滚动事件，实现无限加载
const handleScroll = (event: Event) => {
  const target = event.target as HTMLElement
  const { scrollTop, clientHeight, scrollHeight } = target

  // 当滚动到距离底部20px时，触发加载更多
  if (scrollHeight - scrollTop - clientHeight < 20 && props.hasMorePhotos && !props.loadingPhotos) {
    emit('load-more')
  }
}

const selectPhoto = (photoId: number) => {
  selectedCoverId.value = photoId
}

const handleSetAsCover = () => {
  if (props.album && selectedCoverId.value) {
    showConfirmDialog.value = true
  }
}

const cancelConfirm = () => {
  showConfirmDialog.value = false
}

const confirmSetCover = () => {
  showConfirmDialog.value = false
  if (props.album && selectedCoverId.value) {
    emit('confirm-cover', props.album.id, selectedCoverId.value)
  }
}

const confirmDeletePhoto = (photo: Photo) => {
  photoToDelete.value = photo
}

const handleDeletePhoto = () => {
  if (photoToDelete.value) {
    emit('delete-photo', photoToDelete.value.id)
    photoToDelete.value = null
  }
}

const triggerImageUpload = () => {
  if (imageInputRef.value) {
    imageInputRef.value.click()
  }
}

const handleSelectedImageUpload = (event: Event) => {
  emit('upload', event)
}
</script>
