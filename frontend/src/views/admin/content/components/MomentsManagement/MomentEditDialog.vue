<template>
  <GenericDialog
    :open="open"
    :title="props.moment?.id ? '编辑动态' : '添加动态'"
    @cancel="closeDialog"
    :loading="loading"
  >
    <template #content>
      <div class="h-full flex flex-col">
        <div class="mb-3">
          <label class="block text-sm font-medium text-gray-700 mb-1">状态</label>
          <div class="flex items-center bg-gray-100 rounded-lg p-1">
            <button
              type="button"
              @click="localMoment.isPublic = true"
              class="flex-1 py-2 text-sm font-medium rounded-md transition-all duration-200"
              :class="localMoment.isPublic ? 'bg-white text-blue-600 shadow-sm' : 'text-gray-600'"
              :disabled="loading"
            >
              公开
            </button>
            <button
              type="button"
              @click="localMoment.isPublic = false"
              class="flex-1 py-2 text-sm font-medium rounded-md transition-all duration-200"
              :class="!localMoment.isPublic ? 'bg-white text-blue-600 shadow-sm' : 'text-gray-600'"
              :disabled="loading"
            >
              私密
            </button>
          </div>
        </div>
        <div class="mb-3">
          <label class="block text-sm font-medium text-gray-700 mb-1"
            >内容 <span class="text-red-500">*</span></label
          >
          <textarea
            v-model="localMoment.content"
            class="w-full win11-input"
            rows="4"
            placeholder="请输入动态内容"
            :disabled="loading"
          ></textarea>
        </div>

        <!-- 图片上传区域 -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1"
            >图片 <span class="text-xs text-gray-400">支持多张图片，JPG/PNG 格式</span></label
          >
          <div class="min-h-0 mt-1 flex-1 overflow-x-hidden overflow-y-auto">
            <div class="px-2 py-4 grid grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-3">
              <!-- 已上传图片列表 -->
              <div
                v-for="image in localMoment.images"
                :key="image.id"
                class="relative aspect-square group"
              >
                <img
                  :src="image.file?.url || ''"
                  class="w-full h-full object-cover rounded border border-gray-300"
                />
                <!-- 删除按钮：手机常显，PC hover 显 -->
                <button
                  type="button"
                  @click.stop="removeImage(image.id)"
                  class="absolute -top-1.5 -right-1.5 bg-red-500 text-white rounded-full w-6 h-6 flex items-center justify-center text-xs z-10 shadow-md opacity-100 sm:opacity-0 sm:group-hover:opacity-100 transition-opacity duration-200"
                  :disabled="loading"
                >
                  ×
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
          </div>
        </div>

        <!-- 隐藏的文件输入框 -->
        <input
          ref="imageInputRef"
          type="file"
          accept="image/*"
          @change="uploadImage"
          class="hidden"
          multiple
        />
      </div>
    </template>
    <template #actions>
      <div class="w-full flex">
        <div class="flex-1 text-center cursor-pointer" @click="closeDialog">取消</div>
        <div
          class="w-1/2 border-l border-gray-300 text-center cursor-pointer text-blue-500"
          @click="handleSave()"
        >
          确认
        </div>
      </div>
    </template>
  </GenericDialog>

  <!-- 确认对话框 -->
  <GenericDialog
    :open="showConfirmDialog"
    title="确认保存"
    :loading="loading"
    size-class="max-w-md"
  >
    <template #content>
      <p class="text-gray-700">
        {{ props.moment?.id ? '您确定要保存对这条动态的更改吗？' : '您确定要添加这条新动态吗？' }}
      </p>
    </template>
    <template #actions>
      <div class="w-full flex">
        <div class="flex-1 text-center cursor-pointer" @click="cancelConfirm">取消</div>
        <div
          class="w-1/2 border-l border-gray-300 text-center cursor-pointer text-blue-500"
          @click="confirmSave"
        >
          确定
        </div>
      </div>
    </template>
  </GenericDialog>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue'

import GenericDialog from '@/components/ui/GenericDialog.vue'
import { type Moment } from '@/services/momentApi'
import { useToast } from '@/utils/toastUtils'

const showToast = useToast()

const props = defineProps({
  open: {
    type: Boolean,
    required: true,
  },
  loading: {
    type: Boolean,
    default: false,
  },
  moment: {
    type: Object as () => Moment,
    default: null,
  },
})

interface Emits {
  (e: 'update:open', open: boolean): void
  (e: 'confirm', moment: Moment): void
  (e: 'upload', event: Event): void
  (e: 'cancel'): void
}

const emit = defineEmits<Emits>()
const closeDialog = () => {
  emit('update:open', false)
  emit('cancel')
}

// 定义默认动态对象
const DEFAULT_MOMENT = {
  id: 0,
  content: '',
  isPublic: true,
  images: [],
  likes: 0,
  author: { name: '系统用户', avatar: '' },
  createdAt: new Date().toISOString(),
}

// 创建本地响应式副本
const localMoment = reactive({ ...(props.moment || { ...DEFAULT_MOMENT }) })

// 监听对话框打开状态，只在打开时同步数据
watch(
  () => props.open,
  isOpen => {
    if (isOpen) {
      Object.assign(localMoment, props.moment || { ...DEFAULT_MOMENT })
    }
  }
)

// 单独监听图片变化，只更新图片字段（用于上传图片后同步）
watch(
  () => props.moment?.images,
  newImages => {
    if (newImages && props.open) {
      localMoment.images = [...newImages]
    }
  },
  { deep: true }
)

const imageInputRef = ref<HTMLInputElement>()
const showConfirmDialog = ref(false)

const triggerImageUpload = () => {
  if (imageInputRef.value) {
    imageInputRef.value.click()
  }
}

const uploadImage = (event: Event) => {
  emit('upload', event)
}

const removeImage = (imageId: number) => {
  localMoment.images = (localMoment.images || []).filter(img => img.id !== imageId)
}

const handleSave = async () => {
  // 校验必填项
  if (!localMoment.content.trim()) {
    showToast('请输入动态内容', 'error')
    return
  }
  showConfirmDialog.value = true
}

const cancelConfirm = () => {
  showConfirmDialog.value = false
}

const confirmSave = async () => {
  showConfirmDialog.value = false
  // 传递确认的数据给父组件
  emit('confirm', localMoment)
}
</script>
