import api from './api'

// 统一的文件信息接口
export interface FileInfo {
  id: number
  url: string
  thumbnail: string
  name?: string
  size?: number
  mime_type?: string
}

export interface UploadResponse {
  file: FileInfo
}
// 用户相关API接口
export const uploadApi = {
  /**
   * 上传图片
   */
  uploadImage(avatar: FormData) {
    return api.post<{ code: number; data: UploadResponse; message: string }>(
      `/file/upload`,
      avatar,
      {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      }
    )
  },
}
