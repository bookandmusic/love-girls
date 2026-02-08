# 配置管理系统

本文档介绍 love-girl 项目的配置管理系统，包括技术架构设计、配置项说明以及使用示例。

## 技术架构设计

结合 Gin 与 Viper 的特性，配置管理系统采用自动化配置管理架构。核心在于利用 fsnotify 监听文件、Viper 管理多层配置，以及使用优雅重启机制实现连接无损切换。

### 1. 静态初始化与自愈阶段 (Startup)

当服务调用 `main()` 启动时，配置中心按照以下顺序建立内存镜像：

#### Step 1：默认值注入
使用 `viper.SetDefault` 在代码中硬编码基础配置（如端口 8182），作为系统底座。

#### Step 2：文件检测与创建
检查 `config.yaml` 文件是否存在。若不存在，执行 `viper.SafeWriteConfigAs("config.yaml")`。此举确保了即使在空目录下，Gin 也能拿到合法的端口和模式配置。

#### Step 3：加载与环境变量覆盖
调用 `viper.ReadInConfig()` 读取磁盘文件，随后立即调用 `viper.AutomaticEnv()`。优先级细化：Viper 会建立一个 Key 索引。查询 port 时，它会优先看系统变量 `APP_PORT` 是否有值，若无则看 config.yaml。

### 2. 自动监听与逻辑验证阶段 (Watch & Validate)

配置文件的修改是异步且不可控的，因此需要建立"缓冲区校验"机制：

#### 热监听 (Hot Watch)
调用 `viper.WatchConfig()`，并注册 `OnConfigChange` 回调函数。

#### 原子性校验逻辑
当回调触发时，不要立即应用。应先创建一个 NewViper 实例去尝试 Unmarshal 新文件：
- **语法检查**：YAML 格式是否正确？
- **业务检查**：修改后的数据库连接是否通畅？端口是否被系统占用？

#### 重启触发器
只有校验通过，主进程使用 `syscall.Exec` 替换自身进程，实现平滑重启。

### 3. 优雅重启执行阶段 (Graceful Restart)

这是保证 Gin 正在处理的请求（如大文件上传或长轮询）不中断的关键。

#### 重启机制
**本项目使用 `syscall.Exec` 实现零停机重启。**

相比 endless 库的 fork 机制，Go 的 `syscall.Exec` 能够更可靠地实现进程替换，直接继承父进程的所有文件描述符（包括监听 socket）。

#### 进程替换细化
当配置验证通过时，系统会：
1. 获取当前可执行文件路径
2. 使用 `syscall.Exec` 替换当前进程为新进程
3. 新进程重新加载配置，启动 Gin 服务
4. 监听 socket 被正确继承，所有连接无缝切换

### 4. 方案集成架构

| 组件 | 职责 | 交互细节 |
|------|------|----------|
| Viper 实例 | 配置存储与监听 | 负责 config.yaml 的读写及 fsnotify 监控 |
| Validation Layer | 准入控制 | 在 OnConfigChange 触发时，阻断错误配置进入重启流程 |
| Exec/Syscall | 进程替换 | 使用 syscall.Exec 替换进程，继承所有文件描述符 |
| Gin Engine | 业务请求处理 | 专注于路由分发，不感知底层的重启逻辑 |

### 5. 关键边界行为

#### 配置错误处理
如果用户改错了配置：
```
Viper 监控到变更 -> 校验失败 -> 打印 Error 日志 -> 服务继续按旧配置运行
```
用户修复文件后，再次保存，系统重新触发校验流程。

#### 环境变量的持久性
重启后的子进程会重新加载环境变量。如果用户在 config.yaml 改了 `port: 9090`，但环境变量 `APP_PORT` 一直是 `8080`，那么重启后 Gin 依然会监听在 `8080`。这符合"环境变量最高优先"的设计原则。

---

## 配置项说明

后端配置项通过环境变量设置，环境变量名使用下划线 `_` 替换配置路径中的点 `.`。

### 应用配置

