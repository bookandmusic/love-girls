# 后端优化计划

> 本文档记录后端项目的优化计划，按优先级和功能模块分类。

## 概述

**分析时间**: 2026-03-13

**项目技术栈**: Go 1.25+ / Gin / GORM / Wire / Viper

**当前状态**: 功能可用，存在安全风险和架构改进空间

---

## ⚠️ 重要约束

> **请严格遵守以下规则**

1. **只修改后端代码**：本次优化仅涉及 `backend/` 目录下的 Go 代码，禁止修改前端代码、文档、配置文件模板等其他内容

2. **修改完成后必须执行代码校验**：
   ```bash
   cd backend
   ./tools/check.sh all
   ```

3. **校验脚本说明** (`backend/tools/check.sh`)：
   
   | 命令 | 说明 |
   |------|------|
   | `./tools/check.sh all` | 执行全部检查（推荐） |
   | `./tools/check.sh format` | 代码格式化 |
   | `./tools/check.sh lint` | 静态检查 (golangci-lint) |
   | `./tools/check.sh complexity` | 复杂度检查 (gocyclo) |
   | `./tools/check.sh gen` | 生成代码 (wire + swag) |

4. **确保校验通过后方可提交**：`check.sh all` 会依次执行：
   - 格式化代码
   - 生成依赖注入代码
   - 生成 API 文档
   - 静态检查
   - 复杂度检查 (函数复杂度 < 15)

---

## TODO 进度跟踪

> 完成后请在对应项打勾 `[x]`

### 一、安全优化 (SEC)

- [ ] **SEC-001** 密码字段暴露给前端 `严重` `internal/model/user.go`
- [ ] **SEC-002** JWT Secret 硬编码默认值 `严重` `internal/config/config.go`
- [ ] **SEC-003** 文件上传接口无认证 `高` `internal/handler/file.go`
- [ ] **SEC-004** 缺少请求速率限制 `中` `中间件层`
- [ ] **SEC-005** 系统初始化接口无速率限制 `高` `internal/handler/system.go`
- [ ] **SEC-006** 缺少安全头部中间件 `中`
- [ ] **SEC-007** 缺少 CORS 配置 `中`
- [ ] **SEC-008** 缺少密码强度校验 `中`
- [ ] **SEC-009** WebDAV 密码明文存储 `中`
- [ ] **SEC-010** Swagger 生产环境暴露 `低`

### 二、基础设施优化 (INF)

- [x] **INF-001** RequestID 追踪 `高`
- [x] **INF-002** 统一错误处理 `高`
- [x] **INF-003** 结构化日志 `中`
- [x] **INF-004** 请求超时控制 `中`

### 三、配置管理 (CFG)

- [ ] **CFG-001** 数据库连接池未配置 `中`
- [ ] **CFG-002** 配置监听使用 fmt.Printf `低`

### 四、API 功能增强 (API)

- [x] **API-001** 查询接口排序支持 `中`
- [x] **API-002** 查询接口过滤支持 `中`
- [x] **API-003** 排序字段白名单 `中`
- [x] **API-004** 过滤字段白名单 `中`
- [x] **API-005** 分页上限限制 `中`

### 五、图片 URL 动态域名 (IMG)

- [x] **IMG-001** 图片 URL 域名写死配置 `中`
- [x] **IMG-002** 缺少多域名部署支持 `中`
- [x] **IMG-003** 移除冗余 gin_proxy.enabled 配置 `中`
- [x] **IMG-004** 四种场景 URL 动态生成逻辑 `中`

### 六、性能优化 (PERF)

- [ ] **PERF-001** N+1 查询问题 `高` `service/system.go`
- [ ] **PERF-002** 点赞功能无并发控制 `高` `repo/moment.go`
- [ ] **PERF-003** 缺少数据库索引优化 `中`
- [ ] **PERF-004** 文件上传无大小业务校验 `中`

### 七、代码质量优化 (CODE)

- [ ] **CODE-001** 零测试覆盖率 `高`
- [ ] **CODE-002** 包名与内置类型冲突 (`internal/error`) `中`
- [ ] **CODE-003** 魔法数字散落各处 `低`
- [ ] **CODE-004** `MustGetAuthClaims` 使用 panic `低`
- [ ] **CODE-005** 变量命名不清晰 (`v`, `v2`) `低`

### 进度统计

| 分类 | 总数 | 已完成 | 进度 |
|------|------|--------|------|
| 安全优化 | 10 | 0 | 0% |
| 基础设施 | 4 | 4 | 100% |
| 配置管理 | 2 | 0 | 0% |
| API 增强 | 5 | 5 | 100% |
| 图片域名 | 4 | 4 | 100% |
| 性能优化 | 4 | 0 | 0% |
| 代码质量 | 5 | 0 | 0% |
| **合计** | **34** | **13** | **38%** |

