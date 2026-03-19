# AGENTS.md - 代码库操作指南

> 本文档是 AI 开发助手的工作手册，包含项目结构索引、命令速查和文档索引。

本项目是一个全栈应用，包含 Vue 3 + TypeScript 前端和 Go + Gin 后端，支持 Tauri 桌面/移动端客户端。

**技术版本**:
- Go 1.25+
- Node 24+
- Vue 3
- TypeScript
- Tauri 2.0

## 1. 项目结构

```
love-girl/
├── backend/              # Go 后端 (Gin + GORM)
│   ├── internal/         # 内部包 (auth/config/db/handler/middleware/model/repo/service/storage/utils...)
│   ├── provider/         # 依赖注入 (Wire)
│   ├── tools/            # 开发工具脚本
│   ├── main.go           # 入口文件
│   └── go.mod            # Go 模块定义
├── frontend/             # Vue 3 前端 (Vite + TypeScript)
│   ├── src/              # 源代码
│   │   ├── assets/       # 静态资源 (images/icons/menu/lottie)
│   │   ├── components/   # 组件 (ui/business)
│   │   ├── layouts/      # 布局组件
│   │   ├── router/       # 路由配置
│   │   ├── stores/       # Pinia 状态管理
│   │   ├── views/        # 页面视图 (admin/frontend)
│   │   ├── services/     # API 服务
│   │   ├── utils/        # 工具函数
│   │   └── types/        # TypeScript 类型定义
│   ├── src-tauri/        # Tauri 桌面/移动端配置
│   │   ├── src/          # Rust 源代码
│   │   ├── icons/        # 应用图标 (各平台)
│   │   ├── capabilities/ # 权限配置
│   │   ├── Cargo.toml    # Rust 依赖
│   │   └── tauri.conf.json # Tauri 主配置
│   ├── public/           # 静态资源
│   └── package.json      # Node 依赖
├── scripts/              # 开发脚本
│   ├── generate-android-keystore.sh      # Android 签名密钥生成 (本地)
│   └── generate-keystore-docker.sh       # Android 签名密钥生成 (Docker)
├── docs/                 # 项目文档
│   ├── user/             # 用户文档 (部署、配置)
│   ├── dev/              # 开发文档 (api、client、plan)
│   └── assets/           # 文档资源
└── .github/              # GitHub 配置
    └── workflows/        # CI/CD 工作流
```

## 2. 常用命令

### 前端 (frontend/)

| 命令 | 描述 |
|-----|------|
| `pnpm install` | 安装依赖 |
| `pnpm dev` | 启动开发服务器 |
| `pnpm type-check` | 类型检查 |
| `pnpm build` | 构建生产版本 |
| `pnpm build:desktop` | 构建桌面端版本 |
| `pnpm lint` | Lint 并自动修复 |
| `pnpm format` | 代码格式化 |

### Tauri 客户端 (frontend/)

| 命令 | 描述 |
|-----|------|
| `pnpm tauri dev` | 启动 Tauri 开发模式 |
| `pnpm tauri build` | 构建桌面客户端 |
| `pnpm tauri android init` | 初始化 Android 项目 |
| `pnpm tauri android build` | 构建 Android APK |

### 后端 (backend/)

| 命令 | 描述 |
|-----|------|
| `go mod tidy` | 安装 Go 依赖 |
| `go run main.go` | 启动开发服务器 |
| `go test ./...` | 运行所有测试 |
| `./tools/check.sh all` | 运行所有检查 |

**端口**: 前端 5173，后端 8182

## 3. 文档索引

| 文档 | 描述 |
|-----|------|
| [docs/README.md](./docs/README.md) | 文档目录与索引 |
| [docs/user/deployment.md](./docs/user/deployment.md) | 部署指南 |
| [docs/user/config.md](./docs/user/config.md) | 配置项说明 |
| [docs/dev/api/README.md](./docs/dev/api/README.md) | API 文档索引 |
| [docs/dev/client/android-keystore.md](./docs/dev/client/android-keystore.md) | Android 签名配置 |
| [docs/dev/plan/backend-optimization.md](./docs/dev/plan/backend-optimization.md) | 后端优化计划 |

---

祝你开发愉快！🚀