| 环境变量 | 说明 | 默认值 | 可选值 |
|---------|------|--------|--------|
| `APP_NAME` | 应用名称 | `love-girl` | - |
| `APP_ENV` | 运行环境 | `prod` | `dev`, `test`, `prod` |
| `APP_INITIALIZED` | 系统是否已初始化 | `false` | `true`, `false` |

### 服务器配置

| 环境变量 | 说明 | 默认值 | 可选值 |
|---------|------|--------|--------|
| `SERVER_ADDR` | 服务监听地址 | `:8182` | 如 `:8181`, `0.0.0.0:8181` |
| `SERVER_SCHEMA` | 协议类型 | `http` | `http`, `https` |
| `SERVER_HOST_NAME` | 主机名 | `localhost:8182` | 如 `example.com:8181` |
| `SERVER_MODE` | Gin 运行模式 | `debug` | `debug`, `release`, `test` |

### 日志配置

| 环境变量 | 说明 | 默认值 | 可选值 |
|---------|------|--------|--------|
| `LOG_LEVEL` | 日志级别 | `debug` | `debug`, `info`, `warn`, `error` |

### 数据库配置

| 环境变量 | 说明 | 默认值 | 可选值 |
|---------|------|--------|--------|
| `DATASOURCE_DATABASE_DRIVER` | 数据库驱动 | `sqlite` | `sqlite`, `mysql`, `postgres` |
| `DATASOURCE_DATABASE_DSN` | 数据库连接字符串 | `./data/love-girl.db` | 根据驱动不同而不同 |

**DSN 示例**:
- SQLite: `file:./data/love-girl.db?_fk=1`
- MySQL: `user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local`
- PostgreSQL: `host=localhost port=5432 user=postgres password=secret dbname=lovegirl sslmode=disable`

### JWT 配置

| 环境变量 | 说明 | 默认值 | 说明 |
|---------|------|--------|------|
| `JWT_SECRET` | JWT 签名密钥 | `love-girl-123456789012345678901234` | 建议至少 32 字符 |
| `JWT_ISSUER` | JWT 签发者 | `love-girl` | - |
| `JWT_EXPIRE` | JWT 过期时间（秒） | `900` | 如 `3600` (1小时) |

### 存储配置

#### 通用存储配置

| 环境变量 | 说明 | 默认值 | 可选值 |
|---------|------|--------|--------|
| `STORAGE_BACKEND` | 存储后端类型 | `local` | `local`, `s3`, `webdav` |

#### 访问策略配置

| 环境变量 | 说明 | 默认值 | 说明 |
|---------|------|--------|------|
| `STORAGE_ACCESS_GIN_PROXY_ENABLED` | Gin 代理访问 | `true` | 通过本地 HTTP 代理访问文件 |
| `STORAGE_ACCESS_IMAGE_PROXY_ENABLED` | 图片代理开关 | `false` | 如 `true` 则需配置 base_url |
| `STORAGE_ACCESS_IMAGE_PROXY_BASE_URL` | 图片代理服务地址 | - | 如 `http://imgproxy:8080` |

#### 本地存储 (STORAGE_BACKEND=local)

| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| `STORAGE_LOCAL_ROOT` | 本地存储根路径 | `./data/uploads` |

#### S3 兼容存储 (STORAGE_BACKEND=s3)

| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| `STORAGE_S3_USE_SSL` | 是否使用 SSL | `false` |
| `STORAGE_S3_ENDPOINT` | S3 端点地址 | - |
| `STORAGE_S3_REGION` | S3 区域 | `us-east-1` |
| `STORAGE_S3_BUCKET` | S3 存储桶名称 | - |
| `STORAGE_S3_CREDENTIALS_ACCESS_KEY_ID` | S3 访问密钥 ID | - |
| `STORAGE_S3_CREDENTIALS_SECRET_ACCESS_KEY` | S3 秘密访问密钥 | - |
| `STORAGE_S3_PUBLIC_BASE_URL` | S3 公共访问地址 | - |
| `STORAGE_S3_PRESIGN_ENABLED` | 启用预签名 URL | `true` |
| `STORAGE_S3_PRESIGN_EXPIRE` | 预签名 URL 过期时间（秒） | `3600` |