---

## 一、安全优化

### 1.1 高优先级

| 编号 | 问题 | 风险等级 | 位置 | 状态 |
|------|------|----------|------|------|
| SEC-001 | 密码字段暴露给前端 | 严重 | `internal/model/user.go` | 待修复 |
| SEC-002 | JWT Secret 硬编码默认值 | 严重 | `internal/config/config.go` | 待修复 |
| SEC-003 | 文件上传接口无认证 | 高 | `internal/handler/file.go` | 待修复 |
| SEC-004 | 缺少请求速率限制 | 中 | 中间件层 | 待实现 |
| SEC-005 | 系统初始化接口无速率限制 | 高 | `internal/handler/system.go` | 待修复 |

#### SEC-001: 密码字段暴露给前端

**问题描述**:
```go
// 当前代码
type User struct {
    Password string `gorm:"size:128;not null" json:"password"` // 密码会暴露给前端!
}
```

**修复方案**:
```go
type User struct {
    Password string `gorm:"size:128;not null" json:"-"` // 隐藏敏感字段
}
```

**影响范围**: 用户相关 API 响应

---

#### SEC-002: JWT Secret 硬编码默认值

**问题描述**:
```go
// config/config.go
v.SetDefault("jwt.secret", "love-girl-123456789012345678901234")
```

**修复方案**:
- 移除默认值，生产环境强制要求配置
- 添加启动时校验，生产环境禁止使用默认值
- 配置文件不存在时提示用户设置

```go
// 建议实现
func validateJWTConfig(cfg *JWTConfig, env string) error {
    if env == "prod" && cfg.Secret == "love-girl-123456789012345678901234" {
        return errors.New("生产环境禁止使用默认 JWT Secret")
    }
    return nil
}
```

---

#### SEC-003: 文件上传接口无认证

**问题描述**:
```go
// handler/file.go
func (h *FileHandler) RegisterRoutes(...) {
    fileGroup := apiGroup.Group("/file")
    {
        fileGroup.POST("/upload", h.SaveFile) // 无认证！
    }
}
```

**修复方案**:
```go
func (h *FileHandler) RegisterRoutes(...) {
    fileGroup := apiGroup.Group("/file")
    {
        fileGroup.GET("/:id", h.GetFile) // 公开访问
    }
    
    // 需要认证的路由
    authGroup := apiGroup.Group("/file")
    authGroup.Use(authMiddleware.Handle())
    {
        authGroup.POST("/upload", h.SaveFile) // 需要认证
    }
}
```

---

#### SEC-004: 请求速率限制

**实现方案**:

```go
// middleware/ratelimit.go
package middleware

import (
    "net/http"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
)

type RateLimiter struct {
    limiters sync.Map
    rps      int
    burst    int
}

func NewRateLimiter(rps, burst int) *RateLimiter {
    return &RateLimiter{rps: rps, burst: burst}
}

func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
    if limiter, ok := rl.limiters.Load(key); ok {
        return limiter.(*rate.Limiter)
    }
    limiter := rate.NewLimiter(rate.Limit(rl.rps), rl.burst)
    rl.limiters.Store(key, limiter)
    return limiter
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        key := c.ClientIP()
        if !rl.getLimiter(key).Allow() {
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                "code":    1,
                "message": "请求过于频繁，请稍后再试",
            })
            return
        }
        c.Next()
    }
}
```

---

### 1.2 中优先级

| 编号 | 问题 | 风险等级 | 建议 |
|------|------|----------|------|
| SEC-006 | 缺少安全头部中间件 | 中 | 添加 X-Frame-Options, X-Content-Type-Options 等 |
| SEC-007 | 缺少 CORS 配置 | 中 | 添加 CORS 中间件 |
| SEC-008 | 缺少密码强度校验 | 中 | 添加密码复杂度验证 |
| SEC-009 | WebDAV 密码明文存储 | 中 | 使用环境变量或密钥管理 |
| SEC-010 | Swagger 生产环境暴露 | 低 | 条件启用或限制访问 |

#### SEC-006: 安全头部中间件

```go
// middleware/security.go
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        c.Next()
    }
}
```

---

## 二、基础设施优化

### 2.1 中间件增强

| 编号 | 功能 | 优先级 | 状态 |
|------|------|--------|------|
| INF-001 | RequestID 追踪 | 高 | 待实现 |
| INF-002 | 统一错误处理 | 高 | 待实现 |
| INF-003 | 结构化日志 | 中 | 待实现 |
| INF-004 | 请求超时控制 | 中 | 待实现 |

