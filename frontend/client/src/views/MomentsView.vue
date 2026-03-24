<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import VueEasyLightbox from "vue-easy-lightbox";

import BaseIcon from "@/components/ui/BaseIcon.vue";
import ActionSheet, {
  type ActionSheetAction,
} from "@/components/ui/ActionSheet.vue";
import FloatingAddButton from "@/components/ui/FloatingAddButton.vue";
import { useLongPress } from "@/composables/useLongPress";
import MainLayout from "@/layouts/MainLayout.vue";
import { type Moment, momentApi } from "@/services/momentApi";
import { useAuthStore } from "@/stores/auth";
import { useSystemStore } from "@/stores/system";
import { useUIStore } from "@/stores/ui";
import { useToast } from "@/utils/toastUtils";
import { useAutoFillPage } from "@/utils/useAutoFillPage";

import MomentEditDialog from "./components/dialogs/MomentEditDialog.vue";
import DeleteConfirmDialog from "./components/dialogs/DeleteConfirmDialog.vue";

const authStore = useAuthStore();
const uiStore = useUIStore();
const systemStore = useSystemStore();

const systemInfo = computed(() => systemStore.getSystemInfo);

const showToast = useToast();

const moments = ref<Moment[]>([]);
const currentPage = ref(1);
const totalPages = ref(0);
const pageSize = ref(8);
const loadingMore = ref(false);
const hasMore = computed(() => currentPage.value < totalPages.value);

const scrollContainer = ref<HTMLElement | null>(null);

const fetchMoments = async (page: number, append = false) => {
  if (loadingMore.value) return;
  loadingMore.value = true;

  try {
    const response = await momentApi.getMoments(page, pageSize.value);
    if (append) {
      moments.value = [...moments.value, ...response.data.moments];
    } else {
      moments.value = response.data.moments;
    }
    totalPages.value = response.data.totalPages;
    currentPage.value = page;
  } catch {
    showToast("获取动态列表失败", "error");
  } finally {
    loadingMore.value = false;
    uiStore.setLoading(false);
    checkAndAutoLoadMore();
  }
};

const loadNextPage = () => {
  if (hasMore.value && !loadingMore.value) {
    fetchMoments(currentPage.value + 1, true);
  }
};

const { checkAndAutoLoadMore } = useAutoFillPage(
  scrollContainer,
  hasMore,
  loadingMore,
  loadNextPage,
);

const handleScroll = (e: Event) => {
  const target = e.target as HTMLElement;
  if (!target || loadingMore.value || !hasMore.value) return;

  const bottomDistance =
    target.scrollHeight - target.scrollTop - target.clientHeight;
  if (bottomDistance < 100) {
    loadNextPage();
  }
};

const likeMoment = async (momentId: number) => {
  try {
    const response = await momentApi.likeMoment(momentId);
    if (response.code === 0) {
      const moment = moments.value.find((m) => m.id === momentId);
      if (moment) {
        moment.likes = response.data.likes;
      }
      showToast("点赞成功", "success");
    } else {
      showToast(response.msg || "点赞失败", "error");
    }
  } catch (error: unknown) {
    const axiosError = error as { response?: { data?: { message?: string } } };
    const message = axiosError.response?.data?.message || "点赞失败";
    showToast(message, "error");
  }
};

const visibleRef = ref(false);
const indexRef = ref(0);
const imgsRef = ref("");
const onShow = () => (visibleRef.value = true);
const onHide = () => (visibleRef.value = false);
const viewImage = (imageUrl: string) => {
  imgsRef.value = imageUrl;
  onShow();
};

const showActionSheet = ref(false);
const selectedMoment = ref<Moment | null>(null);

const { onPointerDown, onPointerUp, onPointerLeave, onPointerCancel } =
  useLongPress({
    duration: 500,
    onFinish: () => {
      if (selectedMoment.value) {
        showActionSheet.value = true;
      }
    },
  });

const handlePointerDown = (moment: Moment, event: PointerEvent) => {
  selectedMoment.value = moment;
  onPointerDown(event);
};

const actionSheetActions = computed<ActionSheetAction[]>(() => [
  {
    label: selectedMoment.value?.isPublic ? "设为私密" : "设为公开",
    handler: () => handleTogglePublic(selectedMoment.value),
  },
  {
    label: "编辑",
    handler: () => openEditDialog(selectedMoment.value),
  },
  {
    label: "删除",
    destructive: true,
    handler: () => openDeleteDialog(selectedMoment.value),
  },
]);

