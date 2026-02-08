import api from './api'
import type { FileInfo } from './upload'

export interface Photo {
  id: number
  placeId: number
  file?: FileInfo
}

export interface Place {
  id: number
  name: string
  latitude: number
  longitude: number
  image?: Photo
  description: string
  date: string
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
  data: Place
  msg?: string
}

interface DeletePlaceResponse {
  code: number
  msg?: string
}

export const placeApi = {
  // 获取地点列表
  async getPlaces(page: number, size: number) {
    const response = await api.get<GetPlacesResponse>('/places', {
      params: {
        page,
        size,
      },
    })
    return response.data
  },

  // 创建地点
  async createPlace(placeData: Omit<Place, 'id'>) {
    const response = await api.post<CreateOrUpdatePlaceResponse>('/places', placeData)
    return response.data
  },

  // 更新地点
  async updatePlace(id: number, placeData: Partial<Place>) {
    const response = await api.put<CreateOrUpdatePlaceResponse>(`/places/${id}`, placeData)
    return response.data
  },

  // 删除地点
  async deletePlace(id: number) {
    const response = await api.delete<DeletePlaceResponse>(`/places/${id}`)
    return response.data
  },
}
