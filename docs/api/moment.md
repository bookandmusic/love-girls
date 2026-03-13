# Moment API 文档

## 概述

Moment API 提供动态（时光动态）管理功能，支持发布动态、点赞、设置公开状态等操作。

---

## 1. 获取动态列表

获取动态列表，支持分页、排序和过滤查询。

### 请求信息

- **接口路径**: `GET /api/v1/moments`
- **Content-Type**: `application/json`
- **需要认证**: 否（未登录用户只能查看公开动态，已登录用户可以查看自己的所有动态）

### 请求参数

#### 分页参数

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| page | int | 否 | 1 | 页码，从 1 开始 |
| size | int | 否 | 10 | 每页数量，最大 100 |

#### 排序参数

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| sort_by | string | 否 | created_at | 排序字段，可选值：`created_at`、`likes` |
| order | string | 否 | desc | 排序方向，可选值：`asc`、`desc` |

#### 过滤参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| filter | []string | 否 | 过滤条件，格式：`field:op:value`，可传多个 |

**支持的过滤字段和操作符**：

| 字段 | 支持的操作符 | 说明 |
|------|-------------|------|
| is_public | eq | 按公开状态过滤，值：`true`/`false` |
| user_id | eq | 按用户ID过滤 |
| likes | eq, gt, lt, gte, lte | 按点赞数过滤 |

**过滤示例**：
- `filter=is_public:eq:true` - 只查看公开动态
- `filter=likes:gt:10` - 点赞数大于10的动态
- `filter=likes:gte:5&filter=likes:lte:20` - 点赞数在5到20之间的动态

### 请求示例（curl）

```bash
# 基础分页查询
curl -X GET "http://localhost:8182/api/v1/moments?page=1&size=10"

# 按点赞数降序排序
curl -X GET "http://localhost:8182/api/v1/moments?page=1&size=10&sort_by=likes&order=desc"

# 只查看公开动态
curl -X GET "http://localhost:8182/api/v1/moments?page=1&size=10&filter=is_public:eq:true"

# 查看点赞数大于10的动态
curl -X GET "http://localhost:8182/api/v1/moments?page=1&size=10&sort_by=likes&order=desc&filter=likes:gt:10"
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| message | string | 响应消息 |
| data | object | 响应数据 |
| data.moments | array | 动态列表 |
| data.moments[].id | uint64 | 动态 ID |
| data.moments[].content | string | 动态内容 |
| data.moments[].images | array | 图片列表 |
| data.moments[].images[].id | uint64 | 图片 ID |
| data.moments[].images[].momentId | uint64 | 动态 ID |
| data.moments[].images[].file | object | 文件信息 |
| data.moments[].images[].file.id | uint64 | 文件 ID |
| data.moments[].images[].file.url | string | 文件访问 URL |
| data.moments[].images[].file.thumbnailUrl | string | 缩略图 URL |
| data.moments[].likes | int | 点赞数 |
| data.moments[].createdAt | string | 创建时间 |
| data.moments[].author | object | 作者信息 |
| data.moments[].author.name | string | 作者姓名 |
| data.moments[].author.avatar | object | 作者头像文件信息 |
| data.moments[].isPublic | boolean | 是否公开 |
| data.page | int | 当前页码 |
| data.size | int | 每页数量 |
| data.total | int64 | 总数量 |
| data.totalPages | int | 总页数 |

### 成功响应示例

```json
{
  "code": 0,
  "message": "查询成功",
  "data": {
    "moments": [
      {
        "id": 1,
        "content": "今天是个美好的日子",
        "images": [
          {
            "id": 101,
            "momentId": 1,
            "file": {
              "id": 101,
              "url": "http://localhost:8182/api/v1/file/101",
              "thumbnailUrl": "http://localhost:8182/api/v1/file/101?w=200&h=200"
            }
          }
        ],
        "likes": 5,
        "createdAt": "2024-01-20 14:30:05",
        "author": {
          "name": "鹿",
          "avatar": {
            "id": 10,
            "url": "http://localhost:8182/api/v1/file/10",
            "thumbnailUrl": "http://localhost:8182/api/v1/file/10?w=200&h=200"
          }
        },
        "isPublic": true
      }
    ],
    "page": 1,
    "size": 10,
    "total": 15,
    "totalPages": 2
  }
}
```

### 错误响应

**400 Bad Request - 参数错误**
```json
{
  "code": 1,
  "message": "无效的页码参数",
  "data": null
}
```

**500 Internal Server Error - 查询失败**
```json
{
  "code": 1,
  "message": "系统内部错误",
  "data": null
}
```

---

## 2. 创建动态

创建一个新的动态。

### 请求信息

- **接口路径**: `POST /api/v1/moments`
- **Content-Type**: `application/json`
- **需要认证**: 是

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| content | string | 是 | 动态内容 |
| imageIds | array | 否 | 图片ID列表 |
| isPublic | boolean | 否 | 是否公开，默认 true |
| userId | uint64 | 是 | 用户ID |

### 请求示例

```json
{
  "content": "今天是个美好的日子",
  "imageIds": [1001],
  "isPublic": true,
  "userId": 1
}
```

### 请求示例（curl）

```bash
curl -X POST "http://localhost:8080/api/v1/moments" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "content": "今天是个美好的日子",
    "imageIds": [1001],
    "isPublic": true,
    "userId": 1
  }'
