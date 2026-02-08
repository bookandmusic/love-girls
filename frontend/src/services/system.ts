import api from './api'
import type { FileInfo } from './upload'

// 定义站点设置接口
export interface SiteSettings {
  siteTitle: string
  siteDescription: string
}

// 定义系统初始化请求参数
export interface InitSystemRequest {
  siteName: string
  siteDescription: string
  startDate: string
  userAName: string
  userARole: 'boy' | 'girl'
  userAEmail?: string
  userAPhone?: string
  userBName: string
  userBRole: 'boy' | 'girl'
  userBEmail?: string
  userBPhone?: string
  sitePassword: string
  sitePasswordConfirm: string
}

// 定义系统初始化响应
export interface InitSystemResponse {
  code: number
  msg?: string
  initialized?: boolean
  data?: Record<string, unknown>
}

// 定义检查初始化响应
export interface CheckInitializedResponse {
  code: number
  data: {
    initialized: boolean
  }
}

// 定义系统信息接口
export interface SystemInfo {
  site: {
    name: string
    description: string
    startDate: string
  }
  couple: {
    boy: {
      name: string
      avatar?: FileInfo
    }
    girl: {
      name: string
      avatar?: FileInfo
    }
  }
}

// 定义获取系统信息响应
export interface GetSystemInfoResponse {
  code: number
  data: SystemInfo
}

// 定义系统初始化相关API
export const systemApi = {
  /**
   * 初始化系统
   * @param data 初始化数据，包含站点信息、用户信息和密码
   */
  initSystem(data: InitSystemRequest) {
    return api.post<InitSystemResponse>('/system/init', data)
  },

  /**
   * 检查系统是否已初始化
   */
  checkInitialized() {
    return api.get<CheckInitializedResponse>('/system/init')
  },

  /**
   * 获取系统信息
   */
  getSystemInfo() {
    return api.get<GetSystemInfoResponse>('/system/info')
  },

  /**
   * 获取站点设置
   */
  getSiteSettings() {
    return api.get<{ code: number; data: SiteSettings }>('/system/settings/site')
  },

  /**
   * 保存站点设置
   */
  saveSiteSettings(settings: SiteSettings) {
    return api.post<{ code: number; data: SiteSettings }>('/system/settings/site', settings)
  },
}
