# 构建、测试和 Lint 命令

## 前端 (frontend/)

```bash
# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev

# 类型检查
pnpm type-check

# 构建生产版本
pnpm build

# 预览构建
pnpm preview

# 运行 Lint 并自动修复
pnpm lint

# 代码格式化
pnpm format
```

## 后端 (backend/)

```bash
# 安装 Go 依赖
cd backend && go mod tidy

# 运行单个测试文件
go test -v ./internal/model/place_test.go

# 运行单个测试函数
go test -v -run TestPlaceImage ./internal/model/

# 运行所有测试
go test ./...

# 运行所有检查 (格式化、lint、复杂度、生成代码)
# 确保脚本有执行权限: chmod +x tools/check.sh
./tools/check.sh all

# 单独执行各项检查
./tools/check.sh format      # 代码格式化
./tools/check.sh lint        # 静态检查
./tools/check.sh complexity  # 复杂度检查
./tools/check.sh wire        # 生成依赖注入代码
./tools/check.sh swag        # 生成 API 文档
./tools/check.sh gen         # wire + swag
```
