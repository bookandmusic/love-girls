<template>
  <div class="w-full h-full flex flex-col overflow-hidden">
    <!-- 地点列表 -->
    <div ref="scrollContainer" class="flex-1 min-h-0 overflow-y-auto pr-2">
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        <PlaceItem
          v-for="place in places"
          :key="place.id"
          :place="place"
          @edit="openEditDialog"
          @delete="confirmDelete = $event"
        />
      </div>

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
        <template v-else-if="!hasMore && places.length > 0">
          <span>已经到底啦</span>
        </template>
      </div>
    </div>

    <!-- 地点编辑对话框 -->
    <PlaceEditDialog
      v-model:open="showDialog"
      :place="currentPlace"
      :loading="uiStore.loading"
      @confirm="savePlace"
      @cancel="closeDialog"
      @upload="handleImageUpload"
    />

    <!-- 删除确认对话框 -->
    <GenericDialog
      variant="admin"
      :open="!!confirmDelete"
      title="删除确认"
      :loading="uiStore.loading"
      size-class="max-w-md"
      @cancel="cancelDelete"
    >
      <template #content>
        <p class="text-gray-700">您确定要删除这个地点吗？此操作不可恢复。</p>
      </template>
      <template #actions>
        <div class="w-full flex">
          <div class="flex-1 text-center cursor-pointer" @click="cancelDelete">取消</div>
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
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'

import GenericDialog from '@/components/ui/GenericDialog.vue'
import { type Photo, type Place, placeApi } from '@/services/placeApi'
import { uploadApi } from '@/services/upload'
import { useUIStore } from '@/stores/ui'
import { calculateFileHash } from '@/utils/fileUtils'
import { useToast } from '@/utils/toastUtils'
import { useAutoFillPage } from '@/utils/useAutoFillPage'

import PlaceEditDialog from './PlacesManagement/PlaceEditDialog.vue'
import PlaceItem from './PlacesManagement/PlaceItem.vue'

// 定义默认地点对象
const DEFAULT_PLACE: Place = {
  id: 0,
  name: '',
  latitude: 0,
  longitude: 0,
  date: new Date().toISOString().substring(0, 10),
  image: undefined,
  description: '',
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

const places = ref<Place[]>([])
const totalPlaces = ref(0)
const currentPage = ref(1)
const pageSize = ref(5)
const totalPages = computed(() => Math.ceil(totalPlaces.value / pageSize.value) || 1)
const hasMore = computed(() => currentPage.value < totalPages.value)

// 滚动容器引用
const scrollContainer = ref<HTMLElement | null>(null)

// 观测点相关
const loadMoreTrigger = ref<HTMLElement | null>(null)
let observer: IntersectionObserver | null = null

// 编辑相关
const showDialog = ref(false)
const currentPlace = ref<Place | null>(null)

// 删除确认
const confirmDelete = ref<Place | null>(null)

// 取消删除
const cancelDelete = () => {
  confirmDelete.value = null
}

// 处理图片上传
const handleImageUpload = async (event: Event): Promise<void> => {
  const target = event.target as HTMLInputElement
  if (!target.files || target.files.length === 0) return

  const file = target.files[0]
  if (!file) return

  uiStore.setLoading(true)
  try {
    // 计算文件哈希值
    const hash = await calculateFileHash(file)

    // 生成路径: places/%Y/%m
    const now = new Date()
    const year = now.getFullYear()
    const month = String(now.getMonth() + 1).padStart(2, '0')
    const path = `places/${year}/${month}`

    const formData = new FormData()
    formData.append('file', file)
    formData.append('hash', hash)
    formData.append('path', path)

    // 上传文件
    const response = await uploadApi.uploadImage(formData)

    if (response.data.code !== 0) {
      throw new Error(response.data.message || '图片上传失败')
    }

    // 返回上传后的图片信息
    const photo: Photo = {
      id: response.data.data.file.id,
      placeId: 0,
      file: response.data.data.file,
    }

    // 更新当前地点的图片
    if (currentPlace.value) {
      currentPlace.value.image = photo
    }

    showToast('图片上传成功', 'success')
  } catch {
    showToast('图片上传失败', 'error')
  } finally {
    uiStore.setLoading(false)
    // 清空 input，允许重复选择相同文件
    target.value = ''
  }
}

// 加载地点列表
const loadPlaces = async (append = false) => {
  uiStore.setLoading(true)
  try {
    const response = await placeApi.getPlaces(currentPage.value, pageSize.value)

    if (append) {
      places.value = [...places.value, ...response.data.places]
    } else {
      places.value = response.data.places
    }

    totalPlaces.value =
      response.data.total || response.data.totalCount || response.data.places.length
  } catch {
    showToast('加载地点失败', 'error')
  } finally {
    uiStore.setLoading(false)
    // 检查是否需要自动加载更多以填充页面
    checkAndAutoLoadMore()
  }
}

// 加载下一页
const handleNextPage = () => {
  if (hasMore.value) {
    currentPage.value++
    loadPlaces(true)
  }
}

// 自动填充页面逻辑
const { checkAndAutoLoadMore } = useAutoFillPage(
  scrollContainer,
  hasMore,
  computed(() => uiStore.loading),
  handleNextPage
)

// 处理交叉观测
const handleIntersect = (entries: IntersectionObserverEntry[]) => {
  const entry = entries[0]
  if (entry && entry.isIntersecting && hasMore.value && !uiStore.loading) {
    handleNextPage()
  }
}

// 删除地点
const performDelete = async () => {
  if (!confirmDelete.value) return

  const placeId = confirmDelete.value.id

  uiStore.setLoading(true)
  try {
    await placeApi.deletePlace(placeId)

    // 删除后重置到第一页
    currentPage.value = 1
    await loadPlaces(false)
    showToast('地点删除成功', 'success')
    confirmDelete.value = null
  } catch {
    showToast('删除地点失败', 'error')
  } finally {
    uiStore.setLoading(false)
  }
}

// 打开添加对话框
const openAddDialog = () => {
  currentPlace.value = { ...DEFAULT_PLACE }
  showDialog.value = true
}

// 打开编辑对话框
const openEditDialog = (place: Place) => {
  currentPlace.value = { ...place }
  showDialog.value = true
}

// 保存地点
const savePlace = async (placeData: Partial<Place>) => {
  uiStore.setLoading(true)

  try {
    let response

    if (placeData.id && placeData.id > 0) {
      // 更新现有地点
      response = await placeApi.updatePlace(placeData.id, placeData as Place)
      if (response.code === 0 && response.data) {
        // 更新本地列表中的对应项
        const index = places.value.findIndex(p => p.id === placeData.id)
        if (index !== -1) {
          places.value[index] = response.data
        }
        showToast('地点更新成功', 'success')
      }
    } else {
      // 创建新地点
      response = await placeApi.createPlace(placeData as Omit<Place, 'id'>)
      if (response.code === 0 && response.data) {
        // 重置到第一页，确保新添加的地点显示
        currentPage.value = 1
        await loadPlaces(false)
        showToast('地点添加成功', 'success')
      }
    }

    if (response.code !== 0) {
      throw new Error(response.msg || '操作失败')
    }

    closeDialog()
  } catch {
    showToast('保存地点失败', 'error')
  } finally {
    uiStore.setLoading(false)
  }
}

// 关闭对话框
const closeDialog = () => {
  showDialog.value = false
  currentPlace.value = null
}

onMounted(() => {
  loadPlaces()
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
