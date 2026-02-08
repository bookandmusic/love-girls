// 相册相关接口
import type { MockMethod } from 'vite-plugin-mock'

export interface Photo {
  id: number
  albumId: number
  url: string
  thumbnailUrl?: string
}
// 定义相册类型
interface Album {
  id: number
  name: string
  description: string
  coverImage?: Photo
  createdAt: string // ISO 日期字符串，如 '2023-05-01'
  photoCount: number
}

// 更新相册时允许的字段（不包含 id、createdAt、photoCount 等只读字段）
interface UpdateAlbumDTO {
  name?: string
  description?: string
  coverImage?: Photo
}

// 模拟数据
const mockAlbums: Album[] = [
  {
    id: 1,
    name: '我们的第一次旅行',
    description: '北京之旅',
    coverImage: {
      id: 1,
      url: 'https://picsum.photos/600/400?random=2001',
      thumbnailUrl: 'https://picsum.photos/300/200?random=2001',
      albumId: 1,
    },
    createdAt: '2023-05-01',
    photoCount: 24,
  },
  {
    id: 2,
    name: '浪漫的夜晚',
    description: '上海外滩',
    coverImage: {
      id: 2,
      url: 'https://picsum.photos/600/400?random=2002',
      thumbnailUrl: 'https://picsum.photos/300/200?random=2002',
      albumId: 2,
    },
    createdAt: '2023-08-15',
    photoCount: 18,
  },
  {
    id: 3,
    name: '春天的约会',
    description: '杭州西湖',
    coverImage: {
      id: 3,
      url: 'https://picsum.photos/600/400?random=2003',
      thumbnailUrl: 'https://picsum.photos/300/200?random=2003',
      albumId: 3,
    },
    createdAt: '2024-03-20',
    photoCount: 32,
  },
  {
    id: 4,
    name: '周末的周末',
    description: '苏州',
    coverImage: undefined,
    createdAt: '2024-07-01',
    photoCount: 10,
  },
  {
    id: 5,
    name: '周末的周末',
    description: '上海',
    coverImage: undefined,
    createdAt: '2024-07-01',
    photoCount: 10,
  },
]

const albums: MockMethod[] = [
  // 获取相册列表
  {
    url: '/api/v1/albums',
    method: 'get',
    response: (params: { query: Record<string, string> }) => {
      const { query } = params
      const page = parseInt(query?.page || '1', 10)
      const size = parseInt(query?.size || '5', 10)

      const startIndex = (page - 1) * size
      const endIndex = startIndex + size
      const paginatedAlbums = mockAlbums.slice(startIndex, endIndex)

      return {
        code: 0,
        data: {
          albums: paginatedAlbums,
          totalPages: Math.ceil(mockAlbums.length / size),
          total: mockAlbums.length,
        },
        msg: '获取成功',
      }
    },
  },
  // 更新相册
  {
    url: '/api/v1/albums/:id',
    method: 'put',
    response: (params: { query: Record<string, string>; body: UpdateAlbumDTO }) => {
      const idParam = params.query.id
      if (idParam === undefined) {
        return {
          code: 1,
          data: null,
          msg: 'ID参数缺失',
        }
      }
      const id = parseInt(idParam, 10)
      const albumIndex = mockAlbums.findIndex(a => a.id === id)

      if (albumIndex === -1) {
        return {
          code: 1,
          data: null,
          msg: '相册不存在',
        }
      }

      // 构造要更新的字段（可能包含 undefined）
      const updatedFields: Partial<Album> = {
        name: params.body.name,
        description: params.body.description,
        coverImage: params.body.coverImage,
      }

      // 合并并断言为 Album（安全，因为基底是完整对象）
      mockAlbums[albumIndex] = {
        ...mockAlbums[albumIndex],
        ...updatedFields,
      } as Album

      return {
        code: 0,
        data: mockAlbums[albumIndex],
        msg: '更新成功',
      }
    },
  },
  // 删除相册
  {
    url: '/api/v1/albums/:id',
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
      const albumIndex = mockAlbums.findIndex(a => a.id === id)

      if (albumIndex === -1) {
        return {
          code: 1,
          data: null,
          msg: '相册不存在',
        }
      }

      mockAlbums.splice(albumIndex, 1)

      return {
        code: 0,
        data: null,
        msg: '删除成功',
      }
    },
  },
  // 创建相册
  {
    url: '/api/v1/albums',
    method: 'post',
    response: (params: { body: Omit<Album, 'id' | 'createdAt' | 'photoCount'> }) => {
      const newAlbum = {
        id: mockAlbums.length + 1,
        createdAt: new Date().toISOString(),
        photoCount: 0,
        ...params.body,
      }
      mockAlbums.push(newAlbum)
      return {
        code: 0,
        data: newAlbum,
        msg: '创建成功',
      }
    },
  },
]

export default albums
