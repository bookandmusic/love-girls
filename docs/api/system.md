# System API 文档

## 概述

System API 提供系统初始化、状态检查、系统信息查询和站点设置管理功能。这些接口主要用于情侣纪念站点的初始化设置和配置管理。

---

## 1. 初始化系统

初始化系统，创建站点信息、用户信息和设置访问密码。

### 请求信息

- **接口路径**: `POST /api/v1/system/init`
- **Content-Type**: `application/json`
- **需要认证**: 否

### 请求参数

| 参数名 | 类型 | 必填 | 说明 | 限制条件 |
|--------|------|------|------|----------|
| siteName | string | 是 | 站点名称 | 不能为空 |
| siteDescription | string | 否 | 站点描述 | 可选 |
| startDate | string | 是 | 故事开始的日期 | 格式：YYYY-MM-DD |
| userAName | string | 是 | 用户 A 昵称 | 不能为空 |
| userARole | string | 是 | 用户 A 角色 | 枚举值：boy / girl |
| userAPhone | string | 否 | 用户 A 手机号 | 可选 |
| userBName | string | 是 | 用户 B 昵称 | 不能为空 |
| userBRole | string | 是 | 用户 B 角色 | 枚举值：boy / girl |
| userBEmail | string | 否 | 用户 B 邮箱 | 可选，需符合邮箱格式 |
| userBPhone | string | 否 | 用户 B 手机号 | 可选 |
| sitePassword | string | 是 | 站点访问密码 | 最小长度：6 |
| sitePasswordConfirm | string | 是 | 确认密码 | 必须与 sitePassword 一致 |

### 请求示例

