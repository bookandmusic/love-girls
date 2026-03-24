<template>
  <GenericDialog
    :open="open"
    :title="moment?.id ? '编辑动态' : '发布动态'"
    @cancel="closeDialog"
    :loading="loading"
  >
    <template #content>
      <div class="space-y-4 h-full">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1.5">
            状态
          </label>
          <div class="flex items-center bg-gray-100 rounded-lg p-1">
            <button
              type="button"
              @click="localMoment.isPublic = true"
              class="flex-1 py-2 text-sm font-medium rounded-md transition-all duration-200"
              :class="
                localMoment.isPublic
                  ? 'bg-white text-[var(--fe-primary)] shadow-sm'
                  : 'text-gray-500'
              "
              :disabled="loading"
            >
              公开
            </button>
            <button
              type="button"
              @click="localMoment.isPublic = false"
              class="flex-1 py-2 text-sm font-medium rounded-md transition-all duration-200"
              :class="
                !localMoment.isPublic
                  ? 'bg-white text-[var(--fe-primary)] shadow-sm'
                  : 'text-gray-500'
              "
              :disabled="loading"
            >
              私密
            </button>
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1.5">
            内容 <span class="text-red-500">*</span>
          </label>
          <textarea
            v-model="localMoment.content"
            class="w-full win11-input resize-none"
            rows="4"
            placeholder="记录这一刻..."
            :disabled="loading"
          ></textarea>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1.5">
            图片
          </label>
          <div class="grid grid-cols-4 gap-2">
            <div
              v-for="image in localMoment.images"
              :key="image.id"
              class="relative aspect-square"
            >
              <img
                :src="image.file?.url || ''"
                class="w-full h-full object-cover rounded-lg border border-gray-200"
              />
              <button
                type="button"
                @click.stop="removeImage(image.id)"
                class="absolute -top-1.5 -right-1.5 bg-red-500 text-white rounded-full w-5 h-5 flex items-center justify-center text-xs shadow-md"
                :disabled="loading"
              >
                ×
              </button>
            </div>

            <div
              @click="triggerImageUpload"
              class="aspect-square border-2 border-dashed border-gray-300 rounded-lg flex items-center justify-center cursor-pointer hover:border-[var(--fe-primary)] hover:bg-gray-50 transition-all"
              :disabled="loading"
            >
              <span class="text-xl text-gray-400">+</span>
            </div>
          </div>
          <input
            ref="imageInputRef"
            type="file"
            accept="image/*"
            @change="uploadImage"
            class="hidden"
            multiple
          />
        </div>
      </div>
    </template>
    <template #actions>
      <button
        class="flex-1 py-3.5 text-center text-gray-500 font-medium hover:bg-gray-50 active:bg-gray-100 transition-colors"
        @click="closeDialog"
      >
        取消
      </button>
      <button
        class="flex-1 py-3.5 text-center text-[var(--fe-primary)] font-semibold border-l border-gray-100 hover:bg-gray-50 active:bg-gray-100 transition-colors"
        @click="handleSave"
      >
        确认
      </button>
    </template>
  </GenericDialog>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from "vue";

import GenericDialog from "@/components/ui/GenericDialog.vue";
import { type Moment, type Photo } from "@/services/momentApi";
import { uploadApi } from "@/services/upload";
import { calculateFileHash } from "@/utils/fileUtils";
import { useToast } from "@/utils/toastUtils";

const showToast = useToast();

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
});

interface Emits {
  (e: "update:open", open: boolean): void;
  (e: "confirm", moment: Moment): void;
  (e: "cancel"): void;
}

const emit = defineEmits<Emits>();

const closeDialog = () => {
  emit("update:open", false);
  emit("cancel");
};

const formatLocalDateTime = () => {
  const now = new Date();
  const year = now.getFullYear();
  const month = String(now.getMonth() + 1).padStart(2, "0");
  const day = String(now.getDate()).padStart(2, "0");
  const hours = String(now.getHours()).padStart(2, "0");
  const minutes = String(now.getMinutes()).padStart(2, "0");
  const seconds = String(now.getSeconds()).padStart(2, "0");
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
};

const DEFAULT_MOMENT = {
  id: 0,
  content: "",
  isPublic: true,
  images: [],
  likes: 0,
  author: { name: "系统用户", avatar: "" },
  createdAt: "",
};

const localMoment = reactive({ ...(props.moment || { ...DEFAULT_MOMENT }) });

watch(
  () => props.open,
  (isOpen) => {
    if (isOpen) {
      if (props.moment && props.moment.id) {
        Object.assign(localMoment, props.moment);
      } else {
        Object.assign(localMoment, {
          ...DEFAULT_MOMENT,
          createdAt: formatLocalDateTime(),
        });
      }
    }
  },
);

const imageInputRef = ref<HTMLInputElement>();

const triggerImageUpload = () => {
  if (imageInputRef.value) {
    imageInputRef.value.click();
  }
};

const uploadImage = async (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (!target.files || target.files.length === 0) return;

  const files = Array.from(target.files);

  try {
    for (const file of files) {
      const hash = await calculateFileHash(file);

      const now = new Date();
      const year = now.getFullYear();
      const month = String(now.getMonth() + 1).padStart(2, "0");
      const path = `moments/${year}/${month}`;

      const formData = new FormData();
      formData.append("file", file);
      formData.append("hash", hash);
      formData.append("path", path);

      const response = await uploadApi.uploadImage(formData);
      if (response.data.code === 0) {
        const newImage: Photo = {
          id: response.data.data.file.id,
          momentId: localMoment.id || 0,
          file: response.data.data.file,
        };
        if (!localMoment.images) {
          localMoment.images = [];
        }
        if (!localMoment.images.some((img) => img.id === newImage.id)) {
          localMoment.images.push(newImage);
        }
      }
    }
    showToast("图片上传成功", "success");
  } catch {
    showToast("图片上传失败", "error");
  }

  target.value = "";
};

const removeImage = (imageId: number) => {
  localMoment.images = (localMoment.images || []).filter(
    (img) => img.id !== imageId,
  );
};

const handleSave = async () => {
  if (!localMoment.content.trim()) {
    showToast("请输入动态内容", "error");
    return;
  }
  emit("confirm", localMoment);
};
</script>