#### INF-001: RequestID 追踪

**目的**: 便于日志追踪和问题定位

**实现方案**:
```go
// middleware/requestid.go
import "github.com/google/uuid"

func RequestID() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        c.Set("request_id", requestID)
        c.Header("X-Request-ID", requestID)
        c.Next()
    }
}
```

---

#### INF-002: 统一错误处理

**目的**: 减少 Handler 层重复代码

**实现方案**:
```go
// handler/response.go
type AppError struct {
    Code    int
    Message string
    Detail  string
    Cause   error
}

func (e *AppError) Error() string {
    return e.Message
}

func Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, Response{
        Code:    0,
        Message: "success",
        Data:    data,
    })
}

func BadRequest(c *gin.Context, message string, err error, logger *log.Logger) {
    if logger != nil && err != nil {
        logger.Error(message, "error", err, "request_id", c.GetString("request_id"))
    }
    c.JSON(http.StatusBadRequest, Response{
        Code:    1,
        Message: message,
    })
}

func InternalError(c *gin.Context, message string, err error, logger *log.Logger) {
    if logger != nil {
        logger.Error(message, "error", err, "request_id", c.GetString("request_id"))
    }
    c.JSON(http.StatusInternalServerError, Response{
        Code:    1,
        Message: "系统内部错误",
    })
}
```

---

### 2.2 配置管理

| 编号 | 问题 | 优先级 | 状态 |
|------|------|--------|------|
| CFG-001 | 数据库连接池未配置 | 中 | 待实现 |
| CFG-002 | 配置监听使用 fmt.Printf | 低 | 待优化 |

#### CFG-001: 数据库连接池配置

**实现方案**:
```go
// config/config.go
type DatabaseConfig struct {
    Driver       string `mapstructure:"driver"`
    DSN          string `mapstructure:"dsn"`
    MaxIdleConns int    `mapstructure:"max_idle_conns"` // 新增
    MaxOpenConns int    `mapstructure:"max_open_conns"` // 新增
    ConnMaxLife  int    `mapstructure:"conn_max_life"`  // 新增，单位秒
}

// db/db.go
func NewDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
    // ... 现有连接代码 ...
    
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    
    if cfg.MaxIdleConns > 0 {
        sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
    }
    if cfg.MaxOpenConns > 0 {
        sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
    }
    if cfg.ConnMaxLife > 0 {
        sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLife) * time.Second)
    }
    
    return db, nil
}
```

**默认值建议**:
```go
v.SetDefault("datasource.database.max_idle_conns", 10)
v.SetDefault("datasource.database.max_open_conns", 100)
v.SetDefault("datasource.database.conn_max_life", 3600)
```

---

## 三、功能增强

### 3.1 查询接口优化

| 编号 | 功能 | 优先级 | 状态 |
|------|------|--------|------|
| API-001 | 查询接口排序支持 | 中 | 待实现 |
| API-002 | 查询接口过滤支持 | 中 | 待实现 |
| API-003 | 排序字段白名单 | 中 | 待实现 |
| API-004 | 过滤字段白名单 | 中 | 待实现 |
| API-005 | 分页上限限制 | 中 | 待实现 |

#### API-001 & API-002: 排序和过滤支持

**当前状态**:
- Repo 层已支持 `WithOrder()` 和 `FilterCondition`
- Handler 层未暴露给 API

**改造方案**:

