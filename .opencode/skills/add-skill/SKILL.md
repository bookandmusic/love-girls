---
name: add-skill
description: 创建新的 skill 文件，帮助快速定义和标准化新的自动化工作流
license: MIT
compatibility: opencode
---

# add-skill

创建新的 skill 文件，帮助快速定义和标准化新的自动化工作流。

## 前置条件

- 了解需要创建的 skill 的功能范围
- 熟悉项目结构和已有 skill 的格式

## Skill 文件结构

每个 skill 需要包含以下文件：

```
.opencode/skills/<skill-name>/
└── SKILL.md              # Skill 定义文件
```

## SKILL.md 格式

```markdown
---
name: <skill-name>
description: <一句话描述 skill 功能>
license: MIT
compatibility: opencode
---

# <Skill 名称>

<Sskill 的详细描述>

## 前置条件

- 条件1
- 条件2

## 执行步骤

### 步骤1：<步骤名称>

<步骤说明>

### 步骤2：<步骤名称>

<步骤说明>

## 使用示例

```
请执行 add-skill skill
```

## 注意事项

- 注意事项1
- 注意事项2
```

## 执行步骤

### 1. 收集 Skill 基本信息

向用户询问以下信息：

```
请提供以下信息以创建新的 skill：

1. Skill 名称（英文，用于目录和调用）：
   - 建议使用 kebab-case（如 code-check、git-commit）
   - 用于调用命令：<skill-name>

2. 功能描述（中文，一句话）：
   - 简要说明这个 skill 做什么
   - 示例：代码质量检查、生成规范 Commit Message

3. 使用场景：
   - 什么时候应该调用这个 skill
   - 解决什么问题
```

### 2. 收集执行步骤

```
请描述这个 skill 的执行步骤：

步骤1：<第一步做什么>
步骤2：<第二步做什么>
...（根据需要添加更多步骤）
```

### 3. 收集其他信息

询问可选信息：

- 是否有前置条件需要安装或配置？
- 是否有特殊注意事项？
- 是否需要示例代码？

### 4. 生成 Skill 文件

在 `.opencode/skills/<skill-name>/` 目录下创建 SKILL.md：

```bash
mkdir -p .opencode/skills/<skill-name>
```

根据收集的信息生成 SKILL.md 文件。

### 5. 确认生成结果

```
已创建新的 skill：

- 路径：.opencode/skills/<skill-name>/SKILL.md
- 名称：<skill-name>
- 描述：<description>

调用方式：
  请执行 <skill-name> skill

是否需要修改？
1. 确认完成
2. 修改内容
3. 重新输入
```

## 常用 Skill 模板

### 模板一：代码检查类

```markdown
---
name: <skill-name>
description: 对代码进行检查，确保符合项目规范
license: MIT
compatibility: opencode
---

# <Skill 名称>

对代码进行检查，确保符合项目规范。

## 前置条件

- 安装必要的工具
- 配置正确

## 执行步骤

### 1. 检测检查范围

首先检查本次修改涉及的文件范围。

### 2. 执行检查

运行相应的检查命令：

```bash
# 命令1
# 命令2
```

## 错误处理

- 错误1：解决方法
- 错误2：解决方法

## 注意事项

- 注意事项
```

### 模板二：代码生成类

```markdown
---
name: <skill-name>
description: 自动生成代码文件
license: MIT
compatibility: opencode
---

# <Skill 名称>

自动生成代码文件。

## 前置条件

- 了解需要生成的代码结构
- 准备好模板或配置

## 执行步骤

### 1. 收集参数

向用户收集必要的信息。

### 2. 生成代码

根据模板生成代码文件。

### 3. 验证结果

检查生成的文件是否正确。

## 使用示例

```
请执行 <skill-name> skill
```

## 注意事项

- 注意事项
```

### 模板三：文档生成类

```markdown
---
name: <skill-name>
description: 生成项目文档
license: MIT
compatibility: opencode
---

# <Skill 名称>

生成项目文档。

## 前置条件

- 了解文档结构
- 准备好模板

## 执行步骤

### 1. 收集信息

收集生成文档所需的信息。

### 2. 生成文档

创建文档文件。

### 3. 格式化文档

格式化生成的文档。

## 输出格式

描述生成的文档格式和结构。

## 注意事项

- 注意事项
```

## 注意事项

- Skill 名称使用 kebab-case（短横线分隔）
- 描述使用中文，简明扼要
- 前置条件必须清晰，避免执行时出错
- 执行步骤要有条理，便于 AI 执行
- 包含错误处理，提高容错性
- 生成后建议调用 code-check 检查格式
