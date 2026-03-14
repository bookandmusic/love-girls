# 基础设施优化计划

> 本文档是后端优化计划第二部分"基础设施优化"的详细实施方案。

## 概述

**目标**: 增强后端基础设施，提升可观测性、错误处理能力和系统稳定性

**涉及模块**:
- `internal/middleware/` - 中间件层
- `internal/handler/` - 处理器层
- `internal/config/` - 配置层
- `internal/log/` - 日志层（已存在）

**预计工时**: 1-2 天

---

## ⚠️ 重要约束

> **请严格遵守以下规则**

1. **只修改后端代码**：本次优化仅涉及 `backend/` 目录下的 Go 代码

2. **修改完成后必须执行代码校验**：
   ```bash
   cd backend
   ./tools/check.sh all
   ```

3. **已存在的组件（避免重复创建）**：
   - `internal/log/log.go` - 已有 slog 日志封装
   - `internal/handler/base.go` - 已有 Response 结构体
   - `internal/server/server.go` - 已使用 gin.Recovery()

---

## 任务清单

| 编号 | 任务 | 优先级 | 状态 | 备注 |
|------|------|--------|------|------|
| INF-001 | RequestID 追踪中间件 | 高 | ✅ 已完成 | 新建 requestid.go |
| INF-002 | 统一错误处理机制 | 高 | ✅ 已完成 | 新建 response.go、recovery.go |
| INF-003 | 结构化日志增强 | 中 | ✅ 已完成 | 扩展 log.go，新建 logging.go |
| INF-004 | 请求超时控制 | 中 | ✅ 已完成 | 新建 timeout.go，扩展 ServerConfig |

---

## INF-001: RequestID 追踪中间件

### 目标

为每个请求生成唯一标识符，便于日志追踪和问题定位。

### 实现方案

#### 1. 创建中间件文件

**文件**: `internal/middleware/requestid.go`

```go
package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

const (
    // RequestIDKey 是 Context 中存储 RequestID 的键
    RequestIDKey = "request_id"
    // RequestIDHeader 是 HTTP 头部名称
    RequestIDHeader = "X-Request-ID"
)

// RequestID 返回一个中间件，为每个请求生成或传递 RequestID
func RequestID() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 优先使用上游传递的 RequestID
        requestID := c.GetHeader(RequestIDHeader)
        if requestID == "" {
            requestID = uuid.New().String()
        }

        // 存储到 Context 和 Response Header
        c.Set(RequestIDKey, requestID)
        c.Header(RequestIDHeader, requestID)

        c.Next()
    }
}
```

#### 2. 注册中间件

**文件**: `internal/handler/router.go` (或主路由文件)

```go
// 在其他中间件之前注册
r.Use(middleware.RequestID())
```

#### 3. 在日志中使用

```go
// 在 Handler 或 Service 中获取 RequestID
func (h *SomeHandler) Handle(c *gin.Context) {
    requestID := c.GetString(middleware.RequestIDKey)
    log.Printf("[%s] Processing request", requestID)
}
```

### 改动文件清单

| 文件 | 操作 |
|------|------|
| `internal/middleware/requestid.go` | 新建 |
| `internal/handler/router.go` | 修改 - 注册中间件 |
| `go.mod` | 修改 - 添加 `github.com/google/uuid` 依赖 |

### 验证方式

```bash
# 发送请求，检查响应头
curl -I http://localhost:8182/api/v1/health
# 应返回 X-Request-ID: <uuid>

# 传递已有 RequestID
curl -H "X-Request-ID: test-123" http://localhost:8182/api/v1/health
# 应返回 X-Request-ID: test-123
```

---

## INF-002: 统一错误处理机制

### 目标

减少 Handler 层重复代码，统一错误响应格式，便于日志记录。

### 现有组件

项目已存在 `internal/handler/base.go` 中的 `Response` 结构体。

### 实现方案

#### 1. 定义错误类型

**文件**: `internal/handler/response.go` (新建)

