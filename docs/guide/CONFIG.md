# 配置说明

## 数据目录

所有数据统一存放在 `DATA_DIR` 目录，只需配置顶级目录，子目录路径由程序自动计算：

```
{DATA_DIR}/                       # 默认 ./data/
├── configs/                      # 配置文件目录（自动创建）
│   └── config.yaml               # 主配置文件（自动生成）
├── uploads/                      # 上传文件目录（自动创建）
└── love-girl.db                  # SQLite 数据库（自动创建）
```

---

## 配置方式

### 环境变量（推荐）

所有配置项支持环境变量，格式为**全大写 + 下划线**：

```bash
# 设置数据目录
export DATA_DIR=/var/lib/love-girl

# 使用 MySQL 数据库
export DATABASE_DRIVER=mysql
export DATABASE_DSN="user:pass@tcp(localhost:3306)/lovegirl?charset=utf8mb4&parseTime=True"

# 使用 S3 存储
export STORAGE_BACKEND=s3
export STORAGE_S3_ENDPOINT=s3.amazonaws.com
export STORAGE_S3_BUCKET=my-bucket
export STORAGE_S3_ACCESS_KEY_ID=your-access-key
export STORAGE_S3_SECRET_ACCESS_KEY=your-secret-key
```

### 配置文件

配置文件位置：`{DATA_DIR}/configs/config.yaml`

```yaml
# ===========================================
# 应用配置
# ===========================================
app:
  name: love-girl          # 应用名称
  env: prod                # 环境: dev / test / prod

# ===========================================
# 数据目录（只需配置顶级目录）
# ===========================================
data_dir: ./data           # 所有数据存放根目录

# ===========================================
# HTTP 服务器配置
# ===========================================
server:
  addr: :8182              # 监听地址
  mode: release            # Gin 模式: debug / release / test
  request_timeout: 0       # 请求超时（秒），0 表示不限制

# ===========================================
# 日志配置
# ===========================================
log:
  level: info              # 日志级别: debug / info / warn / error
  format: text             # 输出格式: text / json
  output: stdout           # 输出位置: stdout 或文件路径

# ===========================================
# 数据库配置
# ===========================================
datasource:
  database:
    driver: sqlite         # 数据库类型: sqlite / mysql / postgres
    dsn: ""                # 连接串（SQLite 自动使用 {DATA_DIR}/love-girl.db）

# ===========================================
# JWT 认证配置
# ===========================================
jwt:
  secret: ""               # JWT 密钥（留空自动生成 64 字符随机密钥）
  issuer: love-girl        # 签发者
  expire: 86400            # Token 有效期（秒），默认 24 小时

# ===========================================
# 存储配置
# ===========================================
storage:
  backend: local           # 存储类型: local / s3 / webdav
  # Local 存储路径由 DATA_DIR 自动计算: {DATA_DIR}/uploads

  # --- S3 存储（backend: s3 时配置）---
  s3:
    use_ssl: true          # 是否使用 SSL
    endpoint: ""           # S3 端点地址
    region: ""             # 区域
    bucket: ""             # 存储桶名称
    credentials:
      access_key_id: ""        # Access Key
      secret_access_key: ""    # Secret Key
    public_url: ""        # 公开访问地址（可选）
    presign_enable: false # 是否启用预签名 URL
    presign_expire: 3600  # 预签名 URL 有效期（秒）

  # --- WebDAV 存储（backend: webdav 时配置）---
  webdav:
    endpoint: ""           # WebDAV 端点地址
    base_path: ""          # 基础路径
    public_url: ""         # 公开访问地址（可选）
    auth:
      username: ""         # 认证用户名
      password: ""         # 认证密码

# ===========================================
# 图片代理配置（可选）
# ===========================================
image_proxy:
  internal_url: ""         # 内网地址，Gin 转发缩略图用
  public_url: ""           # 公开地址，前端直接访问
```

### 配置优先级

1. **环境变量**（最高优先级）
2. **配置文件**
3. **默认值**

---

## 环境变量一览

### 基础配置

| 环境变量 | 默认值 | 说明 |
|----------|--------|------|
| `DATA_DIR` | `./data` | 数据目录根路径，所有子目录自动计算 |
| `APP_NAME` | `love-girl` | 应用名称 |
| `APP_ENV` | `prod` | 环境：`dev` / `test` / `prod` |

### 服务器配置

| 环境变量 | 默认值 | 说明 |
|----------|--------|------|
| `SERVER_ADDR` | `:8182` | 监听地址 |
| `SERVER_MODE` | `debug` | Gin 模式：`debug` / `release` / `test` |
| `SERVER_REQUEST_TIMEOUT` | `0` | 请求超时（秒），0 不限制 |

### 日志配置

