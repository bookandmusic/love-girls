// 文件上传相关接口
import type { MockMethod } from 'vite-plugin-mock'

const upload: MockMethod[] = [
  // 通用文件上传接口
  {
    url: '/api/v1/upload',
    method: 'post',
    response: (params: { body: Record<string, unknown>; query: Record<string, string> }) => {
      // 模拟上传处理，实际mock环境中无法真正接收文件
      // 我们假设文件上传成功并生成一个随机的URL

      const fileType = params.query.type || 'image' // 默认为图片类型
      const timestamp = Date.now()
      const randomNum = Math.floor(Math.random() * 10000)

      let fileUrl = ''
      let fileName = ''

      switch (fileType) {
        case 'image':
          fileUrl = `https://picsum.photos/800/600?random=${timestamp}${randomNum}`
          fileName = `image_${timestamp}.jpg`
          break
        case 'avatar':
          fileUrl = `https://picsum.photos/200/200?random=${timestamp}${randomNum}`
          fileName = `avatar_${timestamp}.jpg`
          break
        default:
          fileUrl = `https://picsum.photos/800/600?random=${timestamp}${randomNum}`
          fileName = `file_${timestamp}.${fileType || 'jpg'}`
      }

      return {
        code: 0,
        data: {
          url: fileUrl,
          fileName: fileName,
          size: Math.floor(Math.random() * 1000000) + 1024, // 随机大小 1KB - 1MB
          type: fileType,
        },
        msg: '上传成功',
      }
    },
  },

  // 图片上传专用接口
  {
    url: '/api/v1/upload/image',
    method: 'post',
    response: () => {
      const timestamp = Date.now()
      const randomNum = Math.floor(Math.random() * 100000)

      return {
        code: 0,
        data: {
          type: 'image',
          id: timestamp,
          url: `https://picsum.photos/800/600?random=${timestamp}${randomNum}`,
          thumbnailUrl: `https://picsum.photos/300/200?random=${timestamp}${randomNum}`,
          fileName: `image_${timestamp}.jpg`,
          size: Math.floor(Math.random() * 500000) + 1024, // 随机大小 1KB - 500KB
        },
        msg: '图片上传成功',
      }
    },
  },
]

export default upload
