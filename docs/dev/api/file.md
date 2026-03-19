# File API 文档

## 概述

File API 提供文件上传和下载功能，主要用于上传图片（如用户头像、相册照片等）以及文件访问。

---

## 1. 上传文件

上传单个文件。

### 请求信息

- **接口路径**: `POST /api/v1/file/upload`
- **Content-Type**: `multipart/form-data`
- **需要认证**: 否

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| file | file | 是 | 文件 |
| path | string | 否 | 存储路径前缀 |
| hash | string | 是 | 文件内容MD5哈希值 |
| thumbnailWidth | int | 否 | 缩略图宽度 |
| thumbnailHeight | int | 否 | 缩略图高度 |

### 请求示例（curl）

```bash
curl -X POST "http://localhost:8080/api/v1/file/upload" \
  -F "file=@photo.jpg" \
  -F "hash=sha256_hash_value"
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| message | string | 响应消息 |
| data | object | 响应数据 |
| data.id | uint64 | 文件 ID |
| data.url | string | 原始文件 URL |
| data.thumbnailUrl | string | 缩略图 URL |

### 成功响应示例

```json
{
  "code": 0,
  "message": "文件保存成功",
  "data": {
    "id": 101,
    "url": "https://example.com/api/v1/file/101",
    "thumbnailUrl": "https://example.com/api/v1/file/101?width=400&height=300"
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

**400 Bad Request - 缩略图尺寸参数非法**
```json
{
  "code": 1,
  "message": "缩略图尺寸参数非法",
  "data": null
}
```

**500 Internal Server Error - 上传失败**
```json
{
  "code": 1,
  "message": "系统内部错误",
  "data": null
}
```

---

## 2. 获取文件

根据文件 ID 获取文件内容。

### 请求信息

- **接口路径**: `GET /api/v1/file/:id`
- **需要认证**: 否

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | uint64 | 是 | 文件 ID |

### 请求示例（curl）

```bash
curl -X GET "http://localhost:8080/api/v1/file/101"
```

### 响应

直接返回文件内容。

### 错误响应

**400 Bad Request - 无效的文件ID格式**
```json
{
  "code": 1,
  "message": "无效的文件ID格式",
  "data": null
}
```

**500 Internal Server Error - 文件读取失败**
```json
{
  "code": 1,
  "message": "系统内部错误",
  "data": null
}
```

---

## 注意事项

1. **文件大小**: 建议限制在合理范围内。
2. **文件格式**: 支持常见的图片格式如 jpg, png, gif 等。
3. **认证**: 上传操作不需要登录认证，但实际应用中可能需要添加认证。
4. **文件ID**: 使用 uint64 类型的数值ID作为文件标识。
5. **哈希计算**: `hash` 参数必须使用MD5算法计算文件内容的哈希值，前端可使用FileReader API读取文件内容后计算。