# 常见问题

## 后端启动失败

**问题**: 端口被占用或配置文件错误

```bash
# 1. 检查端口占用
lsof -i :8182

# 2. 查看详细错误日志
cd backend && go run main.go

# 3. 检查配置文件语法
# backend/data/configs/config.yaml
```

**解决方案**:
- 修改端口或关闭占用进程
- 确认 `config.yaml` 格式正确
- 首次运行会自动创建必要目录

---

## 前端无法连接后端

**问题**: API 请求失败或 CORS 错误

**解决方案**:
1. 确认后端服务已启动 (`http://localhost:8182`)
2. 检查前端 API_BASE_URL 环境变量配置
3. 验证后端 CORS 配置：

```yaml
server:
  cors:
    allow_origins:
      - "http://localhost:5173"
```

---

## 数据库连接失败

**问题**: SQLite 数据库权限或路径错误

```bash
# 确保数据目录有写权限
chmod 755 backend/data

# 检查数据库文件是否存在
ls -la backend/data/love-girl.db
```

---

## 前端构建失败

**问题**: 依赖安装不完整或 Node 版本不匹配

```bash
# 1. 删除依赖重新安装
rm -rf node_modules package-lock.json
pnpm install

# 2. 检查 Node 版本
node -v  # 应在 24.x 以上

# 3. 清除 Vite 缓存
rm -rf node_modules/.vite
```