```bash
curl -X POST "http://localhost:8080/api/v1/system/init" \
  -H "Content-Type: application/json" \
  -d '{
    "siteName": "鹿与星的纪念站",
    "siteDescription": "记录我们的美好时光",
    "startDate": "2024-01-01",
    "userAName": "鹿",
    "userARole": "boy",
    "userAEmail": "lu@example.com",
    "userAPhone": "",
    "userBName": "星",
    "userBRole": "girl",
    "userBEmail": "xing@example.com",
    "userBPhone": "",
    "sitePassword": "123456",
    "sitePasswordConfirm": "123456"
  }'
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| message | string | 响应消息 |
| data | object | 响应数据 |
| data.initialized | boolean | 系统是否已初始化 |

### 成功响应

**200 OK - 初始化成功**
```json
{
  "code": 0,
  "message": "系统初始化成功",
  "data": {
    "initialized": true
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

**500 Internal Server Error - 系统初始化失败**
```json
{
  "code": 1,
  "message": "系统初始化失败",
  "data": null
}
```

常见错误原因：
- 参数校验失败：必填字段为空、密码长度不足6位、两次密码不一致等
- 系统已经初始化，无法重复初始化
- 创建初始数据失败（数据库错误）
- 更新初始化状态失败（配置文件写入错误）

---

## 2. 检查系统初始化状态

检查系统是否已完成初始化。

### 请求信息

- **接口路径**: `GET /api/v1/system/init`
- **Content-Type**: `application/json`
- **需要认证**: 否

### 请求参数

无

### 请求示例

```bash
curl -X GET "http://localhost:8080/api/v1/system/init"
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| message | string | 响应消息 |
| data | object | 响应数据 |
| data.initialized | boolean | 系统是否已初始化 |

### 成功响应

**200 OK - 查询成功**
```json
{
  "code": 0,
  "message": "查询成功",
  "data": {
    "initialized": true
  }
}
```

**200 OK - 查询成功（未初始化）**
```json
{
  "code": 0,
  "message": "查询成功",
  "data": {
    "initialized": false
  }
}
```

### 错误响应

**500 Internal Server Error - 检查失败**
```json
{
  "code": 1,
  "message": "检查初始化状态失败",
  "data": null
}
```

---

## 3. 获取系统信息

获取系统配置的基本信息。

### 请求信息

- **接口路径**: `GET /api/v1/system/info`
- **Content-Type**: `application/json`
- **需要认证**: 否

### 请求参数

无

### 请求示例

```bash
curl -X GET "http://localhost:8080/api/v1/system/info"
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| message | string | 响应消息 |
| data | object | 响应数据 |
| data.site | object | 站点信息对象 |
| data.site.name | string | 站点名称 |
| data.site.description | string | 站点描述 |
| data.site.startDate | string | 故事开始的日期 |
| data.couple | object | 情侣信息对象 |
| data.couple.boy | object | 男友信息 |
| data.couple.boy.name | string | 男友昵称 |
| data.couple.boy.avatar | string | 男友头像URL |
| data.couple.girl | object | 女友信息 |
| data.couple.girl.name | string | 女友昵称 |
| data.couple.girl.avatar | string | 女友头像URL |

### 成功响应

**200 OK - 查询成功**
```json
{
  "code": 0,
  "message": "查询成功",
  "data": {
    "site": {
      "name": "鹿与星的纪念站",
      "description": "记录我们的美好时光",
      "startDate": "2024-01-01"
    },
    "couple": {
      "boy": {
        "name": "鹿",
        "avatar": ""
      },
      "girl": {
        "name": "星",
        "avatar": ""
      }
    }
  }
}
```

### 错误响应

**404 Not Found - 系统未初始化或获取失败**
```json
{
  "code": 1,
  "message": "系统未初始化或获取信息失败",
  "data": null
}
```

**500 Internal Server Error - 获取失败**
```json
{
  "code": 1,
  "message": "系统内部错误",
  "data": null
}
```

常见错误原因：
- 系统未初始化（找不到站点设置）
- 数据库查询失败
- 用户信息不完整

---

## 4. 获取站点设置

获取当前站点的设置信息。

### 请求信息

- **接口路径**: `GET /api/v1/system/settings/site`
- **Content-Type**: `application/json`
- **需要认证**: 是（建议仅管理员可用）

### 请求参数

无

### 请求示例（curl）

```bash
curl -X GET "http://localhost:8080/api/v1/system/settings/site" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| message | string | 响应消息 |
| data | object | 响应数据，包含所有通用设置项 |
| data.siteTitle | string | 站点标题 |
| data.siteDescription | string | 站点描述 |
| data.startDate | string | 故事开始日期 |

### 成功响应

**200 OK - 查询成功**
```json
{
  "code": 0,
  "message": "查询成功",
  "data": {
    "siteTitle": "鹿与星的纪念站",
    "siteDescription": "记录我们的美好时光",
    "startDate": "2024-01-01"
  }
}
```

### 错误响应

**401 Unauthorized - 未认证**
```json
{
  "code": 1,
  "message": "未授权访问",
  "data": null
}
```

**500 Internal Server Error - 获取失败**
```json
{
  "code": 1,
  "message": "系统内部错误",
  "data": null
}
```

---

## 5. 保存站点设置

更新站点的设置信息。

### 请求信息

- **接口路径**: `POST /api/v1/system/settings/site`
- **Content-Type**: `application/json`
- **需要认证**: 是（建议仅管理员可用）

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| siteTitle | string | 是 | 站点标题 |
| siteDescription | string | 否 | 站点描述 |

### 请求示例

```json
{
  "siteTitle": "鹿与星的纪念站",
  "siteDescription": "记录我们的美好时光"
}
```

### 请求示例（curl）

```bash
curl -X POST "http://localhost:8080/api/v1/system/settings/site" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "siteTitle": "鹿与星的纪念站",
    "siteDescription": "记录我们的美好时光"
  }'
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| message | string | 响应消息 |
| data | object | 响应数据，包含保存的设置项 |
| data.siteTitle | string | 站点标题 |
| data.siteDescription | string | 站点描述 |

### 成功响应

**200 OK - 保存成功**
```json
{
  "code": 0,
  "message": "保存成功",
  "data": {
    "siteTitle": "鹿与星的纪念站",
    "siteDescription": "记录我们的美好时光"
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

**401 Unauthorized - 未认证**
```json
{
  "code": 1,
  "message": "未授权访问",
  "data": null
}
```

**500 Internal Server Error - 保存失败**
```json
{
  "code": 1,
  "message": "系统内部错误",
  "data": null
}
```

---

## 5. 获取仪表盘统计数据

获取仪表盘的统计数据，用于展示系统的整体运营情况和数据概览。

### 请求信息

- **接口路径**: `GET /api/v1/system/dashboard/stats`
- **Content-Type**: `application/json`
- **需要认证**: 是

### 请求参数

无

### 请求示例（curl）

```bash
curl -X GET "http://localhost:8080/api/v1/system/dashboard/stats" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| message | string | 响应消息 |
| data | object | 响应数据 |
| data.albumStats | object | 相册统计 |
| data.albumStats.total | int | 相册总数 |
| data.albumStats.totalPhotos | int | 照片总数 |
| data.placeStats | object | 地点统计 |
| data.placeStats.total | int | 地点总数 |
| data.momentStats | object | 动态统计 |
| data.momentStats.total | int | 动态总数 |
| data.wishStats | object | 愿望统计 |
| data.wishStats.total | int | 愿望总数 |
| data.wishStats.pending | int | 待审核愿望数 |

### 响应示例

```json
{
  "code": 0,
  "message": "查询成功",
  "data": {
    "albumStats": {
      "total": 10,
      "totalPhotos": 256
    },
    "placeStats": {
      "total": 15
    },
    "momentStats": {
      "total": 42
    },
    "wishStats": {
      "total": 20,
      "pending": 5
    }
  }
}
```

---

## 注意事项

1. **初始化限制**: 系统只能初始化一次，初始化完成后再次调用该接口会返回错误。
2. **密码安全**: 站点访问密码会被加密存储，请妥善保管。
3. **设置权限**: 获取和修改站点设置都需要认证。
4. **全局影响**: 站点设置会影响所有用户的访问体验，修改前请确认内容。
5. **缓存策略**: 仪表盘数据可能会缓存一段时间，实时性要求不高的场景可以适当延长缓存时间。
6. **活动排序**: 最近活动按时间倒序排列，最新的活动显示在前面。
7. **性能优化**: 数据统计可能涉及聚合查询，建议在非高峰时段更新缓存数据。
8. **数据范围**: 统计数据基于当前登录用户可见的数据范围。