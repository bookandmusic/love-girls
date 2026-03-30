<script setup lang="ts">
import "leaflet/dist/leaflet.css";

import L from "leaflet";
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from "vue";
import { PullRefresh as VanPullRefresh } from "vant";

import BaseIcon from "@/components/ui/BaseIcon.vue";
import ActionSheet, {
  type ActionSheetAction,
} from "@/components/ui/ActionSheet.vue";
import FloatingAddButton from "@/components/ui/FloatingAddButton.vue";
import { useLongPress } from "@/composables/useLongPress";
import MainLayout from "@/layouts/MainLayout.vue";
import type { Place } from "@/services/placeApi";
import { placeApi } from "@/services/placeApi";
import { usePlacesStore } from "@/stores/places";
import { useSystemStore } from "@/stores/system";
import { useUIStore } from "@/stores/ui";
import { useToast } from "@/utils/toastUtils";

import PlaceEditDialog from "./components/dialogs/PlaceEditDialog.vue";
import DeleteConfirmDialog from "./components/dialogs/DeleteConfirmDialog.vue";

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

const uiStore = useUIStore();
const systemStore = useSystemStore();
const placesStore = usePlacesStore();

const systemInfo = computed(() => systemStore.getSystemInfo);
const places = computed(() => {
  const list = placesStore.getPlaces;
  return [...list].sort(
    (a, b) => new Date(b.date).getTime() - new Date(a.date).getTime(),
  );
});

const mapRef = ref<HTMLDivElement | null>(null);
const showToast = useToast();
const isRefreshing = ref(false);

let map: L.Map | null = null;
const markerMap = new Map<number, L.Marker>();

const fetchPlaces = async () => {
  try {
    await placesStore.fetchPlaces();
    updateMarkers();
  } catch {
    showToast("获取地点数据失败", "error");
  }
};

onMounted(async () => {
  uiStore.setLoading(true);

  try {
    await Promise.all([systemStore.fetchSystemInfo(), fetchPlaces()]);

    if (places.value && places.value.length > 0) {
      await nextTick();
      initMap();
    }
  } catch {
    showToast("获取地点数据失败，稍后重试", "error");
  } finally {
    uiStore.setLoading(false);
  }
});

onBeforeUnmount(() => {
  if (map) {
    map.remove();
    map = null;
    markerMap.clear();
  }
});

const initialCenter: [number, number] = [35, 105];
const initialZoom = 4;

function resetMap() {
  if (!map) return;
  map.setView(initialCenter, initialZoom);
}

function addResetButtonTogether() {
  if (!map) return;

  const zoomControlContainer = map.zoomControl.getContainer();
  if (!zoomControlContainer) return;

  const resetBtn = L.DomUtil.create(
    "a",
    "leaflet-control-zoom-reset",
    zoomControlContainer,
  );
  resetBtn.title = "重置地图";
  resetBtn.style.width = "30px";
  resetBtn.style.height = "30px";
  resetBtn.style.lineHeight = "30px";
  resetBtn.style.textAlign = "center";
  resetBtn.style.display = "block";
  resetBtn.style.cursor = "pointer";
  resetBtn.style.userSelect = "none";
  resetBtn.style.fill = "currentColor";
  resetBtn.style.color = "currentColor";
  resetBtn.style.boxSizing = "border-box";

  const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 512 512"><path fill="currentColor" fill-rule="evenodd" d="M426.667 106.667v42.666L358 149.33c36.077 31.659 58.188 77.991 58.146 128.474c-.065 78.179-53.242 146.318-129.062 165.376s-154.896-15.838-191.92-84.695C58.141 289.63 72.637 204.42 130.347 151.68a85.33 85.33 0 0 0 33.28 30.507a124.59 124.59 0 0 0-46.294 97.066c1.05 69.942 58.051 126.088 128 126.08c64.072 1.056 118.71-46.195 126.906-109.749c6.124-47.483-15.135-92.74-52.236-118.947L320 256h-42.667V106.667zM202.667 64c23.564 0 42.666 19.103 42.666 42.667s-19.102 42.666-42.666 42.666S160 130.231 160 106.667S179.103 64 202.667 64" stroke-width="13" stroke="currentColor"/></svg>`;
  resetBtn.innerHTML = `<div style="display: flex; align-items: center; justify-content: center; height: 100%; width: 100%;">${svg}</div>`;

  L.DomEvent.disableClickPropagation(resetBtn);

  L.DomEvent.on(resetBtn, "click", (e) => {
    L.DomEvent.stopPropagation(e);
    resetMap();
  });
}

