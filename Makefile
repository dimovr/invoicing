# Variables
BINARY_NAME=fakture
VERSION=1.0.0
BUILD_DIR=build
PLATFORMS=darwin/amd64 darwin/arm64 windows/amd64

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Formatting and linting
GOFMT=$(GOCMD) fmt
GOLINT=golangci-lint

# Default target
.DEFAULT_GOAL := help

.PHONY: all clean deps fmt lint test help package package-darwin package-windows run

all: clean fmt lint test build package ## Run a complete build cycle

clean: ## Clean build directories
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	$(GOCLEAN)

deps: ## Install dependencies
	@echo "Ensuring dependencies are up to date..."
	$(GOMOD) tidy
	$(GOGET) -v ./...

fmt: ## Format source code
	@echo "Formatting code..."
	$(GOFMT) ./...

lint: ## Lint source code
	@echo "Linting code..."
	$(GOLINT) run

test: ## Run tests
	@echo "Running tests..."
	$(GOCMD) test -v ./...

run: build ## Run the application locally
	@echo "Running application..."
	@cp -r templates $(BUILD_DIR)/
	cd $(BUILD_DIR) && ./$(BINARY_NAME)

package: package-darwin package-windows ## Package for all platforms

package-darwin: ## Package for macOS (Intel and Apple Silicon)
	@echo "Packaging for macOS..."
	@mkdir -p $(BUILD_DIR)/darwin_arm64
	
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/darwin_arm64/$(BINARY_NAME) -v
	
	@cp -r templates $(BUILD_DIR)/darwin_arm64/
	@cd $(BUILD_DIR) && zip -r $(BINARY_NAME)_$(VERSION)_darwin_arm64.zip darwin_arm64
	@echo "macOS package created at $(BUILD_DIR)/$(BINARY_NAME)_$(VERSION)_darwin_arm64.zip"

package-windows: ## Package for Windows
	@echo "Packaging for Windows..."
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/windows_amd64/$(BINARY_NAME).exe -v
	
	@cp -r templates $(BUILD_DIR)/windows_amd64/
	
	@cd $(BUILD_DIR) && zip -r $(BINARY_NAME)_$(VERSION)_windows_amd64.zip windows_amd64
	@echo "Windows package created at $(BUILD_DIR)/$(BINARY_NAME)_$(VERSION)_windows_amd64.zip"

help: ## Display available commands
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'