<script setup lang="ts">
import { computed, onMounted, ref } from "vue";

import BaseIcon from "@/components/ui/BaseIcon.vue";
import ActionSheet, {
  type ActionSheetAction,
} from "@/components/ui/ActionSheet.vue";
import FloatingAddButton from "@/components/ui/FloatingAddButton.vue";
import MainLayout from "@/layouts/MainLayout.vue";
import { type Album, albumApi, type Photo } from "@/services/albumApi";
import { useSystemStore } from "@/stores/system";
import { useUIStore } from "@/stores/ui";
import { useToast } from "@/utils/toastUtils";

import AlbumEditDialog from "./components/dialogs/AlbumEditDialog.vue";
import DeleteConfirmDialog from "./components/dialogs/DeleteConfirmDialog.vue";
import AlbumList from "./components/AlbumList.vue";
import PhotoList from "./components/PhotoList.vue";

const uiStore = useUIStore();
const systemStore = useSystemStore();

const systemInfo = computed(() => systemStore.getSystemInfo);

const showToast = useToast();

const albums = ref<Album[]>([]);
const currentAlbum = ref<Album | null>(null);
const currentPage = ref(1);
const totalPages = ref(0);
const pageSize = ref(9);
const loadingMoreAlbums = ref(false);
const hasMoreAlbums = computed(() => currentPage.value < totalPages.value);

const photos = ref<Photo[]>([]);
const currentAlbumId = ref<number | null>(null);
const currentPhotoPage = ref(1);
const totalPhotoPages = ref(0);
const photoPageSize = ref(12);
const loadingMorePhotos = ref(false);
const hasMorePhotos = computed(
  () => currentPhotoPage.value < totalPhotoPages.value,
);

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

const fetchPhotos = async (albumId: number, page: number, append = false) => {
  if (loadingMorePhotos.value) return;
  loadingMorePhotos.value = true;

  try {
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

const handleLoadMoreAlbums = () => {
  if (hasMoreAlbums.value) {
    fetchAlbums(currentPage.value + 1, true);
  }
};

const handleLoadMorePhotos = () => {
  if (currentAlbumId.value && hasMorePhotos.value) {
    fetchPhotos(currentAlbumId.value, currentPhotoPage.value + 1, true);
  }
};

const backToAlbums = () => {
  currentAlbumId.value = null;
  currentAlbum.value = null;
  photos.value = [];
  currentPhotoPage.value = 1;
  totalPhotoPages.value = 0;
};

const handleSelectAlbum = (album: Album) => {
  uiStore.setLoading(true);
  fetchPhotos(album.id, 1);
};

const handleBack = () => {
  backToAlbums();
};

const showActionSheet = ref(false);
const selectedAlbum = ref<Album | null>(null);
const actionSheetActions = computed<ActionSheetAction[]>(() => [
  {
    label: "编辑",
    handler: () => openEditDialog(selectedAlbum.value),
  },
  {
    label: "删除",
    destructive: true,
    handler: () => openDeleteDialog(selectedAlbum.value),
  },
]);

const handleLongPressAlbum = (album: Album) => {
  selectedAlbum.value = album;
  showActionSheet.value = true;
};

const showEditDialog = ref(false);
const editingAlbum = ref<Album | null>(null);
const savingAlbum = ref(false);

const openAddDialog = () => {
  editingAlbum.value = null;
  showEditDialog.value = true;
};

const openEditDialog = (album: Album | null) => {
  editingAlbum.value = album;
  showEditDialog.value = true;
};

const handleSaveAlbum = async (album: Album) => {
  savingAlbum.value = true;
  try {
    if (album.id) {
      await albumApi.updateAlbum(album.id, album);
      showToast("相册更新成功", "success");
    } else {
      await albumApi.createAlbum(album);
      showToast("相册创建成功", "success");
    }
    showEditDialog.value = false;
    await fetchAlbums(1);
  } catch {
    showToast("操作失败", "error");
  } finally {
    savingAlbum.value = false;
  }
};

const showDeleteDialog = ref(false);
const deletingAlbum = ref<Album | null>(null);
const deleting = ref(false);

const openDeleteDialog = (album: Album | null) => {
  deletingAlbum.value = album;
  showDeleteDialog.value = true;
};

const handleDeleteAlbum = async () => {
  if (!deletingAlbum.value) return;
  deleting.value = true;
  try {
    await albumApi.deleteAlbum(deletingAlbum.value.id);
    showToast("相册删除成功", "success");
    showDeleteDialog.value = false;
    await fetchAlbums(1);
  } catch {
    showToast("删除失败", "error");
  } finally {
    deleting.value = false;
  }
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
      <div
        class="flex flex-col items-center justify-center py-20 text-[var(--fe-text-secondary)]"
      >
        <BaseIcon name="camera" size="w-24" />
        <p class="text-xl font-bold mt-4 text-[var(--fe-text-primary)]">
          暂无相册
        </p>
      </div>
    </template>

    <template #main-content>
      <div v-if="!currentAlbumId" class="flex flex-col h-full">
        <AlbumList
          :albums="albums"
          :loading="loadingMoreAlbums"
          :has-more="hasMoreAlbums"
          @select-album="handleSelectAlbum"
          @load-more="handleLoadMoreAlbums"
          @long-press="handleLongPressAlbum"
        />
      </div>

      <div v-else class="flex flex-col h-full">
        <PhotoList
          :photos="photos"
          :loading="loadingMorePhotos"
          :has-more="hasMorePhotos"
          @back="handleBack"
          @load-more="handleLoadMorePhotos"
        />
      </div>

      <FloatingAddButton
        v-if="!currentAlbumId"
        :loading="savingAlbum"
        @click="openAddDialog"
      />
    </template>
  </MainLayout>

  <ActionSheet
    v-model="showActionSheet"
    title="相册操作"
    :actions="actionSheetActions"
  />

  <AlbumEditDialog
    v-model:open="showEditDialog"
    :album="editingAlbum"
    :loading="savingAlbum"
    @confirm="handleSaveAlbum"
  />

  <DeleteConfirmDialog
    v-model:open="showDeleteDialog"
    :loading="deleting"
    title="删除相册"
    :message="`确定要删除「${deletingAlbum?.name || ''}」吗？删除后无法恢复。`"
    @confirm="handleDeleteAlbum"
  />
</template>