```go
// handler/query.go
package handler

import (
    "github.com/gin-gonic/gin"
    "github.com/bookandmusic/love-girl/internal/repo"
)

type QueryParams struct {
    Page    int
    Size    int
    SortBy  string
    Order   string // "asc" or "desc"
    Filters []repo.FilterCondition
}

var AllowedSortFields = map[string][]string{
    "moments":      {"created_at", "likes"},
    "places":       {"created_at", "name"},
    "wishes":       {"created_at"},
    "anniversaries": {"date", "created_at"},
    "albums":       {"created_at", "name"},
}

var AllowedFilterFields = map[string]map[string][]string{
    "moments": {
        "is_public": {"eq"},
        "user_id":   {"eq"},
        "likes":     {"eq", "gt", "lt", "gte", "lte"},
    },
    "wishes": {
        "approved": {"eq"},
    },
}

func ParseQueryParams(c *gin.Context, resource string) *QueryParams {
    params := &QueryParams{
        Page:  1,
        Size:  10,
        Order: "desc",
    }
    
    // 分页参数
    if page := c.Query("page"); page != "" {
        if p, err := strconv.Atoi(page); err == nil && p > 0 {
            params.Page = p
        }
    }
    if size := c.Query("size"); size != "" {
        if s, err := strconv.Atoi(size); err == nil && s > 0 && s <= 100 {
            params.Size = s
        }
    }
    
    // 排序参数
    if sortBy := c.Query("sort_by"); sortBy != "" {
        if isAllowedField(AllowedSortFields[resource], sortBy) {
            params.SortBy = sortBy
        }
    }
    if order := c.Query("order"); order == "asc" || order == "desc" {
        params.Order = order
    }
    
    // 过滤参数: ?filter=likes:gt:10&filter=is_public:eq:true
    if filters := c.QueryArray("filter"); len(filters) > 0 {
        for _, f := range filters {
            if cond := parseFilter(f, resource); cond != nil {
                params.Filters = append(params.Filters, *cond)
            }
        }
    }
    
    return params
}

func isAllowedField(allowed []string, field string) bool {
    for _, f := range allowed {
        if f == field {
            return true
        }
    }
    return false
}

func parseFilter(filterStr, resource string) *repo.FilterCondition {
    parts := strings.Split(filterStr, ":")
    if len(parts) != 3 {
        return nil
    }
    
    field, op, value := parts[0], parts[1], parts[2]
    
    // 检查字段和操作符是否允许
    allowedOps, ok := AllowedFilterFields[resource][field]
    if !ok {
        return nil
    }
    if !isAllowedField(allowedOps, op) {
        return nil
    }
    
    return &repo.FilterCondition{
        Field:    field,
        Operator: op,
        Value:    parseValue(value),
    }
}
```

**API 使用示例**:
```
GET /api/v1/moments?page=1&size=20&sort_by=likes&order=desc&filter=is_public:eq:true

GET /api/v1/wishes?page=1&size=10&sort_by=created_at&order=desc&filter=approved:eq:true
```

---

#### API-005: 分页上限限制

**实现方案**:
```go
// handler/query.go
const (
    DefaultPageSize = 10
    MaxPageSize     = 100
)

func ParsePagination(page, size int) (int, int) {
    if page < 1 {
        page = 1
    }
    if size < 1 {
        size = DefaultPageSize
    }
    if size > MaxPageSize {
        size = MaxPageSize
    }
    return page, size
}
```

---

### 3.2 图片 URL 动态域名

| 编号 | 问题 | 优先级 | 状态 |
|------|------|--------|------|
| IMG-001 | 图片 URL 域名写死配置 | 中 | 待实现 |
| IMG-002 | 缺少多域名部署支持 | 中 | 待实现 |
| IMG-003 | image_proxy.base_url 无法动态调整 | 中 | 待实现 |

#### 当前架构

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   前端      │────▶│    Gin       │────▶│   存储      │
│  (Browser)  │     │  (代理访问)   │     │ (local/S3) │
└─────────────┘     └──────────────┘     └─────────────┘
       │                   │
       │                   ▼
       │           ┌──────────────┐
       └──────────▶│  ImageProxy  │
                   │ (imgproxy)   │
                   └──────────────┘
```

#### 核心原则

1. **原始图片**: 如果存储系统有公开访问链接，优先使用存储本身的链接
2. **缩略图**: 统一由 ImageProxy 生成，只是 ImageProxy 是否有公开链接的区别

#### 四种部署场景详解

#### 配置简化

**原有配置（冗余）**:
```yaml
storage:
  access:
    gin_proxy:
      enabled: true/false      # 冗余！可动态决定
    image_proxy:
      enabled: true/false      # 冗余！
      base_url: "..."
```

**简化后配置**:
```yaml
storage:
  backend: s3/local/webdav
  access:
    image_proxy:
      base_url: "..."          # 空 = ImageProxy 不公开，非空 = 公开
  s3:
    public_base_url: "..."     # 空 = 存储不公开，非空 = 公开
  webdav:
    public_base_url: "..."     # 同上
  # local 存储无 public_base_url，永远不公开
```

**动态决策逻辑**:
```
是否需要 Gin 代理原始图片？
  └─ 存储有公开链接？ 否 → 需要 Gin 代理

缩略图走哪里？
  └─ ImageProxy 有公开链接？
       ├─ 是 → ImageProxy 公开链接
       └─ 否 → Gin 代理转发到内网 ImageProxy
```

---

#### 四种部署场景详解

##### 场景一：ImageProxy 公开 + 存储公开

**配置**:
```yaml
storage:
  backend: s3
  access:
    image_proxy:
      base_url: https://img.example.com  # 非空 = ImageProxy 公开
  s3:
    public_base_url: https://cdn.example.com  # 非空 = 存储公开
