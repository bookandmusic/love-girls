# 快速开始

## 1. 环境要求

### 必需

- **Go**: 1.25 或更高版本
- **Node.js**: 24 或更高版本
- **pnpm**: 包管理器（推荐）或 npm/yarn

### 可选

- **Docker**: 用于容器化部署
- **Git**: 用于版本控制

## 2. 克隆项目

```bash
git clone https://github.com/bookandmusic/love-girl.git
cd love-girl
```

## 3. 后端启动

```bash
cd backend

# 安装 Go 依赖
go mod tidy

# 启动开发服务器
go run main.go
```

后端服务将在 `http://localhost:8182` 启动

**首次启动会自动创建**：
- `./data/configs/config.yaml` - 配置文件
- `./data/love-girl.db` - SQLite 数据库
- `./data/uploads/` - 文件上传目录

## 4. 前端启动

```bash
cd frontend

# 安装 Node 依赖
pnpm install

# 启动开发服务器
pnpm dev
```

前端服务将在 `http://localhost:5173` 启动

## 5. 访问应用

打开浏览器访问: `http://localhost:5173`
