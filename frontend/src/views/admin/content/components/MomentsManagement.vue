<template>
  <div class="w-full h-full flex flex-col overflow-hidden">
    <!-- 主体内容区域：包含列表和分页，允许内部滚动 -->
    <ul class="flex-1 min-h-0 overflow-y-auto pr-2">
      <MomentItem
        v-for="moment in moments"
        :key="moment.id"
        :moment="moment"
        @toggle-public="confirmToggle = $event"
        @edit="editMoment"
        @delete="confirmDelete = $event"
      />
    </ul>

    <!-- 分页：固定在底部，不随内容滚动 -->
    <div class="pt-3 flex-shrink-0">
      <Pagination
        :currentPage="currentPage"
        :totalPages="totalPages"
        @prev="handlePrevPage"
        @next="handleNextPage"
      />
    </div>

    <!-- 编辑对话框 -->
    <MomentEditDialog
      v-model:open="showEditDialog"
      :moment="editingMoment"
      :loading="uiStore.loading"
      @confirm="saveMoment"
      @cancel="closeEditDialog"
      @upload="handleImageUpload"
    />

    <!-- 删除确认对话框 -->
    <GenericDialog
      :open="!!confirmDelete"
      title="删除确认"
      :loading="uiStore.loading"
      @cancel="cancelDelete"
      size-class="max-w-md"
    >
      <template #content>
        <p class="text-gray-700">您确定要删除这条动态吗？此操作不可恢复。</p>
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

    <!-- 切换公开状态确认对话框 -->
    <GenericDialog
      :open="!!confirmToggle"
      :title="confirmToggle?.isPublic ? '设为私密确认' : '设为公开确认'"
      :loading="uiStore.loading"
      @cancel="cancelToggle"
      size-class="max-w-md"
    >
      <template #content>
        <p class="text-gray-700">
          {{
            confirmToggle?.isPublic
              ? '您确定要将这条动态设为私密吗？只有你自己可见。'
              : '您确定要将这条动态设为公开吗？所有人都将可以看到。'
          }}
        </p>
      </template>
      <template #actions>
        <div class="w-full flex">
          <div class="flex-1 text-center cursor-pointer" @click="cancelToggle">取消</div>
          <div
            class="w-1/2 border-l border-gray-300 text-center cursor-pointer text-red-500"
            @click="performTogglePublic"
          >
            {{ confirmToggle?.isPublic ? '设为私密' : '设为公开' }}
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
import { type Moment, momentApi, type Photo } from '@/services/momentApi'
import { uploadApi } from '@/services/upload'
import { useAuthStore } from '@/stores/auth'
import { useUIStore } from '@/stores/ui'
import { calculateFileHash } from '@/utils/fileUtils'
import { useToast } from '@/utils/toastUtils'

import MomentEditDialog from './MomentsManagement/MomentEditDialog.vue'
import MomentItem from './MomentsManagement/MomentItem.vue'

// 提取默认的moment对象为常量
const DEFAULT_MOMENT: Moment = {
  id: 0,
  content: '',
  isPublic: true,
  images: [],
  likes: 0,
  author: { name: '系统用户' },
  createdAt: new Date().toISOString(),
}

const props = defineProps<{ triggerAdd: boolean }>()
watch(
  () => props.triggerAdd,
  val => {
    if (val) addNewMoment()
  }
)

const uiStore = useUIStore()
const authStore = useAuthStore()
const showToast = useToast()

const moments = ref<Moment[]>([])
const totalMoments = ref(0)
const currentPage = ref(1)
const pageSize = ref(5)
const totalPages = computed(() => Math.ceil(totalMoments.value / pageSize.value) || 1)

// 加载动态列表
const loadMoments = async () => {
  uiStore.setLoading(true)
  try {
    const response = await momentApi.getMoments(currentPage.value, pageSize.value)

    moments.value = response.data.moments.map((moment: Moment) => ({
      ...moment,
      isPublic: moment.isPublic !== undefined ? moment.isPublic : true, // 默认为公开
    }))

    totalMoments.value =
      response.data.total || response.data.totalCount || response.data.moments.length
  } catch {
    showToast('加载动态失败', 'error')
  } finally {
    uiStore.setLoading(false)
  }
}

onMounted(() => {
  loadMoments()
})

const handlePrevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    loadMoments()
  }
}

const handleNextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    loadMoments()
  }
}

const confirmToggle = ref<Moment | null>(null)

// 取消切换公开状态
const cancelToggle = () => {
  confirmToggle.value = null
}

// 执行切换公开状态
const performTogglePublic = async () => {
  if (!confirmToggle.value) return

  uiStore.setLoading(true)
  try {
    // 发送API请求更新动态的公开状态
    const response = await momentApi.updateMomentPublic(confirmToggle.value.id, {
      isPublic: !confirmToggle.value.isPublic,
    })
    if (response.code === 0) {
      confirmToggle.value.isPublic = !confirmToggle.value.isPublic
      const momentIndex = moments.value.findIndex(u => u.id === confirmToggle!.value!.id)
      if (momentIndex !== -1) {
        moments.value[momentIndex] = {
          ...confirmToggle.value,
        }
      }

      // 更新本地状态
      showToast(confirmToggle.value.isPublic ? '动态已设为公开' : '动态已设为私密', 'success')
    } else {
      showToast('状态更新失败', 'error')
    }
  } catch {
    showToast('状态更新失败', 'error')
  } finally {
    uiStore.setLoading(false)
    confirmToggle.value = null
  }
}

