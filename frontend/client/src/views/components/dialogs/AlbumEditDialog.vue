<template>
  <GenericDialog
    :open="open"
    :title="album?.id ? '编辑相册' : '新建相册'"
    @cancel="closeDialog"
    :loading="loading"
  >
    <template #content>
      <form class="space-y-4 h-full">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1.5">
            相册名称 <span class="text-red-500">*</span>
          </label>
          <input
            v-model="form.name"
            type="text"
            class="w-full win11-input"
            placeholder="请输入相册名称"
            :disabled="loading"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1.5">
            描述
          </label>
          <textarea
            v-model="form.description"
            rows="3"
            class="w-full win11-input resize-none"
            placeholder="请输入相册描述"
            :disabled="loading"
          ></textarea>
        </div>
      </form>
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
import { ref, watch } from "vue";

import GenericDialog from "@/components/ui/GenericDialog.vue";
import { type Album } from "@/services/albumApi";
import { useToast } from "@/utils/toastUtils";

const showToast = useToast();

interface Props {
  open: boolean;
  loading?: boolean;
  album?: Album | null;
}

const props = withDefaults(defineProps<Props>(), {
  album: null,
  loading: false,
});

const emit = defineEmits<{
  (e: "update:open", open: boolean): void;
  (e: "confirm", album: Album): void;
  (e: "cancel"): void;
}>();

const DEFAULT_ALBUM: Album = {
  id: 0,
  name: "",
  description: "",
  coverImage: undefined,
  photoCount: 0,
  createdAt: new Date().toISOString(),
};

const form = ref<Album>({ ...DEFAULT_ALBUM });

const closeDialog = () => {
  emit("update:open", false);
  emit("cancel");
};

watch(
  () => props.album,
  (newAlbum) => {
    if (newAlbum) {
      form.value = { ...newAlbum };
    } else {
      form.value = { ...DEFAULT_ALBUM };
    }
  },
  { deep: true, immediate: true },
);

watch(
  () => props.open,
  (isOpen) => {
    if (isOpen) {
      if (props.album) {
        form.value = { ...props.album };
      } else {
        form.value = { ...DEFAULT_ALBUM };
      }
    }
  },
);

const handleSave = async () => {
  if (!form.value.name.trim()) {
    showToast("请输入相册名称", "error");
    return;
  }
  emit("confirm", form.value);
};
</script>
