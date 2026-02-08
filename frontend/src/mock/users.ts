// 用户相关的mock数据
import type { MockMethod } from 'vite-plugin-mock'

interface RequestParams {
  body: Record<string, unknown>
  headers: Record<string, string>
  url: string
}

// 定义用户类型
interface User {
  id: number
  name: string
  email: string
  role: string
  joinDate: string
  avatar: string
  avatarId: number
}

// 模拟用户数据
const mockUsers: User[] = [
  {
    id: 1,
    name: '小伍',
    email: 'xiaowu@example.com',
    role: '男朋友',
    joinDate: '2021-01-01',
    avatar: '',
    avatarId: 0,
  },
  {
    id: 2,
    name: '小陆',
    email: 'xiaolu@example.com',
    role: '女朋友',
    joinDate: '2021-01-01',
    avatar: '',
    avatarId: 0,
  },
]

const users: MockMethod[] = [
  // 登录接口
  {
    url: '/api/v1/user/token',
    method: 'post',
    response: (params: RequestParams) => {
      const { body } = params
      const { username, password } = body

      // 模拟验证逻辑
      if (username === 'admin' && password === '123456') {
        // 登录成功，返回token
        return {
          access_token: 'fake-jwt-token-for-admin',
          token_type: 'bearer',
          expires_in: 3600,
        }
      } else {
        // 登录失败
        return {
          code: 1,
          message: '用户名或密码错误',
          data: null,
        }
      }
    },
  },

  // 验证token接口 (使用获取用户信息接口)
  {
    url: '/api/v1/user',
    method: 'get',
    response: (params: RequestParams) => {
      // 从请求头获取token
      const headers = params.headers
      const authHeader = headers?.authorization
      const token = authHeader ? authHeader.replace('Bearer ', '') : ''

      if (token === 'fake-jwt-token-for-admin') {
        return {
          code: 0,
          message: '查询成功',
          data: {
            userName: 'admin',
            userId: 1,
            userEmail: 'admin@example.com',
          },
        }
      } else {
        return {
          code: 1,
          message: '未授权访问',
          data: null,
        }
      }
    },
  },

  // 获取用户列表
  {
    url: '/api/v1/users',
    method: 'get',
    response: () => {
      return {
        code: 0,
        message: '查询成功',
        data: mockUsers,
      }
    },
  },

  // 更新用户信息
  {
    url: '/api/v1/users/:id',
    method: 'put',
    response: (params: RequestParams) => {
      const urlParts = params.url.split('/')
      const userId = parseInt(urlParts[urlParts.length - 1] || '0')

      const userIndex = mockUsers.findIndex(u => u.id === userId)
      if (userIndex === -1) {
        return {
          code: 1,
          message: '用户不存在',
          data: null,
        }
      }

      // 更新用户信息
      const user = mockUsers[userIndex]!
      const updatedUser: User = {
        ...user,
        name: (params.body.name as string) || user.name,
        email: (params.body.email as string) || user.email,
        avatar: (params.body.avatar as string) || user.avatar,
        avatarId: (params.body.avatarId as number) || user.avatarId,
        role: (params.body.role as string) || user.role,
        joinDate: user.joinDate, // 保持原有的加入日期
        id: user.id, // 保持原有的id，不被更新
      }

      mockUsers[userIndex] = updatedUser

      // 如果提供了新密码，则更新密码
      if (params.body.newPassword) {
        // 这里可以模拟密码更新逻辑
        console.log(`用户 ${mockUsers[userIndex]!.name} 的密码已修改`)
      }

      return {
        code: 0,
        message: '更新成功',
        data: updatedUser,
      }
    },
  },
]

export default users