```

**动态判断**:
- 存储公开 ✓ → 原始图片直连 CDN
- ImageProxy 公开 ✓ → 缩略图走 ImageProxy 公开链接

**URL 生成逻辑**:

| 图片类型 | URL 来源 | 示例 |
|----------|----------|------|
| 原始图片 | 存储公开链接 (CDN) | `https://cdn.example.com/path/file.jpg` |
| 缩略图 | ImageProxy + 存储链接 | `https://img.example.com/200x200/https://cdn.example.com/path/file.jpg` |

**数据流**:
```
原始图片: 前端 ──直连──▶ CDN (存储)
缩略图:   前端 ──▶ ImageProxy ──▶ CDN (存储)
```

**适用场景**: 生产环境，CDN + 图片处理服务，性能最优

---

##### 场景二：ImageProxy 公开 + 存储不公开

**配置**:
```yaml
storage:
  backend: s3
  access:
    image_proxy:
      base_url: https://img.example.com  # 非空 = ImageProxy 公开
  s3:
    public_base_url: ""                  # 空 = 存储不公开（内网 MinIO）
```

**动态判断**:
- 存储不公开 ✗ → 原始图片走 Gin 代理
- ImageProxy 公开 ✓ → 缩略图走 ImageProxy 公开链接

**URL 生成逻辑**:

| 图片类型 | URL 来源 | 示例 |
|----------|----------|------|
| 原始图片 | Gin 代理链接 | `https://api.example.com/api/v1/file/123` |
| 缩略图 | ImageProxy + Gin 链接 | `https://img.example.com/200x200/https://api.example.com/api/v1/file/123` |

**数据流**:
```
原始图片: 前端 ──▶ Gin ──▶ 存储 (内网)
缩略图:   前端 ──▶ ImageProxy ──▶ Gin ──▶ 存储 (内网)
```

**适用场景**: 内网存储 + 外部图片处理服务

---

##### 场景三：ImageProxy 不公开 + 存储公开

**配置**:
```yaml
storage:
  backend: s3
  access:
    image_proxy:
      base_url: ""                       # 空 = ImageProxy 不公开（内网部署）
  s3:
    public_base_url: https://cdn.example.com  # 非空 = 存储公开
```

**动态判断**:
- 存储公开 ✓ → 原始图片直连 CDN
- ImageProxy 不公开 ✗ → 缩略图走 Gin 代理转发

**URL 生成逻辑**:

| 图片类型 | URL 来源 | 示例 |
|----------|----------|------|
| 原始图片 | 存储公开链接 (CDN) | `https://cdn.example.com/path/file.jpg` |
| 缩略图 | Gin 代理（转发到内网 ImageProxy） | `https://api.example.com/api/v1/file/123?w=200&h=200` |

**数据流**:
```
原始图片: 前端 ──直连──▶ CDN (存储)
缩略图:   前端 ──▶ Gin ──▶ ImageProxy (内网) ──▶ CDN (存储)
```

**适用场景**: CDN 直连原图，内网图片处理服务生成缩略图

**实现要点**: Gin 收到带 `w/h` 参数的请求时，转发到内网 ImageProxy

---

##### 场景四：ImageProxy 不公开 + 存储不公开

**配置**:
```yaml
storage:
  backend: local  # 或内网 s3
  access:
    image_proxy:
      base_url: ""                       # 空 = ImageProxy 不公开
  local:
    root: ./data/uploads                 # local 存储无 public_base_url
```

**动态判断**:
- 存储不公开 (local) → 原始图片走 Gin 代理
- ImageProxy 不公开 ✗ → 缩略图走 Gin 代理转发

**URL 生成逻辑**:

| 图片类型 | URL 来源 | 示例 |
|----------|----------|------|
| 原始图片 | Gin 代理链接 | `https://api.example.com/api/v1/file/123` |
| 缩略图 | Gin 代理（转发到内网 ImageProxy） | `https://api.example.com/api/v1/file/123?w=200&h=200` |

**数据流**:
```
原始图片: 前端 ──▶ Gin ──▶ 存储 (本地/内网)
缩略图:   前端 ──▶ Gin ──▶ ImageProxy (内网) ──▶ 存储
```

**适用场景**: 开发环境、小型部署，全链路走 Gin 代理

---

#### 四种场景汇总表

| 场景 | ImageProxy | 存储公开 | 原始图片 | 缩略图 | Gin 代理 |
|------|-----------|---------|---------|--------|---------|
| 1 | ✅ 公开 | ✅ 公开 | CDN 直连 | ImageProxy 公开 | 仅缩略图回源 |
| 2 | ✅ 公开 | ❌ 不公开 | Gin 代理 | ImageProxy 公开 | 原始图必需 |
| 3 | ❌ 不公开 | ✅ 公开 | CDN 直连 | Gin 转发 | 缩略图必需 |
| 4 | ❌ 不公开 | ❌ 不公开 | Gin 代理 | Gin 转发 | 全部必需 |

