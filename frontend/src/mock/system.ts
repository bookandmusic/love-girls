// 系统初始化相关接口
import type { MockMethod } from 'vite-plugin-mock'

// 使用sessionStorage模拟后端存储
let systemInitialized = false
let systemData: Record<string, unknown> | null = null

const systemInit: MockMethod[] = [
  {
    url: '/api/v1/system/init',
    method: 'post',
    // 模拟初始化请求
    response: (request: { body: Record<string, unknown> }) => {
      try {
        // 解析请求数据
        const formData = request.body

        // 检查formData是否为FormData对象，否则直接使用request.body
        let getFormDataValue: (key: string) => unknown
        if (
          formData &&
          typeof (formData as { get?: (...args: unknown[]) => unknown }).get === 'function'
        ) {
          // 类型断言，确保formData是FormData类型
          const formDataAsFormData = formData as unknown as FormData
          getFormDataValue = (key: string) => formDataAsFormData.get(key)
        } else {
          // 如果不是FormData对象，直接从对象中取值
          getFormDataValue = (key: string) => formData[key]
        }

        if (formData) {
          // 模拟处理初始化数据
          systemData = {
            siteInfo: {
              siteName: getFormDataValue('siteName'),
              siteDescription: getFormDataValue('siteDescription') || '',
              startDate: getFormDataValue('startDate'),
            },
            users: [
              {
                id: 1,
                name: (getFormDataValue('userAName') as string) || '用户A',
                role: (getFormDataValue('userARole') as string) || '角色A',
                email: (getFormDataValue('userAEmail') as string) || '',
                phone: (getFormDataValue('userAPhone') as string) || '',
                avatar: getFormDataValue('avatarA') ? 'data:image/avatar-a' : null,
              },
              {
                id: 2,
                name: (getFormDataValue('userBName') as string) || '用户B',
                role: (getFormDataValue('userBRole') as string) || '角色B',
                email: (getFormDataValue('userBEmail') as string) || '',
                phone: (getFormDataValue('userBPhone') as string) || '',
                avatar: getFormDataValue('avatarB') ? 'data:image/avatar-b' : null,
              },
            ],
            password: (getFormDataValue('sitePassword') as string) || '',
          }

          // 设置系统为已初始化状态
          systemInitialized = true

          return {
            code: 200,
            msg: '项目初始化成功',
            initialized: systemInitialized,
            data: systemData,
          }
        }

        return { code: 1, msg: '初始化数据不能为空' }
      } catch (error) {
        console.error('初始化失败:', error)
        return { code: 1, msg: '初始化失败：' + (error as Error).message }
      }
    },
  },
  {
    url: '/api/v1/system/init',
    method: 'get',
    response: () => {
      return {
        code: 0,
        data: {
          initialized: true,
        },
      }
    },
  },
  {
    url: '/api/v1/system/info',
    method: 'get',
    response: () => {
      return {
        code: 0,
        data: {
          site: {
            name: '小伍与小陆',
            description: '始于心动，守于日常。往后余生，皆是彼此。',
            startDate: '2021-01-01',
          },
          couple: {
            boy: {
              name: '小伍',
              avatar: '',
            },
            girl: {
              name: '小陆',
              avatar: '',
            },
          },
        },
      }
    },
  },
]

export default [
  ...systemInit,
  {
    url: '/api/v1/system/settings/site',
    method: 'get',
    response: () => {
      return {
        code: 0,
        data: {
          siteTitle: '爱的纪念册',
          siteDescription: '记录我们美好时光的地方',
          siteFooter: '© 2023 爱的纪念册. 保留所有权利.',
        },
      }
    },
  },
  {
    url: '/api/v1/system/settings/site',
    method: 'post',
    response: () => {
      return {
        code: 0,
        message: '设置保存成功',
      }
    },
  },
]
