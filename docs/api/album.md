# Album API 文档

## 概述

Album API 提供相册管理功能，包括相册的创建、查询、更新和删除等操作，以及相册内照片的管理功能。

---

## 1. 获取相册列表

获取相册列表，支持分页查询。

### 请求信息

- **接口路径**: `GET /api/v1/albums`
- **Content-Type**: `application/json`
- **需要认证**: 否

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，从 1 开始，默认 1 |
| size | int | 否 | 每页数量，默认 10 |

### 请求示例（curl）

```bash
curl -X GET "http://localhost:8080/api/v1/albums?page=1&size=10"
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| msg | string | 响应消息 |
| data | object | 响应数据 |
| data.albums | array | 相册列表 |
| data.albums[].id | int | 相册 ID |
| data.albums[].name | string | 相册名称 |
| data.albums[].description | string | 相册描述 |
| data.albums[].coverImage | object | 封面图片 |
| data.albums[].coverImage.id | int | 图片 ID |
| data.albums[].coverImage.url | string | 原图 URL |
| data.albums[].coverImage.thumbnailUrl | string | 缩略图 URL |
| data.albums[].photoCount | int | 照片数量 |
| data.albums[].createdAt | string | 创建时间（ISO 8601） |
| data.totalPages | int | 总页数 |
| data.total | int | 总数量 |
| data.totalCount | int | 总记录数 |
| data.page | int | 当前页码 |
| data.size | int | 每页数量 |

### 响应示例

```json
{
  "code": 0,
  "msg": "查询成功",
  "data": {
    "albums": [
      {
        "id": 1,
        "name": "我们的回忆",
        "description": "记录我们的美好时光",
        "coverImage": {
          "id": 101,
          "albumId": 1,
          "url": "https://example.com/photos/cover.jpg",
          "thumbnailUrl": "https://example.com/photos/thumb_cover.jpg"
        },
        "photoCount": 25,
        "createdAt": "2024-01-15T10:30:00Z"
      }
    ],
    "totalPages": 5,
    "total": 50,
    "totalCount": 50,
    "page": 1,
    "size": 10
  }
}
```

---



## 2. 创建相册

创建一个新的相册。

### 请求信息

- **接口路径**: `POST /api/v1/albums`
- **Content-Type**: `application/json`
- **需要认证**: 是

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 相册名称 |
| description | string | 否 | 相册描述 |

### 请求示例

```json
{
  "name": "我们的回忆",
  "description": "记录我们的美好时光"
}
```

### 请求示例（curl）

```bash
curl -X POST "http://localhost:8080/api/v1/albums" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "name": "我们的回忆",
    "description": "记录我们的美好时光"
  }'
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| msg | string | 响应消息 |
| data | object | 响应数据 |
| data.id | int | 相册 ID |
| data.name | string | 相册名称 |
| data.description | string | 相册描述 |
| data.photoCount | int | 照片数量 |
| data.createdAt | string | 创建时间（ISO 8601） |

### 响应示例

```json
{
  "code": 0,
  "msg": "创建成功",
  "data": {
    "id": 1,
    "name": "我们的回忆",
    "description": "记录我们的美好时光",
    "photoCount": 0,
    "createdAt": "2024-01-15T10:30:00Z"
  }
}
```

---

## 3. 更新相册

更新相册的信息。

### 请求信息

- **接口路径**: `PUT /api/v1/albums/:id`
- **Content-Type**: `application/json`
- **需要认证**: 是

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 相册 ID |

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 否 | 相册名称 |
| description | string | 否 | 相册描述 |

### 请求示例

```json
{
  "name": "我们的回忆（更新版）",
  "description": "记录我们的美好时光，每一天"
}
```

### 请求示例（curl）

```bash
curl -X PUT "http://localhost:8080/api/v1/albums/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "name": "我们的回忆（更新版）",
    "description": "记录我们的美好时光，每一天"
  }'
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| msg | string | 响应消息 |
| data | object | 响应数据 |
| data.id | int | 相册 ID |
| data.name | string | 相册名称 |
| data.description | string | 相册描述 |
| data.photoCount | int | 照片数量 |
| data.createdAt | string | 创建时间（ISO 8601） |

### 响应示例

```json
{
  "code": 0,
  "msg": "更新成功",
  "data": {
    "id": 1,
    "name": "我们的回忆（更新版）",
    "description": "记录我们的美好时光，每一天",
    "photoCount": 25,
    "createdAt": "2024-01-15T10:30:00Z"
  }
}
```

---

## 4. 删除相册

删除指定的相册。

### 请求信息

- **接口路径**: `DELETE /api/v1/albums/:id`
- **Content-Type**: `application/json`
- **需要认证**: 是

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 相册 ID |

### 请求示例（curl）

```bash
curl -X DELETE "http://localhost:8080/api/v1/albums/1" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| msg | string | 响应消息 |
| data | null | 响应数据，删除操作无数据返回 |

### 响应示例

```json
{
  "code": 0,
  "msg": "删除成功",
  "data": null
}
```

---

## 5. 获取相册中的照片列表

获取指定相册的照片列表，支持分页查询。

### 请求信息

- **接口路径**: `GET /api/v1/albums/:id/photos`
- **Content-Type**: `application/json`
- **需要认证**: 否

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 相册 ID |

### 查询参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 是 | 页码，从 1 开始 |
| size | int | 是 | 每页数量 |