> **Gin 代理是否需要** 由存储是否有公开链接 **自动决定**，无需手动配置

---

#### 动态域名支持

**问题**: 当前域名写死在配置中，无法根据前端请求域名动态返回对应 URL

**需求场景**:
- 内网访问: `http://192.168.1.100:8182/api/v1/file/123`
- 外网访问: `https://api.example.com/api/v1/file/123`
- 同一服务需要支持多域名访问

**实现方案**:

```go
// service/file.go

// URLConfig 用于存储动态域名配置
type URLConfig struct {
    // ImageProxy 是否有公开可访问的 URL
    ImageProxyPublicURL string
    // 存储是否有公开可访问的 URL
    StoragePublicURL string
}

// GetURLConfig 根据请求动态获取 URL 配置
func (s *FileService) GetURLConfig(c *gin.Context) *URLConfig {
    host := s.getRequestHost(c)
    
    // 检查多域名配置
    if cfg, ok := s.domainConfigs[host]; ok {
        return cfg
    }
    
    // 使用默认配置
    return &URLConfig{
        ImageProxyPublicURL: s.access.ImageProxy.BaseURL,
        StoragePublicURL:    s.getStoragePublicURL(),
    }
}

// getRequestHost 从请求中获取主机名，支持代理头
func (s *FileService) getRequestHost(c *gin.Context) string {
    // 优先使用 X-Forwarded-Host（反向代理场景）
    if host := c.GetHeader("X-Forwarded-Host"); host != "" {
        return host
    }
    return c.Request.Host
}

// getRequestScheme 从请求中获取协议，支持代理头
func (s *FileService) getRequestScheme(c *gin.Context) string {
    // 优先使用 X-Forwarded-Proto
    if proto := c.GetHeader("X-Forwarded-Proto"); proto != "" {
        return proto
    }
    if c.Request.TLS != nil {
        return "https"
    }
    return "http"
}

// BuildGinProxyURL 构建 Gin 代理 URL（动态域名）
func (s *FileService) BuildGinProxyURL(c *gin.Context, fileID uint64) string {
    scheme := s.getRequestScheme(c)
    host := s.getRequestHost(c)
    return fmt.Sprintf("%s://%s/api/v1/file/%d", scheme, host, fileID)
}

// GetOriginalImageURL 获取原始图片 URL
func (s *FileService) GetOriginalImageURL(ctx context.Context, c *gin.Context, file *model.File) (string, error) {
    cfg := s.GetURLConfig(c)
    
    // 存储有公开链接，优先使用存储链接
    if cfg.StoragePublicURL != "" {
        return s.Storage.URL(ctx, file.ID, file.Path, 0, 0, nil)
    }
    
    // 否则使用 Gin 代理
    ginURL := s.BuildGinProxyURL(c, file.ID)
    return s.Storage.URL(ctx, file.ID, file.Path, 0, 0, func(id uint64) string {
        return ginURL
    })
}

// GetThumbnailURL 获取缩略图 URL
func (s *FileService) GetThumbnailURL(ctx context.Context, c *gin.Context, file *model.File, width, height int) (string, error) {
    cfg := s.GetURLConfig(c)
    
    // 先获取原始图片 URL
    originalURL, err := s.GetOriginalImageURL(ctx, c, file)
    if err != nil {
        return "", err
    }
    
    // ImageProxy 有公开链接
    if cfg.ImageProxyPublicURL != "" {
        escapedURL := url.QueryEscape(originalURL)
        return fmt.Sprintf("%s/%dx%d/%s", cfg.ImageProxyPublicURL, width, height, escapedURL), nil
    }
    
    // ImageProxy 无公开链接，返回 Gin 代理 URL（带尺寸参数）
    // Gin 会负责转发到内网 ImageProxy
    ginURL := s.BuildGinProxyURL(c, file.ID)
    return fmt.Sprintf("%s?w=%d&h=%d", ginURL, width, height), nil
}
```

**多域名配置示例**:
```yaml
server:
  schema: http
  host_name: localhost:8182

storage:
  backend: s3
  access:
    image_proxy:
      base_url: https://img.example.com  # 默认：ImageProxy 公开
  s3:
    public_base_url: https://cdn.example.com  # 默认：存储公开

  # 新增：多域名配置（覆盖默认值）
  domains:
    # 外网访问
    - host: api.example.com
      image_proxy_public_url: https://img.example.com
      storage_public_url: https://cdn.example.com
    # 内网访问（无公开链接，自动走 Gin 代理）
    - host: 192.168.1.100:8182
      image_proxy_public_url: ""
      storage_public_url: ""
    # 本地开发
    - host: localhost:8182
      image_proxy_public_url: ""
      storage_public_url: ""
```

