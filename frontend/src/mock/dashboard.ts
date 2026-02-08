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
            total: 12,
            totalPhotos: 142,
          },
          // 足迹统计
          placeStats: {
            total: 24,
          },
          // 纪念日统计
          momentStats: {
            total: 8,
          },
          // 祝福统计
          wishStats: {
            total: 42,
            pending: 3,
          },
        },
      }
    },
  },
]

export default dashboard
