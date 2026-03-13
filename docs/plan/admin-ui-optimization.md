# 后台管理系统 iOS Glass UI & 交互优化计划 (V2 - 视觉隔离版)

## 1. 概述
本项目旨在将 `love-girl` 后台管理系统重构为具有 iOS 原生感的 “Glass UI”。

**核心原则：**
- **视觉隔离 (Visual Isolation)**：所有样式重构必须限制在后台私有样式中，禁止修改全局公共 CSS（如 `src/assets/main.css` 等），确保前台展示效果不受任何影响。
- **组件安全**：对于前后台共用的组件（如 `GenericDialog`, `BaseIcon`），应通过 Prop 传参或在后台父容器中进行 CSS 覆盖来实现优化，而非直接修改组件内部的全局定义。

---

## 2. 自动化验证流程 (Quality Gate)
每完成一个功能模块的修改，必须在 `frontend/` 目录下顺序执行以下命令：
1. `pnpm lint` - 代码规范检查
2. `pnpm type-check` - TypeScript 类型检查
3. `pnpm format` - 代码格式化
4. `pnpm build` - 生产环境构建测试（确保前台与后台均能正常编译）

---

## 3. 实施时间线与功能分解

### 第一阶段：视觉基础与样式隔离 (预计 1-2 天)
**目标：建立后台独立的样式变量系统。**
- **优化 `admin-theme.css`**:
  - 在 `.admin-layout` 容器下定义后台专用的 CSS 变量（如 `--admin-glass-thin`, `--admin-accent-color`）。
  - 使用特定的类名前缀（如 `.admin-*`）隔离所有样式定义。
  - 引入 SF Pro 备用字体栈，并仅作用于后台根容器。
- **材质定义**: 规范后台专用的圆角（20px）、模糊度（12px-20px）和中性色调。

### 第二阶段：布局与沉浸式导航 (预计 2 天)
**目标：重构后台核心骨架。**
- **`AdminLayout.vue` 重构**:
  - 顶部导航栏改用 `.admin-header` 的 `Ultra-thin` 材质，实现内容穿透。
  - 侧边栏菜单项使用后台私有的选中态样式。
- **移动端交互**:
  - 将后台移动端菜单重构为从底部弹出的 Sheet 样式。

### 第三阶段：基础组件的后台定制 (预计 2-3 天)
**目标：在不影响前台的前提下，优化后台共用组件表现。**
- **`GenericDialog.vue` 后台适配**:
  - 通过后台私有 CSS 覆盖 `GenericDialog` 在管理端弹出时的边框、模糊度和按钮布局，使其具备 iOS Alert/Sheet 的质感。
- **分页器与空状态**:
  - 优化 `AdminPagination.vue` 和 `AdminEmptyState.vue`。

### 第四阶段：管理页面深度重构 (预计 3 天)
**目标：完善各业务页面的 iOS 交互链路。**
- **内容管理交互**: 
  - 列表项 (`MomentItem`, `AlbumItem` 等) 增加按下反馈与圆角优化。
  - 将“添加/修改”表单容器升级为 Sheet 模式。
- **仪表盘与设置**:
  - `DashboardView.vue` 统计卡片重构。
  - `SettingsView.vue` 采用 iOS Grouped List 布局。

---

## 4. 关键文件修改清单 (严格受限)
- `frontend/src/views/admin/styles/admin-theme.css` (**核心修改地**)
- `frontend/src/views/admin/AdminLayout.vue`
- `frontend/src/views/admin/components/*.vue` (后台专用组件)
- `frontend/src/views/admin/content/components/*.vue` (后台列表项)
