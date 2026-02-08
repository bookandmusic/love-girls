// 地点相关接口
import type { MockMethod } from 'vite-plugin-mock'

export interface Photo {
  id: number
  placeId: number
  url: string
  thumbnailUrl?: string
}
// 定义类型
export interface Place {
  id: number
  name: string
  latitude: number
  longitude: number
  image?: Photo
  description: string
  date: string // YYYY-MM-DD
}

interface GetPlacesResponse {
  code: number
  data: {
    places: Place[]
    totalPages: number
    total?: number
    totalCount?: number
    page?: number
    size?: number
  }
  msg?: string
}

interface CreateOrUpdatePlaceResponse {
  code: number
  data?: Place
  msg?: string
}

interface DeletePlaceResponse {
  code: number
  msg?: string
}

// 模拟数据
const mockPlaces: Place[] = [
  {
    id: 1,
    name: '北京天安门',
    latitude: 39.9092,
    longitude: 116.3975,
    image: {
      url: 'https://picsum.photos/200/300?random=1',
      thumbnailUrl: 'https://picsum.photos/200/300?random=1',
      id: 1,
      placeId: 1,
    },
    description: '我们的第一次旅行',
    date: '2023-05-01',
  },
  {
    id: 2,
    name: '上海外滩',
    latitude: 31.2363,
    longitude: 121.4903,
    image: {
      url: 'https://picsum.photos/200/300?random=2',
      thumbnailUrl: 'https://picsum.photos/200/300?random=2',
      id: 2,
      placeId: 2,
    },
    description: '浪漫的夜晚',
    date: '2023-08-15',
  },
  {
    id: 3,
    name: '杭州西湖',
    latitude: 30.2424,
    longitude: 120.1495,
    image: undefined,
    description: '春天的约会',
    date: '2024-03-20',
  },
  {
    id: 4,
    name: '成都宽窄巷子',
    latitude: 30.6701,
    longitude: 104.0661,
    image: undefined,
    description: '美食之旅',
    date: '2024-05-01',
  },
  {
    id: 5,
    name: '西安兵马俑',
    latitude: 34.3839,
    longitude: 109.2719,
    image: undefined,
    description: '历史探索',
    date: '2024-07-15',
  },
]

const places: MockMethod[] = [
  // 获取地点列表（支持分页）
  {
    url: '/api/v1/places',
    method: 'get',
    response: (params: { query: Record<string, string> }): GetPlacesResponse => {
      const { query } = params
      const page = parseInt(query?.page || '1', 10)
      const size = parseInt(query?.size || '5', 10)

      const startIndex = (page - 1) * size
      const endIndex = startIndex + size
      const paginatedPlaces = mockPlaces.slice(startIndex, endIndex)

      return {
        code: 0,
        data: {
          places: paginatedPlaces,
          totalPages: Math.ceil(mockPlaces.length / size),
          total: mockPlaces.length,
        },
        msg: '获取成功',
      }
    },
  },

  // 创建地点
  {
    url: '/api/v1/places',
    method: 'post',
    response: (params: { body: Omit<Place, 'id'> }): CreateOrUpdatePlaceResponse => {
      const newPlace: Place = {
        id: mockPlaces.length > 0 ? Math.max(...mockPlaces.map(p => p.id)) + 1 : 1,
        ...params.body,
      }
      mockPlaces.push(newPlace)
      return {
        code: 0,
        data: newPlace,
        msg: '创建成功',
      }
    },
  },

  // 更新地点信息
  {
    url: '/api/v1/places/:id',
    method: 'put',
    response: (params: {
      query: Record<string, string>
      body: Partial<Place>
    }): CreateOrUpdatePlaceResponse => {
      const idParam = params.query.id
      if (idParam === undefined) {
        return {
          code: 1,
          msg: 'ID参数缺失',
        }
      }
      const id = parseInt(idParam, 10)
      const placeIndex = mockPlaces.findIndex(p => p.id === id)

      if (placeIndex === -1) {
        return {
          code: 1,
          msg: '地点不存在',
        }
      }

      const current = mockPlaces[placeIndex]
      if (!current) {
        return {
          code: 1,
          msg: '地点不存在',
        }
      }

      // 安全合并：只覆盖非 undefined 字段
      const updated: Place = {
        id: current.id, // 确保id始终存在且不为undefined
        name: params.body.name ?? current.name,
        latitude: params.body.latitude ?? current.latitude,
        longitude: params.body.longitude ?? current.longitude,
        image: params.body.image ?? current.image,
        description: params.body.description ?? current.description,
        date: params.body.date ?? current.date,
      }

      mockPlaces[placeIndex] = updated

      return {
        code: 0,
        data: updated,
        msg: '更新成功',
      }
    },
  },

  // 删除地点
  {
    url: '/api/v1/places/:id',
    method: 'delete',
    response: (params: { query: Record<string, string> }): DeletePlaceResponse => {
      const idParam = params.query.id
      if (idParam === undefined) {
        return {
          code: 1,
          msg: 'ID参数缺失',
        }
      }
      const id = parseInt(idParam, 10)
      const placeIndex = mockPlaces.findIndex(p => p.id === id)

      if (placeIndex === -1) {
        return {
          code: 1,
          msg: '地点不存在',
        }
      }

      mockPlaces.splice(placeIndex, 1)

      return {
        code: 0,
        msg: '删除成功',
      }
    },
  },
]

export default places
