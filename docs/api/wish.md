# Wish API 文档

## 概述

Wish API 提供愿望管理功能，支持愿望的创建、查询、删除和审核操作。

---

## 1. 获取祝福列表

获取祝福列表，支持分页查询。

### 请求信息

- **接口路径**: `GET /api/v1/wishes`
- **Content-Type**: `application/json`
- **需要认证**: 否

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 是 | 页码，从 1 开始 |
| size | int | 是 | 每页数量 |

### 请求示例（curl）

```bash
curl -X GET "http://localhost:8080/api/v1/wishes?page=1&size=10" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| message | string | 响应消息 |
| data | object | 响应数据 |
| data.wishes | array | 祝福列表 |
| data.wishes[].id | int | 祝福 ID |
| data.wishes[].content | string | 祝福内容 |
| data.wishes[].authorName | string | 作者名称 |
| data.wishes[].email | string | 作者邮箱 |
| data.wishes[].createdAt | string | 创建时间（ISO 8601） |
| data.wishes[].approved | boolean | 是否已审核 |
| data.totalPages | int | 总页数 |
| data.total | int | 总数量 |
| data.totalCount | int | 总记录数 |
| data.page | int | 当前页码 |
| data.size | int | 每页数量 |

### 响应示例

```json
{
  "code": 0,
  "message": "查询成功",
  "data": {
    "wishes": [
      {
        "id": 1,
        "content": "祝我们永远幸福快乐",
        "authorName": "鹿",
        "email": "lu@example.com",
        "createdAt": "2024-01-15T10:30:00Z",
        "approved": true
      },
      {
        "id": 2,
        "content": "愿我们的爱情天长地久",
        "authorName": "星",
        "email": "xing@example.com",
        "createdAt": "2024-01-16T11:00:00Z",
        "approved": false
      }
    ],
    "totalPages": 1,
    "total": 2,
    "totalCount": 2,
    "page": 1,
    "size": 10
  }
}
```

---


## 2. 创建祝福

创建一个新的祝福。

### 请求信息

- **接口路径**: `POST /api/v1/wishes`
- **Content-Type**: `application/json`
- **需要认证**: 否

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| content | string | 是 | 祝福内容 |
| authorName | string | 是 | 作者名称 |
| email | string | 否 | 作者邮箱 |

### 请求示例

```json
{
  "content": "祝我们永远幸福快乐",
  "authorName": "鹿",
  "email": "lu@example.com"
}
```

### 请求示例（curl）

```bash
curl -X POST "http://localhost:8080/api/v1/wishes" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "content": "祝我们永远幸福快乐",
    "authorName": "鹿",
    "email": "lu@example.com"
  }'
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| message | string | 响应消息 |
| data | object | 响应数据 |
| data.id | int | 祝福 ID |
| data.content | string | 祝福内容 |
| data.authorName | string | 作者名称 |
| data.email | string | 作者邮箱 |
| data.createdAt | string | 创建时间（ISO 8601） |
| data.approved | boolean | 是否已审核 |

### 响应示例

```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "id": 1,
    "content": "祝我们永远幸福快乐",
    "authorName": "鹿",
    "email": "lu@example.com",
    "createdAt": "2024-01-15T10:30:00Z",
    "approved": false
  }
}
```

---

## 3. 删除祝福

删除指定的祝福。

### 请求信息

- **接口路径**: `DELETE /api/v1/wishes/:id`
- **Content-Type**: `application/json`
- **需要认证**: 是

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 祝福 ID |

### 请求示例（curl）

```bash
curl -X DELETE "http://localhost:8080/api/v1/wishes/1" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| message | string | 响应消息 |
| data | null | 响应数据，删除操作无数据返回 |

### 响应示例

```json
{
  "code": 0,
  "message": "删除成功",
  "data": null
}
```

---

## 4. 批准祝福

批准指定的祝福，使其可以公开显示。

### 请求信息

- **接口路径**: `PUT /api/v1/wishes/:id/approve`
- **Content-Type**: `application/json`
- **需要认证**: 是（仅管理员可用）

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 祝福 ID |

### 请求示例（curl）

```bash
curl -X PUT "http://localhost:8080/api/v1/wishes/1/approve" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| message | string | 响应消息 |
| data | null | 响应数据，审核操作无数据返回 |

### 响应示例

```
{
  "code": 0,
  "message": "审核成功",
  "data": null
}
```

---

---

## 注意事项

1. **审核机制**: 祝福创建后需要管理员审核通过才能公开显示。
2. **权限控制**: 列表查询和创建祝福不需要认证，删除和批准操作需要认证，批准操作仅管理员可用。
3. **邮箱格式**: 邮箱地址必须符合标准邮箱格式。
4. **内容限制**: 祝福内容长度建议在 10-500 字符之间。
5. **隐私保护**: 未审核的祝福仅管理员可见，审核通过后公开显示。

---

## 版本历史

| 版本 | 日期 | 说明 |
|------|------|------|
| 1.3.0 | 2026-02-02 | 更新：列表查询和创建祝福不需要认证，删除单个祝福查询接口 |
| 1.2.0 | 2026-01-30 | 重命名：将Blessing API重命名回Wish API，以匹配前端实现 |
| 1.1.0 | 2026-01-29 | 重命名：将Wish API重命名为Blessing API |
| 1.0.0 | 2026-01-28 | 初始版本，支持祝福的增删查和审核功能 |