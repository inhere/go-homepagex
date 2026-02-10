#!/bin/bash

# Home Dashboard 启动脚本

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_FILE="${1:-$SCRIPT_DIR/config/config.yaml}"

echo "================================"
echo "  Home Dashboard Starter"
echo "================================"
echo ""

# 检查 Go 是否安装
check_go() {
    if ! command -v go &> /dev/null; then
        if [ -f "/usr/local/go/bin/go" ]; then
            export PATH=$PATH:/usr/local/go/bin
        else
            echo "❌ Go is not installed!"
            echo ""
            echo "Please install Go first:"
            echo "  wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz"
            echo "  sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz"
            echo "  export PATH=\$PATH:/usr/local/go/bin"
            echo ""
            exit 1
        fi
    fi
    echo "✓ Go version: $(go version)"
}

# 编译后端
build_backend() {
    echo ""
    echo "Building backend..."
    cd "$SCRIPT_DIR/backend"
    go mod tidy
    go build -o home-dashboard
    echo "✓ Backend built successfully"
}

# 启动服务
start_server() {
    echo ""
    echo "Starting server..."
    echo "Config file: $CONFIG_FILE"
    echo ""
    echo "================================"
    echo ""
    
    cd "$SCRIPT_DIR"
    exec "$SCRIPT_DIR/backend/home-dashboard" "$CONFIG_FILE"
}

# 主流程
main() {
    check_go
    build_backend
    start_server
}

main "$@"
