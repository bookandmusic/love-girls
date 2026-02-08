import api from './api'
import type { FileInfo } from './upload'

export interface Photo {
  id: number
  albumId: number
  file?: FileInfo
  alt?: string
  createdAt: string
}

export interface Album {
  id: number
  name: string
  description: string
  coverImage?: Photo
  photoCount: number
  createdAt: string
}

interface GetAlbumsResponse {
  code: number
  data: {
    albums: Album[]
    totalPages: number
    total?: number
    totalCount?: number
    page?: number
    size?: number
  }
  message?: string
}

interface CreateOrUpdateAlbumResponse {
  code: number
  data: Album
  message?: string
}

interface DeleteAlbumResponse {
  code: number
  message?: string
}

interface GetPhotosResponse {
  code: number
  data: {
    photos: Photo[]
    totalPages: number
    total?: number
    totalCount?: number
    page?: number
    size?: number
  }
  message?: string
}

export const albumApi = {
  // 获取相册列表
  async getAlbums(page: number, size: number) {
    const response = await api.get<GetAlbumsResponse>('/albums', {
      params: {
        page,
        size,
      },
    })
    return response.data
  },

  // 创建相册
  async createAlbum(albumData: Omit<Album, 'id'>) {
    const response = await api.post<CreateOrUpdateAlbumResponse>('/albums', albumData)
    return response.data
  },

  // 更新相册
  async updateAlbum(id: number, albumData: Partial<Album>) {
    const response = await api.put<CreateOrUpdateAlbumResponse>(`/albums/${id}`, albumData)
    return response.data
  },

  // 删除相册
  async deleteAlbum(id: number) {
    const response = await api.delete<DeleteAlbumResponse>(`/albums/${id}`)
    return response.data
  },

  // 获取相册下的照片列表
  async getPhotos(albumId: number, page: number, size: number) {
    const response = await api.get<GetPhotosResponse>(`/albums/${albumId}/photos`, {
      params: {
        page,
        size,
      },
    })
    return response.data
  },

  // 设置相册封面
  async setCover(albumId: number, photoId: number) {
    const response = await api.put<CreateOrUpdateAlbumResponse>(`/albums/${albumId}/cover`, {
      photoId,
    })
    return response.data
  },

  // 添加照片到相册
  async addPhotos(albumId: number, photoIds: number[]) {
    const response = await api.post<{ code: number; data: Photo[]; message?: string }>(
      `/albums/${albumId}/photos`,
      { photoIds }
    )
    return response.data
  },

  // 从相册中删除照片
  async deletePhoto(albumId: number, photoId: number) {
    const response = await api.delete<DeleteAlbumResponse>(`/albums/${albumId}/photos/${photoId}`)
    return response.data
  },
}
