import api from './api'

// 定义仪表盘相关的类型
export interface AlbumStats {
  total: number
  totalPhotos: number
}

export interface PlaceStats {
  total: number
}

export interface MomentStats {
  total: number
}

export interface WishStats {
  total: number
  pending: number
}

export interface DashboardData {
  albumStats: AlbumStats
  placeStats: PlaceStats
  momentStats: MomentStats
  wishStats: WishStats
}

// 仪表盘相关API接口
export const dashboardApi = {
  /**
   * 获取仪表盘统计数据
   */
  async getDashboardStats() {
    const response = await api.get<{ code: number; data: DashboardData }>('/system/dashboard/stats')
    return response.data
  },
}
