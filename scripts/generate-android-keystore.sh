#!/bin/bash

set -e

KEYSTORE_FILE="lovegirl.keystore"
KEY_ALIAS="lovegirl"
VALIDITY=10000
KEY_SIZE=2048
PASSWORD_LENGTH=32

echo "=== Android Keystore 生成工具 ==="
echo ""

if ! command -v keytool &> /dev/null; then
    echo "错误: keytool 未找到，请安装 JDK"
    echo "  Ubuntu/Debian: sudo apt install openjdk-17-jdk"
    echo "  macOS: brew install openjdk@17"
    exit 1
fi

if [ -f "$KEYSTORE_FILE" ]; then
    echo "警告: $KEYSTORE_FILE 已存在"
    read -p "是否覆盖? (y/N): " OVERWRITE
    if [ "$OVERWRITE" != "y" ] && [ "$OVERWRITE" != "Y" ]; then
        echo "已取消"
        exit 0
    fi
    rm -f "$KEYSTORE_FILE"
fi

# 生成随机密码
generate_password() {
    openssl rand -base64 32 | tr -d '/+=' | head -c $PASSWORD_LENGTH
}

KEYSTORE_PASS=$(generate_password)
KEY_PASS=$(generate_password)

echo "正在生成 keystore..."
echo ""

keytool -genkey -v \
    -keystore "$KEYSTORE_FILE" \
    -alias "$KEY_ALIAS" \
    -keyalg RSA \
    -keysize $KEY_SIZE \
    -validity $VALIDITY \
    -storepass "$KEYSTORE_PASS" \
    -keypass "$KEY_PASS" \
    -dname "CN=LoveGirl, OU=Development, O=LoveGirl Team, L=Unknown, ST=Unknown, C=CN" \
    2>/dev/null

if [ $? -eq 0 ]; then
    KEYSTORE_BASE64=$(cat "$KEYSTORE_FILE" | base64 | tr -d '\n')
    
    echo ""
    echo "=========================================="
    echo "              生成成功"
    echo "=========================================="
    echo ""
    echo "Keystore 文件: $(pwd)/$KEYSTORE_FILE"
    echo ""
    echo "GitHub Secrets 配置值:"
    echo ""
    echo "[ANDROID_KEYSTORE_BASE64]"
    echo "$KEYSTORE_BASE64"
    echo ""
    echo "[KEYSTORE_PASS]"
    echo "$KEYSTORE_PASS"
    echo ""
    echo "[KEY_ALIAS]"
    echo "$KEY_ALIAS"
    echo ""
    echo "[KEY_PASS]"
    echo "$KEY_PASS"
    echo ""
    echo "配置步骤:"
    echo "1. 打开 GitHub 仓库页面"
    echo "2. 点击 Settings 标签"
    echo "3. 左侧菜单选择 Secrets and variables → Actions"
    echo "4. 点击 New repository secret 按钮"
    echo "5. 分别添加上述 4 个 Secret"
else
    echo "生成失败"
    exit 1
fi