const showEditDialog = ref(false);
const editingMoment = ref<Moment | null>(null);
const savingMoment = ref(false);

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

const DEFAULT_MOMENT: Moment = {
  id: 0,
  content: "",
  isPublic: true,
  images: [],
  likes: 0,
  author: { name: "系统用户" },
  createdAt: "",
};

const openAddDialog = () => {
  editingMoment.value = { ...DEFAULT_MOMENT, createdAt: formatLocalDateTime() };
  showEditDialog.value = true;
};

const openEditDialog = (moment: Moment | null) => {
  editingMoment.value = moment;
  showEditDialog.value = true;
};

const handleSaveMoment = async (moment: Moment) => {
  savingMoment.value = true;
  try {
    const imageIds = moment.images?.map((img) => img.id) || [];

    if (moment.id) {
      await momentApi.updateMoment(moment.id, {
        content: moment.content,
        isPublic: moment.isPublic,
        imageIds: imageIds,
        createdAt: moment.createdAt,
      });
      showToast("动态更新成功", "success");
    } else {
      await momentApi.createMoment({
        content: moment.content,
        isPublic: moment.isPublic,
        imageIds: imageIds,
        likes: 0,
        author: { name: "系统用户" },
        createdAt: moment.createdAt,
        userId: authStore.userInfo?.userId || 1,
      });
      showToast("动态发布成功", "success");
    }
    showEditDialog.value = false;
    await fetchMoments(1);
  } catch {
    showToast("操作失败", "error");
  } finally {
    savingMoment.value = false;
  }
};

const handleTogglePublic = async (moment: Moment | null) => {
  if (!moment) return;
  try {
    await momentApi.updateMomentPublic(moment.id, {
      isPublic: !moment.isPublic,
    });
    moment.isPublic = !moment.isPublic;
    showToast(moment.isPublic ? "已设为公开" : "已设为私密", "success");
  } catch {
    showToast("操作失败", "error");
  }
};

const showDeleteDialog = ref(false);
const deletingMoment = ref<Moment | null>(null);
const deleting = ref(false);

const openDeleteDialog = (moment: Moment | null) => {
  deletingMoment.value = moment;
  showDeleteDialog.value = true;
};

const handleDeleteMoment = async () => {
  if (!deletingMoment.value) return;
  deleting.value = true;
  try {
    await momentApi.deleteMoment(deletingMoment.value.id);
    showToast("动态删除成功", "success");
    showDeleteDialog.value = false;
    await fetchMoments(1);
  } catch {
    showToast("删除失败", "error");
  } finally {
    deleting.value = false;
  }
};

onMounted(async () => {
  uiStore.setLoading(true);
  await systemStore.fetchSystemInfo();
  await fetchMoments(1);
});
</script>

