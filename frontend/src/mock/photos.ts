// 照片相关接口
import type { MockMethod } from 'vite-plugin-mock'

// 模拟照片数据存储
const albumPhotosStore: Record<
  number,
  Array<{
    id: number
    albumId: number
    url: string
    thumbnailUrl: string
    description: string
    createdAt: string
  }>
> = {}

// 生成相册照片数据
function generateAlbumPhotos(albumId: number) {
  if (albumPhotosStore[albumId]) {
    return albumPhotosStore[albumId]
  }

  let photos = []
  switch (albumId) {
    case 1: // 第一次旅行
      photos = Array.from({ length: 12 }, (_, i) => ({
        id: albumId * 100 + i + 1,
        albumId: albumId,
        url: `https://picsum.photos/800/600?random=${i + 100}`,
        thumbnailUrl: `https://picsum.photos/300/200?random=${i + 100}`,
        description: `旅行中的美好时光 ${i + 1}`,
        createdAt: new Date(2023, 4, 1 + i).toISOString(),
      }))
      break

    case 2: // 浪漫的夜晚
      photos = Array.from({ length: 8 }, (_, i) => ({
        id: albumId * 100 + i + 1,
        albumId: albumId,
        url: `https://picsum.photos/800/600?random=${i + 200}`,
        thumbnailUrl: `https://picsum.photos/300/200?random=${i + 200}`,
        description: `浪漫时刻 ${i + 1}`,
        createdAt: new Date(2023, 7, 15 + i).toISOString(),
      }))
      break

    case 3: // 春天的约会
      photos = Array.from({ length: 15 }, (_, i) => ({
        id: albumId * 100 + i + 1,
        albumId: albumId,
        url: `https://picsum.photos/800/600?random=${i + 300}`,
        thumbnailUrl: `https://picsum.photos/300/200?random=${i + 300}`,
        description: `春日踏青 ${i + 1}`,
        createdAt: new Date(2024, 2, 20 + i).toISOString(),
      }))
      break

    case 4: // 美食之旅
      photos = Array.from({ length: 22 }, (_, i) => ({
        id: albumId * 100 + i + 1,
        albumId: albumId,
        url: `https://picsum.photos/800/600?random=${i + 400}`,
        thumbnailUrl: `https://picsum.photos/300/200?random=${i + 400}`,
        description: `美味食物 ${i + 1}`,
        createdAt: new Date(2024, 4, 1 + i).toISOString(),
      }))
      break

    case 5: // 历史探索
      photos = Array.from({ length: 18 }, (_, i) => ({
        id: albumId * 100 + i + 1,
        albumId: albumId,
        url: `https://picsum.photos/800/600?random=${i + 500}`,
        thumbnailUrl: `https://picsum.photos/300/200?random=${i + 500}`,
        description: `历史文物 ${i + 1}`,
        createdAt: new Date(2024, 6, 15 + i).toISOString(),
      }))
      break

    case 6: // 海边度假
      photos = Array.from({ length: 25 }, (_, i) => ({
        id: albumId * 100 + i + 1,
        albumId: albumId,
        url: `https://picsum.photos/800/600?random=${i + 600}`,
        thumbnailUrl: `https://picsum.photos/300/200?random=${i + 600}`,
        description: `海滩美景 ${i + 1}`,
        createdAt: new Date(2024, 7, 10 + i).toISOString(),
      }))
      break

    case 7: // 雪山之旅
      photos = Array.from({ length: 14 }, (_, i) => ({
        id: albumId * 100 + i + 1,
        albumId: albumId,
        url: `https://picsum.photos/800/600?random=${i + 700}`,
        thumbnailUrl: `https://picsum.photos/300/200?random=${i + 700}`,
        description: `雪山风光 ${i + 1}`,
        createdAt: new Date(2024, 11, 1 + i).toISOString(),
      }))
      break

    case 8: // 城市漫步
      photos = Array.from({ length: 10 }, (_, i) => ({
        id: albumId * 100 + i + 1,
        albumId: albumId,
        url: `https://picsum.photos/800/600?random=${i + 800}`,
        thumbnailUrl: `https://picsum.photos/300/200?random=${i + 800}`,
        description: `城市街景 ${i + 1}`,
        createdAt: new Date(2024, 11, 10 + i).toISOString(),
      }))
      break

    default: // 默认相册
      photos = Array.from({ length: 15 }, (_, i) => ({
        id: albumId * 100 + i + 1,
        albumId: albumId,
        url: `https://picsum.photos/800/600?random=${i + 900}`,
        thumbnailUrl: `https://picsum.photos/300/200?random=${i + 900}`,
        description: `照片 ${i + 1}`,
        createdAt: new Date().toISOString(),
      }))
  }

  albumPhotosStore[albumId] = photos
  return photos
}

