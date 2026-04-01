<template>
  <van-pull-refresh
    v-model="isRefreshing"
    :disabled="!isAtTop"
    @refresh="handleRefresh"
  >
    <div
      ref="scrollContainer"
      class="overflow-y-auto custom-scrollbar"
      @scroll="handleScroll"
    >
      <Waterfall
        :list="albums"
        row-key="id"
        img-selector="coverImage.file.thumbnail"
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
            class="rounded-[var(--fe-radius-card)] overflow-hidden shadow-lg border border-white/40 cursor-pointer tap-feedback ios-transition group"
            :style="getPlaceholderStyle(item.id)"
            @click="onSelectAlbum(item)"
            @pointerdown="handlePointerDown(item, $event)"
            @pointermove="handlePointerMove"
            @pointerup="handlePointerUp"
            @pointerleave="handlePointerLeave"
            @pointercancel="handlePointerCancel"
          >
            <div class="overflow-hidden">
              <LazyImg :url="url" :alt="item.name" />
            </div>
            <div class="p-4 border-t border-white/20">
              <div class="flex justify-between items-center mb-1">
                <h3
                  class="text-lg font-bold text-[var(--fe-text-primary)] truncate"
                >
                  {{ item.name }}
                </h3>
                <div
                  class="flex items-center text-xs font-bold text-[var(--fe-primary)]"
                >
                  <BaseIcon name="photo-heart" size="w-4 h-4" class="mr-1" />
                  {{ item.photoCount }}
                </div>
              </div>
              <p class="text-xs text-[var(--fe-text-secondary)] line-clamp-1">
                {{ item.description || "这一刻，永恒。" }}
              </p>
            </div>
          </div>
        </template>
      </Waterfall>

      <div v-if="loading || hasMore" class="py-10 flex justify-center">
        <div
          v-if="loading"
          class="flex items-center space-x-2 text-[var(--fe-text-secondary)]"
        >
          <div
            class="w-5 h-5 border-2 border-[var(--fe-primary)] border-t-transparent rounded-full animate-spin"
          ></div>
          <span class="text-xs font-bold uppercase tracking-widest"
            >加载中...</span
          >
        </div>
        <div
          v-else-if="!hasMore && albums.length > 0"
          class="text-xs font-bold text-[var(--fe-text-secondary)] uppercase tracking-widest opacity-30"
        >
          已显示全部相册
        </div>
      </div>

      <div class="h-20 md:hidden"></div>
    </div>
  </van-pull-refresh>
</template>

<script setup lang="ts">
import { nextTick, ref, watch } from "vue";
import { PullRefresh as VanPullRefresh } from "vant";
import { LazyImg, Waterfall } from "vue-waterfall-plugin-next";
import "vue-waterfall-plugin-next/dist/style.css";

import BaseIcon from "@/components/ui/BaseIcon.vue";
import { useLongPress } from "@/composables/useLongPress";
import type { Album } from "@/services/albumApi";

interface Props {
  albums: Album[];
  loading?: boolean;
  hasMore?: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "select-album", album: Album): void;
  (e: "load-more"): void;
  (e: "long-press", album: Album): void;
  (e: "refresh"): void;
}>();

const isRefreshing = ref(false);
const isAtTop = ref(true);

const handleRefresh = async () => {
  emit("refresh");
  isRefreshing.value = false;
};

const {
  onPointerDown,
  onPointerUp,
  onPointerLeave,
  onPointerCancel,
  onPointerMove,
} = useLongPress({
  duration: 500,
  onFinish: () => {
    if (longPressAlbum.value) {
      isLongPressTriggered.value = true;
      emit("long-press", longPressAlbum.value);
    }
  },
});

const longPressAlbum = ref<Album | null>(null);
const isLongPressTriggered = ref(false);

const handlePointerDown = (album: Album, event: PointerEvent) => {
  longPressAlbum.value = album;
  isLongPressTriggered.value = false;
  onPointerDown(event);
};

const handlePointerMove = (event: PointerEvent) => {
  onPointerMove(event);
};

const handlePointerUp = () => {
  onPointerUp();
};

const handlePointerLeave = () => {
  longPressAlbum.value = null;
  onPointerLeave();
};

const handlePointerCancel = () => {
  longPressAlbum.value = null;
  onPointerCancel();
};

const onSelectAlbum = (album: Album) => {
  if (!isLongPressTriggered.value) {
    emit("select-album", album);
  }
  longPressAlbum.value = null;
  isLongPressTriggered.value = false;
};

const scrollContainer = ref<HTMLElement | null>(null);

const gradients = [
  "linear-gradient(135deg, #f0ada0 0%, #f8c9c0 100%)",
  "linear-gradient(135deg, #a1c4fd 0%, #c2e9fb 100%)",
  "linear-gradient(135deg, #ff9a9e 0%, #fecfef 100%)",
  "linear-gradient(135deg, #84fab0 0%, #8fd3f4 100%)",
  "linear-gradient(135deg, #cfd9df 0%, #e2ebf0 100%)",
  "linear-gradient(135deg, #a6c0fe 0%, #f68084 100%)",
];

const getPlaceholderStyle = (id: number) => {
  const index = id % gradients.length;
  return {
    background: gradients[index],
  };
};

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

const breakpoints = {
  1200: {
    rowPerView: 4,
  },
  800: {
    rowPerView: 3,
  },
  500: {
    rowPerView: 2,
  },
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

watch(
  () => props.albums,
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