<template>
  <div class="h-full w-full">
    <MainLayout
      title="时光动态"
      subtitle="记录我们的点点滴滴"
      :start-date="systemInfo?.site.startDate"
      :show-empty-state="moments.length === 0 && !loadingMore"
    >
      <template #empty-state>
        <BaseIcon
          name="moment"
          size="w-24"
          style="color: var(--fe-text-secondary)"
        />
        <p class="font-bold text-xl mt-4 text-[var(--fe-text-primary)]">
          暂无动态
        </p>
        <p class="text-md mt-2 text-[var(--fe-text-secondary)]">
          期待分享第一条动态
        </p>
      </template>

      <template #main-content>
        <vue-easy-lightbox
          :visible="visibleRef"
          :imgs="imgsRef"
          :index="indexRef"
          @hide="onHide"
          teleport="body"
        ></vue-easy-lightbox>

        <div class="flex flex-col h-full glass-regular">
          <div
            ref="scrollContainer"
            class="flex-grow overflow-y-auto p-4 md:p-8 space-y-0 custom-scrollbar"
            @scroll="handleScroll"
          >
            <div
              v-for="moment in moments"
              :key="moment.id"
              class="py-6 border-b border-black/5 last:border-0 ios-transition"
              @pointerdown="handlePointerDown(moment, $event)"
              @pointerup="onPointerUp"
              @pointerleave="onPointerLeave"
              @pointercancel="onPointerCancel"
            >
              <div class="flex items-start">
                <div
                  class="w-12 h-12 rounded-lg overflow-hidden bg-white/50 border border-white/60 flex items-center justify-center text-[var(--fe-primary)] font-bold mr-4 flex-shrink-0"
                >
                  <img
                    v-if="
                      moment.author.avatar?.thumbnail ||
                      moment.author.avatar?.url
                    "
                    :src="
                      moment.author.avatar?.thumbnail ||
                      moment.author.avatar?.url
                    "
                    :alt="moment.author.name"
                    class="w-full h-full object-cover"
                  />
                  <span v-else>{{ moment.author.name.substring(0, 1) }}</span>
                </div>

                <div class="flex-grow min-w-0">
                  <div class="mb-1">
                    <h3 class="font-bold text-[#576b95] text-base truncate">
                      {{ moment.author.name }}
                    </h3>
                  </div>

                  <p
                    class="text-[var(--fe-text-primary)] leading-relaxed mb-3 text-sm md:text-base"
                  >
                    {{ moment.content }}
                  </p>

                  <div
                    v-if="moment.images && moment.images.length > 0"
                    class="mb-3"
                  >
                    <div
                      class="grid gap-1.5"
                      :class="{
                        'grid-cols-1 w-max': moment.images.length === 1,
                        'grid-cols-2 w-full md:max-w-[280px]':
                          moment.images.length === 2 ||
                          moment.images.length === 4,
                        'grid-cols-3 w-full max-w-[320px] md:max-w-[420px]':
                          moment.images.length === 3 ||
                          moment.images.length >= 5,
                      }"
                    >
                      <div
                        v-for="(image, index) in moment.images"
                        :key="index"
                        class="overflow-hidden cursor-pointer tap-feedback ios-transition"
                        :class="[
                          moment.images.length === 1
                            ? 'rounded-lg max-w-[240px] max-h-[320px]'
                            : 'w-full aspect-square rounded-md',
                        ]"
                        @click.stop="viewImage(image.file?.url || '')"
                      >
                        <img
                          :src="image.file?.thumbnail || image.file?.url"
                          alt="动态图片"
                          class="w-full h-full object-cover"
                          loading="lazy"
                        />
                      </div>
                    </div>
                  </div>

                  <div class="flex justify-between items-center mt-3">
                    <span
                      class="text-xs font-medium text-[var(--fe-text-secondary)] opacity-60"
                    >
                      {{ moment.createdAt }}
                    </span>
                    <button
                      @click.stop="likeMoment(moment.id)"
                      class="flex items-center space-x-1.5 px-2.5 py-1 rounded-md bg-black/5 tap-feedback ios-transition"
                    >
                      <BaseIcon
                        name="like"
                        size="w-3.5 h-3.5"
                        color="var(--fe-primary)"
                      />
                      <span
                        class="text-xs font-bold text-[var(--fe-text-primary)]"
                      >
                        {{ moment.likes }}
                      </span>
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <div
              v-if="loadingMore || hasMore"
              class="py-10 flex justify-center"
            >
              <div
                v-if="loadingMore"
                class="flex items-center space-x-2 text-[var(--fe-text-secondary)]"
              >
                <div
                  class="w-5 h-5 border-2 border-[var(--fe-primary)] border-t-transparent rounded-full animate-spin"
                ></div>
                <span class="text-xs font-bold uppercase tracking-widest"
                  >正在加载更多...</span
                >
              </div>
              <div
                v-else-if="!hasMore && moments.length > 0"
                class="text-xs font-bold text-[var(--fe-text-secondary)] uppercase tracking-widest opacity-30"
              >
                没有更多动态了
              </div>
            </div>

            <div class="h-20 md:hidden"></div>
          </div>
        </div>
      </template>
    </MainLayout>

    <FloatingAddButton :loading="savingMoment" @click="openAddDialog" />

    <ActionSheet
      v-model="showActionSheet"
      title="动态操作"
      :actions="actionSheetActions"
    />

    <MomentEditDialog
      v-model:open="showEditDialog"
      :moment="editingMoment ?? undefined"
      :loading="savingMoment"
      @confirm="handleSaveMoment"
    />

    <DeleteConfirmDialog
      v-model:open="showDeleteDialog"
      :loading="deleting"
      title="删除动态"
      message="确定要删除这条动态吗？删除后无法恢复。"
      @confirm="handleDeleteMoment"
    />
  </div>
</template>