### 请求示例（curl）

```bash
curl -X GET "http://localhost:8080/api/v1/albums/1/photos?page=1&size=10" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| msg | string | 响应消息 |
| data | object | 响应数据 |
| data.photos | array | 照片列表 |
| data.photos[].id | int | 照片 ID |
| data.photos[].albumId | int | 所属相册 ID |
| data.photos[].url | string | 原图 URL |
| data.photos[].thumbnailUrl | string | 缩略图 URL |
| data.photos[].alt | string | 图片描述（可选） |
| data.photos[].createdAt | string | 创建时间（ISO 8601） |
| data.totalPages | int | 总页数 |
| data.total | int | 总数量 |
| data.totalCount | int | 总记录数 |
| data.page | int | 当前页码 |
| data.size | int | 每页数量 |

### 响应示例

```json
{
  "code": 0,
  "msg": "查询成功",
  "data": {
    "photos": [
      {
        "id": 101,
        "albumId": 1,
        "url": "https://example.com/photos/photo1.jpg",
        "thumbnailUrl": "https://example.com/photos/thumb_photo1.jpg",
        "alt": "美好时光",
        "createdAt": "2024-01-15T10:30:00Z"
      },
      {
        "id": 102,
        "albumId": 1,
        "url": "https://example.com/photos/photo2.jpg",
        "thumbnailUrl": "https://example.com/photos/thumb_photo2.jpg",
        "alt": "快乐时光",
        "createdAt": "2024-01-16T11:00:00Z"
      }
    ],
    "totalPages": 3,
    "total": 25,
    "totalCount": 25,
    "page": 1,
    "size": 10
  }
}
```

---

## 6. 添加照片到相册

将已上传的照片添加到相册。

### 请求信息

- **接口路径**: `POST /api/v1/albums/:id/photos`
- **Content-Type**: `application/json`
- **需要认证**: 是

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 相册 ID |

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| photoIds | array | 是 | 照片 ID 列表 |

### 请求示例

```json
{
  "photoIds": [101, 102, 103]
}
```

### 请求示例（curl）

```bash
curl -X POST "http://localhost:8080/api/v1/albums/1/photos" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "photoIds": [101, 102, 103]
  }'
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| msg | string | 响应消息 |
| data | array | 添加成功的照片列表 |
| data[].id | int | 照片 ID |
| data[].albumId | int | 相册 ID |
| data[].url | string | 原图 URL |
| data[].thumbnailUrl | string | 缩略图 URL |
| data[].createdAt | string | 创建时间 |

### 响应示例

```json
{
  "code": 0,
  "msg": "添加成功",
  "data": [
    {
      "id": 101,
      "albumId": 1,
      "url": "https://example.com/photos/photo1.jpg",
      "thumbnailUrl": "https://example.com/photos/thumb_photo1.jpg",
      "createdAt": "2024-01-15T10:30:00Z"
    }
  ]
}
```

---

## 7. 设置相册封面

设置相册的封面图片。

### 请求信息

- **接口路径**: `PUT /api/v1/albums/:id/cover`
- **Content-Type**: `application/json`
- **需要认证**: 是

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 相册 ID |

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| photoId | int | 是 | 照片 ID |

### 请求示例

```json
{
  "photoId": 101
}
```

### 请求示例（curl）

```bash
curl -X PUT "http://localhost:8080/api/v1/albums/1/cover" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "photoId": 101
  }'
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| msg | string | 响应消息 |
| data | object | 相册对象 |

### 响应示例

```json
{
  "code": 0,
  "msg": "设置成功",
  "data": {
    "id": 1,
    "name": "我们的回忆",
    "coverImage": {
      "id": 101,
      "url": "https://example.com/photos/cover.jpg"
    }
  }
}
```

---

## 8. 从相册删除照片

从相册中移除照片（注意：这通常是解除关联或删除照片，具体取决于后端实现，但前端语义为移除）。

### 请求信息

- **接口路径**: `DELETE /api/v1/albums/:id/photos/:photoId`
- **Content-Type**: `application/json`
- **需要认证**: 是

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 相册 ID |
| photoId | int | 是 | 照片 ID |

### 请求示例（curl）

```bash
curl -X DELETE "http://localhost:8080/api/v1/albums/1/photos/101" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| msg | string | 响应消息 |
| data | null | 无返回数据 |

### 响应示例

```json
{
  "code": 0,
  "msg": "删除成功",
  "data": null
}
```

---

## 注意事项

1. **权限控制**: 相册列表查询和照片列表查询不需要认证；创建、更新、删除相册以及相册内照片管理操作需要认证，请求头必须包含有效的 Authorization token。
2. **分页参数**: 获取相册列表和照片列表时，page 和 size 参数为可选，不传时使用默认值（page=1, size=10）。
2. **删除注意**: 删除相册会同时删除相册中的所有照片，操作不可恢复。
3. **封面图片**: 相册创建时没有封面图片，需要通过上传照片后手动设置。
4. **照片添加**: 照片需要先通过 Upload API 上传，然后通过 ID 添加到相册。

---

## 版本历史

| 版本 | 日期 | 说明 |
|------|------|------|
| 1.0.0 | 2026-01-28 | 初始版本，支持相册的增删改查 |
| 1.1.0 | 2026-01-30 | 更新：调整照片添加和删除接口，新增设置封面接口 |