// 编辑相关
const showEditDialog = ref(false)
const editingMoment = ref<Moment>({ ...DEFAULT_MOMENT }) // 使用常量

// 编辑动态
const editMoment = (moment: Moment) => {
  editingMoment.value = { ...moment }
  showEditDialog.value = true
}

// 保存动态（新增或编辑）
const saveMoment = async (momentData: Moment) => {
  uiStore.setLoading(true)
  try {
    if (momentData.id) {
      // 提取图片 ID 数组
      const imageIds = momentData.images?.map(img => img.id) || []

      // 编辑现有动态
      const response = await momentApi.updateMoment(momentData.id, {
        content: momentData.content,
        isPublic: momentData.isPublic,
        imageIds: imageIds,
        likes: momentData.likes,
      })

      // 更新本地数据
      const index = moments.value.findIndex(m => m.id === momentData.id)
      if (index !== -1) {
        moments.value[index] = {
          ...response.data,
          content: momentData.content,
          isPublic: momentData.isPublic,
          images: momentData.images,
        } // 保存图片信息
      }
      showToast('动态更新成功', 'success')
    } else {
      // 提取图片 ID 数组
      const imageIds = momentData.images?.map(img => img.id) || []

      // 添加新动态
      const response = await momentApi.createMoment({
        content: momentData.content || '',
        isPublic: momentData.isPublic || false,
        imageIds: imageIds, // 发送图片 ID 数组
        likes: 0,
        author: { name: '系统用户' },
        createdAt: new Date().toISOString(),
        userId: authStore.userInfo?.userId || 1, // 使用当前登录用户 ID，默认为 1
      })

      // 添加到列表开头
      moments.value.unshift({
        ...response.data,
        images: momentData.images || [], // 保存图片信息
        isPublic: response.data.isPublic !== undefined ? response.data.isPublic : true,
      })
      totalMoments.value++

      // 重置到第一页，确保新添加的动态显示
      currentPage.value = 1
      showToast('动态添加成功', 'success')
    }
  } catch {
    showToast('保存动态失败', 'error')
  } finally {
    uiStore.setLoading(false)
    // 关闭对话框并重置表单
    closeEditDialog()
  }
}

// 添加新动态
const addNewMoment = () => {
  editingMoment.value = { ...DEFAULT_MOMENT }
  showEditDialog.value = true
}

const closeEditDialog = () => {
  showEditDialog.value = false
  editingMoment.value = { ...DEFAULT_MOMENT }
}

// 处理图片上传
const handleImageUpload = async (event: Event): Promise<void> => {
  const target = event.target as HTMLInputElement
  if (!target.files || target.files.length === 0) return

  const files = Array.from(target.files)
  const uploadPromises = files.map(async file => {
    try {
      // 计算 MD5 哈希值
      const hash = await calculateFileHash(file)

      // 生成 path 参数：moments/年/月
      const now = new Date()
      const year = now.getFullYear()
      const month = String(now.getMonth() + 1).padStart(2, '0')
      const path = `moments/${year}/${month}`

      const formData = new FormData()
      formData.append('file', file)
      formData.append('hash', hash)
      formData.append('path', path)

      uiStore.setLoading(true)
      const response = await uploadApi.uploadImage(formData)

      if (response.data.code === 0) {
        const { file } = response.data.data

        // 创建 Photo 对象
        return {
          id: file.id,
          momentId: 0,
          file: file,
        } as Photo
      } else {
        showToast(`上传图片 ${file.name} 失败`, 'error')
        return null
      }
    } catch (error) {
      console.error('Upload error:', error)
      showToast(`上传图片 ${file.name} 失败`, 'error')
      return null
    }
  })

  try {
    const uploadedPhotos = await Promise.all(uploadPromises)
    const validPhotos = uploadedPhotos.filter(photo => photo !== null) as Photo[]

    if (validPhotos.length > 0) {
      const newImages = [...(editingMoment.value.images || [])]
      validPhotos.forEach(photo => {
        if (!newImages.some(existingImg => existingImg.id === photo.id)) {
          newImages.push(photo)
        }
      })
      editingMoment.value.images = newImages
      showToast(`${validPhotos.length} 张图片上传成功`, 'success')
    }
  } finally {
    uiStore.setLoading(false)
    // 重置文件输入，允许重复上传相同文件
    target.value = ''
  }
}

// 删除确认
const confirmDelete = ref<Moment | null>(null)

// 取消删除
const cancelDelete = () => {
  confirmDelete.value = null
}

// 删除动态
const performDelete = async () => {
  if (!confirmDelete.value || !confirmDelete.value.id) return

  const momentId = confirmDelete.value.id
  uiStore.setLoading(true)
  try {
    // 发送API请求删除动态
    await momentApi.deleteMoment(momentId)

    const index = moments.value.findIndex(m => m.id === momentId)
    if (index !== -1) {
      moments.value.splice(index, 1)
      totalMoments.value--

      // 检查当前页是否为空，如果是，则跳转到前一页（如果存在）
      if (moments.value.length === 0 && currentPage.value > 1) {
        currentPage.value--
      }
      // 重新加载当前页数据
      await loadMoments()
    }
    showToast('动态删除成功', 'success')
  } catch {
    showToast('删除失败', 'error')
  } finally {
    uiStore.setLoading(false)
    confirmDelete.value = null
  }
}
</script>
