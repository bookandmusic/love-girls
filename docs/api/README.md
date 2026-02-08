# API 文档目录

本目录包含情侣纪念站点所有 API 接口的详细文档。

## 文档列表

### 核心功能

- **[System API](./system.md)** - 系统初始化、配置和站点设置
- **[Auth API](./auth.md)** - 用户认证
- **[User API](./user.md)** - 用户管理
- **[Dashboard API](./dashboard.md)** - 仪表盘统计

### 内容管理

- **[Album API](./album.md)** - 相册管理（包含照片功能）
- **[Anniversary API](./anniversary.md)** - 纪念日管理
- **[Moment API](./moment.md)** - 动态管理
- **[Place API](./place.md)** - 地点管理
- **[Wish API](./wish.md)** - 愿望管理

## 公共约定

### 基础路径

所有 API 接口的基础路径为：`/api/v1`

### 请求格式

所有 API 请求（除文件上传外）都使用 JSON 格式：

```http
Content-Type: application/json
```

文件上传使用 multipart/form-data 格式：

```http
Content-Type: multipart/form-data
```

### 认证方式

大部分 API 需要认证，使用 Bearer Token 方式：

```http
Authorization: Bearer {token}
```

### 响应结构

所有 API 接口使用统一的响应格式：

```json
{
  "code": 0,
  "msg": "操作成功",
  "data": {}
}
```

#### 响应字段说明

| 字段名 | 类型 | 说明 |
|--------|------|------|
| code | int | 响应码，0 表示成功，1 表示失败 |
| msg | string | 响应消息，描述操作结果 |
| data | any | 响应数据，根据接口不同返回不同结构 |

### 日期格式

日期格式统一使用 ISO 8601 标准：

- 完整日期时间：`2024-01-26T10:30:00Z`
- 仅日期：`2024-01-26`
- 月-日：`02-14`

## 常见问题

### Q: 如何获取 token？

A: 调用登录接口 `POST /api/v1/login`，成功后会在响应中返回 token。

### Q: Token 过期了怎么办？

A: Token 过期后会返回 code 为 0 的响应，需要重新登录获取新的 token。

### Q: 上传图片的接口是什么？

A: 使用 `POST /api/v1/upload/image` 接口上传图片，Content-Type 为 `multipart/form-data`。

### Q: 如何处理分页？

A: 在请求参数中添加 `page` 和 `size`参数，响应中会包含 `totalPages` 和 `total` 信息。

### Q: 日期格式有什么要求？

A:
- 完整日期时间：ISO 8601 格式，如 `2024-01-26T10:30:00Z`
- 仅日期：YYYY-MM-DD 格式，如 `2024-01-26`
- 月-日：MM-DD 格式，如 `02-14`

## 版本历史

| 版本 | 日期 | 说明 |
|------|------|------|
| 1.0.0 | 2026-01-28 | 初始版本，完整的 API 文档 |
| 1.1.0 | 2026-01-29 | 修改：将Memory API重命名为Anniversary API，将Photo API合并到Album API |
| 1.2.0 | 2026-01-30 | 修改：将Blessing API复原为Wish API |

---

**最后更新时间**: 2026-01-30