```go
package handler

import (
    "errors"
    "net/http"

    "github.com/gin-gonic/gin"
)

// AppError 应用层错误
type AppError struct {
    Code    int    // HTTP 状态码
    Message string // 用户可见消息
    Detail  string // 详细错误信息（可选）
    Cause   error  // 原始错误（可选）
}

func (e *AppError) Error() string {
    if e.Cause != nil {
        return e.Message + ": " + e.Cause.Error()
    }
    return e.Message
}

// 预定义错误
var (
    ErrBadRequest   = &AppError{Code: http.StatusBadRequest, Message: "请求参数错误"}
    ErrUnauthorized = &AppError{Code: http.StatusUnauthorized, Message: "未授权访问"}
    ErrForbidden    = &AppError{Code: http.StatusForbidden, Message: "禁止访问"}
    ErrNotFound     = &AppError{Code: http.StatusNotFound, Message: "资源不存在"}
    ErrInternal     = &AppError{Code: http.StatusInternalServerError, Message: "系统内部错误"}
)

// NewAppError 创建应用错误
func NewAppError(code int, message string, cause error) *AppError {
    return &AppError{
        Code:    code,
        Message: message,
        Cause:   cause,
    }
}
```

#### 2. 定义响应辅助函数

**文件**: `internal/handler/response.go` (续)

> **注意**: 使用 `base.go` 中已有的 `Response` 结构体

```go
// Success 成功响应
func Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, Response{
        Code:    0,
        Message: "success",
        Data:    data,
    })
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
    c.JSON(http.StatusOK, Response{
        Code:    0,
        Message: message,
        Data:    data,
    })
}

// Fail 失败响应
func Fail(c *gin.Context, err error) {
    var appErr *AppError
    if errors.As(err, &appErr) {
        c.JSON(appErr.Code, Response{
            Code:    1,
            Message: appErr.Message,
        })
        return
    }

    // 未知错误，返回 500
    c.JSON(http.StatusInternalServerError, Response{
        Code:    1,
        Message: "系统内部错误",
    })
}

// FailWithDetail 失败响应（带详情）
func FailWithDetail(c *gin.Context, err error, detail string) {
    var appErr *AppError
    if errors.As(err, &appErr) {
        appErr.Detail = detail
        c.JSON(appErr.Code, Response{
            Code:    1,
            Message: appErr.Message,
        })
        return
    }

    c.JSON(http.StatusInternalServerError, Response{
        Code:    1,
        Message: "系统内部错误",
    })
}

// BadRequest 400 错误
func BadRequest(c *gin.Context, message string) {
    c.JSON(http.StatusBadRequest, Response{
        Code:    1,
        Message: message,
    })
}

// Unauthorized 401 错误
func Unauthorized(c *gin.Context) {
    c.JSON(http.StatusUnauthorized, Response{
        Code:    1,
        Message: "未授权访问",
    })
}

// Forbidden 403 错误
func Forbidden(c *gin.Context) {
    c.JSON(http.StatusForbidden, Response{
        Code:    1,
        Message: "禁止访问",
    })
}

// NotFound 404 错误
func NotFound(c *gin.Context, message string) {
    if message == "" {
        message = "资源不存在"
    }
    c.JSON(http.StatusNotFound, Response{
        Code:    1,
        Message: message,
    })
}

// InternalError 500 错误
func InternalError(c *gin.Context, message string) {
    if message == "" {
        message = "系统内部错误"
    }
    c.JSON(http.StatusInternalServerError, Response{
        Code:    1,
        Message: message,
    })
}
```

#### 3. 增强 Recovery 中间件

> **注意**: 项目已在 `internal/server/server.go` 中使用 `gin.Recovery()`，此处是增强版本，添加日志记录

**文件**: `internal/middleware/recovery.go` (新建)