function updateMarkers() {
  if (!map) return;

  markerMap.forEach((marker) => marker.remove());
  markerMap.clear();

  places.value.forEach((place: Place) => {
    const imgHtml = place.image?.file
      ? `<img src="${place.image.file.thumbnail || place.image.file.url}" class="w-full h-28 object-cover rounded" />`
      : "";

    const popupHtml = `
    <div class="w-56">
      ${imgHtml}
      <span class="font-bold text-lg"> ${place.name} </span>
      <br/>
      ${place.date}<br/>
      ${place.description}
    </div>
    `;

    const marker = L.marker([place.latitude, place.longitude])
      .addTo(map!)
      .bindPopup(popupHtml);

    markerMap.set(place.id, marker);
  });
}

function initMap() {
  if (!mapRef.value || map) return;

  map = L.map(mapRef.value, {
    zoomControl: true,
    attributionControl: false,
  }).setView([35, 105], 4);

  L.tileLayer("https://map.lw1314.site/{z}/{x}/{y}.png", {
    maxZoom: 18,
  }).addTo(map);

  updateMarkers();

  requestAnimationFrame(() => {
    map?.invalidateSize();
  });
  addResetButtonTogether();
}

function flyToPlace(place: Place) {
  if (!map) return;

  const marker = markerMap.get(place.id);
  if (!marker) return;

  map.invalidateSize(false);

  map.flyTo([place.latitude, place.longitude], 12, {
    duration: 1.5,
    animate: true,
  });
}

const showActionSheet = ref(false);
const selectedPlace = ref<Place | null>(null);

const {
  onPointerDown,
  onPointerUp,
  onPointerLeave,
  onPointerCancel,
  onPointerMove,
} = useLongPress({
  duration: 500,
  onFinish: () => {
    if (selectedPlace.value) {
      showActionSheet.value = true;
    }
  },
});

const handlePointerDown = (place: Place, event: PointerEvent) => {
  selectedPlace.value = place;
  onPointerDown(event);
};

const handlePlaceClick = (place: Place) => {
  flyToPlace(place);
};

const actionSheetActions = computed<ActionSheetAction[]>(() => [
  {
    label: "编辑",
    handler: () => openEditDialog(selectedPlace.value),
  },
  {
    label: "删除",
    destructive: true,
    handler: () => openDeleteDialog(selectedPlace.value),
  },
]);

const showEditDialog = ref(false);
const editingPlace = ref<Place | null>(null);
const savingPlace = ref(false);

const DEFAULT_PLACE: Place = {
  id: 0,
  name: "",
  latitude: 0,
  longitude: 0,
  date: new Date().toISOString().substring(0, 10),
  image: undefined,
  description: "",
};

const openAddDialog = () => {
  editingPlace.value = { ...DEFAULT_PLACE };
  showEditDialog.value = true;
};

const openEditDialog = (place: Place | null) => {
  editingPlace.value = place;
  showEditDialog.value = true;
};

const handleSavePlace = async (place: Place) => {
  savingPlace.value = true;
  try {
    if (place.id) {
      await placeApi.updatePlace(place.id, place);
      showToast("地点更新成功", "success");
    } else {
      await placeApi.createPlace(place);
      showToast("地点添加成功", "success");
    }
    showEditDialog.value = false;
    await fetchPlaces();
  } catch {
    showToast("操作失败", "error");
  } finally {
    savingPlace.value = false;
  }
};

const showDeleteDialog = ref(false);
const deletingPlace = ref<Place | null>(null);
const deleting = ref(false);

const openDeleteDialog = (place: Place | null) => {
  deletingPlace.value = place;
  showDeleteDialog.value = true;
};

