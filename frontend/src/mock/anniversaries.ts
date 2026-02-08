// 纪念日相关接口
import type { MockMethod } from 'vite-plugin-mock'

// 与API端保持一致的类型定义
export interface Anniversary {
  id: number
  title: string
  date: string // MM-DD
  description: string
  calendar: 'solar' | 'lunar'
}

interface GetAnniversariesResponse {
  code: number
  data: {
    anniversaries: Anniversary[]
    totalPages: number
    total?: number
    totalCount?: number
    page?: number
    size?: number
  }
  msg?: string
}

interface CreateOrUpdateAnniversaryResponse {
  code: number
  data?: Anniversary
  msg?: string
}

interface DeleteAnniversaryResponse {
  code: number
  msg?: string
}

// 模拟数据
const mockAnniversaries: Anniversary[] = [
  {
    id: 1,
    title: '初次相遇',
    date: '2021-03-15',
    description: '在那个阳光明媚的下午，我们第一次相遇',
    calendar: 'solar',
  },
  {
    id: 2,
    title: '第一次约会',
    date: '2021-04-02',
    description: '一起去看电影，紧张又兴奋',
    calendar: 'solar',
  },
  {
    id: 3,
    title: '确立关系',
    date: '2022-05-20',
    description: '表白成功的日子，永远难忘',
    calendar: 'solar',
  },
  {
    id: 4,
    title: '第一次旅行',
    date: '2022-07-10',
    description: '一起去海边，留下了美好的回忆',
    calendar: 'solar',
  },
  {
    id: 5,
    title: '第一次见家长',
    date: '2024-01-01',
    description: '重要的里程碑，双方父母都很喜欢我们',
    calendar: 'lunar',
  },
  {
    id: 6,
    title: '订婚',
    date: '2024-10-10',
    description: '许下承诺，决定共度余生',
    calendar: 'lunar',
  },
  {
    id: 7,
    title: '结婚',
    date: '2025-05-20',
    description: '人生最重要的时刻，终于成为彼此的唯一',
    calendar: 'lunar',
  },
]

const anniversaries: MockMethod[] = [
  {
    url: '/api/v1/anniversaries',
    method: 'get',
    response: (params: { query: Record<string, string> }): GetAnniversariesResponse => {
      const { query } = params
      const page = parseInt(query?.page || '1', 10)
      const size = parseInt(query?.size || '10', 10)

      const startIndex = (page - 1) * size
      const endIndex = startIndex + size
      const paginatedAnniversaries = mockAnniversaries.slice(startIndex, endIndex)

      return {
        code: 0,
        data: {
          anniversaries: paginatedAnniversaries,
          totalPages: Math.ceil(mockAnniversaries.length / size),
          total: mockAnniversaries.length,
        },
        msg: '获取成功',
      }
    },
  },
  // 创建纪念日
  {
    url: '/api/v1/anniversaries',
    method: 'post',
    response: (params: { body: Omit<Anniversary, 'id'> }): CreateOrUpdateAnniversaryResponse => {
      const newAnniversary: Anniversary = {
        id: mockAnniversaries.length > 0 ? Math.max(...mockAnniversaries.map(m => m.id)) + 1 : 1,
        ...params.body,
      }
      mockAnniversaries.push(newAnniversary)
      return {
        code: 0,
        data: newAnniversary,
        msg: '创建成功',
      }
    },
  },
  // 更新纪念日
  {
    url: '/api/v1/anniversaries/:id',
    method: 'put',
    response: (params: {
      query: Record<string, string>
      body: Partial<Anniversary>
    }): CreateOrUpdateAnniversaryResponse => {
      const idParam = params.query.id
      if (idParam === undefined) {
        return {
          code: 1,
          msg: 'ID参数缺失',
        }
      }
      const id = parseInt(idParam, 10)
      const index = mockAnniversaries.findIndex(m => m.id === id)
      if (index === -1) {
        return {
          code: 1,
          msg: '纪念日不存在',
        }
      }

      const current = mockAnniversaries[index]
      if (!current) {
        return {
          code: 1,
          msg: '纪念日不存在',
        }
      }
      mockAnniversaries[index] = {
        id: current.id, // 确保id始终存在且不为undefined
        title: params.body.title ?? current.title,
        date: params.body.date ?? current.date,
        description: params.body.description ?? current.description,
        calendar: params.body.calendar ?? current.calendar,
      }
      return {
        code: 0,
        data: mockAnniversaries[index],
        msg: '更新成功',
      }
    },
  },
  // 删除纪念日
  {
    url: '/api/v1/anniversaries/:id',
    method: 'delete',
    response: (params: { query: Record<string, string> }): DeleteAnniversaryResponse => {
      const idParam = params.query.id
      if (idParam === undefined) {
        return {
          code: 1,
          msg: 'ID参数缺失',
        }
      }
      const id = parseInt(idParam, 10)
      const index = mockAnniversaries.findIndex(m => m.id === id)
      if (index === -1) {
        return {
          code: 1,
          msg: '纪念日不存在',
        }
      }
      mockAnniversaries.splice(index, 1)
      return {
        code: 0,
        msg: '删除成功',
      }
    },
  },
]

export default anniversaries