```go
package middleware

import (
    "log/slog"
    "net/http"
    "runtime/debug"

    "github.com/gin-gonic/gin"
    "github.com/bookandmusic/love-girl/internal/handler"
)

// Recovery 错误恢复中间件（增强版）
func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                // 记录堆栈信息
                stack := string(debug.Stack())
                requestID := c.GetString(RequestIDKey)

                slog.Error("panic recovered",
                    "request_id", requestID,
                    "error", err,
                    "stack", stack,
                )

                // 生产环境不暴露详细错误
                handler.InternalError(c, "系统内部错误")
                c.Abort()
            }
        }()
        c.Next()
    }
}
```

#### 4. 改造现有 Handler

**改造前**:
```go
func (h *UserHandler) GetUser(c *gin.Context) {
    id := c.Param("id")
    user, err := h.service.GetByID(c.Request.Context(), id)
    if err != nil {
        c.JSON(500, gin.H{"code": 1, "message": "获取用户失败"})
        return
    }
    c.JSON(200, gin.H{"code": 0, "message": "success", "data": user})
}
```

**改造后**:
```go
func (h *UserHandler) GetUser(c *gin.Context) {
    id := c.Param("id")
    user, err := h.service.GetByID(c.Request.Context(), id)
    if err != nil {
        handler.Fail(c, handler.NewAppError(http.StatusInternalServerError, "获取用户失败", err))
        return
    }
    handler.Success(c, user)
}
```

### 改动文件清单

| 文件 | 操作 |
|------|------|
| `internal/handler/response.go` | 新建 - 错误类型和辅助函数 |
| `internal/middleware/recovery.go` | 新建 - 增强版 Recovery |
| `internal/handler/*.go` | 修改 - 使用统一响应函数 |

### 验证方式

```bash
# 测试错误响应格式
curl http://localhost:8182/api/v1/users/99999
# 应返回: {"code":1,"message":"资源不存在"}

# 测试成功响应格式
curl http://localhost:8182/api/v1/users/1
# 应返回: {"code":0,"message":"success","data":{...}}
```

---

## INF-003: 结构化日志增强

### 目标

增强现有日志系统，集成 RequestID 追踪，支持 JSON 格式输出，便于日志收集和分析。

### 现有组件

项目已存在 `internal/log/log.go`，使用 `slog` + `tint` 实现彩色文本日志。

### 实现方案

#### 1. 扩展 LogConfig

**文件**: `internal/config/app.go`

```go
// LogConfig 已有字段：Level
// 新增字段：
type LogConfig struct {
    Level  string `mapstructure:"level" validate:"required,oneof=debug info warn error"`
    Format string `mapstructure:"format"` // json, text (新增)
    Output string `mapstructure:"output"` // stdout, file path (新增)
}
```

#### 2. 扩展日志包

**文件**: `internal/log/log.go` (修改现有文件)

```go
package log

import (
    "io"
    "log/slog"
    "os"

    "github.com/bookandmusic/love-girl/internal/config"
    "github.com/lmittmann/tint"
)

// Logger 封装 slog.Logger，便于扩展和实现 io.Writer
type Logger struct {
    *slog.Logger
}

// Init 初始化 slog 日志（支持 JSON 和 Text 格式）
func NewLogger(cfg config.LogConfig) *Logger {
    level := getLoggerLevel(cfg.Level)
    
    var handler slog.Handler
    opts := &slog.HandlerOptions{Level: level, AddSource: true}
    
    // 根据配置选择格式
    if cfg.Format == "json" {
        // JSON 格式 - 用于日志收集
        var writer io.Writer = os.Stdout
        if cfg.Output != "" && cfg.Output != "stdout" {
            file, _ := os.OpenFile(cfg.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
            writer = file
        }
        handler = slog.NewJSONHandler(writer, opts)
    } else {
        // 彩色文本格式 - 现有行为
        handler = tint.NewHandler(os.Stdout, &tint.Options{
            Level:      level,
            AddSource:  true,
            TimeFormat: "2006-01-02 15:04:05",
        })
    }

    logger := slog.New(handler)
    slog.SetDefault(logger)

    return &Logger{logger}
}
```

#### 3. 创建日志中间件

**文件**: `internal/middleware/logging.go`

