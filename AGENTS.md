# AGENTS.md - 代码库操作指南

> 本文档是 AI 开发助手的工作手册，包含项目结构索引、命令速查、Skills 和文档索引。

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
│   ├── data/             # 数据目录 (config.yaml/love-girl.db/uploads/)
│   ├── main.go           # 入口文件
│   └── go.mod            # Go 模块定义
├── frontend/             # Vue 3 前端 (Vite + TypeScript)
│   ├── src/              # 源代码 (components/router/stores/views...)
│   ├── public/           # 静态资源
│   └── package.json      # Node 依赖
└── docs/                 # 项目文档 (见第 3 节)
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

## 3. Skills

Skills 是预定义的自动化工作流。

### 可用的 Skills

| Skill | 功能描述 |
|-------|---------|
| **code-check** | 代码质量检查 |
| **git-commit** | 生成规范 Commit Message |
| **requirement-dev** | 需求开发工作流 |
| **add-skill** | 创建新的 Skill |

### 调用方式

```
请执行 <skill-name> skill
```

### 使用场景

| 场景 | 命令 |
|-----|------|
| 代码修改后检查 | `请执行 code-check skill` |
| 提交代码前 | `请执行 git-commit skill` |
| 开始新功能前 | `请执行 requirement-dev skill` |
| 创建新的 Skill | `请执行 add-skill skill` |

**详细说明**: 参见 `.opencode/skills/<skill-name>/SKILL.md`

## 4. 文档索引

### 快速开始

| 文档 | 描述 |
|-----|------|
| [docs/guides/quickstart.md](./docs/guides/quickstart.md) | 环境要求、项目启动 |

### 开发指南

| 文档 | 描述 |
|-----|------|
| [docs/guides/commands.md](./docs/guides/commands.md) | 构建、测试、lint 命令 |
| [docs/guides/coding-style.md](./docs/guides/coding-style.md) | 代码风格规范 |
| [docs/guides/debugging.md](./docs/guides/debugging.md) | 调试技巧 |
| [docs/guides/faq.md](./docs/guides/faq.md) | 常见问题解答 |

### 配置与 API

| 文档 | 描述 |
|-----|------|
| [docs/config.md](./docs/config.md) | 项目配置说明 |
| [docs/api/README.md](./docs/api/README.md) | API 文档索引 |

## 5. AI 行为指引

1. **阅读文档**: 遇到问题时，先查阅 `docs/guides/` 下的相关文档
2. **调用 Skills**: 按场景调用合适的 skill
   - 代码修改后 → `code-check`
   - 提交代码前 → `git-commit`
   - 新功能开发前 → `requirement-dev`
   - 创建新的 Skill → `add-skill`
3. **生成设计文档**: 新功能开发使用 `requirement-dev` skill 自动生成
4. **代码检查**: 每次代码修改后必须执行 `code-check` skill

## 6. API 文档

启动后端服务后访问: `http://localhost:8182/swagger/index.html`

---

祝你开发愉快！🚀
