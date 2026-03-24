<template>
  <GenericDialog
    :open="open"
    :title="props.place?.id ? '编辑地点' : '添加地点'"
    @cancel="closeDialog"
    :loading="loading"
    size-class="max-w-3xl h-screen"
  >
    <template #content>
      <div
        class="space-y-4 px-2 overflow-y-auto pr-4"
        style="max-height: calc(100vh - 180px)"
      >
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1"
            >地点名称 <span class="text-red-500">*</span></label
          >
          <div class="flex gap-2">
            <input
              v-model="form.name"
              type="text"
              class="flex-1 admin-input"
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
                  : 'bg-[var(--primary-color)] text-white hover:opacity-90'
              "
            >
              <span
                v-if="geocodingLoading"
                class="animate-spin h-5 w-5 border-2 border-white border-t-transparent rounded-full"
              ></span>
              <svg
                v-else
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                stroke-width="2"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                />
              </svg>
            </button>
          </div>
        </div>

        <div v-if="geocodingResults.length > 0">
          <label class="block text-sm font-medium text-gray-700 mb-1"
            >搜索结果</label
          >
          <div class="border border-gray-200 rounded-lg overflow-hidden">
            <div
              v-for="(result, index) in geocodingResults"
              :key="index"
              @click="selectGeocodingResult(result)"
              class="p-3 cursor-pointer hover:bg-gray-50 border-b last:border-b-0 transition-colors"
              :class="{
                'bg-[var(--primary-color)]/10': selectedResult === result,
              }"
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
          <p class="text-xs text-gray-400 mt-1">
            点击选择地址，或在下方手动输入经纬度
          </p>
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
              class="w-full admin-input"
              placeholder="纬度"
              :disabled="loading"
              @input="onCoordinatesChange"
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
              class="w-full admin-input"
              placeholder="经度"
              :disabled="loading"
              @input="onCoordinatesChange"
            />
          </div>
        </div>

        <div v-if="showMapPreview">
          <label class="block text-sm font-medium text-gray-700 mb-1"
            >地图预览</label
          >
          <div class="h-48 rounded-lg overflow-hidden border border-gray-200">
            <div ref="mapRef" class="w-full h-full"></div>
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1"
            >日期 <span class="text-red-500">*</span></label
          >
          <input
            v-model="form.date"
            type="date"
            class="w-full admin-input"
            :disabled="loading"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1"
            >图片</label
          >
          <div class="mt-1">
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
          <label class="block text-sm font-medium text-gray-700 mb-1"
            >描述</label
          >
          <textarea
            v-model="form.description"
            class="w-full admin-input"
            rows="3"
            placeholder="请输入描述"
            :disabled="loading"
          ></textarea>
        </div>
      </div>
    </template>

    <template #actions>
      <div class="w-full flex">
        <div class="flex-1 text-center cursor-pointer" @click="closeDialog">
          取消
        </div>
        <div
          class="w-1/2 border-l border-gray-300 text-center cursor-pointer text-blue-500"
          @click="handleSave"
        >
          确认
        </div>
      </div>
    </template>
  </GenericDialog>

  <GenericDialog
    :open="showConfirmDialog"
    title="确认保存"
    :loading="loading"
    size-class="w-md h-md"
    @cancel="cancelConfirm"
  >
    <template #content>
      <p class="text-gray-700">
        {{
          props.place?.id
            ? "您确定要保存对这个地点的更改吗？"
            : "您确定要添加这个新地点吗？"
        }}
      </p>
    </template>
    <template #actions>
      <div class="w-full flex">
        <div class="flex-1 text-center cursor-pointer" @click="cancelConfirm">
          取消
        </div>
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
import "leaflet/dist/leaflet.css";

import L from "leaflet";
import { computed, nextTick, onBeforeUnmount, ref, watch } from "vue";

import GenericDialog from "@/components/ui/GenericDialog.vue";
import type { GeocodingResult } from "@/services/geocodingApi";
import { geocodingApi } from "@/services/geocodingApi";
import type { Place } from "@/services/placeApi";
import { useToast } from "@/utils/toastUtils";