```go
package middleware

import (
    "log/slog"
    "time"

    "github.com/gin-gonic/gin"
)

// Logging 请求日志中间件（集成 RequestID）
func Logging() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        query := c.Request.URL.RawQuery

        c.Next()

        latency := time.Since(start)
        status := c.Writer.Status()
        requestID := c.GetString(RequestIDKey)

        slog.Info("request",
            "request_id", requestID,
            "method", c.Request.Method,
            "path", path,
            "query", query,
            "status", status,
            "latency", latency.String(),
            "client_ip", c.ClientIP(),
        )
    }
}
```

#### 4. 配置文件更新

**文件**: `data/config.yaml` (示例)

```yaml
log:
  level: info      # debug, info, warn, error
  format: json     # json, text (新增)
  output: stdout   # stdout 或文件路径 (新增)
```

### 改动文件清单

| 文件 | 操作 |
|------|------|
| `internal/config/app.go` | 修改 - 扩展 LogConfig |
| `internal/log/log.go` | 修改 - 支持 JSON 格式 |
| `internal/middleware/logging.go` | 新建 |
| `internal/server/server.go` | 修改 - 注册 Logging 中间件 |
| `internal/**/*.go` | 修改 - 使用 slog 或注入的 Logger 替代 fmt.Printf |

### 验证方式

```bash
# 启动服务后查看日志格式
# JSON 格式示例:
{"time":"2026-03-14T10:00:00Z","level":"INFO","msg":"request","request_id":"abc-123","method":"GET","path":"/api/v1/users","status":200,"latency":"5.23ms"}

# Text 格式示例 (现有):
time=2026-03-14T10:00:00Z level=INFO msg=request request_id=abc-123 method=GET path=/api/v1/users status=200 latency=5.23ms
```

---

## INF-004: 请求超时控制

### 目标

为请求设置全局超时，防止长时间运行的请求占用资源。

### 现有组件

项目已存在 `internal/config/app.go` 中的 `ServerConfig`。

### 实现方案

#### 1. 创建超时中间件

**文件**: `internal/middleware/timeout.go` (新建)

```go
package middleware

import (
    "context"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/bookandmusic/love-girl/internal/handler"
)

// Timeout 超时中间件
func Timeout(timeout time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 创建带超时的 Context
        ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
        defer cancel()

        // 替换请求 Context
        c.Request = c.Request.WithContext(ctx)

        // 使用 channel 监听完成
        finished := make(chan struct{})
        go func() {
            defer close(finished)
            c.Next()
        }()

        select {
        case <-finished:
            // 请求正常完成
        case <-ctx.Done():
            // 请求超时
            c.JSON(http.StatusRequestTimeout, handler.Response{
                Code:    1,
                Message: "请求超时",
            })
            c.Abort()
        }
    }
}
```

#### 2. 扩展 ServerConfig

**文件**: `internal/config/app.go` (修改现有文件)

```go
// ServerConfig 已有字段：Addr, Mode, InternalURL
// 新增字段：
type ServerConfig struct {
    Addr           string        `mapstructure:"addr" validate:"required"`
    Mode           string        `mapstructure:"mode" validate:"required,oneof=debug release test"`
    InternalURL    string        `mapstructure:"internal_url"`
    RequestTimeout time.Duration `mapstructure:"request_timeout"` // 新增：请求处理超时
}
```

**文件**: `data/config.yaml` (示例)

```yaml
server:
  addr: :8182
  mode: debug
  request_timeout: 30s  # 新增
```

#### 3. 注册中间件

**文件**: `internal/server/server.go` (修改现有文件)

```go
// 在路由组上应用超时中间件
apiGroup := engine.Group("/api/v1")
apiGroup.Use(middleware.Timeout(cfg.Server.RequestTimeout))
```

#### 4. 在 Service 层使用 Context

```go
// 确保所有数据库操作使用 Context
func (s *UserService) GetByID(ctx context.Context, id uint64) (*model.User, error) {
    var user model.User
    err := s.db.WithContext(ctx).First(&user, id).Error
    return &user, err
}
```

### 改动文件清单

