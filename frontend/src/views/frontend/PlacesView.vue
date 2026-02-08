PlacesView.vue
<script setup lang="ts">
import 'leaflet/dist/leaflet.css'

import L from 'leaflet'
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue'

import BaseIcon from '@/components/ui/BaseIcon.vue'
import MainLayout from '@/layouts/MainLayout.vue'
import type { Place } from '@/services/placeApi'
import { usePlacesStore } from '@/stores/places'
import { useSystemStore } from '@/stores/system'
import { useUIStore } from '@/stores/ui'
import { useToast } from '@/utils/toastUtils'

/* ========== 修复 Leaflet 默认图标问题 ========== */
delete (L.Icon.Default.prototype as unknown as { [key: string]: unknown })._getIconUrl
L.Icon.Default.mergeOptions({
  iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon-2x.png',
  iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon.png',
  shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-shadow.png',
})

/* ========== Store ========== */
const uiStore = useUIStore()
const systemStore = useSystemStore()
const placesStore = usePlacesStore()

// 获取系统信息和地点数据
const systemInfo = computed(() => systemStore.getSystemInfo)
const places = computed(() => placesStore.getPlaces)

/* ========== 状态 ========== */
const mapRef = ref<HTMLDivElement | null>(null)
const showToast = useToast()

let map: L.Map | null = null
const markerMap = new Map<number, L.Marker>()

/* ========== 生命周期 ========== */
onMounted(async () => {
  uiStore.setLoading(true)

  try {
    // 获取系统信息和地点数据
    await Promise.all([systemStore.fetchSystemInfo(), placesStore.fetchPlaces()])

    // 只有在当地点数据存在时才初始化地图
    if (places.value && places.value.length > 0) {
      await nextTick()
      initMap()
    }
  } catch {
    showToast('获取地点数据失败，稍后重试', 'error')
  } finally {
    uiStore.setLoading(false)
  }
})

onBeforeUnmount(() => {
  if (map) {
    map.remove()
    map = null
    markerMap.clear()
  }
})

/* ========== 地图初始化 ========== */
const initialCenter: [number, number] = [35, 105]
const initialZoom = 4

function resetMap() {
  if (!map) return
  map.setView(initialCenter, initialZoom)
}

function addResetButtonTogether() {
  if (!map) return

  const zoomControlContainer = map.zoomControl.getContainer()
  if (!zoomControlContainer) return

  const resetBtn = L.DomUtil.create('a', 'leaflet-control-zoom-reset', zoomControlContainer)
  resetBtn.title = '重置地图'
  resetBtn.style.width = '30px'
  resetBtn.style.height = '30px'
  resetBtn.style.lineHeight = '30px'
  resetBtn.style.textAlign = 'center'
  resetBtn.style.display = 'block'
  resetBtn.style.cursor = 'pointer'
  resetBtn.style.userSelect = 'none'
  resetBtn.style.fill = 'currentColor'
  resetBtn.style.color = 'currentColor'
  resetBtn.style.boxSizing = 'border-box'

  // 创建SVG图标
  const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 512 512"><path fill="currentColor" fill-rule="evenodd" d="M426.667 106.667v42.666L358 149.33c36.077 31.659 58.188 77.991 58.146 128.474c-.065 78.179-53.242 146.318-129.062 165.376s-154.896-15.838-191.92-84.695C58.141 289.63 72.637 204.42 130.347 151.68a85.33 85.33 0 0 0 33.28 30.507a124.59 124.59 0 0 0-46.294 97.066c1.05 69.942 58.051 126.088 128 126.08c64.072 1.056 118.71-46.195 126.906-109.749c6.124-47.483-15.135-92.74-52.236-118.947L320 256h-42.667V106.667zM202.667 64c23.564 0 42.666 19.103 42.666 42.667s-19.102 42.666-42.666 42.666S160 130.231 160 106.667S179.103 64 202.667 64" stroke-width="13" stroke="currentColor"/></svg>`
  // 将SVG包装在容器中以确保居中
  resetBtn.innerHTML = `<div style="display: flex; align-items: center; justify-content: center; height: 100%; width: 100%;">${svg}</div>`

  L.DomEvent.disableClickPropagation(resetBtn)

  L.DomEvent.on(resetBtn, 'click', e => {
    L.DomEvent.stopPropagation(e)
    resetMap()
  })
}

