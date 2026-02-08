# Anniversary API 文档

## 概述

Anniversary API 提供纪念日管理功能，支持公历和农历日历，用于记录重要的纪念日和特殊日期。

---

## 1. 获取纪念日列表

获取纪念日列表，支持分页查询。

### 请求信息

- **接口路径**: `GET /anniversaries`
- **Content-Type**: `application/json`
- **需要认证**: 否

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，从 1 开始，默认 1 |
| size | int | 否 | 每页数量，默认 10 |

### 请求示例（curl）

```bash
curl -X GET "http://localhost:8080/api/v1/anniversaries?page=1&size=10" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| msg | string | 响应消息 |
| data | object | 响应数据 |
| data.anniversaries | array | 纪念日列表 |
| data.anniversaries[].id | int | 纪念日 ID |
| data.anniversaries[].title | string | 纪念日标题 |
| data.anniversaries[].date | string | 日期（MM-DD 格式） |
| data.anniversaries[].description | string | 纪念日描述 |
| data.anniversaries[].calendar | string | 日历类型：solar（公历）或 lunar（农历） |
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
    "anniversaries": [
      {
        "id": 1,
        "title": "我们的第一次约会",
        "date": "02-14",
        "description": "在咖啡店的美好时光",
        "calendar": "solar"
      },
      {
        "id": 2,
        "title": "生日",
        "date": "08-20",
        "description": "特殊的日子",
        "calendar": "solar"
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

## 2. 创建纪念日

创建一个新的纪念日。

### 请求信息

- **接口路径**: `POST /anniversaries`
- **Content-Type**: `application/json`
- **需要认证**: 是

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| title | string | 是 | 纪念日标题 |
| date | string | 是 | 日期，格式：YYYY-MM-DD |
| description | string | 否 | 纪念日描述 |
| calendar | string | 是 | 日历类型：solar（公历）或 lunar（农历） |

### 请求示例

```json
{
  "title": "我们的第一次约会",
  "date": "02-14",
  "description": "在咖啡店的美好时光",
  "calendar": "solar"
}
```

### 请求示例（curl）

```bash
curl -X POST "http://localhost:8080/api/v1/anniversaries" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "title": "我们的第一次约会",
    "date": "02-14",
    "description": "在咖啡店的美好时光",
    "calendar": "solar"
  }'
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| msg | string | 响应消息 |
| data | object | 响应数据 |
| data.id | int | 纪念日 ID |
| data.title | string | 纪念日标题 |
| data.date | string | 日期（MM-DD 格式） |
| data.description | string | 纪念日描述 |
| data.calendar | string | 日历类型：solar（公历）或 lunar（农历） |

### 响应示例

```json
{
  "code": 0,
  "msg": "创建成功",
  "data": {
    "id": 1,
    "title": "我们的第一次约会",
    "date": "02-14",
    "description": "在咖啡店的美好时光",
    "calendar": "solar"
  }
}
```

---

## 3. 更新纪念日

更新纪念日的信息。

### 请求信息

- **接口路径**: `PUT /anniversaries/:id`
- **Content-Type**: `application/json`
- **需要认证**: 是

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 纪念日 ID |

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| title | string | 否 | 纪念日标题 |
| date | string | 否 | 日期，格式：MM-DD |
| description | string | 否 | 纪念日描述 |
| calendar | string | 否 | 日历类型：solar（公历）或 lunar（农历） |

### 请求示例

```json
{
  "title": "我们的第一次约会（纪念）",
  "date": "02-14",
  "description": "在咖啡店的美好时光，值得铭记",
  "calendar": "solar"
}
```

### 请求示例（curl）

```bash
curl -X PUT "http://localhost:8080/api/v1/anniversaries/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "title": "我们的第一次约会（纪念）",
    "date": "02-14",
    "description": "在咖啡店的美好时光，值得铭记",
    "calendar": "solar"
  }'
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| msg | string | 响应消息 |
| data | object | 响应数据 |
| data.id | int | 纪念日 ID |
| data.title | string | 纪念日标题 |
| data.date | string | 日期（MM-DD 格式） |
| data.description | string | 纪念日描述 |
| data.calendar | string | 日历类型：solar（公历）或 lunar（农历） |

### 响应示例

```json
{
  "code": 0,
  "msg": "更新成功",
  "data": {
    "id": 1,
    "title": "我们的第一次约会（纪念）",
    "date": "02-14",
    "description": "在咖啡店的美好时光，值得铭记",
    "calendar": "solar"
  }
}
```

---

## 4. 删除纪念日

删除指定的纪念日。

### 请求信息

- **接口路径**: `DELETE /anniversaries/:id`
- **Content-Type**: `application/json`
- **需要认证**: 是

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 纪念日 ID |

### 请求示例（curl）

```bash
curl -X DELETE "http://localhost:8080/api/v1/anniversaries/1" \
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

1. **日期格式**: 日期必须使用 MM-DD 格式，例如 "02-14" 表示 2 月 14 日。
2. **日历类型**: calendar 字段支持 "solar"（公历）和 "lunar"（农历）两种类型。
3. **权限控制**: 所有纪念日操作都需要认证。
4. **农历计算**: 农历日期的计算需要后端支持，确保农历转换准确。
5. **重复日期**: 同一日期可以创建多个纪念日，用于记录不同的事件。

---

## 版本历史

| 版本 | 日期 | 说明 |
|------|------|------|
| 1.0.0 | 2026-01-28 | 初始版本，支持纪念日的增删改查，支持公历和农历 |
| 1.1.0 | 2026-01-29 | 重命名：将Memory API重命名为Anniversary API |
| 1.2.0 | 2026-02-01 | 更新API文档，匹配前端接口结构 |