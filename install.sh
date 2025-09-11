#!/bin/bash

#
# Copyright (C) 2020-2025. IntBoat <intboat@gmail.com> - All Rights Reserved.
# Unauthorized copying or redistribution of this file in source and binary forms via any medium is strictly prohibited.
# Proprietary and confidential
#
# @author    IntBoat <intboat@gmail.com>
# @copyright 2020-2025. IntBoat <intboat@gmail.com>
# @modified  2025-09-11 16:31:49
#

# Rime 字典過濾器安裝腳本

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 記錄訊息
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 檢查 Go 是否已安裝
check_go() {
    if ! command -v go &> /dev/null; then
        log_error "Go 未安裝，請先安裝 Go 1.21+"
        echo ""
        echo "安裝方法："
        echo "  Ubuntu/Debian: sudo apt install golang-go"
        echo "  macOS: brew install go"
        echo "  其他系統: https://golang.org/dl/"
        exit 1
    fi
    
    # 檢查 Go 版本
    go_version=$(go version | grep -oP 'go\d+\.\d+' | sed 's/go//')
    required_version="1.21"
    
    if [ "$(printf '%s\n' "$required_version" "$go_version" | sort -V | head -n1)" != "$required_version" ]; then
        log_error "Go 版本過低，需要 1.21+，當前版本: $go_version"
        exit 1
    fi
    
    log_success "Go 版本檢查通過: $(go version)"
}

# 檢查 fonttools 是否已安裝
check_fonttools() {
    if ! command -v ttx &> /dev/null; then
        log_warning "fonttools 未安裝，正在安裝..."
        if command -v pip3 &> /dev/null; then
            pip3 install fonttools
        elif command -v pip &> /dev/null; then
            pip install fonttools
        else
            log_error "未找到 pip，請手動安裝 fonttools: pip install fonttools"
            exit 1
        fi
    fi
    
    log_success "fonttools 檢查通過"
}

# 編譯程序
build_program() {
    log_info "正在編譯程序..."
    
    if [ -f "main.go" ]; then
        go build -ldflags "-s -w" -o rime-filter main.go
        log_success "編譯完成"
    else
        log_error "未找到 main.go 文件"
        exit 1
    fi
}

# 安裝到系統
install_system() {
    if [ "$1" = "--system" ]; then
        log_info "安裝到系統..."
        sudo cp rime-filter /usr/local/bin/
        log_success "已安裝到 /usr/local/bin/rime-filter"
    else
        log_info "程序已編譯完成，可執行文件: ./rime-filter"
        echo ""
        echo "如需安裝到系統，請運行:"
        echo "  sudo cp rime-filter /usr/local/bin/"
        echo ""
        echo "或者使用 Makefile:"
        echo "  make install"
    fi
}

# 顯示使用說明
show_usage() {
    echo -e "${BLUE}Rime 字典過濾器安裝腳本${NC}"
    echo ""
    echo "用法: $0 [選項]"
    echo ""
    echo "選項:"
    echo "  --system    安裝到系統 (/usr/local/bin)"
    echo "  --help      顯示此說明"
    echo ""
    echo "範例:"
    echo "  $0           # 編譯程序"
    echo "  $0 --system  # 編譯並安裝到系統"
}

# 主函數
main() {
    if [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
        show_usage
        exit 0
    fi
    
    log_info "開始安裝 Rime 字典過濾器..."
    echo ""
    
    # 檢查依賴
    check_go
    check_fonttools
    echo ""
    
    # 編譯程序
    build_program
    echo ""
    
    # 安裝
    install_system "$1"
    echo ""
    
    log_success "安裝完成！"
    echo ""
    echo "使用方法:"
    echo "  ./rime-filter -h  # 查看幫助"
    echo "  ./rime-filter -f font.ttc  # 基本使用"
}

# 執行主函數
main "$@"