**改动文件清单**:
```
internal/config/app.go       # 新增 DomainConfig 结构
internal/config/config.go    # 新增 domains 配置解析
internal/service/file.go     # 动态 URL 构建逻辑
internal/storage/base.go     # URL 方法可选 builder
internal/storage/s3.go       # 适配动态 URL
internal/storage/webdav.go   # 适配动态 URL
internal/storage/local.go    # 适配动态 URL
internal/handler/file.go     # 支持 w/h 参数转发 ImageProxy
```

---

## 四、性能优化

| 编号 | 问题 | 优先级 | 位置 | 状态 |
|------|------|--------|------|------|
| PERF-001 | N+1 查询问题 | 高 | `service/system.go` | 待修复 |
| PERF-002 | 点赞功能无并发控制 | 高 | `repo/moment.go` | 待修复 |
| PERF-003 | 缺少数据库索引优化 | 中 | 迁移文件 | 待实现 |
| PERF-004 | 文件上传无大小业务校验 | 中 | `handler/file.go` | 待修复 |

#### PERF-001: N+1 查询问题

**问题描述**:
```go
// service/system.go - GetSystemInfo
users, _ := s.UserService.ListUsers(ctx)
for _, user := range users {
    // 每次循环查询头像，造成 N+1
    user.Avatar = s.FileService.BuildFileResponse(ctx, user.Avatar)
}
```

**修复方案**:
```go
// 在 ListUsers 时预加载 Avatar
func (r *UserRepo) ListUsers(ctx context.Context) ([]model.User, error) {
    var users []model.User
    err := r.db.WithContext(ctx).
        Preload("Avatar").
        Find(&users).Error
    return users, err
}
```

---

#### PERF-002: 点赞功能并发控制

**问题描述**:
```go
// repo/moment.go
func (r *MomentRepo) UpdateLike(ctx context.Context, id uint64) error {
    return r.db.Model(&model.Moment{}).
        Where("id = ?", id).
        UpdateColumn("likes", gorm.Expr("likes + ?", 1)).Error
    // 高并发时可能丢失更新
}
```

**修复方案 A: 使用行锁**:
```go
func (r *MomentRepo) UpdateLike(ctx context.Context, id uint64) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        var m model.Moment
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
            First(&m, id).Error; err != nil {
            return err
        }
        return tx.Model(&m).Update("likes", m.Likes+1).Error
    })
}
```

**修复方案 B: 使用 Redis 原子计数**:
```go
// service/moment.go
func (s *MomentService) LikeMoment(ctx context.Context, id uint64) (*model.Moment, error) {
    // 先在 Redis 中增加计数
    key := fmt.Sprintf("moment:likes:%d", id)
    count, err := s.redis.Incr(ctx, key).Result()
    if err != nil {
        return nil, err
    }
    
    // 异步同步到数据库
    go s.syncLikesToDB(id, count)
    
    return &model.Moment{ID: id, Likes: int(count)}, nil
}
```

---

## 五、代码质量优化

| 编号 | 问题 | 优先级 | 状态 |
|------|------|--------|------|
| CODE-001 | 零测试覆盖率 | 高 | 待实现 |
| CODE-002 | 包名与内置类型冲突 (`internal/error`) | 中 | 待修复 |
| CODE-003 | 魔法数字散落各处 | 低 | 待优化 |
| CODE-004 | `MustGetAuthClaims` 使用 panic | 低 | 待修复 |
| CODE-005 | 变量命名不清晰 (`v`, `v2`) | 低 | 待优化 |

#### CODE-001: 测试体系建设

**测试优先级**:
1. `internal/utils/password.go` - 密码加密测试
2. `internal/auth/jwt.go` - JWT 生成和解析测试
3. `internal/service/` - 业务逻辑单元测试
4. `internal/handler/` - HTTP 处理器测试

**示例测试**:
```go
// internal/utils/password_test.go
package utils

import "testing"

func TestHashPassword(t *testing.T) {
    password := "test123456"
    hash, err := HashPassword(password)
    if err != nil {
        t.Fatalf("HashPassword failed: %v", err)
    }
    if hash == "" {
        t.Fatal("hash should not be empty")
    }
    if hash == password {
        t.Fatal("hash should not equal to plain password")
    }
}

func TestCheckPassword(t *testing.T) {
    password := "test123456"
    hash, _ := HashPassword(password)
    
    if !CheckPassword(password, hash) {
        t.Fatal("password should match")
    }
    if CheckPassword("wrong", hash) {
        t.Fatal("wrong password should not match")
    }
}
```

