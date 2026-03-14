# AGENTS.md - 代码库操作指南

> 本文档是 AI 开发助手的工作手册，包含项目结构索引、命令速查和文档索引。

本项目是一个全栈应用，包含 Vue 3 + TypeScript 前端和 Go + Gin 后端。

**技术版本**:
- Go 1.25+
- Node 24+
- Vue 3
- TypeScript

## 1. 项目结构

```
love-girl/
├── backend/              # Go 后端 (Gin + GORM)
│   ├── internal/         # 内部包 (auth/config/db/handler/model/repo/service...)
│   ├── provider/         # 依赖注入 (Wire)
│   ├── tools/            # 开发工具脚本
│   ├── main.go           # 入口文件
│   └── go.mod            # Go 模块定义
├── frontend/             # Vue 3 前端 (Vite + TypeScript)
│   ├── src/              # 源代码 (components/router/stores/views...)
│   ├── public/           # 静态资源
│   └── package.json      # Node 依赖
├── docs/                 # 项目文档
│   ├── api/              # API 文档
│   ├── guide/            # 配置与部署指南
│   └── assets/           # 文档资源
└── .opencode/            # Skills 配置
```

## 2. 常用命令

### 前端 (frontend/)

| 命令 | 描述 |
|-----|------|
| `pnpm install` | 安装依赖 |
| `pnpm dev` | 启动开发服务器 |
| `pnpm type-check` | 类型检查 |
| `pnpm build` | 构建生产版本 |
| `pnpm lint` | Lint 并自动修复 |
| `pnpm format` | 代码格式化 |

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
| [docs/guide/CONFIG.md](./docs/guide/CONFIG.md) | 配置项说明 |
| [docs/guide/DEPLOYMENT.md](./docs/guide/DEPLOYMENT.md) | 部署指南 |
| [docs/api/README.md](./docs/api/README.md) | API 文档索引 |

---

祝你开发愉快！🚀