# Love Girl

> 一个为情侣设计的全栈应用，记录美好时光，珍藏珍贵回忆

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue-3.5+-4FC08D?style=flat&logo=vue.js)](https://vuejs.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## ✨ 特性

- 🎨 **现代化界面** - 基于 Vue 3 + TypeScript 构建的精美前端
- 🔒 **安全认证** - JWT 令牌认证，保障数据安全
- 💾 **多种存储** - 支持本地存储、S3、WebDAV 等多种存储后端
- ⚡ **零配置启动** - 无需任何配置即可运行，自动生成所有必要配置
- 🔥 **热更新** - 配置文件修改后自动热更新，零停机重启
- 📦 **容器化** - 提供 Docker 部署方案，一键部署
- 🌐 **响应式设计** - 完美适配桌面和移动设备

## 🏗️ 技术架构

![](docs/assets/arch.png)

### 前端技术栈

- **核心框架**: Vue 3.5 + TypeScript + Vite + Pinia + Vue Router
- **样式方案**: TailwindCSS 4
- **HTTP 客户端**: Axios
- **视觉媒体**: vue-easy-lightbox + p5.js + vue3-lottie
- **日期相关**: chinese-days
- **地图**: Leaflet

### 后端技术栈

- **Web 框架**: Gin - 高性能 HTTP Web 框架
- **数据库 ORM**: GORM - 功能强大的 ORM 库
- **配置管理**: Viper - 配置文件解析和管理
- **依赖注入**: Wire - 编译时依赖注入
- **认证**: JWT - JSON Web Token 认证

## 🚀 快速开始

```bash
# Docker 一键启动
docker run -d --name love-girl -p 8182:8182 -v $(pwd)/data:/app/data bookandmusic/love-girl:latest

# 访问
open http://localhost:8182
```

## 📚 文档

- [配置说明](docs/guide/CONFIG.md) - 完整的配置项说明
- [部署指南](docs/guide/DEPLOYMENT.md) - Docker、生产环境部署方案
- [开发指南](AGENTS.md) - 开发环境搭建和代码规范

## ⚙️ 配置

**零配置启动**：直接运行即可，自动创建数据目录、数据库、JWT 密钥。

所有数据存放在 `./data/` 目录：

```
./data/
├── configs/config.yaml    # 配置文件（自动生成）
├── uploads/               # 上传文件
└── love-girl.db           # SQLite 数据库
```

## 🤝 贡献

欢迎贡献代码、报告问题或提出建议！

## 📄 许可证

[MIT License](LICENSE)

## ⭐ Star History

[![Star History Chart](https://api.star-history.com/svg?repos=bookandmusic/love-girl&type=Date)](https://star-history.com/#bookandmusic/love-girl&Date)

---

Made with ❤️ for couples