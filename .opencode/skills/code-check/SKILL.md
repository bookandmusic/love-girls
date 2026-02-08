---
name: code-check
description: 对代码修改后进行质量检查，确保代码符合项目规范。后端执行 bash tools/check.sh all，前端执行 pnpm lint && pnpm type-check && pnpm format
license: MIT
compatibility: opencode
---

# code-check

对代码修改后进行质量检查，确保代码符合项目规范。

## 前置条件

- 后端检查需要安装 Go 工具链
- 前端检查需要安装 Node.js 和 pnpm

## 执行步骤

### 1. 检测修改的文件范围

首先检查本次修改涉及了哪些文件（backend/ 或 frontend/ 目录）

### 2. 后端代码检查

如果修改了后端代码，执行：

```bash
bash tools/check.sh all
```

该命令会依次执行：
- 代码格式化（goimports-reviser）
- 生成依赖注入代码（wire）
- 生成 API 文档（swag）
- 静态检查（golangci-lint）
- 复杂度检查（gocyclo）

### 3. 前端代码检查

如果修改了前端代码，按顺序执行：

```bash
cd frontend
pnpm lint      # 运行 Lint 并自动修复
pnpm type-check # 类型检查
pnpm format    # 代码格式化
```

## 错误处理

### 后端检查失败

1. **goimports-reviser 错误**：检查导入语句是否正确
2. **wire 错误**：检查依赖注入配置
3. **swag 错误**：检查 API 注释语法
4. **golangci-lint 错误**：修复 lint 报告的问题
5. **gocyclo 错误**：考虑拆分复杂函数

### 前端检查失败

1. **lint 错误**：修复 ESLint 报告的问题（通常可自动修复）
2. **type-check 错误**：修复 TypeScript 类型错误
3. **format 错误**：检查 Prettier 配置

## 报告格式

检查完成后，报告以下内容：

1. **检查范围**：本次修改涉及的文件目录
2. **执行结果**：各项检查是否通过
3. **发现问题**：列出需要手动修复的问题
4. **自动修复**：列出已自动修复的问题

## 注意事项

- **每次代码修改完成后必须执行此 skill**
- 只有所有检查都通过后才能提交代码
- 如果无法自动修复的问题，记录下来并告知用户
