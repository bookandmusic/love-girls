# Place API 文档

## 概述

Place API 提供地点管理功能，用于记录情侣共同去过的地方，包含地理位置、图片和日期信息。

---

## 1. 获取地点列表

获取地点列表，支持分页、排序和过滤查询。

### 请求信息

- **接口路径**: `GET /api/v1/places`
- **Content-Type**: `application/json`
- **需要认证**: 否

### 请求参数

#### 分页参数

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| page | int | 否 | 1 | 页码，从 1 开始 |
| size | int | 否 | 10 | 每页数量，最大 100 |

#### 排序参数

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| sort_by | string | 否 | created_at | 排序字段，可选值：`created_at`、`name` |
| order | string | 否 | desc | 排序方向，可选值：`asc`、`desc` |

#### 过滤参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| filter | []string | 否 | 过滤条件，格式：`field:op:value`，可传多个 |

**支持的过滤字段和操作符**：

| 字段 | 支持的操作符 | 说明 |
|------|-------------|------|
| name | like | 按地点名称模糊搜索 |

**过滤示例**：
- `filter=name:like:北京` - 搜索名称包含"北京"的地点

### 请求示例（curl）

```bash
# 基础分页查询
curl -X GET "http://localhost:8182/api/v1/places?page=1&size=10"

# 按名称升序排序
curl -X GET "http://localhost:8182/api/v1/places?page=1&size=10&sort_by=name&order=asc"

# 搜索名称包含"西湖"的地点
curl -X GET "http://localhost:8182/api/v1/places?page=1&size=10&filter=name:like:西湖"
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| msg | string | 响应消息 |
| data | object | 响应数据 |
| data.places | array | 地点列表 |
| data.places[].id | int | 地点 ID |
| data.places[].name | string | 地点名称 |
| data.places[].latitude | number | 纬度 |
| data.places[].longitude | number | 经度 |
| data.places[].image | object | 地点图片信息 |
| data.places[].image.id | int | 图片 ID |
| data.places[].image.placeId | int | 地点 ID |
| data.places[].image.file | object | 文件信息 |
| data.places[].description | string | 地点描述 |
| data.places[].date | string | 日期（YYYY-MM-DD） |
| data.page | int | 当前页码 |
| data.size | int | 每页数量 |
| data.total | int64 | 总数量 |
| data.totalPages | int | 总页数 |

### 响应示例

```json
{
  "code": 0,
  "message": "查询成功",
  "data": {
    "places": [
      {
        "id": 1,
        "name": "西湖",
        "latitude": 30.2592,
        "longitude": 120.1302,
        "image": {
          "id": 201,
          "placeId": 1,
          "file": {
            "id": 201,
            "url": "http://localhost:8182/api/v1/file/201",
            "thumbnailUrl": "http://localhost:8182/api/v1/file/201?w=200&h=200"
          }
        },
        "description": "第一次一起看湖",
        "date": "2024-01-15"
      }
    ],
    "page": 1,
    "size": 10,
    "total": 1,
    "totalPages": 1
  }
}
```

---



## 2. 创建地点

创建一个新的地点。

### 请求信息

- **接口路径**: `POST /api/v1/places`
- **Content-Type**: `application/json`
- **需要认证**: 是

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 地点名称 |
| latitude | number | 是 | 纬度 |
| longitude | number | 是 | 经度 |
| description | string | 否 | 地点描述 |
| date | string | 是 | 日期，格式：YYYY-MM-DD |
| image | object | 否 | 地点图片 |
| image.url | string | 是 | 图片 URL |
| image.thumbnailUrl | string | 是 | 缩略图 URL |

### 请求示例

```json
{
  "name": "西湖",
  "latitude": 30.2592,
  "longitude": 120.1302,
  "description": "第一次一起看湖",
  "date": "2024-01-15",
  "image": {
    "url": "https://example.com/photos/place1.jpg",
    "thumbnailUrl": "https://example.com/photos/thumb_place1.jpg"
  }
}
```

### 请求示例（curl）

```bash
curl -X POST "http://localhost:8080/api/v1/places" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "name": "西湖",
    "latitude": 30.2592,
    "longitude": 120.1302,
    "description": "第一次一起看湖",
    "date": "2024-01-15",
    "image": {
      "url": "https://example.com/photos/place1.jpg",
      "thumbnailUrl": "https://example.com/photos/thumb_place1.jpg"
    }
  }'
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| msg | string | 响应消息 |
| data | object | 响应数据 |
| data.id | int | 地点 ID |
| data.name | string | 地点名称 |
| data.latitude | number | 纬度 |
| data.longitude | number | 经度 |
| data.description | string | 地点描述 |
| data.date | string | 日期（YYYY-MM-DD） |

