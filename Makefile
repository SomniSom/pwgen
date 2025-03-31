# Определяем имя исполняемого файла
APP_NAME = pw.gen

# Путь к исходному коду
SRC_DIR = cmd

# Компилятор Go
GO_CMD = go

# Целевые платформы
PLATFORMS ?= darwin/amd64 linux/amd64 windows/amd64

# Имена исполняемых файлов для каждой платформы
BUILD_TARGETS := $(addsuffix -$(APP_NAME), $(basename $(PLATFORMS)))

# Функция для сборки и упаковки приложения для конкретной платформы
build-%:
	@GOOS=$* GOARCH=${platform##*/} $(GO_CMD) build -o $*.$(APP_NAME) $(SRC_DIR)/main.go

# Упаковка в дистрибутив (создание архива)
package:
	@for platform in $(PLATFORMS); do \
		ARCHIVE_NAME=$(APP_NAME)-$$platform.tar.gz; \
		tar -czf $$ARCHIVE_NAME $$(build-$$platform); \
		rm build-$$platform.$(APP_NAME); \
	done

# Сборка для всех платформ
build_and_package: $(BUILD_TARGETS) package
	@echo "All platforms built and packaged."

# Отдельные цели сборки для каждой платформы
.PHONY: build-darwin-amd64
build-darwin-amd64:
	$(GO_CMD) build -o build-macos/$(APP_NAME) $(SRC_DIR)/main.go

.PHONY: build-linux-amd64
build-linux-amd64:
	$(GO_CMD) build -o build-linux/$(APP_NAME) $(SRC_DIR)/main.go

.PHONY: build-windows-amd64
build-windows-amd64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO_CMD) build -o build-win/$(APP_NAME).exe $(SRC_DIR)/main.go

# Удаление всех собранных файлов
clean:
	rm -rf build-macos build-linux build-win 2>/dev/null || true

# Помощь: отображает доступные команды
help:
	@echo "Available commands:"
	@echo "  make build-darwin-amd64    - Build for macOS"
	@echo "  make build-linux-amd64     - Build for Linux"
	@echo "  make build-windows-amd64   - Build for Windows"
	@echo "  make package               - Package the application into distribution archives"
	@echo "  make clean                 - Remove all built files"

.PHONY: $(BUILD_TARGETS) build_and_package package help
