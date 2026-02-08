---
name: git-commit
description: 生成符合 Conventional Commits 规范的 git commit 信息
license: MIT
Compatibility: opencode
---

# git-commit

生成符合 Conventional Commits 规范的 git commit 信息。

## 前置条件

- 安装 Git
- 了解本次代码修改的内容

## Commit Message 格式

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

## Type 类型

| Type | 说明 |
|------|------|
| `feat` | 新功能 |
| `fix` | Bug 修复 |
| `docs` | 文档更新 |
| `style` | 代码格式调整（不影响功能） |
| `refactor` | 代码重构 |
| `perf` | 性能优化 |
| `test` | 测试相关 |
| `chore` | 构建/工具/辅助工具修改 |
| `revert` | 回滚提交 |

## 执行步骤

### 1. 分析修改内容

首先检查本次修改涉及的文件和内容：

```bash
git status          # 查看已修改的文件
git diff --stat     # 查看修改统计
git diff            # 查看具体修改内容
```

### 2. 确定 Type 类型

根据修改内容选择合适的 type：

- **feat**: 新增功能、页面、组件
- **fix**: 修复 bug、错误
- **docs**: 更新 README、文档、注释
- **style**: 格式化代码、修改空格/缩进
- **refactor**: 重构代码、不改变功能
- **perf**: 优化代码性能
- **test**: 添加/修改测试
- **chore**: 修改构建脚本、依赖、配置文件
- **revert**: 回滚之前的提交

### 3. 确定 Scope（可选）

Scope 用于标识修改的影响范围：

**后端常见 Scope**：
- `auth` - 认证相关
- `api` - API 接口
- `db` - 数据库相关
- `model` - 数据模型
- `repo` - 数据访问层
- `service` - 业务逻辑
- `handler` - 处理器
- `middleware` - 中间件
- `config` - 配置相关

**前端常见 Scope**：
- `router` - 路由相关
- `store` - 状态管理
- `component` - 组件相关
- `view` - 页面视图
- `api` - 接口封装
- `style` - 样式相关
- `config` - 构建配置

### 4. 编写 Description

遵循以下规则：

- 使用中文描述
- 简洁明了，不超过 50 字符
- 使用动词开头的祈使句
- 不要大写首字母
- 句末不加句号

**正确示例**：
- `feat: 添加用户登录功能`
- `fix: 修复图片上传失败的 bug`
- `docs: 更新 API 文档`

**错误示例**：
- `feat: 添加了用户登录功能`（使用了"了"）
- `feat: USER LOGIN`（大写）
- `feat: 添加用户登录功能。`（句末有句号）

### 5. 编写 Body（可选）

当描述不足以说明修改时，添加详细的说明：

- 说明为什么需要这个修改
- 说明修改的细节
- 列出重要的改动点
- 每行不超过 72 字符

### 6. 编写 Footer（可选）

用于：

- 标记不兼容的改动（以 `BREAKING CHANGE:` 开头）
- 关联 Issue（以 `Closes #123` 格式）
- 标记贡献者

## 生成 Commit Message

根据分析结果，生成符合规范的 commit message：

```bash
git commit -m "<type>[optional scope]: <description>"
```

## 示例

### 简单提交

```bash
git commit -m "feat: 添加用户注册功能"
```

### 带 Scope 的提交

```bash
git commit -m "feat(auth): 添加用户注册功能"
git commit -m "fix(api): 修复用户查询接口的权限验证"
```

### 详细提交

```bash
git commit -m "feat(user): 新增用户头像上传功能

- 支持 JPG/PNG 格式
- 最大支持 5MB
- 上传成功后自动裁剪为正方形

Closes #123"
```

### 不兼容改动

```bash
git commit -m "refactor(api): 重构用户认证模块

- 改用 JWT 认证方式
- 移除旧的 session 机制

BREAKING CHANGE: API 认证方式变更，需要更新客户端代码"
```

## 注意事项

- 提交前先执行 `git add <files>` 暂存修改的文件
- 确保只提交相关的修改，避免混合提交
- 提交信息使用中文，与项目语言保持一致
- 每次提交只做一件事情
