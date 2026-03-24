import { defineStore } from "pinia";

import { systemApi, type SystemInfo } from "@/services/system";

const DEFAULT_SYSTEM_INFO: SystemInfo = {
  site: {
    name: "Love Story",
    description: "我们的故事",
    startDate: "2024-01-01",
  },
  couple: {
    boy: {
      name: "他",
    },
    girl: {
      name: "她",
    },
  },
};

interface SystemState {
  systemInfo: SystemInfo | null;
  loadError: boolean;
}

export const useSystemStore = defineStore("system", {
  state: (): SystemState => ({
    systemInfo: null,
    loadError: false,
  }),

  getters: {
    getSystemInfo(state): SystemInfo {
      return state.systemInfo || DEFAULT_SYSTEM_INFO;
    },
  },

  actions: {
    async fetchSystemInfo(): Promise<SystemInfo> {
      if (this.systemInfo) {
        return this.systemInfo;
      }

      try {
        const response = await systemApi.getSystemInfo();
        if (response.data.code === 0) {
          this.systemInfo = response.data.data;
          this.loadError = false;
          return this.systemInfo;
        }
        this.loadError = true;
        return DEFAULT_SYSTEM_INFO;
      } catch (error) {
        console.error("获取系统信息失败:", error);
        this.loadError = true;
        return DEFAULT_SYSTEM_INFO;
      }
    },

    clearCache() {
      this.systemInfo = null;
      this.loadError = false;
    },
  },
});
