# AGENTS.md - 代码库操作指南

> 本文档是 AI 开发助手的工作手册，包含项目结构索引、命令速查和文档索引。

本项目是一个全栈应用，包含 Vue 3 + TypeScript 前端和 Go + Gin 后端，支持 Tauri 桌面/移动端客户端。

## 0. 工作规则

### 开发需求处理流程

用户提出开发需求时，必须遵循以下流程：

1. **需求确认**：使用 `question` 工具询问用户，确定所有不确定的内容
2. **生成计划**：在 `docs/dev/plan/` 目录生成计划文档，命名规则：`YYYY-MM-DD-功能名称.md`（优先使用中文）
3. **等待确认**：计划必须等待用户确认后才能执行
4. **执行检查**：计划最后需明确需要执行的校验工具（如 `./tools/check.sh all`）

### 错误修复处理流程

用户报告错误时，必须遵循以下流程：

1. **分析代码**：优先分析相关代码，定位问题根因
2. **生成计划**：在 `docs/dev/plan/` 目录生成修复计划
3. **等待确认**：修复计划必须等待用户确认后才能执行

### 计划文档模板

```markdown
# [功能名称] 计划

> 本文档记录 [功能名称] 的实现计划。

## 概述

**创建时间**: YYYY-MM-DD
**优先级**: 高/中/低
**影响范围**: [影响模块列表]

---

## 需求描述

[详细说明需求内容]

---

## 实现方案

### 步骤一：[步骤名称]

**文件**: `[文件路径]`

**改动说明**:
- [具体改动点]

### 步骤二：...

---

## 校验命令

完成后执行以下命令进行验证：

```bash
# 后端检查
cd backend && ./tools/check.sh all

# 前端检查
cd frontend && pnpm lint && pnpm type-check && pnpm format
```

---

## TODO 进度跟踪

- [ ] [任务一]
- [ ] [任务二]

---

## 变更记录

| 日期 | 版本 | 变更内容 |
|------|------|----------|
| YYYY-MM-DD | v1.0 | 初始版本 |
```

**技术版本**:
- Go 1.25+
- Node 24+
- Vue 3
- TypeScript
- Tauri 2.0

## 1. 项目结构

```
love-girl/
├── frontend/                    # pnpm workspace 根目录
│   ├── package.json            # workspace 配置和脚本
│   ├── pnpm-workspace.yaml    # workspace 成员定义
│   ├── node_modules/           # 共享依赖
│   ├── web-frontend/           # Web 前台项目
│   │   ├── src/
│   │   │   ├── views/          # 前台页面 (Home/Albums/Moments/Places/Anniversaries)
│   │   │   ├── router/         # 路由配置
│   │   │   ├── components/     # UI组件和业务组件
│   │   │   ├── services/       # API服务
│   │   │   ├── stores/         # 状态管理
│   │   │   └── assets/         # 前台资源
│   │   └── package.json
│   ├── web-admin/              # Web 后台管理项目
│   │   ├── src/
│   │   │   ├── views/          # 后台页面 (Dashboard/Users/Content/Settings)
│   │   │   ├── router/         # 路由配置
│   │   │   └── assets/         # 后台资源
│   │   └── package.json
│   └── client/                 # Tauri 桌面/移动端客户端
│       ├── src/
│       │   ├── views/          # 所有页面 (前台+后台+client专用)
│       │   ├── router/         # 完整路由配置
│       │   ├── utils/          # client专用工具 (platform.ts)
│       │   └── components/     # client专用组件 (SplashScreen/DesktopMenu)
│       ├── src-tauri/          # Tauri 配置
│       └── package.json
├── backend/                    # Go 后端 (Gin + GORM)
│   ├── internal/               # 内部包
│   ├── provider/               # 依赖注入 (Wire)
│   ├── tools/                  # 开发工具脚本
│   ├── main.go                 # 入口文件
│   └── go.mod                  # Go 模块定义
├── scripts/                    # 开发脚本
└── docs/                       # 项目文档
```

## 2. 常用命令

### frontend/ 目录 (pnpm workspace)

| 命令 | 描述 |
|-----|------|
| `pnpm install` | 安装所有依赖 |
| `pnpm dev:frontend` | 启动前台开发服务器 (端口 5173) |
| `pnpm dev:admin` | 启动后台开发服务器 (端口 5174) |
| `pnpm dev:client` | 启动客户端开发服务器 (端口 5175) |
| `pnpm build:frontend` | 构建前台生产版本 |
| `pnpm build:admin` | 构建后台生产版本 |
| `pnpm build:client` | 构建客户端生产版本 |
| `pnpm build:all` | 构建所有项目 |
| `pnpm lint` | Lint 所有项目 |
| `pnpm type-check` | 类型检查所有项目 |

### 子项目构建 (子路径部署)

```bash
# 前台部署在 /
pnpm build:frontend

# 后台部署在 /admin/
VITE_BASE_URL=/admin/ pnpm build:admin
```

### Tauri 客户端 (frontend/client/)

| 命令 | 描述 |
|-----|------|
| `cd frontend/client && pnpm tauri dev` | 启动 Tauri 开发模式 |
| `cd frontend/client && pnpm tauri build` | 构建桌面客户端 |
| `cd frontend/client && pnpm tauri android init` | 初始化 Android 项目 |
| `cd frontend/client && pnpm tauri android build` | 构建 Android APK |

### 后端 (backend/)

| 命令 | 描述 |
|-----|------|
| `go mod tidy` | 安装 Go 依赖 |
| `go run main.go` | 启动开发服务器 |
| `go test ./...` | 运行所有测试 |
| `./tools/check.sh all` | 运行所有检查 |

**端口**: 前台 5173，后台 5174，客户端 5175，后端 8182

## 3. 部署配置

### Docker 部署说明

Docker 镜像构建时自动将前端打包嵌入后端：
- **前台**：部署在 `/` 路径
- **后台**：部署在 `/admin/` 路径

访问地址：
- 前台：`http://localhost:8182/`
- 后台：`http://localhost:8182/admin/`

## 4. 文档索引

| 文档 | 描述 |
|-----|------|
| [docs/README.md](./docs/README.md) | 文档目录与索引 |
| [docs/user/deployment.md](./docs/user/deployment.md) | 部署指南 |
| [docs/user/config.md](./docs/user/config.md) | 配置项说明 |
| [docs/dev/api/README.md](./docs/dev/api/README.md) | API 文档索引 |
| [docs/dev/client/android-keystore.md](./docs/dev/client/android-keystore.md) | Android 签名配置 |
| [docs/dev/plan/frontend-split-plan.md](./docs/dev/plan/frontend-split-plan.md) | 前端拆分方案 |

---

祝你开发愉快！🚀