const handleDeletePlace = async () => {
  if (!deletingPlace.value) return;
  deleting.value = true;
  try {
    await placeApi.deletePlace(deletingPlace.value.id);
    showToast("地点删除成功", "success");
    showDeleteDialog.value = false;
    await fetchPlaces();
  } catch {
    showToast("删除失败", "error");
  } finally {
    deleting.value = false;
  }
};

const handleRefresh = async () => {
  await fetchPlaces();
  isRefreshing.value = false;
};
</script>

<template>
  <MainLayout
    title="足迹地图"
    subtitle="印刻我们的同行足迹"
    :start-date="systemInfo?.site.startDate"
    :show-empty-state="places.length === 0"
  >
    <template #empty-state>
      <BaseIcon
        name="place"
        size="w-24"
        style="color: var(--fe-text-secondary)"
      />
      <p class="font-bold text-xl mt-4 text-[var(--fe-text-primary)]">
        暂无地点数据
      </p>
      <p class="text-md mt-2 text-[var(--fe-text-secondary)]">
        期待标记第一个足迹
      </p>
    </template>

    <template #main-content>
      <div
        class="flex-grow flex flex-col overflow-hidden bg-[var(--fe-bg-gray)]/30"
      >
        <div class="p-4 md:p-6 flex-shrink-0">
          <div
            class="h-64 md:h-96 rounded-[var(--fe-radius-card)] overflow-hidden shadow-lg border border-white/40 relative"
          >
            <div ref="mapRef" class="w-full h-full z-0"></div>
          </div>
        </div>

        <van-pull-refresh
          v-model="isRefreshing"
          @refresh="handleRefresh"
          class="flex-grow"
        >
          <div class="overflow-y-auto custom-scrollbar px-4 md:px-8 pb-8">
            <h2
              class="text-xs font-bold text-[var(--fe-text-secondary)] uppercase tracking-widest mb-3 ml-1"
            >
              我们去过的地方
            </h2>

            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              <div
                v-for="place in places"
                :key="place.id"
                class="flex items-center p-4 glass-thick rounded-2xl border border-white/40 shadow-sm tap-feedback ios-transition active:scale-[0.98] cursor-pointer"
                @click="handlePlaceClick(place)"
                @pointerdown="handlePointerDown(place, $event)"
                @pointermove="onPointerMove"
                @pointerup="onPointerUp"
                @pointerleave="onPointerLeave"
                @pointercancel="onPointerCancel"
              >
                <div
                  class="w-12 h-12 rounded-xl bg-gradient-to-br from-[var(--fe-primary)] to-[#f8c9c0] flex items-center justify-center text-white font-bold mr-4 shadow-sm flex-shrink-0"
                >
                  {{ place.name.substring(0, 1) }}
                </div>
                <div class="min-w-0">
                  <h3 class="font-bold text-[var(--fe-text-primary)] truncate">
                    {{ place.name }}
                  </h3>
                  <p
                    class="text-xs font-medium text-[var(--fe-text-secondary)] mt-0.5"
                  >
                    {{ place.date }}
                  </p>
                </div>
                <div class="ml-auto">
                  <BaseIcon
                    name="right"
                    size="w-4 h-4"
                    color="var(--fe-text-secondary)"
                  />
                </div>
              </div>
            </div>

            <div class="h-20 md:hidden"></div>
          </div>
        </van-pull-refresh>
      </div>
    </template>
  </MainLayout>

  <FloatingAddButton :loading="savingPlace" @click="openAddDialog" />

  <ActionSheet
    v-model="showActionSheet"
    title="地点操作"
    :actions="actionSheetActions"
  />

  <PlaceEditDialog
    v-model:open="showEditDialog"
    :place="editingPlace"
    :loading="savingPlace"
    @confirm="handleSavePlace"
  />

  <DeleteConfirmDialog
    v-model:open="showDeleteDialog"
    :loading="deleting"
    title="删除地点"
    :message="`确定要删除「${deletingPlace?.name || ''}」吗？删除后无法恢复。`"
    @confirm="handleDeletePlace"
  />
</template>

<style scoped>
:deep(.leaflet-container) {
  font-family: inherit;
}
</style>