const photos: MockMethod[] = [
  // 获取相册照片列表
  {
    url: '/api/v1/albums/:id/photos',
    method: 'get',
    response: (params: { query: Record<string, string>; url: string }) => {
      const query = params.query
      // 从URL中提取albumId
      const urlMatch = params.url.match(/\/albums\/(\d+)\/photos/)
      const albumId = urlMatch && urlMatch[1] ? parseInt(urlMatch[1]) : 0
      const page = parseInt(query?.page || '1')
      const size = parseInt(query?.size || '9')

      // 根据不同相册ID生成不同的照片数据
      const allPhotos = generateAlbumPhotos(albumId)

      // 计算分页数据
      const startIndex = (page - 1) * size
      const endIndex = startIndex + size
      const paginatedPhotos = allPhotos.slice(startIndex, endIndex)

      return {
        code: 0,
        data: {
          photos: paginatedPhotos,
          totalPages: Math.ceil(allPhotos.length / size),
          total: allPhotos.length,
          currentPage: page,
        },
        msg: '获取成功',
      }
    },
  },

  // 设置相册封面
  {
    url: '/api/v1/albums/:id/cover',
    method: 'put',
    response: (params: { body: { photoId: number }; url: string }) => {
      const urlMatch = params.url.match(/\/albums\/(\d+)\/cover/)
      const albumId = urlMatch && urlMatch[1] ? parseInt(urlMatch[1]) : 0
      const { photoId } = params.body

      const allPhotos = generateAlbumPhotos(albumId)
      const coverPhoto = allPhotos.find(p => p.id === photoId)

      return {
        code: 0,
        data: {
          id: albumId,
          name: `相册 ${albumId}`,
          description: '相册描述',
          coverImage: coverPhoto || null,
          photoCount: allPhotos.length,
          createdAt: new Date().toISOString(),
        },
        msg: '封面设置成功',
      }
    },
  },

  // 添加照片到相册
  {
    url: '/api/v1/albums/:id/photos',
    method: 'post',
    response: (params: { body: { photoIds: number[] }; url: string }) => {
      const urlMatch = params.url.match(/\/albums\/(\d+)\/photos/)
      const albumId = urlMatch && urlMatch[1] ? parseInt(urlMatch[1]) : 0
      const { photoIds } = params.body

      // 模拟添加照片
      const newPhotos = photoIds.map((id, index) => ({
        id: id,
        albumId: albumId,
        url: `https://picsum.photos/800/600?random=${Date.now() + index}`,
        thumbnailUrl: `https://picsum.photos/300/200?random=${Date.now() + index}`,
        description: `新添加的照片`,
        createdAt: new Date().toISOString(),
      }))

      // 添加到存储
      if (!albumPhotosStore[albumId]) {
        generateAlbumPhotos(albumId)
      }
      // Ensure array exists before spreading
      const existingPhotos = albumPhotosStore[albumId] || []
      albumPhotosStore[albumId] = [...existingPhotos, ...newPhotos]

      return {
        code: 0,
        data: newPhotos,
        msg: '照片添加成功',
      }
    },
  },

  // 从相册中删除照片
  {
    url: '/api/v1/albums/:albumId/photos/:photoId',
    method: 'delete',
    response: (params: { url: string }) => {
      const urlMatch = params.url.match(/\/albums\/(\d+)\/photos\/(\d+)/)
      const albumId = urlMatch && urlMatch[1] ? parseInt(urlMatch[1]) : 0
      const photoId = urlMatch && urlMatch[2] ? parseInt(urlMatch[2]) : 0

      // 从存储中删除照片
      if (albumPhotosStore[albumId]) {
        albumPhotosStore[albumId] = albumPhotosStore[albumId].filter(p => p.id !== photoId)
      }

      return {
        code: 0,
        msg: '照片删除成功',
      }
    },
  },
]

export default photos