---

#### CODE-002: 包名冲突

**问题描述**: `internal/error` 与 Go 内置 `error` 类型冲突

**修复方案**: 重命名为 `internal/errmsg` 或 `internal/errors`

**改动范围**:
```
internal/error/         → internal/errmsg/
├── db.go              → errmsg/db.go
└── ...

所有引用该包的文件需更新 import 路径
```

---

## 六、实施计划

### Phase 1: 安全修复（立即执行）

**预计工时**: 2-3 天

| 任务 | 优先级 | 预计时间 |
|------|--------|----------|
| SEC-001 密码字段隐藏 | 高 | 0.5h |
| SEC-002 JWT Secret 配置强制 | 高 | 1h |
| SEC-003 文件上传认证 | 高 | 1h |
| SEC-004 速率限制中间件 | 中 | 2h |

---

### Phase 2: 基础设施（短期）

**预计工时**: 3-5 天

| 任务 | 优先级 | 预计时间 |
|------|--------|----------|
| INF-001 RequestID 中间件 | 高 | 2h |
| INF-002 统一错误处理 | 高 | 3h |
| SEC-006 安全头部中间件 | 中 | 1h |
| SEC-007 CORS 中间件 | 中 | 1h |
| CFG-001 数据库连接池配置 | 中 | 1h |
| CODE-001 核心模块测试 | 高 | 4h |

---

### Phase 3: 功能增强（中期）

**预计工时**: 5-7 天

| 任务 | 优先级 | 预计时间 |
|------|--------|----------|
| API-001/002 排序过滤支持 | 中 | 4h |
| API-003/004 字段白名单 | 中 | 2h |
| API-005 分页上限限制 | 中 | 1h |
| IMG-001 动态域名支持 | 中 | 4h |
| IMG-002/003 多域名配置 | 中 | 3h |
| API 文档更新 | 中 | 2h |

---

### Phase 4: 性能优化（中期）

**预计工时**: 3-5 天

| 任务 | 优先级 | 预计时间 |
|------|--------|----------|
| PERF-001 N+1 查询修复 | 高 | 2h |
| PERF-002 点赞并发控制 | 高 | 3h |
| PERF-003 数据库索引优化 | 中 | 2h |
| PERF-004 文件上传校验 | 中 | 1h |

---

### Phase 5: 代码质量（长期）

**预计工时**: 持续进行

| 任务 | 优先级 | 预计时间 |
|------|--------|----------|
| CODE-002 包名重构 | 中 | 2h |
| CODE-003 常量提取 | 低 | 2h |
| CODE-004 panic 改造 | 低 | 1h |
| INF-003 结构化日志 | 中 | 3h |
| 测试覆盖率提升 | 中 | 持续 |

---

## 七、风险与注意事项

### 7.1 兼容性风险

| 变更 | 影响范围 | 缓解措施 |
|------|----------|----------|
| SEC-001 密码字段隐藏 | 前端可能依赖该字段 | 检查前端代码，确认无依赖 |
| API-001/002 新增参数 | 现有 API 调用 | 参数可选，向后兼容 |
| IMG-001 动态域名 | URL 格式可能变化 | 保持默认行为不变 |

### 7.2 测试建议

每次变更后执行:
```bash
# 后端检查
cd backend
./tools/check.sh all

# 运行测试
go test ./...

# 手动测试
# 1. 用户登录/注册
# 2. 文件上传/访问
# 3. 动态列表/创建/点赞
# 4. 系统初始化
```

---

## 八、变更记录

| 日期 | 版本 | 变更内容 |
|------|------|----------|
| 2026-03-14 | v1.4 | 完成基础设施优化 (INF-001~004)：RequestID 追踪、统一错误处理、结构化日志、请求超时控制 |
| 2026-03-14 | v1.3 | 完成图片 URL 动态域名计划 (IMG-001~004)，移除静态域名配置，实现动态 URL 生成 |
| 2026-03-13 | v1.2 | 完成 API 增强计划 (API-001~005)，新增 handler/query.go 查询参数解析，更新各 Handler 支持排序和过滤 |
| 2026-03-13 | v1.1 | 添加 TODO 进度跟踪，优化图片 URL 配置说明 |
| 2026-03-13 | v1.0 | 初始版本，完成全面分析 |

> **更新进度后**：请在"进度统计"表格中更新数字

---

> 本文档将随优化进度持续更新。