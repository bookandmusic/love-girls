// 动态相关接口
import type { MockMethod } from 'vite-plugin-mock'

// 定义嵌套类型
interface Author {
  name: string
  avatar: string
}
interface Photo {
  id: number
  momentId: number
  url: string
  thumbnailUrl: string
  alt?: string
}
interface Moment {
  id: number
  content: string
  images: Photo[]
  likes: number
  createdAt: string // 格式: 'YYYY-MM-DD HH:mm:ss'
  author: Author
  isPublic: boolean
}

// 创建动态的请求体（所有字段都应提供，但可设默认值）
interface CreateMomentDTO {
  content?: string
  images?: Photo[]
  likes?: number
  author?: Partial<Author>
  isPublic?: boolean
}

// 更新动态的请求体（所有字段可选）
interface UpdateMomentDTO {
  content?: string
  images?: Photo[]
  likes?: number
  createdAt?: string
  author?: Partial<Author>
  isPublic?: boolean
}

// 仅更新公开状态
interface UpdatePublicStatusDTO {
  isPublic?: boolean
}

// 模拟数据
const mockMoments: Moment[] = [
  {
    id: 1,
    content: '今天的夕阳真是太美了，和你一起看夕阳是我最幸福的时刻。',
    images: [
      {
        id: 1,
        momentId: 1,
        url: 'https://picsum.photos/600/400?random=1001',
        thumbnailUrl: 'https://picsum.photos/200/150?random=1001',
      },
    ],
    likes: 5,
    createdAt: '2024-01-15',
    author: {
      name: '小陆',
      avatar: '',
    },
    isPublic: true,
  },
  {
    id: 2,
    content: '一起做的第一顿饭，虽然简单但很温暖。生活中的小确幸就是和你在一起的每一天。',
    images: [],
    likes: 8,
    createdAt: '2024-01-22',
    author: {
      name: '小伍',
      avatar: '',
    },
    isPublic: false,
  },
  {
    id: 3,
    content: '周末的公园散步，发现了一个很美的小湖，下次我们带野餐垫来吧！',
    images: [
      {
        id: 2,
        momentId: 3,
        url: 'https://picsum.photos/600/400?random=1002',
        thumbnailUrl: 'https://picsum.photos/200/150?random=1002',
      },
      {
        id: 3,
        momentId: 3,
        url: 'https://picsum.photos/600/400?random=1003',
        thumbnailUrl: 'https://picsum.photos/200/150?random=1003',
      },
    ],
    likes: 3,
    createdAt: '2024-02-05',
    author: {
      name: '小陆',
      avatar: '',
    },
    isPublic: true,
  },
  {
    id: 4,
    content: '今天是我们在一起的第365天，时间过得真快，但每一天都很珍贵。',
    images: [],
    likes: 12,
    createdAt: '2024-02-10',
    author: {
      name: '小伍',
      avatar: '',
    },
    isPublic: true,
  },
  {
    id: 5,
    content: '一起看的电影太好看了，回家的路上还在讨论剧情呢！',
    images: [
      {
        id: 4,
        momentId: 4,
        url: 'https://picsum.photos/600/400?random=1004',
        thumbnailUrl: 'https://picsum.photos/200/150?random=1004',
      },
      {
        id: 5,
        momentId: 4,
        url: 'https://picsum.photos/600/400?random=1005',
        thumbnailUrl: 'https://picsum.photos/200/150?random=1005',
      },
    ],
    likes: 6,
    createdAt: '2024-02-18',
    author: {
      name: '小陆',
      avatar: 'https://picsum.photos/200/200?random=1016',
    },
    isPublic: true,
  },
  {
    id: 6,
    content: '下雨天窝在家里一起煮火锅，外面雨声阵阵，屋内热气腾腾，这就是家的感觉吧。',
    images: [
      {
        id: 6,
        momentId: 5,
        url: 'https://picsum.photos/600/400?random=1006',
        thumbnailUrl: 'https://picsum.photos/200/150?random=1006',
      },
      {
        id: 7,
        momentId: 5,
        url: 'https://picsum.photos/600/400?random=1007',
        thumbnailUrl: 'https://picsum.photos/200/150?random=1007',
      },
    ],
    likes: 9,
    createdAt: '2024-03-02',
    author: {
      name: '小伍',
      avatar: '',
    },
    isPublic: true,
  },
  {
    id: 7,
    content: '一起去爬山，虽然很累但是山顶的风景真的值得！',
    images: [],
    likes: 7,
    createdAt: '2024-03-15',
    author: {
      name: '小陆',
      avatar: '',
    },
    isPublic: true,
  },
  {
    id: 8,
    content: '你送我的花已经开了，粉粉的特别好看，就像你一样温柔。',
    images: [],
    likes: 11,
    createdAt: '2024-03-22',
    author: {
      name: '小伍',
      avatar: '',
    },
    isPublic: true,
  },
]

