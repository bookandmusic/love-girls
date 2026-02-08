# 代码风格指南

## 前端 (Vue 3 + TypeScript)

### 导入顺序 (eslint-plugin-simple-import-sort)

```typescript
// 1. Vue 和框架导入
import { createPinia } from 'pinia'
import { createApp } from 'vue'

// 2. Vue Router
import router from './router'

// 3. 组件导入
import App from './App.vue'

// 4. 路径别名导入
import { useUserStore } from '@/stores/user'

// 5. 相对路径导入
import './assets/main.css'
```

### 格式化 (Prettier)

```json
{
  "semi": false,
  "singleQuote": true,
  "tabWidth": 2,
  "printWidth": 100,
  "trailingComma": "es5",
  "arrowParens": "avoid"
}
```

### 命名规范

- **组件文件名**: PascalCase (如 `UserProfile.vue`, `BaseButton.vue`)
- **组件变量**: PascalCase
- **普通变量/函数**: camelCase
- **常量**: UPPER_SNAKE_CASE
- **Props**: camelCase
- **Store**: PascalCase 后缀 (如 `userStore`, `authStore`)

### 类型定义

- 使用 TypeScript 泛型而非 `any`
- 导出类型定义放在 `src/types/` 目录
- 避免使用 `as any`，优先使用类型断言

### 错误处理

- 使用 Pinia store 管理全局状态错误
- API 请求统一使用 `services/` 封装，捕获异常
- 组件中使用 `try/catch` 处理异步操作

---

## 后端 (Go + Gin)

### 导入顺序

```go
import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "gorm.io/gorm"
    "github.com/bookandmusic/love-girl/internal/model"
)
```

### 代码检查工具

使用 `goimports-reviser` 格式化 (自动处理导入分组和未使用导入)

### 命名规范

- **包名**: 小写单词，简短 (如 `model`, `handler`, `repo`)
- **结构体**: PascalCase
- **函数/变量**: camelCase
- **常量**: CamelCase 或 MixedCaps
- **接口**: er 后缀 (如 `Reader`, `Writer`)
- **错误变量**: `Err` 前缀 (如 `ErrNotFound`)

### 错误处理

- 使用 `errors.New()` 或 `fmt.Errorf()` 创建错误
- 优先返回错误而非打印
- 使用日志库记录错误

### GORM 最佳实践

- 在 `internal/model/` 定义模型
- 使用 `AutoMigrate()` 自动迁移
- 关联使用 `Preload()` 预加载
- 测试使用 SQLite 内存数据库 (`:memory:`)

### 测试规范

- 测试文件命名: `*_test.go`
- 测试函数: `TestXxx(t *testing.T)`
- 使用 `t.Fatalf()` 处理致命错误
- 测试数据库使用 `gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})`
