<template>
  <GenericDialog
    :open="open"
    :title="props.place?.id ? '编辑地点' : '添加地点'"
    @cancel="closeDialog"
    :loading="loading"
  >
    <template #content>
      <div class="space-y-4 px-2 overflow-y-auto h-full">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1"
            >地点名称 <span class="text-red-500">*</span></label
          >
          <input
            v-model="form.name"
            type="text"
            class="w-full win11-input"
            placeholder="请输入地点名称"
            :disabled="loading"
          />
        </div>

        <div class="grid grid-cols-1 gap-x-4 sm:grid-cols-2">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1"
              >纬度 <span class="text-red-500">*</span>
              <span class="text-xs text-gray-400">(-90~90)</span></label
            >
            <input
              v-model.number="form.latitude"
              type="number"
              step="any"
              class="w-full win11-input"
              placeholder="纬度"
              :disabled="loading"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1"
              >经度 <span class="text-red-500">*</span>
              <span class="text-xs text-gray-400">(-180~180)</span></label
            >
            <input
              v-model.number="form.longitude"
              type="number"
              step="any"
              class="w-full win11-input"
              placeholder="经度"
              :disabled="loading"
            />
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1"
            >日期 <span class="text-red-500">*</span></label
          >
          <input v-model="form.date" type="date" class="w-full win11-input" :disabled="loading" />
        </div>

        <!-- 图片上传区域 -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">图片</label>
          <div class="mt-1">
            <!-- 显示已上传图片 -->
            <div v-if="form.image" class="relative aspect-w-16 aspect-h-9 mb-2">
              <img
                :src="form.image?.file?.url || ''"
                class="w-full h-48 object-cover rounded border border-gray-300"
                alt="地点图片"
              />
              <button
                type="button"
                @click="removeImage"
                class="absolute top-2 right-2 bg-red-500 text-white rounded-full w-6 h-6 flex items-center justify-center text-xs z-10 shadow-md"
                :disabled="loading"
              >
                ×
              </button>
            </div>

            <!-- 上传按钮 -->
            <div
              v-if="!form.image"
              @click="triggerImageUpload"
              class="flex items-center justify-center w-full h-32 border-2 border-dashed border-gray-300 rounded-lg cursor-pointer hover:border-[var(--primary-color)] hover:bg-gray-50"
              :disabled="loading"
            >
              <div class="text-center">
                <p class="text-gray-500">点击上传图片</p>
                <p class="text-xs text-gray-400 mt-1">JPG/PNG 格式</p>
              </div>
            </div>

            <!-- 隐藏的文件输入框 -->
            <input
              ref="imageInputRef"
              type="file"
              accept="image/*"
              @change="handleSelectedImageUpload"
              class="hidden"
            />
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">描述</label>
          <textarea
            v-model="form.description"
            class="w-full win11-input"
            rows="3"
            placeholder="请输入描述"
            :disabled="loading"
          ></textarea>
        </div>
      </div>
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
    @cancel="cancelConfirm"
  >
    <template #content>
      <p class="text-gray-700">
        {{ props.place?.id ? '您确定要保存对这个地点的更改吗？' : '您确定要添加这个新地点吗？' }}
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
import { ref, watch } from 'vue'

import GenericDialog from '@/components/ui/GenericDialog.vue'
import type { Place } from '@/services/placeApi'
import { useToast } from '@/utils/toastUtils'

const showToast = useToast()

interface Props {
  open: boolean
  loading: boolean
  place?: Place | null
}

interface Emits {
  (e: 'update:open', open: boolean): void
  (e: 'confirm', place: Place): void
  (e: 'upload', event: Event): void
  (e: 'cancel'): void
}

const props = withDefaults(defineProps<Props>(), {
  place: null,
  loading: false,
})

const emit = defineEmits<Emits>()

const closeDialog = () => {
  emit('update:open', false)
  emit('cancel')
}

const DEFAULT_PLACE: Place = {
  id: 0,
  name: '',
  latitude: 0,
  longitude: 0,
  date: new Date().toISOString().substring(0, 10),
  image: undefined,
  description: '',
}

const form = ref<Place>({ ...DEFAULT_PLACE })

const imageInputRef = ref<HTMLInputElement>()
const showConfirmDialog = ref(false)

// 监听编辑地点的变更
watch(
  () => props.place,
  newPlace => {
    form.value = newPlace ? { ...newPlace } : { ...DEFAULT_PLACE }
  },
  { deep: true, immediate: true }
)

const triggerImageUpload = () => {
  if (imageInputRef.value) {
    imageInputRef.value.click()
  }
}

const handleSelectedImageUpload = async (event: Event) => {
  emit('upload', event)
}

const removeImage = () => {
  form.value.image = undefined
}

const handleSave = async () => {
  // 校验必填项
  if (!form.value.name.trim()) {
    showToast('请输入地点名称', 'error')
    return
  }
  if (form.value.latitude === 0 && form.value.longitude === 0) {
    showToast('请输入纬度和经度', 'error')
    return
  }
  if (form.value.latitude < -90 || form.value.latitude > 90) {
    showToast('纬度范围应为-90到90', 'error')
    return
  }
  if (form.value.longitude < -180 || form.value.longitude > 180) {
    showToast('经度范围应为-180到180', 'error')
    return
  }
  if (!form.value.date) {
    showToast('请选择日期', 'error')
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
