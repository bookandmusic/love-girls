import api from './api'

export interface Wish {
  id: number
  content: string
  authorName: string
  email: string
  createdAt: string
  approved: boolean
}

interface GetWishesResponse {
  code: number
  data: {
    wishes: Wish[]
    totalPages: number
    total: number
    totalCount?: number
    page?: number
    size?: number
  }
  message?: string
}

interface CreateWishRequest {
  content: string
  authorName: string
  email: string
}

interface CreateWishResponse {
  code: number
  data: Wish
  message?: string
}

interface DeleteWishResponse {
  code: number
  message?: string
}

interface ApproveWishResponse {
  code: number
  message?: string
}

export const wishApi = {
  // 获取愿望列表
  async getWishes(page: number, size: number, approved?: boolean) {
    const params: { page: number; size: number; approved?: boolean } = {
      page,
      size,
    }

    // 前台必须传入 approved=true
    if (approved !== undefined) {
      params.approved = approved
    }

    const response = await api.get<GetWishesResponse>('/wishes', {
      params,
    })
    return response.data
  },

  // 创建愿望
  async createWish(wishData: CreateWishRequest) {
    const response = await api.post<CreateWishResponse>('/wishes', wishData)
    return response.data
  },

  // 删除愿望
  async deleteWish(id: number) {
    const response = await api.delete<DeleteWishResponse>(`/wishes/${id}`)
    return response.data
  },

  // 批准愿望
  async approveWish(id: number) {
    const response = await api.put<ApproveWishResponse>(`/wishes/${id}/approve`)
    return response.data
  },
}
