// 模拟祝福数据
import type { MockMethod } from 'vite-plugin-mock'

const mockWishes = [
  {
    id: 1,
    content: '祝你们永远幸福，白头偕老！',
    authorName: '朋友小王',
    email: 'wang@example.com',
    createdAt: '2024-10-15 14:30:22',
    approved: false,
  },
  {
    id: 2,
    content: '看着你们幸福的样子，真让人羡慕，祝你们永远甜蜜！',
    authorName: '同事小李',
    email: 'li@example.com',
    createdAt: '2024-10-16 09:15:33',
    approved: true,
  },
  {
    id: 3,
    content: '爱情最美的样子就是你们这样，祝愿你们的爱情长长久久！',
    authorName: '表姐小张',
    email: 'zhang@example.com',
    createdAt: '2024-10-17 18:45:11',
    approved: false,
  },
  {
    id: 4,
    content: '愿你们的爱情如美酒般越来越醇香！',
    authorName: '同学小刘',
    email: 'liu@example.com',
    createdAt: '2024-10-18 11:20:45',
    approved: true,
  },
  {
    id: 5,
    content: '真心祝愿你们幸福美满，永远相爱！',
    authorName: '邻居小陈',
    email: 'chen@example.com',
    createdAt: '2024-10-19 20:10:12',
    approved: false,
  },
  {
    id: 6,
    content: '看到你们这么恩爱，我也要加油找到另一半！',
    authorName: '室友小赵',
    email: 'zhao@example.com',
    createdAt: '2024-10-20 15:30:55',
    approved: true,
  },
  {
    id: 7,
    content: '你们是天生一对，地设一双！',
    authorName: '长辈小孙',
    email: 'sun@example.com',
    createdAt: '2024-10-21 08:45:22',
    approved: false,
  },
  {
    id: 8,
    content: '爱情的花朵在你们心中绽放，愿它永远不凋零！',
    authorName: '远方朋友小周',
    email: 'zhou@example.com',
    createdAt: '2024-10-22 16:25:30',
    approved: true,
  },
]

const wishes: MockMethod[] = [
  // 获取祝福列表（公开接口）
  {
    url: '/api/v1/wishes',
    method: 'get',
    response: (params: { query: Record<string, string> }) => {
      const { query } = params
      const page = parseInt(query?.page || '1')
      const size = parseInt(query?.size || '5')

      // 计算总页数
      const totalPages = Math.ceil(mockWishes.length / size) || 1

      // 计算起始索引
      const startIndex = (page - 1) * size
      const endIndex = startIndex + size

      // 获取当前页数据
      const data = mockWishes.slice(startIndex, endIndex)

      return {
        code: 0,
        data: {
          wishes: data,
          totalPages: totalPages,
          total: mockWishes.length,
          totalCount: mockWishes.length,
          page: page,
          size: size,
        },
        msg: '获取成功',
      }
    },
  },

  // 提交祝福（公开接口）
  {
    url: '/api/v1/wishes',
    method: 'post',
    response: (params: { body: Record<string, string> }) => {
      try {
        const { body } = params
        // 生成新的整数ID - 使用当前最大ID+1的方式
        const newId =
          mockWishes.length > 0 ? Math.max(...mockWishes.map((w: { id: number }) => w.id)) + 1 : 1

        const newWish = {
          id: newId,
          content: body.content || `${newId}测试`,
          authorName: body.authorName || `${newId}测试`,
          email: body.email || `${newId}测试@example.com`,
          createdAt: new Date().toISOString().slice(0, 19).replace('T', ' '),
          approved: false, // 新提交的祝福默认未批准
        }

        // 添加到列表开头
        mockWishes.unshift(newWish)

        return {
          code: 0,
          data: newWish,
          msg: '祝福提交成功',
        }
      } catch {
        return {
          code: 1,
          data: null,
          msg: '提交失败，请检查数据格式',
        }
      }
    },
  },
  // 删除指定祝福
  {
    url: '/api/v1/wishes/:id',
    method: 'delete',
    response: (params: { query: Record<string, string> }) => {
      const idParam = params.query.id
      if (idParam === undefined) {
        return {
          code: 1,
          data: null,
          msg: 'ID参数缺失',
        }
      }
      const id = parseInt(idParam, 10)
      const wishIndex = mockWishes.findIndex(w => w.id === id)

      if (wishIndex === -1) {
        return {
          code: 1,
          data: null,
          msg: '愿望不存在',
        }
      }

      mockWishes.splice(wishIndex, 1)

      return {
        code: 0,
        data: null,
        msg: '删除成功',
      }
    },
  },
  // 批准指定祝福
  {
    url: '/api/v1/wishes/:id/approve',
    method: 'put',
    response: (params: { query: Record<string, string> }) => {
      const idParam = params.query.id
      if (idParam === undefined) {
        return {
          code: 1,
          data: null,
          msg: 'ID参数缺失',
        }
      }
      const id = parseInt(idParam, 10)
      const wish = mockWishes.find(w => w.id === id)

      if (!wish) {
        return {
          code: 1,
          data: null,
          msg: '愿望不存在',
        }
      }

      // 更新批准状态
      wish.approved = true

      return {
        code: 0,
        data: null,
        msg: '批准成功',
      }
    },
  },
]

export default wishes