### 响应示例

```json
{
  "code": 0,
  "msg": "创建成功",
  "data": {
    "id": 1,
    "name": "西湖",
    "latitude": 30.2592,
    "longitude": 120.1302,
    "description": "第一次一起看湖",
    "date": "2024-01-15"
  }
}
```

---

## 3. 更新地点

更新地点的信息。

### 请求信息

- **接口路径**: `PUT /api/v1/places/:id`
- **Content-Type**: `application/json`
- **需要认证**: 是

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 地点 ID |

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 否 | 地点名称 |
| latitude | number | 否 | 纬度 |
| longitude | number | 否 | 经度 |
| description | string | 否 | 地点描述 |
| date | string | 否 | 日期，格式：YYYY-MM-DD |
| image | object | 否 | 地点图片 |
| image.url | string | 否 | 图片 URL |
| image.thumbnailUrl | string | 否 | 缩略图 URL |

### 请求示例

```json
{
  "name": "西湖（更新）",
  "description": "第一次一起看湖，记忆深刻"
}
```

### 请求示例（curl）

```bash
curl -X PUT "http://localhost:8080/api/v1/places/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "name": "西湖（更新）",
    "description": "第一次一起看湖，记忆深刻"
  }'
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| msg | string | 响应消息 |
| data | object | 响应数据 |
| data.id | int | 地点 ID |
| data.name | string | 地点名称 |
| data.latitude | number | 纬度 |
| data.longitude | number | 经度 |
| data.description | string | 地点描述 |
| data.date | string | 日期（YYYY-MM-DD） |

### 响应示例

```json
{
  "code": 0,
  "msg": "更新成功",
  "data": {
    "id": 1,
    "name": "西湖（更新）",
    "latitude": 30.2592,
    "longitude": 120.1302,
    "description": "第一次一起看湖，记忆深刻",
    "date": "2024-01-15"
  }
}
```

---

## 4. 删除地点

删除指定的地点。

### 请求信息

- **接口路径**: `DELETE /api/v1/places/:id`
- **Content-Type**: `application/json`
- **需要认证**: 是

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 地点 ID |

### 请求示例（curl）

```bash
curl -X DELETE "http://localhost:8080/api/v1/places/1" \
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

---

## 注意事项

1. **权限控制**: 除了获取地点列表外，其他地点操作都需要认证。
2. **坐标格式**: 纬度范围 -90 到 90，经度范围 -180 到 180。
3. **日期格式**: 日期必须使用 YYYY-MM-DD 格式。
4. **图片上传**: 地点的图片需要先通过上传接口上传，然后在创建/更新地点时提供 URL。
5. **删除注意**: 删除地点会同时删除关联的图片，操作不可恢复。

---

## 版本历史

| 版本 | 日期 | 说明 |
|------|------|------|
| 1.2.0 | 2026-03-13 | 新增：列表接口支持分页、排序和过滤功能，新增 page、size、sort_by、order、filter 参数 |
| 1.1.0 | 2026-02-02 | 更新：列表查询不需要认证，删除单个地点查询接口 |
| 1.0.0 | 2026-01-28 | 初始版本，支持地点的增删改查，支持地理位置信息 |