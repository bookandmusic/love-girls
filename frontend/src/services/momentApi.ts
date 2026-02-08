import api from './api'
import type { FileInfo } from './upload'

export interface Photo {
  id: number
  momentId: number
  file?: FileInfo
}
export interface Moment {
  id: number
  content: string
  images?: Photo[]
  imageIds?: number[]
  likes: number
  createdAt: string
  author: {
    name: string
    avatar?: FileInfo
  }
  isPublic: boolean
  userId?: number
}

interface GetMomentsResponse {
  code: number
  data: {
    moments: Moment[]
    totalPages: number
    total?: number
    totalCount?: number
    page?: number
    size?: number
  }
  msg?: string
}

interface UpdateMomentPublicRequest {
  isPublic: boolean
}

interface DeleteMomentResponse {
  code: number
  msg?: string
}

interface CreateOrUpdateMomentResponse {
  code: number
  data: Moment
  msg?: string
}

interface LikeMomentResponse {
  code: number
  data: {
    likes: number
  }
  msg?: string
}

export const momentApi = {
  // 获取动态列表
  async getMoments(page: number, size: number) {
    const response = await api.get<GetMomentsResponse>('/moments', {
      params: {
        page,
        size,
      },
    })
    return response.data
  },

  // 创建动态
  async createMoment(momentData: Omit<Moment, 'id'>) {
    const response = await api.post<CreateOrUpdateMomentResponse>('/moments', momentData)
    return response.data
  },

  // 更新动态
  async updateMoment(id: number, momentData: Partial<Moment>) {
    const response = await api.put<CreateOrUpdateMomentResponse>(`/moments/${id}`, momentData)
    return response.data
  },

  // 更新动态公开状态
  async updateMomentPublic(id: number, data: UpdateMomentPublicRequest) {
    const response = await api.put<CreateOrUpdateMomentResponse>(`/moments/${id}/public`, data)
    return response.data
  },

  // 点赞动态
  async likeMoment(momentId: number) {
    const response = await api.post<LikeMomentResponse>(`/moments/${momentId}/like`)
    return response.data
  },

  // 删除动态
  async deleteMoment(id: number) {
    const response = await api.delete<DeleteMomentResponse>(`/moments/${id}`)
    return response.data
  },
}
