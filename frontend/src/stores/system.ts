// stores/system.ts
import { defineStore } from 'pinia'

import { systemApi, type SystemInfo } from '@/services/system'

// 定义系统信息类型

interface SystemState {
  systemInfo: SystemInfo | null
  initialized: boolean | null // 修改为nullable，与路由中的实现保持一致
}

export const useSystemStore = defineStore('system', {
  state: (): SystemState => ({
    systemInfo: null,
    initialized: null, // 初始值改为null
  }),

  getters: {
    getSystemInfo(state): SystemInfo | null {
      return state.systemInfo
    },
    isInitialized(state): boolean {
      return state.initialized === true
    },
    // 返回原始的可空初始化状态
    isInitializedNullable(state): boolean | null {
      return state.initialized
    },
  },

  actions: {
    async checkInitialization(): Promise<boolean> {
      // 如果已经检查过初始化状态，直接返回结果
      if (this.initialized !== null) {
        return this.isInitialized
      }

      try {
        const response = await systemApi.checkInitialized()
        this.initialized = response.data.data.initialized
        return this.isInitialized
      } catch (error) {
        console.error('检查初始化状态失败:', error)
        // 出错时返回false，确保用户能够访问初始化页面
        this.initialized = false
        return this.isInitialized
      }
    },

    async fetchSystemInfo(): Promise<SystemInfo | null> {
      // 如果已经初始化过，直接返回缓存的数据
      if (this.systemInfo) {
        console.log('从缓存中获取系统信息')
        return this.systemInfo
      }

      try {
        const response = await systemApi.getSystemInfo()
        if (response.data.code === 0) {
          this.systemInfo = response.data.data
          // 获取系统信息成功说明系统已经初始化
          this.initialized = true
          return this.systemInfo
        }
        return null
      } catch (error) {
        console.error('获取系统信息失败:', error)
        return null
      }
    },

    // 更新系统信息
    updateSystemInfo(info: SystemInfo) {
      this.systemInfo = info
    },

    // 设置初始化状态
    setInitialized(status: boolean) {
      this.initialized = status
    },

    // 清除缓存
    clearCache() {
      this.systemInfo = null
      this.initialized = null
    },
  },
})