```

### 成功响应示例

```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "id": 1,
    "content": "今天是个美好的日子",
    "images": [
      {
        "id": 101,
        "momentId": 1,
        "url": "/uploads/moment1.jpg",
        "thumbnailUrl": "/uploads/thumbnail/moment1_thumb.jpg"
      }
    ],
    "likes": 0,
    "createdAt": "2024-01-20 14:30:05",
    "author": {
      "name": "鹿",
      "avatar": "/avatars/user1.jpg"
    },
    "isPublic": true
  }
}
```

### 错误响应

**400 Bad Request - 参数校验失败**
```json
{
  "code": 1,
  "message": "参数校验失败",
  "data": null
}
```

**500 Internal Server Error - 创建失败**
```json
{
  "code": 1,
  "message": "系统内部错误",
  "data": null
}
```

---

## 3. 更新动态

更新动态的信息。

### 请求信息

- **接口路径**: `PUT /api/v1/moments/:id`
- **Content-Type**: `application/json`
- **需要认证**: 是

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | uint64 | 是 | 动态 ID |

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| content | string | 否 | 动态内容 |
| imageIds | array | 否 | 图片ID列表 |
| isPublic | boolean | 否 | 是否公开 |

### 请求示例

```json
{
  "content": "今天是个美好的日子，值得铭记",
  "isPublic": false
}
```

### 请求示例（curl）

```bash
curl -X PUT "http://localhost:8080/api/v1/moments/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "content": "今天是个美好的日子，值得铭记",
    "isPublic": false
  }'
```

### 成功响应示例

```json
{
  "code": 0,
  "message": "更新成功",
  "data": {
    "id": 1,
    "content": "今天是个美好的日子，值得铭记",
    "images": [
      {
        "id": 101,
        "momentId": 1,
        "url": "/uploads/moment1.jpg",
        "thumbnailUrl": "/uploads/thumbnail/moment1_thumb.jpg"
      }
    ],
    "likes": 0,
    "createdAt": "2024-01-20 14:30:05",
    "author": {
      "name": "鹿",
      "avatar": "/avatars/user1.jpg"
    },
    "isPublic": false
  }
}
```

### 错误响应

**400 Bad Request - 参数校验失败**
```json
{
  "code": 1,
  "message": "参数校验失败",
  "data": null
}
```

**400 Bad Request - 无效ID**
```json
{
  "code": 1,
  "message": "无效的动态ID",
  "data": null
}
```

**404 Not Found - 动态不存在**
```json
{
  "code": 1,
  "message": "动态不存在",
  "data": null
}
```

**500 Internal Server Error - 更新失败**
```json
{
  "code": 1,
  "message": "系统内部错误",
  "data": null
}
```

---

## 4. 更新动态公开状态

更新动态的公开状态。

### 请求信息

- **接口路径**: `PUT /api/v1/moments/:id/public`
- **Content-Type**: `application/json`
- **需要认证**: 是

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | uint64 | 是 | 动态 ID |

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| isPublic | boolean | 是 | 是否公开 |

### 请求示例

```json
{
  "isPublic": false
}
```

### 请求示例（curl）

```bash
curl -X PUT "http://localhost:8080/api/v1/moments/1/public" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "isPublic": false
  }'