const moments: MockMethod[] = [
  // 获取动态列表
  {
    url: '/api/v1/moments',
    method: 'get',
    response: (params: { query: Record<string, string> }) => {
      const { query } = params
      const page = parseInt(query?.page || '1', 10)
      const size = parseInt(query?.size || '5', 10)

      const startIndex = (page - 1) * size
      const endIndex = startIndex + size
      const paginatedMoments = mockMoments.slice(startIndex, endIndex)

      return {
        code: 0,
        data: {
          moments: paginatedMoments,
          totalPages: Math.ceil(mockMoments.length / size),
          total: mockMoments.length,
        },
        msg: '获取成功',
      }
    },
  },

  // 创建动态
  {
    url: '/api/v1/moments',
    method: 'post',
    response: (params: { body: CreateMomentDTO }) => {
      const newId = mockMoments.length > 0 ? Math.max(...mockMoments.map(m => m.id)) + 1 : 1

      const defaultAuthor: Author = {
        name: '小陆',
        avatar: '',
      }

      const newMoment: Moment = {
        id: newId,
        content: params.body.content ?? '',
        images: Array.isArray(params.body.images) ? params.body.images : [],
        likes: typeof params.body.likes === 'number' ? params.body.likes : 0,
        createdAt: new Date().toISOString().slice(0, 19).replace('T', ' '),
        author: params.body.author
          ? {
              name: params.body.author.name ?? defaultAuthor.name,
              avatar: params.body.author.avatar ?? defaultAuthor.avatar,
            }
          : defaultAuthor,
        isPublic: typeof params.body.isPublic === 'boolean' ? params.body.isPublic : true,
      }

      mockMoments.push(newMoment)

      return {
        code: 0,
        data: newMoment,
        msg: '创建成功',
      }
    },
  },
  // 更新动态
  {
    url: '/api/v1/moments/:id',
    method: 'put',
    response: (params: { query: Record<string, string>; body: UpdateMomentDTO }) => {
      const idParam = params.query.id
      if (idParam === undefined) {
        return {
          code: 1,
          data: null,
          msg: 'ID参数缺失',
        }
      }
      const id = parseInt(idParam, 10)
      const momentIndex = mockMoments.findIndex(m => m.id === id)

      if (momentIndex === -1) {
        return {
          code: 1,
          data: null,
          msg: '动态不存在',
        }
      }

      // ✅ 安全断言：此时 current 一定存在
      const current = mockMoments[momentIndex]!

      const updated: Moment = {
        ...current,
        content: params.body.content ?? current.content,
        images: params.body.images ?? current.images,
        likes: params.body.likes ?? current.likes,
        createdAt: params.body.createdAt ?? current.createdAt,
        isPublic: params.body.isPublic ?? current.isPublic,
        author: {
          name: params.body.author?.name ?? current.author.name,
          avatar: params.body.author?.avatar ?? current.author.avatar,
        },
      }

      mockMoments[momentIndex] = updated

      return {
        code: 0,
        data: updated,
        msg: '更新成功',
      }
    },
  },
  {
    url: '/api/v1/moments/:id/like',
    method: 'post',
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
      const momentIndex = mockMoments.findIndex(m => m.id === id)

      if (momentIndex === -1) {
        return {
          code: 1,
          data: null,
          msg: '动态不存在',
        }
      }

      const current = mockMoments[momentIndex]!
      mockMoments[momentIndex] = {
        ...current,
        likes: current.likes + 1,
      }

      return {
        code: 0,
        data: mockMoments[momentIndex],
        msg: '更新成功',
      }
    },
  },
  {
    url: '/api/v1/moments/:id/public',
    method: 'put',
    response: (params: { query: Record<string, string>; body: UpdatePublicStatusDTO }) => {
      const idParam = params.query.id
      if (idParam === undefined) {
        return {
          code: 1,
          data: null,
          msg: 'ID参数缺失',
        }
      }
      const id = parseInt(idParam, 10)
      const momentIndex = mockMoments.findIndex(m => m.id === id)

      if (momentIndex === -1) {
        return {
          code: 1,
          data: null,
          msg: '动态不存在',
        }
      }

      const current = mockMoments[momentIndex]!
      mockMoments[momentIndex] = {
        ...current,
        isPublic: params.body.isPublic ?? current.isPublic,
      }

      return {
        code: 0,
        data: mockMoments[momentIndex],
        msg: '更新成功',
      }
    },
  },
  // 删除动态
  {
    url: '/api/v1/moments/:id',
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
      const momentIndex = mockMoments.findIndex(m => m.id === id)

      if (momentIndex === -1) {
        return {
          code: 1,
          data: null,
          msg: '动态不存在',
        }
      }

      mockMoments.splice(momentIndex, 1)

      return {
        code: 0,
        data: null,
        msg: '删除成功',
      }
    },
  },
]

export default moments
