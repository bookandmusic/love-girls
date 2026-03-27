<template>
  <GenericDialog
    :open="open"
    :title="place?.id ? '编辑地点' : '添加地点'"
    @cancel="closeDialog"
    :loading="loading"
  >
    <template #content>
      <div class="space-y-4 h-full">
        <div class="relative" ref="searchResultRef">
          <label class="block text-sm font-medium text-gray-700 mb-1.5">
            地点名称 <span class="text-red-500">*</span>
          </label>
          <div class="flex gap-2">
            <input
              v-model="form.name"
              type="text"
              class="flex-1 win11-input"
              placeholder="请输入地点名称"
              :disabled="loading"
              @keyup.enter="searchAddress"
            />
            <button
              type="button"
              @click="searchAddress"
              :disabled="loading || geocodingLoading || !form.name.trim()"
              class="w-10 h-10 flex items-center justify-center rounded-lg transition-colors"
              :class="
                loading || geocodingLoading || !form.name.trim()
                  ? 'bg-gray-200 text-gray-400 cursor-not-allowed'
                  : 'bg-[var(--fe-primary)] text-white hover:opacity-90'
              "
            >
              <span
                v-if="geocodingLoading"
                class="animate-spin h-5 w-5 border-2 border-white border-t-transparent rounded-full"
              ></span>
              <svg
                v-else
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                />
              </svg>
            </button>
          </div>

          <Transition name="expand">
            <div
              v-if="geocodingResults.length > 0"
              class="absolute top-full left-0 right-0 mt-2 z-20 glass-regular rounded-xl border border-white/40 shadow-xl h-48 overflow-y-auto"
            >
              <div class="p-2 space-y-1">
                <div
                  v-for="(result, index) in geocodingResults"
                  :key="index"
                  @click="selectGeocodingResult(result)"
                  class="rounded-xl p-3 border transition-all cursor-pointer active:scale-[0.98]"
                  :class="
                    selectedResult === result
                      ? 'border-[var(--fe-primary)] bg-[var(--fe-primary)]/5'
                      : 'border-transparent hover:bg-white/30'
                  "
                >
                  <div class="text-sm text-gray-800 line-clamp-2">
                    {{ result.displayName }}
                  </div>
                  <div class="text-xs text-gray-500 mt-1">
                    纬度: {{ parseFloat(result.lat).toFixed(6) }} | 经度:
                    {{ parseFloat(result.lon).toFixed(6) }}
                  </div>
                </div>
              </div>
            </div>
          </Transition>
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">
              纬度 <span class="text-red-500">*</span>
            </label>
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
            <label class="block text-sm font-medium text-gray-700 mb-1.5">
              经度 <span class="text-red-500">*</span>
            </label>
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
          <label class="block text-sm font-medium text-gray-700 mb-1.5">
            日期 <span class="text-red-500">*</span>
          </label>
          <input
            v-model="form.date"
            type="date"
            class="w-full win11-input"
            :disabled="loading"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1.5">
            图片
          </label>
          <div v-if="form.image" class="relative mb-2">
            <img
              :src="form.image?.file?.url || ''"
              class="w-full h-32 object-cover rounded-lg border border-gray-200"
            />
            <button
              type="button"
              @click="removeImage"
              class="absolute top-2 right-2 bg-red-500 text-white rounded-full w-6 h-6 flex items-center justify-center text-xs shadow-md"
              :disabled="loading"
            >
              ×
            </button>
          </div>
          <div
            v-if="!form.image"
            @click="triggerImageUpload"
            class="flex items-center justify-center w-full h-24 border-2 border-dashed border-gray-300 rounded-lg cursor-pointer hover:border-[var(--fe-primary)] hover:bg-gray-50 transition-all"
            :disabled="loading"
          >
            <div class="text-center text-gray-500">
              <p>点击上传图片</p>
            </div>
          </div>
          <input
            ref="imageInputRef"
            type="file"
            accept="image/*"
            @change="handleSelectedImageUpload"
            class="hidden"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1.5">
            描述
          </label>
          <textarea
            v-model="form.description"
            class="w-full win11-input resize-none"
            rows="3"
            placeholder="请输入描述"
            :disabled="loading"
          ></textarea>
        </div>
      </div>
    </template>
    <template #actions>
      <button
        class="flex-1 py-3.5 text-center text-gray-700 font-medium hover:bg-gray-100 active:bg-gray-200 transition-colors"
        @click="closeDialog"
      >
        取消
      </button>
      <button
        class="flex-1 py-3.5 text-center text-[var(--fe-primary-dark)] font-semibold border-l border-gray-100 hover:bg-gray-100 active:bg-gray-200 transition-colors"
        @click="handleSave"
      >
        确认
      </button>
    </template>
  </GenericDialog>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from "vue";

