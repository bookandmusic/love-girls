# 项目文档

本目录包含 Love Girl 项目的完整文档。

## 目录结构

```
docs/
├── assets/           # 文档资源（图片等）
├── user/             # 用户文档（面向应用使用者）
└── dev/              # 开发文档（面向开发者）
    ├── api/          # API 接口文档
    ├── client/       # 客户端构建
    └── plan/         # 开发计划
```

## 文档索引

### 用户文档

面向应用使用者，包含部署和配置说明。

| 文档 | 说明 |
|------|------|
| [部署指南](./user/deployment.md) | Docker 部署、生产环境配置示例 |
| [配置说明](./user/config.md) | 完整的配置项说明、环境变量、配置优先级 |

### 开发文档

面向开发者，包含 API 文档、客户端构建、开发计划等。

#### API 文档

| 文档 | 说明 |
|------|------|
| [API 概览](./dev/api/README.md) | API 公共约定、请求/响应格式、常见问题 |
| [System API](./dev/api/system.md) | 系统初始化、配置和站点设置 |
| [User API](./dev/api/user.md) | 用户管理 |
| [Album API](./dev/api/album.md) | 相册管理（包含照片功能） |
| [Anniversary API](./dev/api/anniversary.md) | 纪念日管理 |
| [Moment API](./dev/api/moment.md) | 动态管理 |
| [Place API](./dev/api/place.md) | 地点管理 |
| [File API](./dev/api/file.md) | 文件上传与管理 |

#### 客户端构建

| 文档 | 说明 |
|------|------|
| [Android 签名配置](./dev/client/android-keystore.md) | GitHub Actions 构建 Android APK 的签名配置 |

#### 开发计划

| 文档 | 说明 |
|------|------|
| [后端优化计划](./dev/plan/backend-optimization.md) | 后端安全、性能、代码质量优化 |

## 快速导航

### 快速部署

1. [部署指南](./user/deployment.md) - Docker 一键部署
2. [配置说明](./user/config.md) - 了解配置项

### 客户端构建

1. [Android 签名配置](./dev/client/android-keystore.md) - 配置 APK 签名
2. GitHub Actions 自动构建 - 推送 tag 触发

### 开发者

1. [API 文档](./dev/api/README.md) - 接口调用参考
2. [后端优化计划](./dev/plan/backend-optimization.md) - 后端开发任务跟踪