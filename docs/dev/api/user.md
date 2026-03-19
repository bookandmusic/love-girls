# User API 文档

## 概述

User API 提供用户管理和认证功能，包括用户登录、获取用户信息、获取用户列表、更新用户信息等操作。

---

## 1. 用户登录

用户登录并获取访问令牌。

### 请求信息

- **接口路径**: `POST /api/v1/user/token`
- **Content-Type**: `application/json`
- **需要认证**: 否

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

### 请求示例（curl）

```bash
curl -X POST "http://localhost:8080/api/v1/user/token" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "your_username",
    "password": "your_password"
  }'
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| access_token | string | 访问令牌 |
| token_type | string | 令牌类型 |
| expires_in | int | 过期时间（秒） |

### 成功响应示例

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "bearer",
  "expires_in": 3600
}
```

### 错误响应

**400 Bad Request - 参数错误**
```json
{
  "code": 1,
  "message": "参数格式错误或字段缺失",
  "data": null
}
```

**401 Unauthorized - 登录失败**
```json
{
  "code": 1,
  "message": "用户名或密码错误",
  "data": null
}
```

---

## 2. 获取用户信息

获取当前认证用户的信息。

### 请求信息

- **接口路径**: `GET /api/v1/user`
- **Content-Type**: `application/json`
- **需要认证**: 是

### 请求参数

通过 Authorization header 传递 token：
```
Authorization: Bearer <token>
```

### 请求示例（curl）

```bash
curl -X GET "http://localhost:8080/api/v1/user" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| message | string | 响应消息 |
| data | object | 响应数据 |
| data.userName | string | 用户名 |
| data.userId | int | 用户 ID |
| data.userEmail | string | 用户邮箱 |

### 成功响应示例

```json
{
  "code": 0,
  "message": "查询成功",
  "data": {
    "userName": "张三",
    "userId": 1,
    "userEmail": "zhangsan@example.com"
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

**500 Internal Server Error - 查询失败**
```json
{
  "code": 1,
  "message": "系统内部错误",
  "data": null
}
```

---

## 3. 获取用户列表

获取所有用户的列表。

### 请求信息

- **接口路径**: `GET /api/v1/users`
- **Content-Type**: `application/json`
- **需要认证**: 是

### 请求参数

通过 Authorization header 传递 token：
```
Authorization: Bearer <token>
```

### 请求示例（curl）

```bash
curl -X GET "http://localhost:8080/api/v1/users" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| message | string | 响应消息 |
| data | array | 用户列表 |
| data[].id | number | 用户 ID |
| data[].name | string | 用户名 |
| data[].email | string | 用户邮箱 |
| data[].role | string | 用户角色 |
| data[].joinDate | string | 加入日期 |
| data[].avatar | string | 用户头像 |
| data[].avatarId | number | 头像 ID |

### 成功响应示例

```json
{
  "code": 0,
  "message": "查询成功",
  "data": [
    {
      "id": 1,
      "name": "用户1",
      "email": "user1@example.com",
      "role": "user",
      "joinDate": "2026-01-01",
      "avatar": "",
      "avatarId": 1
    },
    {
      "id": 2,
      "name": "用户2",
      "email": "user2@example.com",
      "role": "user",
      "joinDate": "2026-01-02",
      "avatar": "",
      "avatarId": 2
    }
  ]
}
```

---

## 4. 更新用户信息

更新指定用户的信息。

### 请求信息

- **接口路径**: `PUT /api/v1/users/{id}`
- **Content-Type**: `application/json`
- **需要认证**: 是

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | number | 是 | 用户 ID（路径参数） |
| name | string | 是 | 用户名 |
| email | string | 是 | 用户邮箱 |
| avatar | string | 否 | 用户头像 |
| avatarId | number | 否 | 头像 ID |
| role | string | 否 | 用户角色 |
| newPassword | string | 否 | 新密码 |

### 请求示例（curl）

```bash
curl -X PUT "http://localhost:8080/api/v1/users/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "name": "新用户名",
    "email": "newemail@example.com",
    "avatarId": 2,
    "newPassword": "newpassword123"
  }'
```

### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| code | int | 状态码，0 表示成功 |
| message | string | 响应消息 |
| data | object | 更新后的用户信息 |
| data.id | number | 用户 ID |
| data.name | string | 用户名 |
| data.email | string | 用户邮箱 |
| data.role | string | 用户角色 |
| data.joinDate | string | 加入日期 |
| data.avatar | string | 用户头像 |
| data.avatarId | number | 头像 ID |

### 成功响应示例

```json
{
  "code": 0,
  "message": "更新成功",
  "data": {
    "id": 1,
    "name": "新用户名",
    "email": "newemail@example.com",
    "role": "user",
    "joinDate": "2026-01-01",
    "avatar": "",
    "avatarId": 2
  }
}
```

### 错误响应

**400 Bad Request - 参数错误**
```json
{
  "code": 1,
  "message": "参数格式错误或字段缺失",
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

**500 Internal Server Error - 更新失败**
```json
{
  "code": 1,
  "message": "系统内部错误",
  "data": null
}
```

---

## 注意事项

1. **权限控制**: 除登录接口外，其他接口都需要通过 Authorization header 传递 Bearer token 进行认证。
2. **登录安全**: 登录接口对用户名和密码进行验证，失败时不透露具体错误原因。
3. **令牌管理**: 访问令牌具有时效性（默认3600秒），过期后需要重新登录获取新令牌。
4. **Token 存储**: Token 应该安全地存储在客户端，推荐使用 localStorage 或 sessionStorage。
5. **安全传输**: 所有 API 请求都应该使用 HTTPS 协议，确保数据传输安全。
6. **请求头格式**: Authorization 头的格式必须是 `Bearer {token}`，注意 Bearer 后面有空格。
7. **用户管理**: 可以获取所有用户的列表和更新用户信息，需要相应的权限。

---

## 版本历史

| 版本 | 日期 | 说明 |
|------|------|------|
| 3.1.0 | 2026-02-02 | 合并Auth API文档到User API文档，统一管理用户认证和管理相关接口 |
| 3.0.0 | 2026-02-01 | 更新为后端实际实现的 API 接口，包括登录、获取用户信息、获取用户列表和更新用户信息 |