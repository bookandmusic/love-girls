<script setup lang="ts">
import { onMounted, onUnmounted, ref } from "vue";
import { useRouter } from "vue-router";
import { PullRefresh as VanPullRefresh } from "vant";

import BaseIcon from "@/components/ui/BaseIcon.vue";
import { useNotificationStore } from "@/stores/notification";

const router = useRouter();
const notificationStore = useNotificationStore();
const isRefreshing = ref(false);

const handleBack = () => {
  router.back();
};

const handleNotificationClick = async (notification: {
  id: number;
  momentId: number;
}) => {
  await notificationStore.markAsRead(notification.id);
  router.push(`/moments?highlight=${notification.momentId}`);
};

const handleMarkAllRead = async () => {
  await notificationStore.markAllAsRead();
};

const handleRefresh = async () => {
  await notificationStore.fetchNotifications(true);
  isRefreshing.value = false;
};

onMounted(() => {
  notificationStore.fetchNotifications(true);
  notificationStore.startPolling();
});

onUnmounted(() => {
  notificationStore.stopPolling();
});
</script>

<template>
  <div
    class="notifications-view fixed inset-0 z-[200] bg-[var(--fe-bg-primary)]"
  >
    <div class="flex flex-col h-full">
      <div
        class="flex-shrink-0 h-14 flex items-center px-4 border-b border-black/5 bg-white/80 backdrop-blur-sm"
      >
        <button
          @click="handleBack"
          class="p-2 -ml-2 rounded-full hover:bg-black/5 transition-colors"
        >
          <svg
            class="w-6 h-6 text-[var(--fe-text-primary)]"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M15 19l-7-7 7-7"
            />
          </svg>
        </button>
        <h1
          class="flex-1 text-center text-lg font-bold text-[var(--fe-text-primary)]"
        >
          通知
        </h1>
        <button
          v-if="notificationStore.hasUnread"
          @click="handleMarkAllRead"
          class="text-sm text-[var(--fe-primary-dark)] font-medium hover:opacity-80 transition-opacity"
        >
          全部已读
        </button>
        <div v-else class="w-14"></div>
      </div>

      <div class="flex-1 overflow-y-auto">
        <van-pull-refresh v-model="isRefreshing" @refresh="handleRefresh">
          <div
            v-if="
              notificationStore.notifications.length === 0 &&
              !notificationStore.loading
            "
            class="py-20 flex flex-col items-center justify-center"
          >
            <BaseIcon
              name="bell"
              size="w-16"
              style="color: var(--fe-text-secondary)"
            />
            <p class="text-sm text-[var(--fe-text-secondary)] mt-4">
              暂无新通知
            </p>
          </div>

          <div v-else class="divide-y divide-black/5">
            <div
              v-for="notification in notificationStore.notifications"
              :key="notification.id"
              @click="handleNotificationClick(notification)"
              class="p-4 cursor-pointer tap-feedback ios-transition hover:bg-black/5 active:bg-black/10"
            >
              <div class="flex items-start gap-3">
                <div
                  class="w-10 h-10 rounded-full overflow-hidden bg-gray-200 flex-shrink-0 flex items-center justify-center"
                >
                  <img
                    v-if="
                      notification.sender.avatar?.thumbnail ||
                      notification.sender.avatar?.url
                    "
                    :src="
                      notification.sender.avatar.thumbnail ||
                      notification.sender.avatar.url
                    "
                    :alt="notification.sender.name"
                    class="w-full h-full object-cover"
                  />
                  <span v-else class="text-sm font-bold text-gray-500">
                    {{ notification.sender.name.substring(0, 1) }}
                  </span>
                </div>

                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2 mb-1">
                    <span class="font-medium text-[#576b95]">
                      {{ notification.sender.name }}
                    </span>
                    <span class="text-xs text-gray-400">
                      {{ notification.createdAt }}
                    </span>
                  </div>

                  <p class="text-sm text-gray-800">
                    <span v-if="notification.type === 'comment'"
                      >评论了动态</span
                    >
                    <span v-else>回复了评论</span>
                    <span class="text-gray-500 mx-1">:</span>
                    <span class="text-gray-600">{{
                      notification.content
                    }}</span>
                  </p>
                </div>

                <div
                  v-if="!notification.isRead"
                  class="w-2 h-2 rounded-full bg-red-500 flex-shrink-0 mt-2"
                ></div>
              </div>
            </div>
          </div>

          <div
            v-if="notificationStore.loading"
            class="py-8 flex justify-center"
          >
            <div
              class="w-6 h-6 border-2 border-[var(--fe-primary)] border-t-transparent rounded-full animate-spin"
            ></div>
          </div>
        </van-pull-refresh>
      </div>
    </div>
  </div>
</template>