function initMap() {
  if (!mapRef.value || map) return

  map = L.map(mapRef.value, {
    zoomControl: true,
    attributionControl: false,
  }).setView([35, 105], 4)

  L.tileLayer('https://map.lw1314.site/{z}/{x}/{y}.png', {
    maxZoom: 18,
  }).addTo(map)

  // 添加标记
  places.value.forEach((place: Place) => {
    // 仅在存在图片时渲染 <img>
    const imgHtml = place.image?.file
      ? `<img src="${place.image.file.thumbnail || place.image.file.url}" class="w-full h-28 object-cover rounded" />`
      : ''

    const popupHtml = `
    <div class="w-56">
      ${imgHtml}
      <span class="font-bold text-lg"> ${place.name} </span>
      <br/>
      ${place.date}<br/>
      ${place.description}
    </div>
    `

    const marker = L.marker([place.latitude, place.longitude]).addTo(map!).bindPopup(popupHtml)

    markerMap.set(place.id, marker)
  })

  // 强制 Leaflet 重新计算尺寸
  requestAnimationFrame(() => {
    map?.invalidateSize()
  })
  addResetButtonTogether()
}

/* ========== 交互 ========== */
function flyToPlace(place: Place) {
  if (!map) return

  const marker = markerMap.get(place.id)
  if (!marker) return

  // 确保地图完全初始化并计算好尺寸
  map.invalidateSize(false)

  // 使用简单的flyTo而不是复杂的坐标计算，避免首次点击时的坐标计算错误
  map.flyTo([place.latitude, place.longitude], 12, {
    duration: 1.5,
    animate: true,
  })

  // 不再自动打开弹窗，让用户点击标记时才显示信息
}
</script>

<template>
  <MainLayout
    title="足迹地图"
    subtitle="印刻我们的同行足迹"
    :start-date="systemInfo?.site.startDate"
    :show-empty-state="places.length === 0"
  >
    <template #empty-state>
      <BaseIcon name="place" size="w-24" />
      <p class="text-xl font-medium mt-4">暂无地点数据</p>
      <p class="text-md mt-2">还没有添加任何地点</p>
    </template>

    <template #main-content>
      <!-- 有地点数据时显示地图和列表 -->
      <div class="flex-grow flex flex-col overflow-hidden">
        <!-- 地图区域 -->
        <div class="h-96 md:h-[500px] flex-shrink-0">
          <div ref="mapRef" class="w-full h-full"></div>
        </div>

        <!-- 列表区域 -->
        <div class="p-4 flex-grow overflow-y-auto">
          <h2 class="text-xl font-bold mb-3 font-[Ma_Shan_Zheng]">我们去过的地方</h2>

          <!-- 有地点数据时显示列表 -->
          <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
            <div
              v-for="place in places"
              :key="place.id"
              @click="flyToPlace(place)"
              class="flex items-center p-3 generic-card hover:bg-white/80 transition-colors cursor-pointer"
            >
              <div
                class="w-10 h-10 rounded-full bg-gradient-to-br from-red-400 to-pink-500 flex items-center justify-center text-white font-bold mr-3"
              >
                {{ place.name.substring(0, 1) }}
              </div>
              <div>
                <h3 class="font-bold">{{ place.name }}</h3>
                <p class="text-sm text-gray-600">{{ place.date }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </MainLayout>
</template>

<style scoped>
/* Leaflet 字体继承 */
:deep(.leaflet-container) {
  font-family: inherit;
}
</style>
