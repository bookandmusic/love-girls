import api from './api'
import type { FileInfo } from './upload'

// 定义认证相关的类型
export interface UserInfo {
  userName: string
  userId: number
  userEmail: string
}

export interface VerifyTokenResponse {
  code: number
  message: string
  data: UserInfo
}

interface LoginRequest {
  username: string
  password: string
}

interface LoginResponse {
  access_token: string
  token_type: string
  expires_in: number
}

// 定义用户类型
export interface User {
  id: number
  name: string
  email: string
  role: string
  joinDate: string
  avatar?: FileInfo
}

export interface UserFormData {
  newPassword?: string
  id: number
  name: string
  email: string
  role: string
  joinDate: string
  url: string
  thumbnailUrl: string
  avatarId: number
}

interface UpdateUserResponse {
  code: number
  data: User
  message: string
}

/**
 * 用户和认证相关API接口
 */
export const userApi = {
  // 认证相关接口
  /**
   * 用户登录
   */
  async login(loginData: LoginRequest): Promise<LoginResponse> {
    const response = await api.post<LoginResponse>('/user/token', loginData)
    return response.data
  },

  /**
   * 验证token是否有效
   */
  async verifyToken(token: string) {
    const response = await api.get<VerifyTokenResponse>('/user', {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
    return response.data
  },

  // 用户相关接口
  /**
   * 获取用户列表
   */
  async getUsers() {
    const response = await api.get<{ code: number; data: User[] }>('/users')
    return response.data
  },

  /**
   * 更新用户信息
   */
  updateUser(
    id: number,
    data: {
      name: string
      email: string
      avatarId: number
      newPassword?: string
      role: string
    }
  ) {
    return api.put<UpdateUserResponse>(`/users/${id}`, data)
  },
}
