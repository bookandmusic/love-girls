# Android 签名密钥配置指南

GitHub Actions 构建 Android APK 需要配置以下 Secrets：

| Secret | 说明 |
|--------|------|
| `ANDROID_KEYSTORE_BASE64` | 签名密钥文件的 Base64 编码 |
| `KEYSTORE_PASS` | 密钥库密码 |
| `KEY_ALIAS` | 密钥别名 |
| `KEY_PASS` | 密钥密码 |

## 生成密钥

### 方法一：本地生成（需要 JDK）

```bash
./scripts/generate-android-keystore.sh
```

### 方法二：Docker 生成（无需本地 JDK）

```bash
./scripts/generate-keystore-docker.sh
```

脚本会自动生成随机密码并输出四个配置值。

## 配置 GitHub Secrets

1. 打开 GitHub 仓库页面
2. 点击 Settings 标签
3. 左侧菜单选择 Secrets and variables → Actions
4. 点击 New repository secret 按钮
5. 分别添加上述 4 个 Secret

## 安全注意事项

- 不要将 keystore 文件提交到 Git 仓库（已在 `.gitignore` 中排除）
- 妥善备份 keystore 文件，丢失将无法更新已发布的应用
- 密码存储在 GitHub Secrets 中，只有仓库管理员可见

## 验证配置

1. 进入 GitHub 仓库的 Actions 页面
2. 选择 Build Client workflow
3. 点击 Run workflow
4. 查看构建日志确认签名成功