import GenericDialog from "@/components/ui/GenericDialog.vue";
import type { GeocodingResult } from "@/services/geocodingApi";
import { geocodingApi } from "@/services/geocodingApi";
import type { Place } from "@/services/placeApi";
import { uploadApi } from "@/services/upload";
import { calculateFileHash } from "@/utils/fileUtils";
import { useToast } from "@/utils/toastUtils";

const showToast = useToast();

interface Props {
  open: boolean;
  loading?: boolean;
  place?: Place | null;
}

const props = withDefaults(defineProps<Props>(), {
  place: null,
  loading: false,
});

interface Emits {
  (e: "update:open", open: boolean): void;
  (e: "confirm", place: Place): void;
  (e: "cancel"): void;
}

const emit = defineEmits<Emits>();

const closeDialog = () => {
  emit("update:open", false);
  emit("cancel");
  geocodingResults.value = [];
  selectedResult.value = null;
};

const DEFAULT_PLACE: Place = {
  id: 0,
  name: "",
  latitude: 0,
  longitude: 0,
  date: new Date().toISOString().substring(0, 10),
  image: undefined,
  description: "",
};

const form = ref<Place>({ ...DEFAULT_PLACE });
const imageInputRef = ref<HTMLInputElement>();
const searchResultRef = ref<HTMLElement | null>(null);
const geocodingResults = ref<GeocodingResult[]>([]);
const selectedResult = ref<GeocodingResult | null>(null);
const geocodingLoading = ref(false);

const handleClickOutside = (event: MouseEvent) => {
  if (
    searchResultRef.value &&
    !searchResultRef.value.contains(event.target as Node)
  ) {
    geocodingResults.value = [];
  }
};

onMounted(() => {
  document.addEventListener("click", handleClickOutside);
});

onUnmounted(() => {
  document.removeEventListener("click", handleClickOutside);
});

watch(
  () => props.open,
  (isOpen) => {
    if (!isOpen) {
      geocodingResults.value = [];
      selectedResult.value = null;
    }
  },
);

watch(
  () => props.place,
  (newPlace) => {
    form.value = newPlace ? { ...newPlace } : { ...DEFAULT_PLACE };
  },
  { deep: true, immediate: true },
);

async function searchAddress() {
  if (!form.value.name.trim()) return;

  geocodingLoading.value = true;
  selectedResult.value = null;

  try {
    const results = await geocodingApi.search(form.value.name);
    geocodingResults.value = results;

    if (results.length === 0) {
      showToast("未找到匹配的地址", "info");
    }
  } catch {
    showToast("地址搜索失败，请稍后重试", "error");
    geocodingResults.value = [];
  } finally {
    geocodingLoading.value = false;
  }
}

function selectGeocodingResult(result: GeocodingResult) {
  selectedResult.value = result;
  form.value.latitude = parseFloat(result.lat);
  form.value.longitude = parseFloat(result.lon);
  geocodingResults.value = [];
}

const triggerImageUpload = () => {
  if (imageInputRef.value) {
    imageInputRef.value.click();
  }
};

const handleSelectedImageUpload = async (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (!target.files || target.files.length === 0) return;

  const file = target.files[0];
  if (!file) return;

  try {
    const hash = await calculateFileHash(file);

    const now = new Date();
    const year = now.getFullYear();
    const month = String(now.getMonth() + 1).padStart(2, "0");
    const path = `places/${year}/${month}`;

    const formData = new FormData();
    formData.append("file", file);
    formData.append("hash", hash);
    formData.append("path", path);

    const response = await uploadApi.uploadImage(formData);
    if (response.data.code === 0) {
      form.value.image = {
        id: response.data.data.file.id,
        placeId: form.value.id || 0,
        file: response.data.data.file,
      };
      showToast("图片上传成功", "success");
    }
  } catch {
    showToast("图片上传失败", "error");
  }

  target.value = "";
};

const removeImage = () => {
  form.value.image = undefined;
};

const handleSave = async () => {
  if (!form.value.name.trim()) {
    showToast("请输入地点名称", "error");
    return;
  }
  if (form.value.latitude === 0 && form.value.longitude === 0) {
    showToast("请输入纬度和经度", "error");
    return;
  }
  if (!form.value.date) {
    showToast("请选择日期", "error");
    return;
  }
  emit("confirm", form.value);
};
</script>

<style scoped>
.expand-enter-active,
.expand-leave-active {
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.expand-enter-from,
.expand-leave-to {
  opacity: 0;
  transform: translateY(-8px) scale(0.95);
}
</style>
