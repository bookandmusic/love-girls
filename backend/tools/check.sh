#!/bin/bash

# 判断命令是否存在 
function exists() {
    command -v "$1" >/dev/null 2>&1
}

# 安装 goimports-reviser
install_goimports_reviser() {
    if ! exists "goimports-reviser"; then
        echo "goimports-reviser 未安装，正在安装..."
        go install -v github.com/incu6us/goimports-reviser/v3@latest
    fi
}

# 安装 golangci-lint（用于静态检查）
install_golangci_lint() {
    if ! exists "golangci-lint"; then
        echo "golangci-lint 未安装，正在安装..."
        go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    fi
}

# 安装 gocyclo（用于复杂度检查）
install_gocyclo() {
    if ! exists "gocyclo"; then
        echo "gocyclo 未安装，正在安装..."
        go install -v github.com/fzipp/gocyclo/cmd/gocyclo@latest
    fi
}

# 安装 swag（用于API文档生成）
install_swag() {
    if ! exists "swag"; then
        echo "swag 未安装，正在安装..."
        go install github.com/swaggo/swag/cmd/swag@latest
    fi
}

# 安装 wire（用于依赖注入代码生成）
install_wire() {
    if ! exists "wire"; then
        echo "wire 未安装，正在安装..."
        go install github.com/google/wire/cmd/wire@latest
    fi
}

# 生成依赖注入代码
generate_wire() {
    echo "正在生成依赖注入代码..."
    cd provider && wire && cd ..
}

# 生成API文档
generate_swag() {
    echo "正在生成API文档..."
    swag init
}

# 格式化代码
format_code() {
    echo "正在格式化代码..."
    goimports-reviser -rm-unused -set-alias -format ./...
}

# 静态检查
static_check() {
    echo "正在执行静态检查..."
    golangci-lint run
}

# 复杂度检查
complexity_check() {
    echo "正在执行复杂度检查..."
    gocyclo -over 15 .
}

# 解析命令行参数
case "${1:-all}" in
    "format")
        install_goimports_reviser
        format_code
        ;;
    "lint")
        install_golangci_lint
        static_check
        ;;
    "complexity")
        install_gocyclo
        complexity_check
        ;;
    "wire")
        install_wire
        generate_wire
        ;;
    "swag")
        install_swag
        generate_swag
        ;;
    "gen")
        install_wire
        install_swag
        generate_wire
        generate_swag
        ;;
    "all")
        install_goimports_reviser
        install_golangci_lint
        install_gocyclo
        install_wire
        install_swag
        
        echo "开始执行全部检查..."
        
        echo "1. 格式化代码"
        format_code
        
        echo "2. 生成依赖注入代码"
        generate_wire
        
        echo "3. 生成API文档"
        generate_swag
        
        echo "4. 静态检查"
        static_check
        
        echo "5. 复杂度检查"
        complexity_check
        
        echo "全部检查完成！"
        ;;
    *)
        echo "用法: $0 [format|lint|complexity|wire|swag|gen|all]"
        echo "  format:     代码格式化"
        echo "  lint:       静态检查"
        echo "  complexity: 复杂度检查"
        echo "  wire:       生成依赖注入代码"
        echo "  swag:       生成API文档"
        echo "  gen:        执行代码生成(wire + swag)"
        echo "  all:        执行全部操作"
        exit 1
        ;;
esac