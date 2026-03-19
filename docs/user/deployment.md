# 部署指南

## 快速开始

**零配置启动**：直接运行程序即可，无需任何配置文件。

```bash
./love-girl
```

首次启动自动完成：
- 创建数据目录 `./data/`
- 创建 SQLite 数据库
- 创建本地存储目录
- 生成并持久化 JWT Secret

---

## 本地开发

### 环境要求

- Go 1.25+
- Node.js 24+
- SQLite / MySQL / PostgreSQL

### 启动服务

```bash
# 克隆项目
git clone https://github.com/bookandmusic/love-girl.git
cd love-girl

# 启动后端
cd backend
go run main.go

# 启动前端
cd ../frontend
npm install && npm run dev
```

---

## Docker 部署

### 使用 docker run

```bash
docker run -d \
  --name love-girl \
  -p 8182:8182 \
  -v $(pwd)/data:/app/data \
  -e TZ=Asia/Shanghai \
  --restart unless-stopped \
  bookandmusic/love-girl:latest
```

### 使用 Docker Compose

```yaml
services:
  app:
    image: bookandmusic/love-girl:latest
    container_name: love-girl
    ports:
      - "8182:8182"
    volumes:
      - ./data:/app/data
    environment:
      - TZ=Asia/Shanghai
    restart: unless-stopped
```

访问 `http://localhost:8182` 即可使用。

---

## 生产部署

### MySQL + S3

```yaml
services:
  app:
    image: bookandmusic/love-girl:latest
    ports:
      - "8182:8182"
    volumes:
      - ./data:/app/data
    environment:
      - TZ=Asia/Shanghai
      - DATABASE_DRIVER=mysql
      - DATABASE_DSN=user:pass@tcp(mysql:3306)/lovegirl?charset=utf8mb4&parseTime=True
      - STORAGE_BACKEND=s3
      - STORAGE_S3_ENDPOINT=s3.amazonaws.com
      - STORAGE_S3_BUCKET=my-bucket
      - STORAGE_S3_ACCESS_KEY_ID=${AWS_ACCESS_KEY}
      - STORAGE_S3_SECRET_ACCESS_KEY=${AWS_SECRET_KEY}
      - JWT_SECRET=your-very-secure-secret-key-at-least-64-chars
    restart: unless-stopped
    depends_on:
      - mysql

  mysql:
    image: mysql:8
    environment:
      - MYSQL_ROOT_PASSWORD=pass
      - MYSQL_DATABASE=lovegirl
    volumes:
      - mysql_data:/var/lib/mysql

volumes:
  mysql_data:
```

### PostgreSQL + WebDAV

```yaml
services:
  app:
    image: bookandmusic/love-girl:latest
    ports:
      - "8182:8182"
    volumes:
      - ./data:/app/data
    environment:
      - TZ=Asia/Shanghai
      - DATABASE_DRIVER=postgres
      - DATABASE_DSN=postgres://user:pass@postgres:5432/lovegirl?sslmode=disable
      - STORAGE_BACKEND=webdav
      - STORAGE_WEBDAV_ENDPOINT=https://webdav.example.com
      - STORAGE_WEBDAV_BASE_PATH=/uploads
      - STORAGE_WEBDAV_AUTH_USERNAME=admin
      - STORAGE_WEBDAV_AUTH_PASSWORD=secret
    restart: unless-stopped
    depends_on:
      - postgres

  postgres:
    image: postgres:16
    environment:
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=pass
      - POSTGRES_DB=lovegirl
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

---

## 注意事项

- **数据持久化**：确保 `./data` 目录正确挂载到持久化存储
- **JWT 密钥**：生产环境务必手动设置 `JWT_SECRET`，不要使用自动生成的密钥
- **时区设置**：通过 `TZ` 环境变量设置时区
- **配置优先级**：环境变量 > 配置文件 > 默认值，详见 [配置说明](CONFIG.md)