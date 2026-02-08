import api from './api'

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
  data: Anniversary
  msg?: string
}

interface DeleteAnniversaryResponse {
  code: number
  msg?: string
}

export const anniversaryApi = {
  // 获取纪念日列表
  async getAnniversaries(page: number, size: number) {
    const response = await api.get<GetAnniversariesResponse>('/anniversaries', {
      params: {
        page,
        size,
      },
    })
    return response.data
  },

  // 创建纪念日
  async createAnniversary(anniversaryData: Omit<Anniversary, 'id'>) {
    const response = await api.post<CreateOrUpdateAnniversaryResponse>(
      '/anniversaries',
      anniversaryData
    )
    return response.data
  },

  // 更新纪念日
  async updateAnniversary(id: number, anniversaryData: Partial<Anniversary>) {
    const response = await api.put<CreateOrUpdateAnniversaryResponse>(
      `/anniversaries/${id}`,
      anniversaryData
    )
    return response.data
  },

  // 删除纪念日
  async deleteAnniversary(id: number) {
    const response = await api.delete<DeleteAnniversaryResponse>(`/anniversaries/${id}`)
    return response.data
  },
}
