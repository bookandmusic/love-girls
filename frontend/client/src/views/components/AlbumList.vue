<template>
  <div
    ref="scrollContainer"
    class="flex-grow overflow-y-auto custom-scrollbar"
    @scroll="handleScroll"
  >
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 p-6">
      <div
        v-for="album in albums"
        :key="album.id"
        class="relative aspect-[4/3] rounded-[var(--fe-radius-card)] overflow-hidden shadow-lg border border-white/40 cursor-pointer tap-feedback ios-transition group"
        @click="onSelectAlbum(album)"
        @pointerdown="handlePointerDown(album, $event)"
        @pointerup="handlePointerUp"
        @pointerleave="handlePointerLeave"
        @pointercancel="handlePointerCancel"
      >
        <!-- 基础渐变背景（作为底层垫片，解决封面加载瞬间的留白问题） -->
        <div
          class="absolute inset-0 w-full h-full"
          :style="getPlaceholderStyle(album.id)"
        ></div>

        <!-- 封面图 (仅在数据存在时渲染) -->
        <img
          v-if="album.coverImage?.file"
          :src="album.coverImage.file.thumbnail || album.coverImage.file.url"
          :alt="album.name"
          class="absolute inset-0 w-full h-full object-cover transition-opacity duration-700"
          @load="
            $event.target &&
            (($event.target as HTMLImageElement).style.opacity = '1')
          "
          style="opacity: 0"
        />

        <!-- 底部毛玻璃信息栏 -->
        <div
          class="absolute inset-x-0 bottom-0 p-4 glass-thick border-t border-white/20 ios-transition rounded-b-[var(--fe-radius-card)]"
        >
          <div class="flex justify-between items-center mb-1">
            <h3
              class="text-lg font-bold text-[var(--fe-text-primary)] truncate"
            >
              {{ album.name }}
            </h3>
            <div
              class="flex items-center text-xs font-bold text-[var(--fe-primary)]"
            >
              <BaseIcon name="photo-heart" size="w-4 h-4" class="mr-1" />
              {{ album.photoCount }}
            </div>
          </div>
          <p class="text-xs text-[var(--fe-text-secondary)] line-clamp-1">
            {{ album.description || "这一刻，永恒。" }}
          </p>
        </div>
      </div>
    </div>

    <!-- 加载状态指示器 -->
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

    <!-- 占位 -->
    <div class="h-20 md:hidden"></div>
  </div>
</template>

<script setup lang="ts">
import { nextTick, ref, watch } from "vue";

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
}>();

const { onPointerDown, onPointerUp, onPointerLeave, onPointerCancel } =
  useLongPress({
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

// 滚动容器引用
const scrollContainer = ref<HTMLElement | null>(null);

// 预设的一组精致渐变色 (iOS 风格)
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
  if (!target || props.loading || !props.hasMore) return;

  const bottomDistance =
    target.scrollHeight - target.scrollTop - target.clientHeight;
  if (bottomDistance < 100) {
    emit("load-more");
  }
};

// 自动填充页面逻辑
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

// 监听数据变化，检查是否需要自动加载
watch(
  () => props.albums,
  () => {
    checkAndAutoLoadMore();
  },
);
</script>
