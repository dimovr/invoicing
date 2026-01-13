# ======================
# Variables
# ======================
BINARY_NAME := fakture
VERSION := 1.0.0
BUILD_DIR := build

DARWIN_ARM64 := $(BUILD_DIR)/darwin_arm64
WINDOWS_AMD64 := $(BUILD_DIR)/windows_amd64

GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOMOD := $(GOCMD) mod
GOFMT := $(GOCMD) fmt

.DEFAULT_GOAL := help

.PHONY: all clean deps fmt lint test build run package help

# ======================
# Core targets
# ======================

all: clean fmt test build package ## Full build cycle

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@$(GOCLEAN)

deps: ## Ensure dependencies
	@$(GOMOD) tidy

fmt: ## Format code
	@$(GOFMT) ./...

lint: ## Lint code (requires golangci-lint)
	@golangci-lint run

test: ## Run tests
	@$(GOCMD) test ./...

# ======================
# Local build & run
# ======================

build: ## Build for local platform
	@echo "Building local binary..."
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)
	@cp -r templates $(BUILD_DIR)/

run: build ## Run locally
ifeq ($(OS),Windows_NT)
	@$(BUILD_DIR)\$(BINARY_NAME).exe
else
	@$(BUILD_DIR)/$(BINARY_NAME)
endif

# ======================
# Packaging
# ======================

package: package-darwin package-windows ## Build all packages

package-darwin: ## Package macOS (Intel + Apple Silicon)
	@echo "Packaging macOS binaries..."

	@mkdir -p $(DARWIN_AMD64) $(DARWIN_ARM64)

	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(DARWIN_AMD64)/$(BINARY_NAME)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(DARWIN_ARM64)/$(BINARY_NAME)

	@cp -r templates $(DARWIN_AMD64)/
	@cp -r templates $(DARWIN_ARM64)/

	@cd $(BUILD_DIR) && \
	zip -r $(BINARY_NAME)_$(VERSION)_darwin_amd64.zip darwin_amd64 && \
	zip -r $(BINARY_NAME)_$(VERSION)_darwin_arm64.zip darwin_arm64

package-windows: ## Package Windows
	@echo "Packaging Windows binary..."

	@mkdir -p $(WINDOWS_AMD64)

	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(WINDOWS_AMD64)/$(BINARY_NAME).exe
	@cp -r templates $(WINDOWS_AMD64)/

	@cd $(BUILD_DIR) && \
	zip -r $(BINARY_NAME)_$(VERSION)_windows_amd64.zip windows_amd64

# ======================
# Help
# ======================

help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "%-20s %s\n", $$1, $$2}'
