import { defineStore } from "pinia";
import { ref, computed } from "vue";

import type { Notification } from "@/services/notificationApi";
import { notificationApi } from "@/services/notificationApi";

export const useNotificationStore = defineStore("notification", () => {
  const notifications = ref<Notification[]>([]);
  const unreadCount = ref(0);
  const loading = ref(false);
  const isPolling = ref(false);

  const hasUnread = computed(() => unreadCount.value > 0);

  let pollInterval: ReturnType<typeof setInterval> | null = null;

  const fetchUnreadCount = async () => {
    try {
      const response = await notificationApi.getUnreadCount();
      if (response.code === 0) {
        unreadCount.value = response.data.count;
      }
    } catch (error) {
      console.error("获取未读通知数量失败:", error);
    }
  };

  const fetchNotifications = async (reset = false) => {
    if (loading.value) return;
    loading.value = true;

    try {
      const response = await notificationApi.getUnreadNotifications();
      if (response.code === 0) {
        if (reset) {
          notifications.value = response.data.notifications || [];
        } else {
          notifications.value = [
            ...notifications.value,
            ...response.data.notifications,
          ];
        }
        unreadCount.value = response.data.total;
      }
    } catch (error) {
      console.error("获取通知失败:", error);
    } finally {
      loading.value = false;
    }
  };

  const markAsRead = async (notificationId: number) => {
    try {
      const response = await notificationApi.markAsRead(notificationId);
      if (response.code === 0) {
        const notification = notifications.value.find(
          (n) => n.id === notificationId,
        );
        if (notification) {
          notification.isRead = true;
        }
        notifications.value = notifications.value.filter((n) => !n.isRead);
        unreadCount.value = Math.max(0, unreadCount.value - 1);
      }
    } catch (error) {
      console.error("标记已读失败:", error);
    }
  };

  const markAllAsRead = async () => {
    try {
      const response = await notificationApi.markAllAsRead();
      if (response.code === 0) {
        notifications.value = [];
        unreadCount.value = 0;
      }
    } catch (error) {
      console.error("全部标记已读失败:", error);
    }
  };

  const startPolling = (intervalMs = 5000) => {
    if (isPolling.value) return;
    isPolling.value = true;
    fetchUnreadCount();
    pollInterval = setInterval(fetchUnreadCount, intervalMs);
  };

  const stopPolling = () => {
    if (pollInterval) {
      clearInterval(pollInterval);
      pollInterval = null;
    }
    isPolling.value = false;
  };

  return {
    notifications,
    unreadCount,
    loading,
    isPolling,
    hasUnread,
    fetchUnreadCount,
    fetchNotifications,
    markAsRead,
    markAllAsRead,
    startPolling,
    stopPolling,
  };
});