#### WebDAV 存储 (STORAGE_BACKEND=webdav)

| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| `STORAGE_WEBDAV_ENDPOINT` | WebDAV 端点地址 | - |
| `STORAGE_WEBDAV_BASE_PATH` | WebDAV 基础路径 | - |
| `STORAGE_WEBDAV_PUBLIC_BASE_URL` | WebDAV 公共访问地址 | - |
| `STORAGE_WEBDAV_AUTH_USERNAME` | WebDAV 认证用户名 | - |
| `STORAGE_WEBDAV_AUTH_PASSWORD` | WebDAV 认证密码 | - |

---

## 环境变量配置示例

### 本地 SQLite + 本地存储

```bash
APP_NAME=love-girl
APP_ENV=production
SERVER_ADDR=:8181
SERVER_HOST_NAME=example.com:8181
SERVER_SCHEMA=https
LOG_LEVEL=info
DATASOURCE_DATABASE_DRIVER=sqlite
DATASOURCE_DATABASE_DSN=file:./data/love-girl.db?_fk=1
JWT_SECRET=your-secret-key-at-least-32-characters-long
JWT_EXPIRE=3600
STORAGE_BACKEND=local
STORAGE_LOCAL_ROOT=./data/uploads
STORAGE_ACCESS_GIN_PROXY_ENABLED=true
```

### MySQL + S3 存储

```bash
APP_ENV=production
DATASOURCE_DATABASE_DRIVER=mysql
DATASOURCE_DATABASE_DSN=root:password@tcp(mysql:3306)/lovegirl?charset=utf8mb4&parseTime=True&loc=Local
STORAGE_BACKEND=s3
STORAGE_S3_ENDPOINT=minio:9000
STORAGE_S3_BUCKET=love-girl
STORAGE_S3_REGION=us-east-1
STORAGE_S3_CREDENTIALS_ACCESS_KEY_ID=minioadmin
STORAGE_S3_CREDENTIALS_SECRET_ACCESS_KEY=minioadmin
STORAGE_S3_PUBLIC_BASE_URL=https://s3.example.com
STORAGE_S3_USE_SSL=false
```

### PostgreSQL + WebDAV 存储

```bash
APP_ENV=production
DATASOURCE_DATABASE_DRIVER=postgres
DATASOURCE_DATABASE_DSN=host=postgres port=5432 user=postgres password=secret dbname=lovegirl sslmode=disable
STORAGE_BACKEND=webdav
STORAGE_WEBDAV_ENDPOINT=https://webdav.example.com
STORAGE_WEBDAV_BASE_PATH=/Resource/Download
STORAGE_WEBDAV_PUBLIC_BASE_URL=https://cdn.example.com
STORAGE_WEBDAV_AUTH_USERNAME=user
STORAGE_WEBDAV_AUTH_PASSWORD=password
```

---

## 配置文件示例

配置文件位置：`./data/configs/config.yaml`

```yaml
app:
  name: love-girl
  env: prod
  initialized: true

server:
  addr: :8182
  schema: http
  host_name: localhost:8182
  mode: debug

log:
  level: debug

datasource:
  database:
    driver: sqlite
    dsn: ./data/love-girl.db

jwt:
  secret: love-girl-123456789012345678901234
  issuer: love-girl
  expire: 900

storage:
  backend: local
  access:
    gin_proxy:
      enabled: true
    image_proxy:
      enabled: false
      base_url: http://localhost:8080
  local:
    root: ./data/uploads
```

---

## 配置热更新

当修改 `config.yaml` 文件后：

1. 系统自动检测到文件变更
2. 进行语法和业务规则验证
3. 如果验证通过：
   - 使用 `syscall.Exec` 替换进程
   - 新进程加载新配置并启动
4. 如果验证失败：
   - 打印错误日志
   - 服务继续按旧配置运行

**注意**：环境变量的优先级高于配置文件。如果同时设置了环境变量和配置文件，环境变量将覆盖配置文件的值。