delete (L.Icon.Default.prototype as unknown as { [key: string]: unknown })
  ._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl:
    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon-2x.png",
  iconUrl:
    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon.png",
  shadowUrl:
    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-shadow.png",
});

const showToast = useToast();

interface Props {
  open: boolean;
  loading: boolean;
  place?: Place | null;
}

interface Emits {
  (e: "update:open", open: boolean): void;
  (e: "confirm", place: Place): void;
  (e: "upload", event: Event): void;
  (e: "cancel"): void;
}

const props = withDefaults(defineProps<Props>(), {
  place: null,
  loading: false,
});

const emit = defineEmits<Emits>();

const closeDialog = () => {
  emit("update:open", false);
  emit("cancel");
  cleanupMap();
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
const showConfirmDialog = ref(false);

const geocodingResults = ref<GeocodingResult[]>([]);
const selectedResult = ref<GeocodingResult | null>(null);
const geocodingLoading = ref(false);

const mapRef = ref<HTMLDivElement | null>(null);
let map: L.Map | null = null;
let marker: L.Marker | null = null;

const showMapPreview = computed(() => {
  return form.value.latitude !== 0 || form.value.longitude !== 0;
});

watch(
  () => props.open,
  (isOpen) => {
    if (!isOpen) {
      cleanupMap();
      geocodingResults.value = [];
      selectedResult.value = null;
    }
  },
);

watch(
  () => props.place,
  (newPlace) => {
    form.value = newPlace ? { ...newPlace } : { ...DEFAULT_PLACE };
    if (newPlace && (newPlace.latitude || newPlace.longitude)) {
      nextTick(() => {
        initMap();
      });
    }
  },
  { deep: true, immediate: true },
);

watch(showMapPreview, (shouldShow) => {
  if (shouldShow) {
    nextTick(() => {
      initMap();
    });
  }
});

function cleanupMap() {
  if (map) {
    map.remove();
    map = null;
    marker = null;
  }
}

function initMap() {
  if (!mapRef.value || map) return;

  const lat = form.value.latitude || 35;
  const lon = form.value.longitude || 105;

  map = L.map(mapRef.value, {
    zoomControl: true,
    attributionControl: false,
  }).setView([lat, lon], form.value.latitude ? 12 : 4);

  L.tileLayer("https://map.lw1314.site/{z}/{x}/{y}.png", {
    maxZoom: 18,
  }).addTo(map);

  if (form.value.latitude && form.value.longitude) {
    marker = L.marker([form.value.latitude, form.value.longitude]).addTo(map);
  }

  requestAnimationFrame(() => {
    map?.invalidateSize();
  });
}

function updateMapMarker() {
  if (!map) return;

  if (marker) {
    map.removeLayer(marker);
  }

  if (form.value.latitude && form.value.longitude) {
    marker = L.marker([form.value.latitude, form.value.longitude]).addTo(map);
    map.setView([form.value.latitude, form.value.longitude], 12);
  }
}

function onCoordinatesChange() {
  selectedResult.value = null;
  updateMapMarker();
}

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

  nextTick(() => {
    initMap();
    updateMapMarker();
  });
}

const triggerImageUpload = () => {
  if (imageInputRef.value) {
    imageInputRef.value.click();
  }
};

const handleSelectedImageUpload = async (event: Event) => {
  emit("upload", event);
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
  if (form.value.latitude < -90 || form.value.latitude > 90) {
    showToast("纬度范围应为-90到90", "error");
    return;
  }
  if (form.value.longitude < -180 || form.value.longitude > 180) {
    showToast("经度范围应为-180到180", "error");
    return;
  }
  if (!form.value.date) {
    showToast("请选择日期", "error");
    return;
  }
  showConfirmDialog.value = true;
};

const cancelConfirm = () => {
  showConfirmDialog.value = false;
};

const confirmSave = async () => {
  showConfirmDialog.value = false;
  emit("confirm", form.value);
};

onBeforeUnmount(() => {
  cleanupMap();
});
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

:deep(.leaflet-container) {
  font-family: inherit;
}
</style>