| 环境变量 | 默认值 | 说明 |
|----------|--------|------|
| `LOG_LEVEL` | `debug` | 日志级别：`debug` / `info` / `warn` / `error` |
| `LOG_FORMAT` | `text` | 输出格式：`text` / `json` |
| `LOG_OUTPUT` | `stdout` | 输出位置：`stdout` 或文件路径 |

### 数据库配置

| 环境变量 | 默认值 | 说明 |
|----------|--------|------|
| `DATABASE_DRIVER` | `sqlite` | 数据库类型：`sqlite` / `mysql` / `postgres` |
| `DATABASE_DSN` | 自动计算 | 连接串（SQLite 自动使用 `{DATA_DIR}/love-girl.db`） |

**DSN 示例**：

```bash
# MySQL
user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True

# PostgreSQL
postgres://user:password@localhost:5432/dbname?sslmode=disable
```

### JWT 配置

| 环境变量 | 默认值 | 说明 |
|----------|--------|------|
| `JWT_SECRET` | 自动生成 | JWT 密钥（64字符随机），首次启动生成并持久化 |
| `JWT_ISSUER` | `love-girl` | 签发者 |
| `JWT_EXPIRE` | `86400` | Token 有效期（秒），默认 24 小时 |

### 存储配置

| 环境变量 | 默认值 | 说明 |
|----------|--------|------|
| `STORAGE_BACKEND` | `local` | 存储类型：`local` / `s3` / `webdav` |

**Local 存储**：路径由 `DATA_DIR` 自动计算为 `{DATA_DIR}/uploads`，不支持单独配置。

#### S3 存储

| 环境变量 | 必须 | 说明 |
|----------|:----:|------|
| `STORAGE_S3_ENDPOINT` | ✅ | S3 端点地址 |
| `STORAGE_S3_BUCKET` | ✅ | 存储桶名称 |
| `STORAGE_S3_ACCESS_KEY_ID` | ✅ | Access Key |
| `STORAGE_S3_SECRET_ACCESS_KEY` | ✅ | Secret Key |
| `STORAGE_S3_USE_SSL` | ❌ | 是否使用 SSL（默认 `false`） |
| `STORAGE_S3_REGION` | ❌ | 区域 |
| `STORAGE_S3_PUBLIC_URL` | ❌ | 公开访问地址 |
| `STORAGE_S3_PRESIGN_ENABLE` | ❌ | 是否启用预签名 URL |
| `STORAGE_S3_PRESIGN_EXPIRE` | ❌ | 预签名 URL 有效期（秒） |

#### WebDAV 存储

| 环境变量 | 必须 | 说明 |
|----------|:----:|------|
| `STORAGE_WEBDAV_ENDPOINT` | ✅ | WebDAV 端点地址 |
| `STORAGE_WEBDAV_BASE_PATH` | ✅ | 基础路径 |
| `STORAGE_WEBDAV_PUBLIC_URL` | ❌ | 公开访问地址 |
| `STORAGE_WEBDAV_AUTH_USERNAME` | ❌ | 认证用户名 |
| `STORAGE_WEBDAV_AUTH_PASSWORD` | ❌ | 认证密码 |

### 图片代理配置

| 环境变量 | 说明 |
|----------|------|
| `IMAGE_PROXY_INTERNAL_URL` | 内网地址，Gin 转发缩略图用 |
| `IMAGE_PROXY_PUBLIC_URL` | 公开地址，前端直接访问 |

---

## 配置热更新

修改 `{DATA_DIR}/configs/config.yaml` 后：

1. 自动检测文件变更
2. 验证配置有效性
3. 验证通过后零停机重启
4. 验证失败则继续使用旧配置

**注意**：环境变量优先级高于配置文件，重启后环境变量会覆盖配置文件的值。

---

## 技术架构

配置系统采用自动化管理架构：

```
┌─────────────────────────────────────────────────────────┐
│                    配置加载流程                          │
├─────────────────────────────────────────────────────────┤
│  默认值 → 环境变量 → 配置文件 → 运行时计算              │
├─────────────────────────────────────────────────────────┤
│  数据目录初始化 → JWT Secret 生成 → 目录结构创建        │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│                    配置热更新流程                        │
├─────────────────────────────────────────────────────────┤
│  文件变更 → fsnotify 监听 → 语法校验 → 业务校验        │
│      ↓                                                  │
│  验证通过 → syscall.Exec 进程替换 → 零停机重启          │
│  验证失败 → 打印错误 → 继续使用旧配置                   │
└─────────────────────────────────────────────────────────┘
```

| 组件 | 职责 |
|------|------|
| Viper | 配置存储、文件监听、环境变量绑定 |
| validator | 配置校验、业务规则验证 |
| fsnotify | 文件变更监听 |
| syscall.Exec | 进程替换、零停机重启 |