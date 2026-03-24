<script setup lang="ts">
import { computed, onMounted, ref } from "vue";

import BaseIcon from "@/components/ui/BaseIcon.vue";
import MainLayout from "@/layouts/MainLayout.vue";
import { type Album, albumApi, type Photo } from "@/services/albumApi";
import { useSystemStore } from "@/stores/system";
import { useUIStore } from "@/stores/ui";
import { useToast } from "@/utils/toastUtils";

import AlbumList from "./components/AlbumList.vue";
import PhotoList from "./components/PhotoList.vue";

const uiStore = useUIStore();
const systemStore = useSystemStore();

// 获取系统信息
const systemInfo = computed(() => systemStore.getSystemInfo);

const showToast = useToast();

// 相册相关状态
const albums = ref<Album[]>([]);
const currentAlbum = ref<Album | null>(null);
const currentPage = ref(1);
const totalPages = ref(0);
const pageSize = ref(9); // 每页显示数量
const loadingMoreAlbums = ref(false);
const hasMoreAlbums = computed(() => currentPage.value < totalPages.value);

// 照片相关状态
const photos = ref<Photo[]>([]);
const currentAlbumId = ref<number | null>(null);
const currentPhotoPage = ref(1);
const totalPhotoPages = ref(0);
const photoPageSize = ref(12); // 每页显示数量
const loadingMorePhotos = ref(false);
const hasMorePhotos = computed(
  () => currentPhotoPage.value < totalPhotoPages.value,
);

// 计算当前显示的标题和副标题
const pageTitle = computed(() => {
  if (currentAlbum.value) {
    return currentAlbum.value.name;
  }
  return "记忆相册";
});

const pageSubtitle = computed(() => {
  if (currentAlbum.value) {
    return currentAlbum.value.description;
  }
  return "珍藏我们的美好瞬间";
});

// 获取相册列表
const fetchAlbums = async (page: number, append = false) => {
  if (loadingMoreAlbums.value) return;
  loadingMoreAlbums.value = true;

  try {
    const response = await albumApi.getAlbums(page, pageSize.value);
    if (append) {
      albums.value = [...albums.value, ...response.data.albums];
    } else {
      albums.value = response.data.albums;
    }
    totalPages.value = response.data.totalPages;
    currentPage.value = page;
  } catch {
    showToast("获取相册列表失败", "error");
  } finally {
    loadingMoreAlbums.value = false;
    uiStore.setLoading(false);
  }
};

// 获取相册中的照片
const fetchPhotos = async (albumId: number, page: number, append = false) => {
  if (loadingMorePhotos.value) return;
  loadingMorePhotos.value = true;

  try {
    // 首次进入相册详情时设置当前相册信息
    if (!append) {
      const album = albums.value.find((a) => a.id === albumId);
      if (album) {
        currentAlbum.value = album;
      }
    }

    const response = await albumApi.getPhotos(
      albumId,
      page,
      photoPageSize.value,
    );
    if (append) {
      photos.value = [...photos.value, ...response.data.photos];
    } else {
      photos.value = response.data.photos;
    }
    totalPhotoPages.value = response.data.totalPages;
    currentPhotoPage.value = page;
    currentAlbumId.value = albumId;
  } catch {
    showToast("获取照片列表失败", "error");
  } finally {
    loadingMorePhotos.value = false;
    uiStore.setLoading(false);
  }
};

// 加载更多相册
const handleLoadMoreAlbums = () => {
  if (hasMoreAlbums.value) {
    fetchAlbums(currentPage.value + 1, true);
  }
};

// 加载更多照片
const handleLoadMorePhotos = () => {
  if (currentAlbumId.value && hasMorePhotos.value) {
    fetchPhotos(currentAlbumId.value, currentPhotoPage.value + 1, true);
  }
};

// 返回相册列表
const backToAlbums = () => {
  currentAlbumId.value = null;
  currentAlbum.value = null;
  photos.value = [];
  currentPhotoPage.value = 1;
  totalPhotoPages.value = 0;
};

// 处理相册选择
const handleSelectAlbum = (album: Album) => {
  uiStore.setLoading(true);
  fetchPhotos(album.id, 1);
};

// 处理返回相册列表
const handleBack = () => {
  backToAlbums();
};

onMounted(async () => {
  uiStore.setLoading(true);
  await systemStore.fetchSystemInfo();
  await fetchAlbums(1);
});
</script>

<template>
  <MainLayout
    :title="pageTitle"
    :subtitle="pageSubtitle"
    :start-date="systemInfo?.site.startDate"
    :show-empty-state="
      !currentAlbumId && albums.length === 0 && !loadingMoreAlbums
    "
  >
    <template #empty-state>
      <BaseIcon name="camera" size="w-24" />
      <p class="text-xl font-bold mt-4 text-[var(--fe-text-primary)]">
        暂无相册
      </p>
      <p class="text-sm mt-2 text-[var(--fe-text-secondary)]">
        还没有添加任何相册
      </p>
    </template>

    <template #main-content>
      <!-- 相册列表视图 -->
      <div
        v-if="!currentAlbumId"
        class="flex flex-col h-full bg-[var(--fe-bg-gray)]/30"
      >
        <AlbumList
          :albums="albums"
          :loading="loadingMoreAlbums"
          :has-more="hasMoreAlbums"
          @select-album="handleSelectAlbum"
          @load-more="handleLoadMoreAlbums"
        />
      </div>

      <!-- 照片列表视图 -->
      <div v-else class="flex flex-col h-full bg-[var(--fe-bg-gray)]/30">
        <PhotoList
          :photos="photos"
          :loading="loadingMorePhotos"
          :has-more="hasMorePhotos"
          @back="handleBack"
          @load-more="handleLoadMorePhotos"
        />
      </div>
    </template>
  </MainLayout>
</template>