| 文件 | 操作 |
|------|------|
| `internal/middleware/timeout.go` | 新建 |
| `internal/config/app.go` | 修改 - 扩展 ServerConfig |
| `internal/server/server.go` | 修改 - 注册超时中间件 |
| `internal/service/*.go` | 检查 - 确保 Context 传递 |

### 验证方式

```bash
# 模拟慢请求（需要后端配合添加测试端点）
# 或通过数据库慢查询验证超时中间件生效
```

---

## 实施顺序

建议按以下顺序实施：

```
Phase 1: 基础设施
├── INF-001 RequestID 中间件 (30 min)
│   └── 新建 internal/middleware/requestid.go
├── INF-002 统一错误处理 (2 h)
│   ├── 新建 internal/handler/response.go
│   └── 新建 internal/middleware/recovery.go (增强版)
└── 验证点：运行 ./tools/check.sh all

Phase 2: 日志增强
├── INF-003 结构化日志增强 (2 h)
│   ├── 扩展 internal/config/app.go (LogConfig)
│   ├── 修改 internal/log/log.go (支持 JSON)
│   ├── 新建 internal/middleware/logging.go
│   └── 修改 internal/server/server.go (注册中间件)
└── 验证点：运行 ./tools/check.sh all

Phase 3: 超时控制
├── INF-004 请求超时 (1 h)
│   ├── 新建 internal/middleware/timeout.go
│   ├── 扩展 internal/config/app.go (ServerConfig)
│   └── 修改 internal/server/server.go (注册中间件)
└── 验证点：运行 ./tools/check.sh all
```

---

## 验证清单

完成所有任务后，执行以下验证：

```bash
# 1. 代码检查
cd backend
./tools/check.sh all

# 2. 运行测试
go test ./...

# 3. 启动服务
go run main.go

# 4. 功能验证
# 4.1 RequestID
curl -I http://localhost:8182/api/v1/health
# 检查响应头 X-Request-ID

# 4.2 错误响应格式
curl http://localhost:8182/api/v1/not-exist
# 检查返回 {"code":1,"message":"..."}

# 4.3 日志格式
# 检查控制台输出是否为结构化日志

# 4.4 超时控制
# 模拟慢请求验证
```

---

## 风险与注意事项

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 统一错误处理改造范围大 | 可能引入回归问题 | 逐个 Handler 改造，改造后测试 |
| 超时中间件可能中断正常请求 | 长时间操作被中断 | 设置合理的超时时间（建议 30s+） |
| 日志格式变更 | 现有日志解析失效 | 提前通知运维，提供过渡期 |

---

## 变更记录

| 日期 | 版本 | 变更内容 |
|------|------|----------|
| 2026-03-14 | v1.2 | 完成实施：创建 requestid.go、response.go、recovery.go、logging.go、timeout.go 中间件 |
| 2026-03-14 | v1.1 | 优化：适配项目现有组件，避免重复创建 |
| 2026-03-14 | v1.0 | 初始版本，从 backend-optimization.md 提取并细化 |

---

## 附录：项目现有组件分析

本计划基于对项目现有代码的分析进行了优化，避免重复创建已存在的组件：

### 已验证的现有组件

| 组件 | 位置 | 说明 |
|------|------|------|
| Response 结构体 | `internal/handler/base.go` | 已有 `Code`, `Message`, `Data` 字段 |
| 日志包 | `internal/log/log.go` | 使用 slog + tint，已有日志级别支持 |
| gin.Recovery | `internal/server/server.go` | 已使用 gin.Recovery() |
| ServerConfig | `internal/config/app.go` | 已有 `Addr`, `Mode`, `InternalURL` |
| LogConfig | `internal/config/app.go` | 已有 `Level` 字段 |

### 优化要点

1. **INF-001**: 新建 `requestid.go`，不依赖现有组件
2. **INF-002**: 复用 `handler/base.go` 中的 Response 结构体
3. **INF-003**: 扩展现有 `internal/log/log.go`，添加 JSON 格式支持
4. **INF-004**: 扩展现有 `internal/config/app.go` 中的 ServerConfig