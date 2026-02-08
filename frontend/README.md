# Love Girl Frontend

这是一个专为情侣设计的浪漫应用程序前端项目。

## 技术栈

- Vue 3
- TypeScript
- Vite
- Tailwind CSS
- Pinia
- Vue Router

## 功能特性

- 响应式设计，适配移动端和桌面端
- 照片墙展示
- 地点记录
- 浪漫动画效果
- PWA支持，可以安装为本地应用

## 开发环境

### 环境要求

- Node.js 版本 ^20.19.0 或 >=22.12.0
- pnpm 包管理器

### 安装依赖

```bash
pnpm install
```

### 启动开发服务器

```bash
pnpm dev
```

### 构建生产版本

```bash
pnpm build
```

### 预览生产构建

```bash
pnpm preview
```

## PWA 支持

该项目支持 Progressive Web App (PWA) 特性，可以安装为本地应用使用。

### 主要特性

- 可安装为桌面/移动应用
- 离线访问能力
- 应用壳模型提升加载速度
- 推送通知支持（待实现）

### 如何安装

1. 在支持 PWA 的浏览器中打开应用（如 Chrome、Edge 等）
2. 点击地址栏中的安装按钮或在菜单中选择"安装应用"
3. 确认安装

安装后应用将作为一个独立的应用程序运行，拥有类似原生应用的体验。

## 项目结构

```
src/
├── assets/           # 静态资源
├── components/       # 公共组件
├── layouts/          # 页面布局
├── mock/             # 模拟数据
├── router/           # 路由配置
├── services/         # API服务
├── stores/           # 状态管理
├── utils/            # 工具函数
├── views/            # 页面视图
└── App.vue           # 根组件
```

## 环境配置

项目支持三种运行环境：
- `mock`: 使用模拟数据
- `dev`: 连接开发服务器
- `prod`: 生产环境配置

通过设置 `VITE_ENV` 环境变量来切换环境。