```

### 成功响应示例

```json
{
  "code": 0,
  "message": "更新成功",
  "data": {
    "id": 1,
    "content": "今天是个美好的日子",
    "images": [],
    "likes": 0,
    "createdAt": "2024-01-20 14:30:05",
    "author": {
      "name": "鹿",
      "avatar": "/avatars/user1.jpg"
    },
    "isPublic": false
  }
}
```

### 错误响应

**400 Bad Request - 参数校验失败**
```json
{
  "code": 1,
  "message": "参数校验失败",
  "data": null
}
```

**400 Bad Request - 无效ID**
```json
{
  "code": 1,
  "message": "无效的动态ID",
  "data": null
}
```

**404 Not Found - 动态不存在**
```json
{
  "code": 1,
  "message": "动态不存在",
  "data": null
}
```

**500 Internal Server Error - 更新失败**
```json
{
  "code": 1,
  "message": "系统内部错误",
  "data": null
}
```

---

## 5. 点赞动态

为动态添加点赞。

### 请求信息

- **接口路径**: `POST /api/v1/moments/:id/like`
- **Content-Type**: `application/json`
- **需要认证**: 是

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | uint64 | 是 | 动态 ID |

### 请求示例（curl）

```bash
curl -X POST "http://localhost:8080/api/v1/moments/1/like" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 成功响应示例

```json
{
  "code": 0,
  "message": "点赞成功",
  "data": {
    "likes": 6
  }
}
```

### 错误响应

**400 Bad Request - 无效ID**
```json
{
  "code": 1,
  "message": "无效的动态ID",
  "data": null
}
```

**404 Not Found - 动态不存在**
```json
{
  "code": 1,
  "message": "动态不存在",
  "data": null
}
```

**500 Internal Server Error - 点赞失败**
```json
{
  "code": 1,
  "message": "系统内部错误",
  "data": null
}
```

---

## 6. 删除动态

删除指定的动态。

### 请求信息

- **接口路径**: `DELETE /api/v1/moments/:id`
- **Content-Type**: `application/json`
- **需要认证**: 是

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | uint64 | 是 | 动态 ID |

### 请求示例（curl）

```bash
curl -X DELETE "http://localhost:8080/api/v1/moments/1" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 成功响应示例

```json
{
  "code": 0,
  "message": "删除成功",
  "data": null
}
```

### 错误响应

**400 Bad Request - 无效ID**
```json
{
  "code": 1,
  "message": "无效的动态ID",
  "data": null
}
```

**404 Not Found - 动态不存在**
```json
{
  "code": 1,
  "message": "动态不存在",
  "data": null
}
```

**500 Internal Server Error - 删除失败**
```json
{
  "code": 1,
  "message": "系统内部错误",
  "data": null
}
```

---

## 注意事项

1. **权限控制**: 创建、更新、删除、点赞和更改公开状态需要认证。
2. **公开状态**: 私有动态仅作者和授权用户可见，公开动态所有人可见。
3. **点赞限制**: 同一用户对同一动态可以点赞多次，每次调用都会增加点赞数。
4. **图片上传**: 动态创建时可以附带图片，图片需要先通过上传接口上传。
5. **删除注意**: 删除动态会同时删除关联的图片，操作不可恢复。

---

## 版本历史

| 版本 | 日期 | 说明 |
|------|------|------|
| 1.1.0 | 2026-03-13 | 新增：列表接口支持排序和过滤功能，新增 sort_by、order、filter 参数 |
| 1.0.0 | 2026-01-31 | 完善版本，支持动态的增删改查、点赞和公开状态设置 |