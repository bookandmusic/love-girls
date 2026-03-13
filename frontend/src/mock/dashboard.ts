// 仪表盘统计数据相关接口
import type { MockMethod } from 'vite-plugin-mock'

const dashboard: MockMethod[] = [
  {
    url: '/api/v1/system/dashboard/stats',
    method: 'get',
    response: () => {
      return {
        code: 0,
        data: {
          // 相册统计
          albumStats: {
            total: 5,
            totalPhotos: 124,
          },
          // 足迹统计
          placeStats: {
            total: 5,
          },
          // 动态统计
          momentStats: {
            total: 8,
            private: 1,
          },
          // 纪念日统计
          anniversaryStats: {
            total: 7,
          },
          // 祝福统计
          wishStats: {
            total: 8,
            pending: 4,
          },
          // 用户统计
          userStats: {
            total: 2,
          },
        },
      }
    },
  },
]

export default dashboard
