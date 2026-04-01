<template>
  <div class="h-full flex flex-col overflow-hidden">
    <div
      class="sticky top-0 z-20 glass-thick p-4 border-b border-white/20 flex items-center"
    >
      <button
        @click="onBack"
        class="flex items-center space-x-1 px-3 py-1.5 rounded-full bg-black/5 tap-feedback ios-transition text-[var(--fe-text-primary)] font-bold text-sm"
      >
        <BaseIcon name="left" size="w-5 h-5" color="var(--fe-text-primary)" />
        <span>返回</span>
      </button>
    </div>

    <van-pull-refresh
      v-model="isRefreshing"
      :disabled="!isAtTop"
      @refresh="handleRefresh"
    >
      <div
        ref="scrollContainer"
        class="overflow-y-auto p-4 md:p-6 custom-scrollbar"
        @scroll="handleScroll"
      >
        <vue-easy-lightbox
          :visible="visibleRef"
          :imgs="imgsRef"
          :index="indexRef"
          @hide="onHide"
          teleport="body"
        ></vue-easy-lightbox>

        <Waterfall
          v-if="photos.length > 0"
          :list="photos"
          row-key="id"
          img-selector="file.thumbnail"
          :breakpoints="breakpoints"
          :gutter="8"
          :has-around-gutter="true"
          :animation-cancel="true"
          :lazyload="true"
          :delay="100"
          background-color="transparent"
        >
          <template #default="{ item, url }">
            <div
              class="overflow-hidden rounded-xl border border-white/40 tap-feedback ios-transition shadow-sm group"
              @click="preview(item.file?.url || '')"
              @pointerdown="handlePointerDown(item, $event)"
              @pointermove="handlePointerMove"
              @pointerup="onPointerUp"
              @pointerleave="onPointerLeave"
              @pointercancel="onPointerCancel"
            >
              <div class="overflow-hidden">
                <LazyImg :url="url" :alt="item.alt || ''" />
              </div>
            </div>
          </template>
        </Waterfall>

        <div
          v-else-if="!loading"
          class="flex-1 flex flex-col items-center justify-center py-20"
        >
          <BaseIcon
            name="photo-heart"
            size="w-24"
            style="color: var(--fe-text-secondary)"
          />
          <p class="text-xl font-bold mt-4 text-[var(--fe-text-primary)]">
            暂无照片
          </p>
          <p class="text-md mt-2 text-[var(--fe-text-secondary)]">
            期待分享第一张照片
          </p>
        </div>

        <div v-if="loading || hasMore" class="py-10 flex justify-center">
          <div
            v-if="loading"
            class="flex items-center space-x-2 text-[var(--fe-text-secondary)]"
          >
            <div
              class="w-5 h-5 border-2 border-[var(--fe-primary)] border-t-transparent rounded-full animate-spin"
            ></div>
            <span class="text-xs font-bold uppercase tracking-widest"
              >正在加载...</span
            >
          </div>
          <div
            v-else-if="!hasMore && photos.length > 0"
            class="text-xs font-bold text-[var(--fe-text-secondary)] uppercase tracking-widest opacity-30"
          >
            已加载全部照片
          </div>
        </div>

        <div class="h-20 md:hidden"></div>
      </div>
    </van-pull-refresh>
  </div>
</template>

<script setup lang="ts">
import { nextTick, ref, watch } from "vue";
import VueEasyLightbox from "vue-easy-lightbox";
import { PullRefresh as VanPullRefresh } from "vant";
import { LazyImg, Waterfall } from "vue-waterfall-plugin-next";
import "vue-waterfall-plugin-next/dist/style.css";

import BaseIcon from "@/components/ui/BaseIcon.vue";
import { useLongPress } from "@/composables/useLongPress";
import type { Photo } from "@/services/albumApi";

interface Props {
  photos: Photo[];
  loading?: boolean;
  hasMore?: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "back"): void;
  (e: "load-more"): void;
  (e: "long-press", photo: Photo): void;
  (e: "refresh"): void;
}>();

const onBack = () => {
  emit("back");
};

const isRefreshing = ref(false);
const isAtTop = ref(true);

const handleRefresh = async () => {
  emit("refresh");
  isRefreshing.value = false;
};

const scrollContainer = ref<HTMLElement | null>(null);

const handleScroll = (e: Event) => {
  const target = e.target as HTMLElement;
  if (!target) return;

  isAtTop.value = target.scrollTop === 0;

  if (props.loading || !props.hasMore) return;

  const bottomDistance =
    target.scrollHeight - target.scrollTop - target.clientHeight;
  if (bottomDistance < 100) {
    emit("load-more");
  }
};

const checkAndAutoLoadMore = async () => {
  await nextTick();

  if (props.loading || !props.hasMore) return;

  const container = scrollContainer.value;
  if (container) {
    const isNotFilled = container.scrollHeight <= container.clientHeight + 10;
    if (isNotFilled) {
      emit("load-more");
    }
  }
};

const selectedPhoto = ref<Photo | null>(null);

const {
  onPointerDown,
  onPointerUp,
  onPointerLeave,
  onPointerCancel,
  onPointerMove,
} = useLongPress({
  duration: 500,
  onFinish: () => {
    if (selectedPhoto.value) {
      emit("long-press", selectedPhoto.value);
    }
  },
});

const handlePointerDown = (photo: Photo, event: PointerEvent) => {
  selectedPhoto.value = photo;
  onPointerDown(event);
};

const handlePointerMove = (event: PointerEvent) => {
  onPointerMove(event);
};

const visibleRef = ref(false);
const indexRef = ref(0);
const imgsRef = ref("");
const onShow = () => {
  visibleRef.value = true;
};
const onHide = () => (visibleRef.value = false);
const preview = (url: string) => {
  imgsRef.value = url;
  onShow();
};

const breakpoints = {
  1200: {
    rowPerView: 5,
  },
  800: {
    rowPerView: 4,
  },
  600: {
    rowPerView: 3,
  },
  400: {
    rowPerView: 2,
  },
};

watch(
  () => props.photos,
  () => {
    checkAndAutoLoadMore();
  },
);
</script>

<style scoped>
:deep(.lazy__img[lazy="loading"]) {
  padding: 5em 0;
  width: 48px;
}

:deep(.lazy__img[lazy="loaded"]) {
  width: 100%;
}

:deep(.lazy__img[lazy="error"]) {
  padding: 5em 0;
  width: 48px;
}
</style>
