import { ref, onUnmounted } from "vue";
import { useAuthStore } from "@/stores/auth";
import { useSystemStore } from "@/stores/system";
import { useNotificationStore } from "@/stores/notification";

const isReady = ref(false);
let refreshIntervalId: ReturnType<typeof setInterval> | null = null;
let initialized = false;

export function useGlobalRefresh() {
  const authStore = useAuthStore();
  const systemStore = useSystemStore();
  const notificationStore = useNotificationStore();

  const initGlobalState = async (): Promise<boolean> => {
    if (initialized) {
      return isReady.value;
    }

    try {
      const isLoggedIn = await authStore.checkAuthStatus();
      if (!isLoggedIn) {
        isReady.value = false;
        initialized = true;
        return false;
      }

      await systemStore.fetchSystemInfo();

      if (authStore.isAuthenticated) {
        notificationStore.startPolling();
      }

      isReady.value = true;
      initialized = true;
      return true;
    } catch (error) {
      console.error("初始化全局状态失败:", error);
      isReady.value = false;
      initialized = true;
      return false;
    }
  };

  const startPeriodicRefresh = (intervalMs: number = 60000) => {
    if (refreshIntervalId) {
      clearInterval(refreshIntervalId);
    }

    refreshIntervalId = setInterval(async () => {
      if (!isReady.value) return;

      try {
        await systemStore.fetchSystemInfo();
      } catch (error) {
        console.error("定时刷新系统信息失败:", error);
      }
    }, intervalMs);
  };

  const stopPeriodicRefresh = () => {
    if (refreshIntervalId) {
      clearInterval(refreshIntervalId);
      refreshIntervalId = null;
    }
  };

  onUnmounted(() => {
    stopPeriodicRefresh();
  });

  return {
    isReady,
    initGlobalState,
    startPeriodicRefresh,
    stopPeriodicRefresh,
  };
}
