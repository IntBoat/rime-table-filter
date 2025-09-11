# Rime 字典過濾器 Makefile

# 編譯器設定
GO = go
BUILD_DIR = build
BINARY_NAME = rime-filter

# 編譯標誌
LDFLAGS = -ldflags "-s -w"
BUILD_FLAGS = -trimpath $(LDFLAGS)

# 預設目標
.PHONY: all build clean install test help

all: build

# 編譯程序
build:
	@echo "正在編譯 $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) main.go
	@echo "編譯完成: $(BUILD_DIR)/$(BINARY_NAME)"

# 快速編譯（用於開發）
dev:
	$(GO) build -o $(BINARY_NAME) main.go

# 清理編譯文件
clean:
	@echo "清理編譯文件..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME)
	@rm -f output.ttx
	@rm -f chinese_characters*.txt
	@rm -f filtered_dict*.yaml
	@rm -f missing_chars*.txt
	@echo "清理完成"

# 安裝到系統
install: build
	@echo "安裝 $(BINARY_NAME) 到系統..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "安裝完成"

# 卸載
uninstall:
	@echo "卸載 $(BINARY_NAME)..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "卸載完成"

# 運行測試
test:
	@echo "運行測試..."
	$(GO) test -v ./...

# 格式化代碼
fmt:
	@echo "格式化代碼..."
	$(GO) fmt ./...

# 檢查代碼
vet:
	@echo "檢查代碼..."
	$(GO) vet ./...

# 顯示幫助
help:
	@echo "可用的目標:"
	@echo "  build     - 編譯程序到 build/ 目錄"
	@echo "  dev       - 快速編譯（用於開發）"
	@echo "  clean     - 清理編譯文件和臨時文件"
	@echo "  install   - 安裝到系統 (/usr/local/bin)"
	@echo "  uninstall - 從系統卸載"
	@echo "  test      - 運行測試"
	@echo "  fmt       - 格式化代碼"
	@echo "  vet       - 檢查代碼"
	@echo "  help      - 顯示此幫助"
