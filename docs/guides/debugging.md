# 调试技巧

## 后端调试

### 使用 Delve 调试器

```bash
# 安装 Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 启动调试服务器
cd backend
dlv debug main.go --headless --listen=:2345 --api-version=2
```

### 日志调试

在开发环境中，使用结构化日志快速定位问题：

```go
import "log/slog"

slog.Info("处理请求", slog.String("path", path), slog.Any("error", err))
```

---

## 前端调试

### 浏览器 DevTools

- **Chrome/Edge**: F12 打开开发者工具
- **Vue DevTools**: Chrome 扩展商店搜索 "Vue.js devtools"

### 调试网络请求

```bash
# 查看 API 请求是否正常
# 1. 打开浏览器 DevTools -> Network 标签
# 2. 过滤 XHR/Fetch 请求
# 3. 检查请求参数和响应数据
```

### HMR 问题排查

当热更新失效时：
1. 检查终端是否有编译错误
2. 尝试重启开发服务器
3. 清除 Vite 缓存：`rm -rf node_modules/.vite`
