<template>
  <GenericDialog :open="open" :title="dialogTitle" @cancel="closeDialog" :loading="loading">
    <template #content>
      <form class="space-y-4">
        <div class="mb-3">
          <label class="block text-sm font-medium text-gray-700 mb-1"
            >相册名称 <span class="text-red-500">*</span></label
          >
          <input
            v-model="form.name"
            type="text"
            class="w-full win11-input"
            placeholder="请输入相册名称"
            :disabled="loading"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">描述</label>
          <textarea
            v-model="form.description"
            rows="3"
            class="w-full win11-input"
            placeholder="请输入相册描述"
            :disabled="loading"
          ></textarea>
        </div>
      </form>
    </template>
    <template #actions>
      <div class="w-full flex">
        <div class="flex-1 text-center cursor-pointer" @click="closeDialog">取消</div>
        <div
          class="w-1/2 border-l border-gray-300 text-center cursor-pointer text-blue-500"
          @click="handleSave"
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
    size-class="w-md h-md"
  >
    <template #content>
      <p class="text-gray-700">
        {{ props.album?.id ? '您确定要保存对这个相册的更改吗？' : '您确定要添加这个新相册吗？' }}
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
import { computed, ref, watch } from 'vue'

import GenericDialog from '@/components/ui/GenericDialog.vue'
import { type Album } from '@/services/albumApi'
import { useToast } from '@/utils/toastUtils'

const showToast = useToast()

interface Props {
  open: boolean
  loading: boolean
  album?: Album | null
}

const props = withDefaults(defineProps<Props>(), {
  album: null,
  loading: false,
})

const emit = defineEmits<{
  (e: 'update:open', open: boolean): void
  (e: 'confirm', album: Album): void
  (e: 'cancel'): void
}>()

const closeDialog = () => {
  emit('update:open', false)
  emit('cancel')
}

// 定义默认相册对象
const DEFAULT_ALBUM: Album = {
  id: 0,
  name: '',
  description: '',
  coverImage: undefined,
  photoCount: 0,
  createdAt: new Date().toISOString(),
}

const form = ref<Album>({ ...DEFAULT_ALBUM })
const showConfirmDialog = ref(false)

// 使用计算属性计算对话框标题
const dialogTitle = computed(() => {
  if (props.album && props.album.id) {
    return '编辑相册'
  }
  return '新建相册'
})

// 监听props.album的变化，更新表单数据
watch(
  () => props.album,
  newAlbum => {
    if (newAlbum) {
      // 确保新专辑存在后再复制其属性
      form.value = { ...newAlbum }
    } else {
      // 如果没有新专辑，清空表单
      form.value = { ...DEFAULT_ALBUM }
    }
  },
  { deep: true, immediate: true }
)

const handleSave = async () => {
  // 校验必填项
  if (!form.value.name.trim()) {
    showToast('请输入相册名称', 'error')
    return
  }
  showConfirmDialog.value = true
}

const cancelConfirm = () => {
  showConfirmDialog.value = false
}

const confirmSave = async () => {
  showConfirmDialog.value = false
  emit('confirm', form.value)
}
</script>
